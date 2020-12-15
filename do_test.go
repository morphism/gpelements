package gpelements

import (
	"bufio"
	"log"
	"os"
	"testing"
)

func TestDo(t *testing.T) {
	f := func(filename string) error {
		log.Printf("Reading %s", filename)
		in, err := os.Open(filename)
		if err != nil {
			return err
		}
		defer in.Close()
		r := bufio.NewReader(in)

		return Do(r, func(e Elements) error {
			// fmt.Printf("%v\n", e)
			return nil
		})
	}

	for _, filename := range []string{"data/test.jsonarray"} {
		if err := f(filename); err != nil {
			t.Fatal(err)
		}
	}

}
