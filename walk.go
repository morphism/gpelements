package gpelements

import (
	"fmt"
	"math/rand"
	"regexp"
)

// Re the NORAD catalog number namespace:
//
// The "Alpha-5" scheme allows for namespace expansion by accepting a
// capital English letter as the leading character.  The resulting
// namespace would hold 360,000 identifiers (assuming I and O are not
// discarded).
//
// According to Celestrak, 18 Space Control Squadron (18 SPCS) is
// currently assigning numbers in the range 7995xxxxx (an "analyst sat
// range").

// NextAlpha5Num generates a new NORAD catalogy number using the "Alpha-5"
// scheme and updates the Name.
func NextAlpha5Num(state int64) (string, int64, error) {
	var (
		blocks = "ABCDEFGHIJKMPQRSTUVWXYZ"
		block  = state / 10000
		rem    = state % 10000
	)

	if int64(len(blocks)) <= block {
		return "", state, fmt.Errorf("NextAlpha5Num at %d is out of blocks (%d)", state, block)
	}

	id := fmt.Sprintf("%c%04d", blocks[block], rem)

	return id, state + 1, nil
}

var idInName = regexp.MustCompile("^#[ ]+ ")

func (e *Elements) UpdateName(id string) {
	s := idInName.ReplaceAllString(e.Name, "")
	e.Name = "#" + id + " " + s
}

func (e *Elements) Reassign(state int64) (int64, error) {
	id, next, err := NextAlpha5Num(state)
	if err != nil {
		return state, err
	}
	e.NoradCatId = NoradCatId(id)
	return next, nil
}

func (e *Elements) IncSetNum() error {
	if 10000 <= e.ElementSet {
		return fmt.Errorf("maximum ElementSet reached: %d", e.ElementSet)
	}
	e.ElementSet++
	return nil
}

// Walk modifies the Elements based on an under-specified random walk.
//
// Name, Id, and ElementSet are not changed.
func (e *Elements) Walk(minSteps, maxSteps int) error {
	if maxSteps < minSteps || maxSteps == 0 {
		return nil
	}

	steps := minSteps + rand.Intn(1+maxSteps-minSteps)

	for i := 0; i < steps; i++ {
		switch rand.Intn(6) {
		case 0:
			e.Inclination = stepDegrees(e.Inclination)
		case 1:
			e.RightAscension = stepDegrees(e.RightAscension)
		case 2:
			e.Eccentricity = stepEcc(e.Eccentricity)
		case 3:
			e.ArgOfPericenter = stepDegrees(e.ArgOfPericenter)
		case 4:
			e.MeanAnomaly = stepDegrees(e.MeanAnomaly)
		case 5:
			e.MeanMotion = stepMotion(e.MeanMotion)
			// ToDo: Somehow update derivatives?

		}
	}

	return nil
}

func rnd() float64 {
	return 1 - 2*rand.Float64()
}

// stepDegrees moves +/- [0,0.01) degree.
func stepDegrees(d float64) float64 {
	return d + rnd()/100
}

// stepEcc moves +/- [0,0.001) percent.
func stepEcc(x float64) float64 {
	y := x * (1 + rnd()/1000)
	if 1 <= y {
		y = 1 - rand.Float64()/1000
	}
	if y <= 0 {
		y = rand.Float64() / 1000
	}

	return y
}

// stepMotion moves +/- [0,0.0001) percent.
func stepMotion(x float64) float64 {
	// ToDo: Something better.

	y := x * (1 + rnd()/10000)

	if y <= 0 {
		y = 1 + rnd()/10
	}

	if y < 0 {
		panic(y)
	}

	return y
}
