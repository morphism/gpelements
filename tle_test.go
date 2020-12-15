package gpelements

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestTLEMany(t *testing.T) {
	in, err := os.Open("test.tle")
	if err != nil {
		t.Fatal(err)
	}
	defer in.Close() // Ignore error.
	r := bufio.NewReader(in)

	f := func(lines []string) error {

		e, err := ParseTLE(lines[0], lines[1], lines[2])
		if err != nil {
			return err
		}

		check := make([]string, 3)
		check[0], check[1], check[2], err = e.MarshalTLE()
		if err != nil {
			return err
		}

		for i, line := range lines {
			if !sameTLELine(line, check[i]) {
				return fmt.Errorf("different\n'%s'\n'%s'\n\n", line, check[i])
			}
		}

		return nil

	}

	if err := DoTLEs(r, 3, f); err != nil {
		t.Fatal(err)
	}

}

func BenchmarkTLEParse(b *testing.B) {
	tle := testTLE

	lines := strings.SplitN(tle, "\n", 3)
	if len(lines) != 3 {
		b.Fatal(len(lines))
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if _, err := ParseTLE(lines[0], lines[1], lines[2]); err != nil {
			b.Fatal(err)
		}
	}

}
