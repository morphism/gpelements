package gpelements

import (
	"fmt"
	"strconv"
	"strings"
)

// NoradCatId exists in an attempt to handle numbers or the new,
// temporary "Alpha-5 scheme" for NORAD catalog identifiers.
//
// Note that JSON serializations have explicit syntax for strings vs
// numbers, so a NoradCatId will try to (de)serialize as a number.
//
// This type also has support for encoding/decoding in a large base in
// order to fit large numbers in five characters in a TLE.
type NoradCatId string

func NewNoradCatId(s string) NoradCatId {
	return NoradCatId(strings.TrimSpace(s))
}

func (s *NoradCatId) UnmarshalJSON(bs []byte) error {
	maybeQuoted := string(bs)
	maybeQuoted = strings.Trim(maybeQuoted, `"`)
	*s = NoradCatId(maybeQuoted)
	return nil
}

func (s NoradCatId) MarshalJSON() ([]byte, error) {
	if _, err := strconv.Atoi(string(s)); err == nil && s[0] != '0' {
		// error calling MarshalJSON for type gpelements.NoradCatId:
		// invalid character '0' after top-level value.
		//
		// That message is why we check for s[0] != '0'.
		return []byte(s), nil
	}
	return []byte(fmt.Sprintf(`"%s"`, string(s))), nil
}

func Base(n int64, alphabet string) string {
	if n < 100000 {
		return strconv.FormatInt(n, 10)
	}

	var (
		b   = int64(len(alphabet))
		acc = ""
	)
	for 0 < n {
		var (
			r = n % b
			c = alphabet[r : r+1]
		)
		acc = acc + c
		if n/b < 0 {
			break
		}
		n /= b
	}
	var (
		pad   = 5 - 1 - len(acc)
		zeros = "00000"
	)
	if pad < 0 {
		pad = 0
	}
	return "z" + zeros[0:pad] + acc

}

func pow(n, e int64) int64 {
	acc := int64(1)
	for i := int64(0); i < e; i++ {
		acc *= n
	}
	return acc
}

func Unbase(s string, alphabet string) (int64, error) {
	if len(s) == 0 || s[0] != 'z' {
		return 0, fmt.Errorf("'%s' not in base '%s'", s, alphabet)
	}
	s = s[1:]
	acc := int64(0)
	base := int64(len(alphabet))
	for i := int64(len(s)) - 1; 0 <= i; i-- {
		c := s[i : i+1]
		r := int64(strings.Index(alphabet, c))
		if r < 0 {
			return 0, fmt.Errorf("'%s' not in base '%s'", s, alphabet)
		}
		acc += r * pow(base, int64(i))
	}
	return acc, nil

}

var noradCatBase = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func (s NoradCatId) Encode() NoradCatId {
	n, err := strconv.ParseInt(string(s), 10, 64)
	if err != nil {
		return s
	}
	return NoradCatId(Base(n, noradCatBase))
}

func (s NoradCatId) Decode() NoradCatId {
	d, err := Unbase(string(s), noradCatBase)
	if err != nil {
		return s
	}
	// ToDo: Error on overflow here?
	return NoradCatId(strconv.FormatInt(d, 10))
}
