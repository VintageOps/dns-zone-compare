package zonecompare

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"reflect"
	"testing"
)

func readFileIntoVariable(filename string) (string, error) {
	var jreport = make(rrMapJzone)

	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		log.Fatal(err)
		return "", err
	}

	err = json.Unmarshal([]byte(content), &jreport)
	if err != nil {
		log.Fatal("Error:", err)
	}

	jsonContent, err := json.Marshal(jreport)
	if err != nil {
		log.Fatal("Error:", err)
	}

	return string(jsonContent), nil
}
func Test_loadMap(t *testing.T) {
	type args struct {
		filename string
		options  Opts
	}
	tests := []struct {
		name string
		args args
		want zoneMap
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := loadMap(tt.args.filename, tt.args.options); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("loadMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_sortDNSSlice(t *testing.T) {
	type args struct {
		x []dnsEntry
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sortDNSSlice(tt.args.x)
		})
	}
}

func Test_removeTabs(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := removeTabs(tt.args.str); got != tt.want {
				t.Errorf("removeTabs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_diffDnsSlices(t *testing.T) {
	type args struct {
		x []dnsEntry
		y []dnsEntry
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := diffDnsSlices(tt.args.x, tt.args.y); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("diffDnsSlices() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_diffDnsEntries(t *testing.T) {
	type args struct {
		origin      []dnsEntry
		destination []dnsEntry
		options     Opts
	}
	tests := []struct {
		name string
		args args
		want map[string][]string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := diffDnsEntries(tt.args.origin, tt.args.destination, tt.args.options); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("diffDnsEntries() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_sliceDnsEntryString(t *testing.T) {
	type args struct {
		x []dnsEntry
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := sliceDnsEntryString(tt.args.x); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("sliceDnsEntryString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_deepSliceAndSort(t *testing.T) {
	type args struct {
		records []dnsEntry
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := deepSliceAndSort(tt.args.records); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("deepSliceAndSort() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_findRepeat(t *testing.T) {
	type args struct {
		slice []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := findRepeat(tt.args.slice); got != tt.want {
				t.Errorf("findRepeat() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_flattenDnsEntrySlice(t *testing.T) {
	type args struct {
		entry []dnsEntry
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := flattenDnsEntrySlice(tt.args.entry); got != tt.want {
				t.Errorf("flattenDnsEntrySlice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_logReport(t *testing.T) {
	type args struct {
		jreport     map[string]map[string][]jzoneDiff
		reportType  string
		name        string
		dnsType     string
		origin      []dnsEntry
		destination []dnsEntry
		options     Opts
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logReport(tt.args.jreport, tt.args.reportType, tt.args.name, tt.args.dnsType, tt.args.origin, tt.args.destination, tt.args.options)
		})
	}
}

func Test_logAndReport(t *testing.T) {
	type args struct {
		status                string
		name                  string
		curType               string
		entrySliceOrigin      []dnsEntry
		entrySliceDestination []dnsEntry
		options               Opts
	}
	tests := []struct {
		name string
		args args
		want []jzoneDiff
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := logAndReport(tt.args.status, tt.args.name, tt.args.curType, tt.args.entrySliceOrigin, tt.args.entrySliceDestination, tt.args.options); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("logAndReport() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestZoneCompare(t *testing.T) {
	type args struct {
		options Opts
	}
	tests := []struct {
		name string
		args args
		want rrMapJzone
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ZoneCompare(tt.args.options); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ZoneCompare() = %v, want %v", got, tt.want)
			}
		})
	}
}
