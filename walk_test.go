package gpelements

import (
	"fmt"
	"strings"
	"testing"
)

func TestWalk(t *testing.T) {
	lines := strings.SplitN(testTLE, "\n", 3)
	if len(lines) != 3 {
		t.Fatal(len(lines))
	}

	e, err := ParseTLE(lines[0], lines[1], lines[2])
	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < 5; i++ {
		err := e.Copy().Walk(2, 4)
		if err != nil {
			t.Fatal(err)
		}
		csv, err := e.MarshalCSV()
		if err != nil {
			t.Fatal(err)
		}
		fmt.Printf("%d,%s", i, csv)
	}

}

func TestNextAlpha5Num(t *testing.T) {
	state := int64(3)
	_, next, err := NextAlpha5Num(state)
	if err != nil {
		t.Fatal(err)
	}
	if next == state {
		t.Fatal(next)
	}
}
