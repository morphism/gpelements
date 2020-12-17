package gpelements

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
)

func ParseTLE(line0, line1, line2 string) (*Elements, error) {
	var (
		e   = NewElements()
		err error
		s   string
	)

	e.Name = strings.TrimRight(line0, " ")

	if s, err = tleExtract(line1, 3, 7); err != nil {
		return nil, wrapErrf(err, "NoradCatId")
	}
	e.NoradCatId = NewNoradCatId(s).Decode()

	if s, err = tleExtract(line1, 8, 8); err != nil {
		return nil, wrapErrf(err, "ClassificationType")
	}
	e.ClassificationType = s

	{
		if s, err = tleExtract(line1, 10, 11); err != nil {
			return nil, wrapErrf(err, "LaunchYear")
		}
		var year int
		if year, err = strconv.Atoi(s); err != nil {
			return nil, wrapErrf(err, "LaunchYear")
		}
		if year < 56 {
			year = 2000 + year
		} else {
			year = 1900 + year
		}
		e.LaunchYear = year
	}

	if s, err = tleExtract(line1, 12, 14); err != nil {
		return nil, wrapErrf(err, "LaunchNum")
	}
	s = strings.TrimLeft(s, "0")
	if 0 < len(s) {
		if e.LaunchNum, err = strconv.Atoi(s); err != nil {
			return nil, wrapErrf(err, "LaunchNum: %s", line1)
		}
	}

	if s, err = tleExtract(line1, 15, 17); err != nil {
		return nil, wrapErrf(err, "LaunchPiece")
	}
	s = strings.TrimRight(s, " ")
	e.LaunchPiece = s

	e.Id = fmt.Sprintf("%04d-%03d%s", e.LaunchYear, e.LaunchNum, e.LaunchPiece)

	{

		var year int
		if s, err = tleExtract(line1, 19, 20); err != nil {
			return nil, wrapErrf(err, "Epoch year")
		}
		if year, err = strconv.Atoi(s); err != nil {
			return nil, wrapErrf(err, "Epoch year")
		}
		if year < 56 {
			year = 2000 + year
		} else {
			year = 1900 + year
		}

		loc := time.FixedZone("UTC", 0)
		t := time.Date(year, 1, 0, 0, 0, 0, 0, loc)

		var day float64
		if s, err = tleExtract(line1, 21, 32); err != nil {
			return nil, wrapErrf(err, "Epoch year")
		}
		if day, err = strconv.ParseFloat(s, 64); err != nil {
			return nil, wrapErrf(err, "Epoch day")
		}
		hours := day * 24
		t0 := Time(t.Add(time.Duration(hours * float64(time.Hour))))
		e.Epoch = &t0
	}

	if s, err = tleExtract(line1, 34, 43); err != nil {
		return nil, wrapErrf(err, "MeanMotionDot")
	}
	s = strings.TrimLeft(s, " ")
	if e.MeanMotionDot, err = tleParseFloat(s); err != nil {
		return nil, wrapErrf(err, "MeanMotionDot '%s'", s)
	}

	if s, err = tleExtract(line1, 45, 52); err != nil {
		return nil, wrapErrf(err, "MeanMotionDDot")
	}
	switch s[0] {
	case ' ':
		s = "0." + s[1:]

	case '-':
		s = "-0." + s[1:]
	}
	if e.MeanMotionDDot, err = tleParseFloat(s); err != nil {
		return nil, wrapErrf(err, "MeanMotionDDot '%s'", s)
	}

	if s, err = tleExtract(line1, 54, 61); err != nil {
		return nil, wrapErrf(err, "BStar")
	}
	switch s[0] {
	case ' ':
		s = "0." + s[1:]

	case '-':
		s = "-0." + s[1:]
	}
	if e.BStar, err = tleParseFloat(s); err != nil {
		return nil, wrapErrf(err, "BStar '%s'", s)
	}

	e.EphemerisType = 0

	if s, err = tleExtract(line1, 65, 68); err != nil {
		return nil, wrapErrf(err, "ElementSet")
	}
	s = strings.TrimLeft(s, " ")
	if e.ElementSet, err = strconv.Atoi(s); err != nil {
		return nil, wrapErrf(err, "ElementSet")
	}

	// Skip checksum?

	if s, err = tleExtract(line2, 3, 7); err != nil {
		return nil, wrapErrf(err, "NoradCatId")
	}
	if id := NewNoradCatId(s).Decode(); id != e.NoradCatId {
		return nil, fmt.Errorf("NoradCatId disagreement: '%s' != '%s'", string(id), string(e.NoradCatId))
	}

	if s, err = tleExtract(line2, 9, 16); err != nil {
		return nil, wrapErrf(err, "Inclination")
	}
	s = strings.TrimLeft(s, " ")
	if e.Inclination, err = tleParseFloat(s); err != nil {
		return nil, wrapErrf(err, "Inclination '%s'", s)
	}

	if s, err = tleExtract(line2, 18, 25); err != nil {
		return nil, wrapErrf(err, "RightAscention")
	}
	s = strings.TrimLeft(s, " ")
	if e.RightAscension, err = tleParseFloat(s); err != nil {
		return nil, wrapErrf(err, "RightAscention '%s'", s)
	}

	if s, err = tleExtract(line2, 27, 33); err != nil {
		return nil, wrapErrf(err, "Eccentricity")
	}
	switch s[0] {
	case ' ':
		s = "0." + s[1:]

	case '-':
		s = "-0." + s[1:]
	default:
		s = "0." + s
	}
	if e.Eccentricity, err = tleParseFloat(s); err != nil {
		return nil, wrapErrf(err, "Eccentricity '%s'", s)
	}

	if s, err = tleExtract(line2, 35, 42); err != nil {
		return nil, wrapErrf(err, "ArgOfPericenter")
	}
	s = strings.TrimLeft(s, " ")
	if e.ArgOfPericenter, err = tleParseFloat(s); err != nil {
		return nil, wrapErrf(err, "ArgOfPericenter '%s'", s)
	}

	if s, err = tleExtract(line2, 44, 51); err != nil {
		return nil, wrapErrf(err, "MeanAnomaly")
	}
	s = strings.TrimLeft(s, " ")
	if e.MeanAnomaly, err = tleParseFloat(s); err != nil {
		return nil, wrapErrf(err, "MeanAnomaly '%s'", s)
	}

	if s, err = tleExtract(line2, 53, 63); err != nil {
		return nil, wrapErrf(err, "MeanMotion")
	}
	s = strings.TrimLeft(s, " ")
	if e.MeanMotion, err = tleParseFloat(s); err != nil {
		return nil, wrapErrf(err, "MeanMotion '%s'", s)
	}

	if s, err = tleExtract(line2, 64, 68); err != nil {
		return nil, wrapErrf(err, "RevAtEpoch '%s'", s)
	}
	s = strings.TrimLeft(s, " ")
	if e.RevAtEpoch, err = tleParseFloat(s); err != nil {
		return nil, wrapErrf(err, "RevAtEpoch '%s'", s)
	}

	return e, nil
}

func wrapErrf(err error, format string, args ...interface{}) error {
	msg := fmt.Sprintf(format, args...)
	return fmt.Errorf("%s: %s", msg, err)
}

func tleParseFloat(s string) (float64, error) {
	s0 := s[0:1] // Maybe a -.
	s = strings.Replace(s[1:], "-", "e-", 1)
	s = strings.Replace(s, "+", "e", 1)
	return strconv.ParseFloat(s0+s, 64)
}

func tleExtract(line string, col0, col1 int) (string, error) {
	var (
		from = col0 - 1
		to   = col1
	)
	if len(line) <= to {
		return "", fmt.Errorf("len('%s') == %d <= %d", line, len(line), to)
	}

	return line[from:to], nil
}

func year56(year int) int {
	if year < 1955 {
		return 0
	}
	if year < 2000 {
		return year - 1900
	}
	return year - 2000
}

func (e *Elements) MarshalTLE() (line0, line1, line2 string, err error) {
	line0 = fmt.Sprintf("% -24s", e.Name)

	// Slowly (very) and clearly (hopefully).
	s := "1 "

	s += fmt.Sprintf("% -5s%s ", e.NoradCatId.Encode(), e.ClassificationType)

	if err := e.UseInternationalDesignator(); err != nil {
		return "", "", "", err
	}
	s += fmt.Sprintf("%02d%03d%-3s ", year56(e.LaunchYear), e.LaunchNum, e.LaunchPiece)

	var (
		year = time.Time(*e.Epoch).Year()
		loc  = time.FixedZone("UTC", 0)
		t0   = time.Date(year, 1, 0, 0, 0, 0, 0, loc) // day == 0 is strange.
		d    = time.Time(*e.Epoch).Sub(t0)
		day  = float64(d.Nanoseconds()) / 1000 / 1000 / 1000 / 60 / 60 / 24
	)

	s += fmt.Sprintf("%02d%012.8f ", year56(year), day)

	{ // 9
		x := fmt.Sprintf("%7.8f", e.MeanMotionDot)
		x = strings.Replace(x, "0.", ".", 1)
		if x[0] != '-' {
			x = " " + x
		}
		s += x + " "
	}

	{ // 10
		s0, oops := formatWeird(e.MeanMotionDDot)
		if err != nil {
			err = oops
			return
		}
		s += s0 + " "
	}

	{ // 11
		s0, oops := formatWeird(e.BStar)
		if err != nil {
			err = oops
			return
		}
		s += "" + s0 + " "
	}

	// 12
	s += "0 "

	// 13
	s += fmt.Sprintf("% 4d", e.ElementSet)

	s += checksum(s)

	line1 = s

	s = "2 " +
		fmt.Sprintf("% -5s ", e.NoradCatId.Encode()) +
		fmt.Sprintf("%8.4f ", e.Inclination) +
		fmt.Sprintf("%8.4f ", e.RightAscension)

	{
		s0 := strconv.FormatFloat(e.Eccentricity, 'f', 8, 64)
		if !strings.HasPrefix(s0, "0.") {
			err = fmt.Errorf("Eccentricity %f is too big", e.Eccentricity)
			return
		}
		s0 += "0000000"
		s += s0[2:9] + " "
	}

	s += fmt.Sprintf("%8.4f ", e.ArgOfPericenter)
	s += fmt.Sprintf("%8.4f ", e.MeanAnomaly)
	s += fmt.Sprintf("%11.8f", e.MeanMotion)

	if 99999 < e.RevAtEpoch {
		err = fmt.Errorf("%f is too many revolutions per day", e.RevAtEpoch)
		return
	}

	s += fmt.Sprintf("% 6d", int(e.RevAtEpoch))[1:]

	// ToDo: Checksum

	s += checksum(s)

	line2 = s

	return
}

func formatWeird(x float64) (string, error) {
	s := strconv.FormatFloat(x*10, 'e', 4, 64)

	s = strings.ReplaceAll(s, ".", "")

	s = strings.ReplaceAll(s, "e", "")
	// s = strings.ReplaceAll(s, "0+0", "0-0")

	e := s[len(s)-1:]
	s = s[0:len(s)-2] + e

	if 0 <= x {
		s = " " + s
	}

	return s, nil
}

func checksum(s string) string {
	// "The checksums for each line are calculated by adding all
	// numerical digits on that line, including the line
	// number. One is added to the checksum for each negative sign
	// (-) on that line. All other non-digit characters are
	// ignored."
	//
	// https://en.wikipedia.org/wiki/Two-line_element_set

	sum := 0
	for _, c := range s {
		inc := 0
		switch c {
		case '0':
			inc = 0
		case '1':
			inc = 1
		case '2':
			inc = 2
		case '3':
			inc = 3
		case '4':
			inc = 4
		case '5':
			inc = 5
		case '6':
			inc = 6
		case '7':
			inc = 7
		case '8':
			inc = 8
		case '9':
			inc = 9
		case '-':
			inc = 1
		}
		sum += inc
	}

	return fmt.Sprintf("%d", sum%10)
}

func DoTLEs(r *bufio.Reader, group int, f func(lines []string) error) error {

	var (
		tle []string
		i   int
	)
	for {
		line, err := r.ReadString('\n')
		line = strings.TrimRight(line, "\n\r")
		if 0 < len(line) {
			if tle == nil {
				tle = make([]string, group)
			}
			tle[i%group] = line
			if (i+1)%group == 0 {
				if err := f(tle); err != nil {
					return err
				}
				tle = nil
			}
			i++
		}
		if err == io.EOF {
			break
		}
	}

	return nil
}
