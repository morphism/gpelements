package gpelements

import (
	"encoding/xml"
	"io/ioutil"
	"testing"
)

func TestXML(t *testing.T) {
	filename := "omm.xml"

	bs, err := ioutil.ReadFile("omm.xml")
	if err != nil {
		t.Skipf("Couldn't read %s: %s", filename, err)
	}

	var es ElementsList
	if err := xml.Unmarshal(bs, &es); err != nil {
		t.Fatal(err)
	}

	for _, e := range es.Es {
		if _, err := e.MarshalCSV(); err != nil {
			t.Fatal(err)
		}
	}
}
