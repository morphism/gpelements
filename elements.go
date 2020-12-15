package gpelements

import (
	"encoding/xml"
	"fmt"
	"strconv"
	"strings"
)

// NoradCatId exists in an attempt to handle numbers or the new,
// temporary "Alpha-5 scheme" for NORAD catalog identifiers.
//
// Note that JSON serializations have explicit syntax for strings vs
// numbers, so a NoradCatId will try to (de)serialize as a number
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
	if _, err := strconv.Atoi(string(s)); err == nil {
		return []byte(s), nil
	}
	return []byte(fmt.Sprintf(`"%s"`, string(s))), nil
}

type Elements struct {
	XMLName xml.Name `json:"-" xml:"omm"`

	OMMId      string `json:"-" xml:"id,attr"`
	OMMVersion string `json:"-" xml:"version,attr"`

	CreationDate *Time `json:",omitempty" xml:"header>CREATION_DATE"`

	Originator string `json:",omitempty" xml:"header>ORIGINATOR"`

	// Name: "OBJECT_NAME": "STARLINK-1329"
	Name string `json:"OBJECT_NAME,omitempty" xml:"body>segment>metadata>OBJECT_NAME"`

	// Id: "OBJECT_ID": "2020-025A",
	Id string `json:"OBJECT_ID,omitempty" xml:"body>segment>metadata>OBJECT_ID"`

	LaunchYear  int    `json:"-" xml:"-"`
	LaunchNum   int    `json:"-" xml:"-"`
	LaunchPiece string `json:"-" xml:"-"`

	// Epoch: "EPOCH": "2020-12-13T03:44:10.927392",
	Epoch *Time `json:"EPOCH" xml:"body>segment>data>meanElements>EPOCH"`

	// MeanMotion: "MEAN_MOTION": 15.05587631,
	MeanMotion float64 `json:"MEAN_MOTION" xml:"body>segment>data>meanElements>MEAN_MOTION"`

	// Eccentricity: "ECCENTRICITY": 0.000133,
	Eccentricity float64 `json:"ECCENTRICITY" xml:"body>segment>data>meanElements>ECCENTRICITY"`

	// Inclination: "INCLINATION": 53.0021,
	Inclination float64 `json:"INCLINATION" xml:"body>segment>data>meanElements>INCLINATION"`

	// RightAscension: "RA_OF_ASC_NODE": 52.3173,
	RightAscention float64 `json:"RA_OF_ASC_NODE" xml:"body>segment>data>meanElements>RA_OF_ASC_NODE"`

	// ArgOfPericenter: "ARG_OF_PERICENTER": 82.0365,
	ArgOfPericenter float64 `json:"ARG_OF_PERICENTER" xml:"body>segment>data>meanElements>ARG_OF_PERICENTER"`

	// MeanAnomaly: "MEAN_ANOMALY": 278.0774,
	MeanAnomaly float64 `json:"MEAN_ANOMALY" xml:"body>segment>data>meanElements>MEAN_ANOMALY"`

	// EphemerisType: "EPHEMERIS_TYPE": 0,
	EphemerisType int `json:"EPHEMERIS_TYPE" xml:"body>segment>data>tleParameters>EPHEMERIS_TYPE"`

	// ClassificationType: "CLASSIFICATION_TYPE": "U",
	ClassificationType string `json:"CLASSIFICATION_TYPE" xml:"body>segment>data>tleParameters>CLASSIFICATION_TYPE"`

	// NoradCatId: "NORAD_CAT_ID": 45531,
	//
	// This value is a string in order to support the "Alpha-5
	// scheme" as well as arbitrary identifiers.
	NoradCatId NoradCatId `json:"NORAD_CAT_ID" xml:"body>segment>data>tleParameters>NORAD_CAT_ID"`

	// ElementSetNum: "ELEMENT_SET_NO": 999,
	ElementSet int `json:"ELEMENT_SET_NO" xml:"body>segment>data>tleParameters>ELEMENT_SET_NO"`

	// RevAtEpoch: "REV_AT_EPOCH": 3623,
	//
	// Usually an int.
	RevAtEpoch float64 `json:"REV_AT_EPOCH" xml:"body>segment>data>tleParameters>REV_AT_EPOCH"`

	// BStar: "BSTAR": 6.5515e-5,
	BStar float64 `json:"BSTAR" xml:"body>segment>data>tleParameters>BSTAR"`

	// MeanMotionDOT: "MEAN_MOTION_DOT": 6.76e-6,
	MeanMotionDot float64 `json:"MEAN_MOTION_DOT" xml:"body>segment>data>tleParameters>MEAN_MOTION_DOT"`

	// MeanMontionDDOT: "MEAN_MOTION_DDOT": 0
	MeanMotionDDot float64 `json:"MEAN_MOTION_DDOT" xml:"body>segment>data>tleParameters>MEAN_MOTION_DDOT"`
}

func (e *Elements) Copy() *Elements {
	acc := *e
	return &acc
}
