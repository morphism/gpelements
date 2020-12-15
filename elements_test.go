package gpelements

import (
	"log"
	"strconv"
	"testing"
)

func TestCopy(t *testing.T) {
	e0 := NewElements()
	e0.Name = "homer"
	e1 := e0.Copy()
	e1.Name = "bart"

	if e0.Name == e1.Name {
		t.Fatal(e0.Name)
	}
}

func TestEncode(t *testing.T) {
	var (
		n  = int64(1234567)
		s  = strconv.FormatInt(n, 10)
		id = NoradCatId(s)
		e  = id.Encode()
		d  = e.Decode()
	)

	log.Printf("DEBUG %s %s %s %s", s, id, e, d)
}
