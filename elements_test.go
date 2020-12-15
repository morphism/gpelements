package gpelements

import "testing"

func TestCopy(t *testing.T) {
	e0 := NewElements()
	e0.Name = "homer"
	e1 := e0.Copy()
	e1.Name = "bart"

	if e0.Name == e1.Name {
		t.Fatal(e0.Name)
	}
}
