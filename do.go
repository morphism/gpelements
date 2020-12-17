package gpelements

import (
	"bufio"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"strings"
)

func (e *Elements) Marshal(how string) (string, error) {
	switch how {
	case "csv":
		return e.MarshalCSV()
	case "xml":
		bs, err := xml.Marshal(e)
		if err != nil {
			return "", err
		}
		return string(bs) + "\n", err
	case "json":
		bs, err := json.Marshal(e)
		if err != nil {
			return "", err
		}
		return string(bs) + "\n", err
	case "kvn":
		return e.MarshalKVN()
	case "tle":
		l0, l1, l2, err := e.MarshalTLE()
		if err != nil {
			return "", err
		}
		return l0 + "\n" + l1 + "\n" + l2 + "\n", nil
	default:
		return "", fmt.Errorf("unknown marshal representation '%s'", how)
	}
}

func MaybeCSV(bs []byte) bool {
	count := 0
	for _, b := range bs {
		if b == ',' {
			count++
		}
	}
	return 4 < count // !
}

func MaybeKVN(bs []byte) bool {
	magic := "CCSDS_OMM_VERS"
	if len(bs) < len(magic) {
		return false
	}
	return magic == string(bs[0:len(magic)])
}

const MinBufferSize = 128

func Do(in io.Reader, bufSize int, f func(Elements) error) error {
	bin := bufio.NewReaderSize(in, bufSize)
	peek, err := bin.Peek(MinBufferSize)
	if err != nil {
		return err
	}

	switch peek[0] {
	case '[', '<': // Just read in the whole thing.
		slurp := func() ([]byte, error) {
			return ioutil.ReadAll(bin)
		}

		var es []Elements

		bs, err := slurp()
		if err != nil {
			return err
		}
		switch peek[0] {
		case '[':
			log.Printf("Detected JSON array input")

			bs, err := slurp()
			if err != nil {
				return err
			}

			bs = []byte(DestringNumbers(string(bs)))
			err = json.Unmarshal(bs, &es)

		case '<':
			log.Printf("Detected XML input")

			list := ElementsList{}
			err = xml.Unmarshal(bs, &list)
			es = list.Es // Hopefully
		}

		for _, e := range es {
			if err = f(e); err != nil {
				// ToDo: Consider toleration.
				return err
			}
		}

	default: // Read line by line
		if peek[0] == '{' {
			err = DoLines(bin, func(s string) error {
				var e Elements
				if err := json.Unmarshal([]byte(s), &e); err != nil {
					return err
				}
				return f(e)
			})
		} else if MaybeKVN(peek) {
			err = DoKVNs(bin, func(s string) error {
				e, _, err := ParseKVN(s)
				if err != nil {
					return err
				}
				return f(*e)
			})

		} else if MaybeCSV(peek) {
			first := true
			err = DoLines(bin, func(s string) error {
				if first && strings.Contains(s, ",EPOCH,") {
					return nil
				}
				first = false
				e, _, err := ParseCSV(s)
				if err != nil {
					return err
				}
				return f(*e)
			})
		} else {
			err = DoTLEs(bin, 3, func(lines []string) error {
				e, err := ParseTLE(lines[0], lines[1], lines[2])
				if err != nil {
					return err
				}
				return f(*e)
			})
		}
	}

	return err
}
