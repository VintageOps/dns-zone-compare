package zonecompare

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"reflect"
	"sort"
	"strings"
	"time"

	"github.com/VintageOps/dns-zone-compare/pkg/utils"
	"github.com/miekg/dns"
)

type Opts struct {
	Domain           string
	Origin           string
	IgnoreTTL        bool
	Ignore           []string
	Deep             []string
	DeepAll          bool
	Found            bool
	Notfound         bool
	Strict           bool
	Json             bool
	PrettyJSON       bool
	Text             bool
	CountText        int
	Destination      string
	Labelorigin      string
	Labeldestination string
}

type dnsEntry struct {
	Hdr     dns.RR_Header
	content string
}

func (entry *dnsEntry) String() string {
	return removeTabs(fmt.Sprintf("%s %s", entry.Hdr.String(), entry.content))
}

type zoneMap map[string]map[string][]dnsEntry

/*
	"mail.example.com.": {
	  "A": [
	    {
	      "differences": {
	        "zone1": [
	          "mail.example.com. 3600 IN A  192.0.2.6",
	          "mail.example.com. 3600 IN A  192.0.2.7"
	        ]
	      },
          "originalRecords":{
		    "zone1": [
			  "mail.example.com. 3600 IN A  192.0.2.4",
			  "mail.example.com. 3600 IN A  192.0.2.6",
			  "mail.example.com. 3600 IN A  192.0.2.7"
		    ],
		    "zone2": [
			  "mail.example.com. 3600 IN A  192.0.2.4"
		    ]
          },
          "status": "different",
	    }
	  ]
	}
*/
/*
zoneDiff
"status":
"origin": dns entry
"destination": dns entry
"differences": map[string]string
*/
type jzoneDiff map[string]interface{}

/*
[name][type]jzoneDiff
*/
type rrMapJzone map[string]map[string][]jzoneDiff

func toDNS(val reflect.Value, entry *dnsEntry) {
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		typeF := val.Type().Field(i)
		if typeF.Type == reflect.TypeOf(dns.RR_Header{}) {
			entry.Hdr = field.Interface().(dns.RR_Header)
		} else {
			switch field.Kind() {
			case reflect.Struct:
				toDNS(field, entry)
			default:
				_entry := strings.Trim(fmt.Sprintf("%v", field.Interface()), "[]")
				entry.content = strings.TrimSpace(entry.content + fmt.Sprintf(" %s", _entry))
			}
		}
	}
}

func checkIfOriginFirstLine(zoneFilePath string) bool {
	fd, err := os.Open(zoneFilePath)
	utils.FatalOnErr(err)
	defer func() {
		err := fd.Close()
		if err != nil {
			utils.FatalOnErr(err)
		}
	}()
	scanner := bufio.NewScanner(fd)
	for scanner.Scan() {
		var line string = scanner.Text()
		if strings.HasPrefix(line, "$ORIGIN") && strings.HasSuffix(line, ".") {
			return true
		}
		break
	}
	return false
}

func loadMap(filename string, options Opts) zoneMap {
	var z *dns.ZoneParser
	var fd *os.File
	var err error
	var rrSlice []dns.RR

	zone := make(zoneMap)

	if strings.Contains(filename, ":") {
		if options.Domain == "" {
			utils.FatalOnErr(fmt.Errorf("when using <address>:<port>, like with %s, "+
				"you must specify a valid fully qualified domain options (e.g. example.com)", filename))
		} else {
			if !strings.HasSuffix(options.Domain, ".") {
				options.Domain += "."
			}
		}
		server := filename
		// TODO: make this timeout an option?
		transfer := new(dns.Transfer)
		transfer.ReadTimeout = time.Duration(10 * time.Second)

		msg := new(dns.Msg)
		msg.SetAxfr(options.Domain)

		axfrChan, err := transfer.In(msg, server)
		if err != nil {
			log.Fatalln(err.Error())
		}

		for x := range axfrChan {
			for _, y := range x.RR {
				y.Header().Rdlength = 0
				rrSlice = append(rrSlice, y)
			}
		}
	} else {
		if options.Domain == "" {
			if !checkIfOriginFirstLine(filename) {
				utils.FatalOnErr(fmt.Errorf("no domain was specified and the provided zonefile %s does "+
					"not have the first line defined with $ORIGIN and ending with a dot('.')\n"+
					"Please either Specify a domain using domain parameters or make sure that $ORIGIN is defined in "+
					"the zonefile (e.g. '$ORIGIN example.com.')", filename))
			}
		} else {
			if !strings.HasSuffix(options.Domain, ".") {
				options.Domain += "."
			}
		}
		fd, err = os.Open(filename)
		utils.FatalOnErr(err)
		defer func() {
			err := fd.Close()
			if err != nil {
				return
			}
		}()
		z = dns.NewZoneParser(fd, options.Domain, "")
		utils.FatalOnErr(z.Err())
		for rr, ok := z.Next(); ok; rr, ok = z.Next() {
			rr.Header().Rdlength = 0
			rrSlice = append(rrSlice, rr)
		}
	}

	for _, rr := range rrSlice {
		var entry dnsEntry
		val := reflect.ValueOf(rr).Elem()
		if options.IgnoreTTL {
			rr.Header().Ttl = 3600
		}
		toDNS(val, &entry)
		if zone[rr.Header().Name] == nil {
			zone[rr.Header().Name] = make(map[string][]dnsEntry)
		}
		zone[rr.Header().Name][val.Type().Name()] =
			append(zone[rr.Header().Name][val.Type().Name()], entry)

	}
	return zone

}

// DNS entry Slice sorter, for proper comparison
func sortDNSSlice(x []dnsEntry) {
	if len(x) > 1 {
		sort.Slice(x, func(i, j int) bool { return x[i].String() < x[j].String() })
	}
}

func removeTabs(str string) string {
	return strings.ReplaceAll(str, "\t", " ")
}

func diffDnsSlices(x, y []dnsEntry) []string {
	var retDiff []string
	for _, _x := range x {
		found := false
		for _, _y := range y {
			if _x.String() == _y.String() {
				found = true
				break
			}
		}
		if !found {
			retDiff = append(retDiff, _x.String())
		}
	}
	return retDiff

}

func logPrint(options Opts, params ...interface{}) {
	if options.Json {
		if options.CountText >= 1 {
			// Both json and Text specified
			log.Println(params...)
		}
	} else {
		log.Println(params...)
	}
}

func diffDnsEntries(origin, destination []dnsEntry, options Opts) map[string][]string {
	retDiff := make(map[string][]string)
	var diffOrigin, diffDestination []string

	// Not sure if we should flatten the difference instead of keeping separated full DNS entries
	diffOrigin = diffDnsSlices(origin, destination)
	diffDestination = diffDnsSlices(destination, origin)
	if len(diffOrigin) > 0 {
		retDiff[options.Labelorigin] = diffOrigin
		// log.Println("different", options.origin, flattenDnsEntrySlice(origin))
		logPrint(options, options.Labelorigin, "different", retDiff[options.Labelorigin])
	}
	if len(diffDestination) > 0 {
		retDiff[options.Labeldestination] = diffDestination
		logPrint(options, options.Labeldestination, "different", retDiff[options.Labeldestination])
	}

	return retDiff
}

func sliceDnsEntryString(x []dnsEntry) []string {
	var returnSlice []string
	for _, entry := range x {
		returnSlice = append(returnSlice, entry.String())
	}
	return returnSlice
}

func deepSliceAndSort(records []dnsEntry) []string {
	var stringSliced []string
	for _, split := range records {
		stringSliced = append(stringSliced, strings.Fields(split.content)...)
	}
	sort.Strings(stringSliced)
	return stringSliced
}

func findRepeat(slice []string) string {
	var output []string
	var outputStr string
	check := make(map[string]struct{})
	for _, entry := range slice {
		if _, found := check[entry]; !found {
			check[entry] = struct{}{}
		} else {
			output = append(output, entry)
		}
	}
	if len(output) > 0 {
		outputStr = strings.Join(output, " ")
	}
	return outputStr
}

func sliceStringDiff(a, b []string) []string {
	mb := make(map[string]struct{}, len(b))
	for _, x := range b {
		mb[x] = struct{}{}
	}
	var diff []string
	for _, x := range a {
		if _, found := mb[x]; !found {
			diff = append(diff, x)
		}
	}
	return diff
}

func logAndReport(status string,
	name string,
	curType string,
	entrySliceOrigin []dnsEntry,
	entrySliceDestination []dnsEntry,
	options Opts) []jzoneDiff {
	var deep = make(map[string]struct{}, len(options.Deep))
	var sliceStringOrigin, sliceStringDestination []string
	var jzoneDiffSlice []jzoneDiff

	for _, i := range options.Deep {
		deep[strings.ToLower(i)] = struct{}{}
	}
	sliceStringOrigin = sliceDnsEntryString(entrySliceOrigin)
	sliceStringDestination = sliceDnsEntryString(entrySliceDestination)
	switch status {
	case "found", "notfound":
		{
			logPrint(options, options.Labelorigin, status, flattenDnsEntrySlice(entrySliceOrigin))
			jzoneDiffSlice = append(jzoneDiffSlice,
				jzoneDiff{
					"status": status,
					"originalRecords": map[string][]string{
						options.Labelorigin: sliceStringOrigin}})
		}
	case "different":
		{
			_jzoneDiff := jzoneDiff{"status": status,
				"originalRecords": map[string][]string{
					options.Labelorigin:      sliceStringOrigin,
					options.Labeldestination: sliceStringDestination},
			}
			if _, found := deep[strings.ToLower(curType)]; found || options.DeepAll {
				_differences := make(map[string][]string)
				_repeats := make(map[string][]string)
				originSlice := deepSliceAndSort(entrySliceOrigin)
				destinationSlice := deepSliceAndSort(entrySliceDestination)
				oriDestSlice := sliceStringDiff(originSlice, destinationSlice)
				destOriSlice := sliceStringDiff(destinationSlice, originSlice)
				oriRepeats := findRepeat(originSlice)
				destRepeats := findRepeat(destinationSlice)
				if len(oriDestSlice) > 0 {
					_differences[options.Labelorigin] = []string{removeTabs(entrySliceOrigin[0].Hdr.String()) + strings.Join(oriDestSlice, " ")}
					logPrint(options, options.Labelorigin, status, removeTabs(entrySliceOrigin[0].Hdr.String())+strings.Join(oriDestSlice, " "))
				}
				if len(destOriSlice) > 0 {
					_differences[options.Labeldestination] = []string{removeTabs(entrySliceDestination[0].Hdr.String()) + strings.Join(destOriSlice, " ")}
					logPrint(options, options.Labeldestination, status, removeTabs(entrySliceDestination[0].Hdr.String())+strings.Join(destOriSlice, " "))
				}
				if len(_differences) > 0 {
					_jzoneDiff["differences"] = _differences
				}
				if len(oriRepeats) > 0 {
					_repeats[options.Labelorigin] = []string{removeTabs(entrySliceDestination[0].Hdr.String()) + oriRepeats}
					logPrint(options, options.Labelorigin, "repeated", removeTabs(entrySliceDestination[0].Hdr.String())+oriRepeats)
				}
				if len(destRepeats) > 0 {
					_repeats[options.Labeldestination] = []string{removeTabs(entrySliceDestination[0].Hdr.String()) + destRepeats}
					logPrint(options, options.Labeldestination, "repeated", removeTabs(entrySliceDestination[0].Hdr.String())+destRepeats)
				}
				if len(_repeats) > 0 {
					_jzoneDiff["repeats"] = _repeats
				}
				if len(destOriSlice) == 0 && len(oriDestSlice) == 0 {
					// there's a repetition or the entries are "deeply" the same
					// if they are the same, change it as found (or don't add it if found reporting is disabled)
					if options.Found && len(_repeats) == 0 {
						_jzoneDiff["status"] = "found"
						delete(_jzoneDiff, options.Labeldestination)
					} else if len(_repeats) == 0 {
						_jzoneDiff = jzoneDiff{}
					}
				}

			} else {

				diffEntries := diffDnsEntries(entrySliceOrigin, entrySliceDestination, options)
				// this happens when strict and the order is different
				if len(diffEntries) > 0 {
					_jzoneDiff["differences"] = diffEntries
				} else if options.Strict {
					_jzoneDiff["differences"] = []string{"Wrong Order"}
				}
			}
			if len(_jzoneDiff) > 0 {
				jzoneDiffSlice = append(jzoneDiffSlice, _jzoneDiff)
			}
		}
	}

	return jzoneDiffSlice
}

func flattenDnsEntrySlice(entry []dnsEntry) string {
	flat := entry[0].String()
	for _, ent := range entry[1:] {
		flat += " " + ent.content
	}
	return flat
}

func logReport(jreport map[string]map[string][]jzoneDiff, reportType string, name string, dnsType string, origin []dnsEntry, destination []dnsEntry, options Opts) {
	if jreport[name][dnsType] == nil {
		jreport[name] = make(map[string][]jzoneDiff)
	}

	switch reportType {
	case "notfound":
		jreport[name][dnsType] = logAndReport(reportType, name, dnsType, origin, []dnsEntry{}, options)
	case "different", "found":
		_logAndReport := logAndReport(reportType, name, dnsType, origin, destination, options)
		if (reportType == "different" && len(_logAndReport) > 0) || reportType == "found" {
			jreport[name][dnsType] = _logAndReport
		}
		// this happens when they are "different" but after deep inspection they are the same.
		if len(jreport[name][dnsType]) == 0 {
			delete(jreport[name], dnsType)
		}
		if len(jreport[name]) == 0 {
			delete(jreport, name)
		}
	default:
		log.Fatalln("We shouldn't reach this point")
	}
}

func ZoneCompare(options Opts) rrMapJzone {
	origin := loadMap(options.Origin, options)
	destination := loadMap(options.Destination, options)
	var jreport = make(rrMapJzone)
	var ignore = make(map[string]struct{}, len(options.Ignore))

	for _, i := range options.Ignore {
		ignore[strings.ToLower(i)] = struct{}{}
	}

	for name, dnsTypes := range origin {

		for dnsType, _ := range dnsTypes {

			if _, found := ignore[strings.ToLower(dnsType)]; found {
				continue
			}
			if destination[name][dnsType] == nil {

				if !options.Notfound {
					logReport(jreport, "notfound", name, dnsType, origin[name][dnsType],
						destination[name][dnsType], options)
				}
				continue
			}
			if !options.Strict {
				sortDNSSlice(origin[name][dnsType])
				sortDNSSlice(destination[name][dnsType])
			}
			if !reflect.DeepEqual(destination[name][dnsType], origin[name][dnsType]) {
				logReport(jreport, "different", name, dnsType, origin[name][dnsType],
					destination[name][dnsType], options)
			} else if options.Found {
				logReport(jreport, "found", name, dnsType, origin[name][dnsType],
					destination[name][dnsType], options)
			}
		}
	}
	return jreport

}
