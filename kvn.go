package gpelements

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
)

var KVNTimeFormat = "2006-01-02T15:04:05.999999999"

func (e *Elements) MarshalKVN() (string, error) {
	f := `CCSDS_OMM_VERS = 2.0
CREATION_DATE  = %s
ORIGINATOR     = %s

OBJECT_NAME    = %s
OBJECT_ID      = %s
CENTER_NAME    = EARTH
REF_FRAME      = TEME
TIME_SYSTEM    = UTC
MEAN_ELEMENT_THEORY = SGP/SGP4

EPOCH          = %s
MEAN_MOTION    = %g
ECCENTRICITY   = %g
INCLINATION    = %g
RA_OF_ASC_NODE = %g
ARG_OF_PERICENTER = %g
MEAN_ANOMALY   = %g

EPHEMERIS_TYPE = %d
CLASSIFICATION_TYPE = %s
NORAD_CAT_ID   = %s
ELEMENT_SET_NO = %d
REV_AT_EPOCH   = %d
BSTAR          = %e
MEAN_MOTION_DOT = %e
MEAN_MOTION_DDOT = %e
`
	return fmt.Sprintf(f,
		e.CreationDate,
		e.Originator,
		e.Name,
		e.Id,
		e.Epoch.Format(KVNTimeFormat),
		e.MeanMotion,
		e.Eccentricity,
		e.Inclination,
		e.RightAscension,
		e.ArgOfPericenter,
		e.MeanAnomaly,
		e.EphemerisType,
		e.ClassificationType,
		e.NoradCatId,
		e.ElementSet,
		int(e.RevAtEpoch),
		e.BStar,
		e.MeanMotionDot,
		e.MeanMotionDDot,
	), nil
}

type String string

func (s *String) Scan(state fmt.ScanState, verb rune) error {
	var acc string
	for {
		r, _, err := state.ReadRune()
		if err != nil {
			return err
		}
		ch := string(r)
		if ch == "\n" {
			state.UnreadRune()
			break
		}
		acc += string(ch)
	}

	*s = String(acc)
	return nil

}

var newlines = regexp.MustCompile("\n+")

// ParseInternationalDesignator tries to parse strings like "1998-067A".
func ParseInternationalDesignator(s string) (y int, n int, p string, err error) {
	y = 0
	n = 0
	p = ""
	if s == "null" {
		return
	}
	if len(s) < 4 {
		return
	}

	m := s[0:4]
	if y, err = strconv.Atoi(m); err != nil {
		err = fmt.Errorf("bad international designator year: '%s' (%s)", m, err)
		return
	}

	if len(s) < 8 {
		return
	}

	m = s[5:8]
	if n, err = strconv.Atoi(m); err != nil {
		err = fmt.Errorf("bad international designator num: '%s' (%s)", m, err)
		return
	}

	p = strings.TrimSpace(s[8:])

	return
}

func (e *Elements) UseInternationalDesignator() error {
	if e.LaunchYear != 0 {
		return nil
	}
	y, n, p, err := ParseInternationalDesignator(e.Id)
	if err != nil {
		return err
	}
	e.LaunchYear = y
	e.LaunchNum = n
	e.LaunchPiece = p
	return nil
}

func ParseKVN(s string) (*Elements, int, error) {

	s = newlines.ReplaceAllString(s, "\n")

	{
		// Scanf has a hard time with zero-length %s input.
		//
		// Should gensym.

		missing := "null"

		s = strings.ReplaceAll(s, "= \n", "= "+missing+"\n")
	}

	f := `CCSDS_OMM_VERS = 2.0
CREATION_DATE  = %s
ORIGINATOR     = %s
OBJECT_NAME    = %s
OBJECT_ID      = %s
CENTER_NAME    = EARTH
REF_FRAME      = TEME
TIME_SYSTEM    = UTC
MEAN_ELEMENT_THEORY = SGP/SGP4
EPOCH          = %s
MEAN_MOTION    = %f
ECCENTRICITY   = %f
INCLINATION    = %f
RA_OF_ASC_NODE = %f
ARG_OF_PERICENTER = %f
MEAN_ANOMALY   = %f
EPHEMERIS_TYPE = %d
CLASSIFICATION_TYPE = %s
NORAD_CAT_ID   = %s
ELEMENT_SET_NO = %d
REV_AT_EPOCH   = %f
BSTAR          = %e
MEAN_MOTION_DOT = %e
MEAN_MOTION_DDOT = %e
`
	f = regexp.MustCompile("\n+").ReplaceAllString(f, "\n")

	var (
		epoch        string
		e            = &Elements{}
		creationDate String
		originator   String
		name         String
	)

	n, err := fmt.Sscanf(s,
		f,
		&creationDate,
		&originator,
		&name,
		&e.Id,
		&epoch,
		&e.MeanMotion,
		&e.Eccentricity,
		&e.Inclination,
		&e.RightAscension,
		&e.ArgOfPericenter,
		&e.MeanAnomaly,
		&e.EphemerisType,
		&e.ClassificationType,
		&e.NoradCatId,
		&e.ElementSet,
		&e.RevAtEpoch,
		&e.BStar,
		&e.MeanMotionDot,
		&e.MeanMotionDDot,
	)

	if err != nil {
		return nil, n, err
	}

	e.Originator = string(originator)
	if e.Originator == "null" {
		e.Originator = ""
	}

	e.Name = string(name)

	{
		y, num, p, err := ParseInternationalDesignator(e.Id)
		if err != nil {
			return nil, n, fmt.Errorf("failed to parse international designator '%s': %s", e.Id, err)
		}
		e.LaunchYear = y
		e.LaunchNum = num
		e.LaunchPiece = p
	}

	if creationDate == "null" {
		creationDate = ""
	}

	if 0 < len(creationDate) {
		e.CreationDate, err = ParseTime(KVNTimeFormat, string(creationDate))
		if err != nil {
			return nil, n, err

		}
	}

	if e.Epoch, err = ParseTime(KVNTimeFormat, epoch); err != nil {
		return nil, n, err
	}

	return e, n, nil
}

func DoKVNs(r *bufio.Reader, f func(s string) error) error {

	var (
		first = "CCSDS_OMM_VERS"
		lines []string
		emit  = func() error {
			if lines != nil {
				return f(strings.Join(lines, "\n"))
			}
			return nil
		}
	)

	for {
		line, err := r.ReadString('\n')
		line = strings.TrimRight(line, "\n\r")

		if strings.HasPrefix(line, first) {
			if err = emit(); err != nil {
				return err
			}
			lines = make([]string, 0, 20)
		}
		lines = append(lines, line)

		if err == io.EOF {
			break
		}
	}

	return emit()
}
