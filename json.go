package gpelements

import (
	"regexp"
)

// This file exists because space-track.com JSON element sets use
// strings for numeric values.  That's maybe a good practice but it
// makes deserialization harder.  Should check the standard.

// DestringNumbers attempts to remove quotation marks from values we
// want to be numeric.

var destringThese = map[string]bool{
	"MEAN_MOTION":       true,
	"ECCENTRICITY":      true,
	"INCLINATION":       true,
	"RA_OF_ASC_NODE":    true,
	"ARG_OF_PERICENTER": true,
	"MEAN_ANOMALY":      true,
	"EPHEMERIS_TYPE":    true,
	"ELEMENT_SET_NO":    true,
	"REV_AT_EPOCH":      true,
	"BSTAR":             true,
	"MEAN_MOTION_DOT":   true,
	"MEAN_MOTION_DDOT":  true,
}

var strPair = regexp.MustCompile(`"([A-Z_]+)":"([-0-9Ee+.]+)"`)

func DestringNumbers(src string) string {
	return strPair.ReplaceAllStringFunc(src, func(s string) string {
		kv := strPair.FindStringSubmatch(s)
		if destringThese[kv[1]] {
			return `"` + kv[1] + `":` + kv[2]
		}
		return s
	})
}
