package gpelements

import (
	"bufio"
	"fmt"
	"os"
	"testing"
	"time"
)

func TestProp(t *testing.T) {
	filename := "test.tle"
	in, err := os.Open(filename)
	if err != nil {
		t.Skipf("Couldn't read %s: %s", filename, err)
	}
	defer in.Close() // Ignore error.

	var (
		r = bufio.NewReader(in)

		now = time.Now().UTC() // ToDo: Use a constant.

		errs  = make([]error, 0, 32)
		count = 0
	)

	f := func(lines []string) error {
		count++

		e, err := ParseTLE(lines[0], lines[1], lines[2])
		if err != nil {
			return err
		}

		o, err := e.SGP4()
		if err != nil {
			return err
		}
		s, err := Prop(o, now)
		if err != nil {
			err = fmt.Errorf("error %s\n%s\n%s\n%s\n", err, lines[0], lines[1], lines[2])
			errs = append(errs, err)
		}

		if false {
			fmt.Printf("%#v\n", s)
		}
		return nil

	}

	if err := DoTLEs(r, 3, f); err != nil {
		t.Fatal(err)
	}

	for i, err := range errs {
		fmt.Printf("%d %s\n", i, err)
	}

	if 0 < len(errs) {
		t.Fatal(errs)
	}

}
