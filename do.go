package gpelements

import (
	"bufio"
	"bytes"
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

func Do(in io.Reader, f func(Elements) error) error {
	// ToDo: Try harder to avoid reading in the entire input!
	bs, err := ioutil.ReadAll(in)
	if err != nil {
		return err
	}

	if len(bs) == 0 {
		return nil
	}

	var es []Elements

	switch bs[0] {
	case '[': // JSON representing and array of Elements.
		log.Printf("Detected JSON array input")
		bs = []byte(DestringNumbers(string(bs)))
		err = json.Unmarshal(bs, &es)

	case '{': // One elements in JSON per line.
		log.Printf("Detected one-per-line JSON input")
		err = DoLines(bufio.NewReader(bytes.NewReader(bs)), func(s string) error {
			var e Elements
			if err := json.Unmarshal([]byte(s), &e); err != nil {
				return err
			}
			es = append(es, e)
			return nil
		})
	case '<': // XML
		log.Printf("Detected XML input")
		list := ElementsList{}
		err = xml.Unmarshal(bs, &list)
		es = list.Es // Hopefully

	default:
		in := bufio.NewReader(bytes.NewReader(bs))
		if MaybeKVN(bs) {
			log.Printf("Detected KVN input")
			err = DoKVNs(in, func(s string) error {
				e, _, err := ParseKVN(s)
				if err != nil {
					return err
				}
				es = append(es, *e)
				return nil
			})

		} else if MaybeCSV(bs) {
			log.Printf("Detected CSV input")
			first := true
			err = DoLines(in, func(s string) error {
				if first && strings.Contains(s, ",EPOCH,") {
					return nil
				}
				first = false
				e, _, err := ParseCSV(s)
				if err != nil {
					return err
				}
				es = append(es, *e)
				return nil
			})
		} else {
			log.Printf("Detected TLE input")
			err = DoTLEs(in, 3, func(lines []string) error {
				e, err := ParseTLE(lines[0], lines[1], lines[2])
				if err != nil {
					return err
				}
				es = append(es, *e)
				return nil
			})
		}
	}

	if err != nil {
		return err
	}

	for i, e := range es {
		if err := f(e); err != nil {
			return fmt.Errorf("%s on set %d", err, i)
		}
	}

	return nil
}
