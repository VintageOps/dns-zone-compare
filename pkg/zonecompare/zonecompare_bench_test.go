package zonecompare

import (
	"testing"
)

func BenchmarkZoneCompare(b *testing.B) {
	options := Opts{
		Domain:      "mail.example.com",
		Origin:      "../../examples/zone1",
		Destination: "../../examples/zone2",
		Ignore:      []string{"ignore1", "ignore2"},
		Notfound:    false,
		Strict:      true,
		Found:       true,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ZoneCompare(options)
	}
}
