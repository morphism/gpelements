package gpelements

import (
	"bufio"
	"fmt"
	"os"
	"testing"
	"time"
)

func TestProp(t *testing.T) {
	filename := "data/test.tle"
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

	// 0 error SGP4 error at ms=1608048467780: code=1: mean elements, ecc >= 1.0 or ecc < -0.001 or a < 0.95 er
	// SIRIUSSAT-1
	// 1 43595U 98067PG  20344.55055662  .07916525  12514-4  56251-3 0  9997
	// 2 43595  51.6588  97.6693 0009437  42.2825  52.9592 16.39998780132780

	// 1 error SGP4 error at ms=1608048467780: code=1: mean elements, ecc >= 1.0 or ecc < -0.001 or a < 0.95 er
	// SIRIUSSAT-2
	// 1 43596U 98067PH  20344.36554324  .08904132  12494-4  57179-3 0  9995
	// 2 43596  51.6062  99.6513 0006261 312.7483  99.6145 16.40921535132742

	if 2 < len(errs) {
		t.Fatal(errs)
	}

}
