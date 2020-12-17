package gpelements

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"time"
)

const (
	CSVHeader = "OBJECT_NAME,OBJECT_ID,EPOCH,MEAN_MOTION,ECCENTRICITY,INCLINATION,RA_OF_ASC_NODE,ARG_OF_PERICENTER,MEAN_ANOMALY,EPHEMERIS_TYPE,CLASSIFICATION_TYPE,NORAD_CAT_ID,ELEMENT_SET_NO,REV_AT_EPOCH,BSTAR,MEAN_MOTION_DOT,MEAN_MOTION_DDOT"
)

var (
	CSVTimeFormat = "2006-01-02T15:04:05.999999999"
)

func (e *Elements) MarshalCSV() (string, error) {
	f := "%q,%q,%q,%g,%g,%g,%g,%g,%g,%d,%q,%q,%d,%g,%g,%g,%g"

	return fmt.Sprintf(f,
		e.Name,
		e.Id,
		e.Epoch.Format(CSVTimeFormat),
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
		e.RevAtEpoch,
		e.BStar,
		e.MeanMotionDot,
		e.MeanMotionDDot,
	), nil
}

// QuoteString will make sure the given line has double-quoted values
// where we (think) we want them.
func QuoteStrings(line string) string {
	ss := strings.Split(line, ",")
	which := []bool{
		true, true, true,
		false, false, false, false, false, false,
		false,
		true, true,
	}
	for i, fix := range which {
		if fix {
			s := strings.Trim(ss[i], `"`)
			ss[i] = `"` + s + `"`
		}
	}

	return strings.Join(ss, ",")
}

func ParseCSV(line string) (*Elements, int, error) {

	line = QuoteStrings(line)
	var (
		e  = &Elements{}
		ts string
		f  = `%q,%q,%q,%f,%f,%f,%f,%f,%f,%d,%q,%q,%d,%f,%f,%f,%f`
	)

	// ToDo: Include Originator, CreationDate?

	n, err := fmt.Sscanf(line, f,
		&e.Name,
		&e.Id,
		&ts,
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

	t, err := time.Parse(CSVTimeFormat, ts)
	if err != nil {
		return nil, n, err

	}
	t0 := Time(t)
	e.Epoch = &t0

	return e, n, nil

}

func DoLines(r *bufio.Reader, f func(s string) error) error {

	i := 0
	for {
		line, err := r.ReadString('\n')
		line = strings.TrimRight(line, "\n\r")
		if 0 < len(line) {
			if err := f(line); err != nil {
				return fmt.Errorf("error parsing line %d: %s", i, err)
			}
		}
		if err == io.EOF {
			break
		}
	}

	return nil
}
