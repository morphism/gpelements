package gpelements

import (
	"time"

	sat "github.com/jsmorph/go-satellite"
	sgp4 "github.com/morphism/sgp4go"
)

func (e *Elements) SGP4() (*sgp4.TLE, error) {
	_, line1, line2, err := e.MarshalTLE()
	if err != nil {
		return nil, err
	}
	return sgp4.NewTLE(line1, line2)
}

// Vect is a 3-vector.
type Vect struct {
	X, Y, Z float32
}

// Ephemeris represents position and velocity.
type Ephemeris struct {
	// V is velocity.
	V Vect

	// C is Cartesian position.
	ECI Vect
}

func Prop(o *sgp4.TLE, t time.Time) (Ephemeris, error) {
	p, v, err := o.PropUnixMillis(t.UnixNano() / 1000 / 1000)
	var e Ephemeris
	if err == nil {
		e = Ephemeris{
			ECI: Vect{float32(p[0]), float32(p[1]), float32(p[2])},
			V:   Vect{float32(v[0]), float32(v[1]), float32(v[2])},
		}
	}
	return e, err
}

func TimeToGST(t time.Time) (float64, float64) {
	var (
		y   = t.Year()
		m   = int(t.Month())
		d   = t.Day()
		h   = t.Hour()
		min = t.Minute()
		sec = t.Second()
		ns  = t.Nanosecond()
	)

	return sat.GSTimeFromDateNano(y, m, d, h, min, sec, ns)
}

type LatLonAlt struct {
	Lat, Lon, Alt float32
}

func ECIToLLA(t time.Time, p Vect) (*LatLonAlt, error) {

	gmst, _ := TimeToGST(t)

	x := sat.Vector3{
		X: float64(p.X),
		Y: float64(p.Y),
		Z: float64(p.Z),
	}

	// sat.ECIToLLA is very slow.
	alt, _, ll := sat.ECIToLLA(x, gmst)

	d, err := sat.LatLongDeg(ll)
	if err != nil {
		return nil, err
	}

	return &LatLonAlt{
		Lat: float32(d.Latitude),
		Lon: float32(d.Longitude),
		Alt: float32(alt),
	}, nil
}
