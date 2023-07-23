package zonecompare

import (
	"encoding/json"
	"io"
	"log"
	"os"
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

func TestZoneCompare(t *testing.T) {
	// Test case 1: Example scenario
	options := Opts{
		Domain:      "mail.example.com",
		Origin:      "../../examples/zone1",
		Destination: "../../examples/zone2",
		Ignore:      []string{"ignore1", "ignore2"},
		Notfound:    false,
		Strict:      true,
		Found:       true,
	}

	result := ZoneCompare(options)
	expectedOutput, err := readFileIntoVariable("expected_output")
	if err != nil {
		log.Fatal("Could not read rawOutput file.")
	}

	if err != nil {
		log.Fatal("Could not read marshalledOutput file.")
	}

	if result != expectedOutput {
		//t.Errorf("Test case 1 failed. Expected: %s, got: %s", expectedOutput, result)
		t.Errorf("Test case 2 failed. Expected: %s, got: %s", expectedOutput, result)
	}

	// Test case 2: Another scenario
	options = Opts{
		Domain:      "mail.example.com",
		Origin:      "../../examples/zone1",
		Destination: "../../examples/zone2",
		Ignore:      []string{"ignore3", "ignore4"},
		Notfound:    true,
		Strict:      false,
		Found:       false,
	}

	result = ZoneCompare(options)

	if result != expectedOutput {
		t.Errorf("Test case 2 failed. Expected: %s, got: %s", expectedOutput, result)
	}

}
