package gpelements

import (
	"bufio"
	"os"
	"testing"
)

func TestKVNs(t *testing.T) {
	filename := "kvn.txt"

	in, err := os.Open(filename)
	if err != nil {
		t.Skipf("Couldn't open %s: %v", filename, err)
	}
	defer in.Close() // Ignore error.
	r := bufio.NewReader(in)

	i := 0
	err = DoKVNs(r, func(s string) error {
		e, _, err := ParseKVN(s)
		if err != nil {
			return err
		}
		if _, err := e.MarshalCSV(); err != nil {
			return err
		}
		i++
		return nil
	})

	if err != nil {
		t.Fatal(err)
	}

	if i == 0 {
		t.Fatal(i)
	}
}

func TestParseInternationalDesignator(t *testing.T) {
	y, n, p, err := ParseInternationalDesignator("1998-067A")
	if err != nil {
		t.Fatal(err)
	}
	if y != 1998 {
		t.Fatal(y)
	}
	if n != 67 {
		t.Fatal(n)
	}
	if p != "A" {
		t.Fatal(p)
	}
}
