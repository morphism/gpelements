package gpelements

import (
	"encoding/xml"
	"fmt"
	"strings"
	"time"
)

type Time time.Time

func NewTime(t time.Time) *Time {
	if t.IsZero() {
		t = time.Now().UTC()
	}
	then := Time(t)
	return &then
}

func (t *Time) String() string {
	if t == nil {
		return ""
	}
	return time.Time(*t).String()
}

func ParseTime(s, layout string) (*Time, error) {
	t, err := time.Parse(s, layout)
	if err != nil {
		return nil, err
	}
	t0 := Time(t)
	return &t0, err
}

func (c *Time) Format(layout string) string {
	return time.Time(*c).Format(layout)
}

func (c *Time) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	tokens := []xml.Token{start}
	s := "null"
	if c != nil {
		s = time.Time(*c).Format(KVNTimeFormat)
	}
	tokens = append(tokens, xml.CharData(s))
	tokens = append(tokens, xml.EndElement{start.Name})
	for _, t := range tokens {
		if err := e.EncodeToken(t); err != nil {
			return err
		}
	}
	return nil
}

func (c *Time) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var v string
	d.DecodeElement(&v, &start)
	if len(v) == 0 {
		return nil
	}
	parse, err := time.Parse(KVNTimeFormat, v)
	if err != nil {
		return err
	}
	*c = Time(parse)
	return nil
}

func (c *Time) UnmarshalJSON(data []byte) error {
	s := strings.Trim(string(data), `"`)
	t, err := time.Parse(KVNTimeFormat, s)
	if err != nil {
		return err
	}
	*c = Time(t)
	return nil
}

func (c *Time) MarshalJSON() ([]byte, error) {
	if c == nil {
		return []byte(`"null"`), nil
	}
	s := fmt.Sprintf(`"%s"`, time.Time(*c).Format(KVNTimeFormat))
	return []byte(s), nil
}
