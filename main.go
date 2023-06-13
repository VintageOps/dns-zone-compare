package main

import (
	"encoding/json"
	"fmt"
	"github.com/miekg/dns"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"reflect"
	"sort"
	"strings"
)

type opts struct {
	domain      string
	origin      string
	destination string
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
	      "status": "different",
	      "zone1": [
	        "mail.example.com. 3600 IN A  192.0.2.4",
	        "mail.example.com. 3600 IN A  192.0.2.6",
	        "mail.example.com. 3600 IN A  192.0.2.7"
	      ],
	      "zone2": [
	        "mail.example.com. 3600 IN A  192.0.2.4"
	      ]
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

func fatalOnErr(err error) {
	if err != nil {
		log.Fatalln(err.Error())
	}
}

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

func loadMap(filename string, options opts) zoneMap {
	var z *dns.ZoneParser
	var fd *os.File
	var err error
	var rrSlice []dns.RR

	zone := make(zoneMap)

	fd, err = os.Open(filename)
	fatalOnErr(err)
	defer func() {
		fd.Close()
	}()

	z = dns.NewZoneParser(fd, options.domain, "")
	fatalOnErr(z.Err())
	for rr, ok := z.Next(); ok; rr, ok = z.Next() {
		rrSlice = append(rrSlice, rr)
	}

	for _, rr := range rrSlice {
		var entry dnsEntry
		val := reflect.ValueOf(rr).Elem()
		toDNS(val, &entry)
		if zone[rr.Header().Name] == nil {
			zone[rr.Header().Name] = make(map[string][]dnsEntry)
		}
		zone[rr.Header().Name][val.Type().Name()] =
			append(zone[rr.Header().Name][val.Type().Name()], entry)

	}
	return zone

}

// DNS entry Slice sorter, for proper comparision
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

func diffDnsEntries(origin, destination []dnsEntry, options opts) map[string][]string {
	retDiff := make(map[string][]string)
	var diffOrigin, diffDestination []string

	// Not sure if we should flatten the difference instead of keeping separated full DNS entries
	diffOrigin = diffDnsSlices(origin, destination)
	diffDestination = diffDnsSlices(destination, origin)
	if len(diffOrigin) > 0 {
		retDiff[options.origin] = diffOrigin
		// log.Println("different", options.origin, flattenDnsEntrySlice(origin))
		log.Println("different", options.origin, retDiff[options.origin])
	}
	if len(diffDestination) > 0 {
		retDiff[options.destination] = diffDestination
		log.Println("different", options.destination, retDiff[options.destination])
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

func logAndReport(status string,
	name string,
	curType string,
	entrySliceOrigin []dnsEntry,
	entrySliceDestination []dnsEntry,
	options opts) []jzoneDiff {
	// TODO: deep diff
	var sliceStringOrigin, sliceStringDestination []string
	var jzoneDiffSlice []jzoneDiff
	sliceStringOrigin = sliceDnsEntryString(entrySliceOrigin)
	sliceStringDestination = sliceDnsEntryString(entrySliceDestination)
	switch status {
	case "found", "notfound":
		{
			log.Println(status, flattenDnsEntrySlice(entrySliceOrigin))
			jzoneDiffSlice = append(jzoneDiffSlice, jzoneDiff{"status": status, options.origin: sliceStringOrigin})
		}
	case "different":
		{
			diffEntries := diffDnsEntries(entrySliceOrigin, entrySliceDestination, options)
			jzoneDiffSlice = append(jzoneDiffSlice, jzoneDiff{"status": status,
				options.origin:      sliceStringOrigin,
				options.destination: sliceStringDestination,
				"differences":       diffEntries})
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

func zoneCompare(origin, destination zoneMap, options opts) string {
	jreport := make(rrMapJzone)
	for name, dnsTypes := range origin {
		for dnsType, _ := range dnsTypes {
			if jreport[name][dnsType] == nil {
				jreport[name] = make(map[string][]jzoneDiff)
			}

			if destination[name][dnsType] == nil {
				jreport[name][dnsType] = logAndReport("notfound", name,
					dnsType,
					origin[name][dnsType],
					[]dnsEntry{},
					options)
				continue
			}
			sortDNSSlice(origin[name][dnsType])
			sortDNSSlice(destination[name][dnsType])
			if !reflect.DeepEqual(destination[name][dnsType], origin[name][dnsType]) {
				jreport[name][dnsType] = logAndReport("different", name,
					dnsType,
					origin[name][dnsType],
					destination[name][dnsType],
					options)
			} else {
				jreport[name][dnsType] = logAndReport("found", name,
					dnsType,
					origin[name][dnsType],
					destination[name][dnsType],
					options)
			}
		}
	}
	jzoneOutput, err := json.Marshal(jreport)
	fatalOnErr(err)
	return string(jzoneOutput)

}

func main() {
	// TODO: use urfave flags
	// var options = opts{domain: "example.com"}
	options := opts{}
	app := &cli.App{
		Name:  "zonecompare",
		Usage: "compare dns zones",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "domain",
				Usage:       "domain to compare, default none",
				Destination: &options.domain,
			},
		},
		Action: func(c *cli.Context) error {
			options.origin = c.Args().Get(0)
			options.destination = c.Args().Get(1)
			origin := loadMap(options.origin, options)
			destination := loadMap(options.destination, options)
			fmt.Println(zoneCompare(origin, destination, options))
			return cli.Exit("", 0)
		},
	}
	err := app.Run(os.Args)
	fatalOnErr(err)
}
