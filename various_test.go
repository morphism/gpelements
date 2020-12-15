package gpelements

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"strings"
	"testing"
)

var (
	testTLE = `ISS (ZARYA)             
1 25544U 98067A   20262.67636574  .00000241  00000+0  12514-4 0  9990
2 25544  51.6432 245.8351 0000884 104.2674 236.9442 15.48952759246507`
)

func sameTLELine(line, check string) bool {

	if strings.TrimSpace(line) == strings.TrimSpace(check) {
		return true
	}

	// '1 08820U 76039A   20262.42169440 -.00000003  00000-0  00000+0 0  9998'
	// '1 08820U 76039A   20262.42169440 -.00000003  00000-0  00000-0 0  9999'

	// Recompute checksum.
	line = line[0 : len(line)-1]
	line = line + checksum(line)

	return strings.TrimSpace(line) == strings.TrimSpace(check)
}

func TestTLEOne(t *testing.T) {
	var (
		tle = testTLE
		e   *Elements
		err error
		js  []byte
	)

	t.Run("TLE", func(t *testing.T) {

		lines := strings.SplitN(tle, "\n", 3)
		if len(lines) != 3 {
			t.Fatal(len(lines))
		}

		if e, err = ParseTLE(lines[0], lines[1], lines[2]); err != nil {
			t.Fatal(err)
		}

		if js, err = json.MarshalIndent(e, "", "  "); err != nil {
			t.Fatal(err)
		}

		check := make([]string, 3)
		check[0], check[1], check[2], err = e.MarshalTLE()
		if err != nil {
			t.Fatal(err)
		}

		for i, line := range lines {
			if !sameTLELine(line, check[i]) {
				err = fmt.Errorf("different\n'%s'\n'%s'\n\n", line, check[i])
				t.Fatal(err)
			}
		}
	})

	t.Run("CSV", func(t *testing.T) {

		csv, err := e.MarshalCSV()
		if err != nil {
			t.Fatal(err)
		}

		check, _, err := ParseCSV(csv)
		if err != nil {
			t.Fatal(err)
		}

		jsc, err := json.MarshalIndent(check, "", "  ")
		if err != nil {
			t.Fatal(err)
		}

		if string(jsc) != string(js) {
			t.Fatal(string(jsc))
		}
	})

	t.Run("KVN", func(t *testing.T) {

		kvn, err := e.MarshalKVN()
		if err != nil {
			t.Fatal(err)
		}

		check, _, err := ParseKVN(kvn)
		if err != nil {
			t.Fatalf("error %s\n%s", err, kvn)
		}

		jsc, err := json.MarshalIndent(check, "", "  ")
		if err != nil {
			t.Fatal(err)
		}

		if string(jsc) != string(js) {
			t.Fatal(string(jsc))
		}
	})

	t.Run("XML", func(t *testing.T) {
		bs, err := xml.MarshalIndent(e, "", "  ")
		if err != nil {
			t.Fatal(err)
		}

		check := NewElements()
		if err := xml.Unmarshal(bs, check); err != nil {
			t.Fatal(err)
		}

		jsc, err := json.MarshalIndent(check, "", "  ")
		if err != nil {
			t.Fatal(err)
		}

		if string(jsc) != string(js) {
			t.Fatal(string(js) + "\n" + string(bs) + "\n" + string(jsc))
		}
	})
}
