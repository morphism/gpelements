package main

import (
	"bufio"
	"encoding/json"
	"encoding/xml"
	"flag"
	"fmt"
	"hash/fnv"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/morphism/gpelements"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {

	var (
		ts = func(t time.Time) string {
			return t.Format(time.RFC3339Nano)
		}

		now = time.Now().UTC()

		transform = flag.NewFlagSet("transform", flag.ExitOnError)
		emit      = transform.String("emit", "csv", "Output represention: csv|csvh|json|jsonarray|tle|kvn|xml")
		tolerate  = transform.Bool("tolerate", false, "Log errors instead of stopping")

		prop         = flag.NewFlagSet("prop", flag.ExitOnError)
		propFrom     = prop.String("from", ts(now), "Propagation start time")
		propTo       = prop.String("to", ts(now.Add(time.Hour)), "Propagation end time")
		propInterval = prop.Duration("interval", 10*time.Minute, "Propagation end time")
		propHigher   = prop.Bool("higher-precision", true, "Higher-precision (as able)")

		orbit         = flag.NewFlagSet("on-orbit", flag.ExitOnError)
		orbitFrom     = orbit.String("from", ts(now), "Propagation start time")
		orbitTo       = orbit.String("to", ts(now.Add(time.Hour)), "Propagation end time")
		orbitInterval = orbit.Duration("interval", 10*time.Minute, "Propagation end time")

		walk       = flag.NewFlagSet("walk", flag.ExitOnError)
		minSteps   = walk.Int("min-steps", 1, "Minimum number of steps")
		maxSteps   = walk.Int("max-steps", 3, "Maximum number of steps")
		incSet     = walk.Bool("inc-set", true, "Increment element set number")
		resetEpoch = walk.Bool("reset-epoch", true, "Set Epoch to now")
		seed       = walk.Int64("seed", time.Now().UTC().UnixNano(), "RNG seed (defaults to current time in ns)")

		rename      = flag.NewFlagSet("rename", flag.ExitOnError)
		renameState = rename.Int64("state", 0, "Next catalog number in Alpha-5 A range")
		renameClear = rename.Bool("clear", false, "Remove original name (suffix)")

		sample    = flag.NewFlagSet("sample", flag.PanicOnError)
		sampleMod = sample.Int("mod", 10, "Sampling hash modulus")
		sampleRem = sample.Int("rem", 0, "Sampling hash remainder")

		random           = flag.NewFlagSet("random", flag.PanicOnError)
		randomPercentage = random.Float64("percent", 0, "Approximate percent of lines to emit")
	)

	usage := func() {
		fmt.Fprintf(os.Stderr, `Usage: %s transform|prop|on-orbit|walk|rename|sample|random ...

Subcommands:

  transform:

`, os.Args[0])
		transform.PrintDefaults()

		fmt.Fprintf(os.Stderr, `

  csvh emits CSV output with a header line.

  json emits each element set as a single line of JSON.

  jsonarray emits an array of element sets as one big blob of JSON.

`)

		fmt.Fprintf(os.Stderr, "\n  prop: Propagate\n\n")
		prop.PrintDefaults()

		fmt.Fprintf(os.Stderr, "\n  on-orbit: Filter for on-orbit\n\n")
		orbit.PrintDefaults()

		fmt.Fprintf(os.Stderr, "\n  walk: Random walk\n\n")
		walk.PrintDefaults()

		fmt.Fprintf(os.Stderr, "\n  rename: Update name, catalog number\n\n")
		rename.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\n")

		fmt.Fprintf(os.Stderr, "\n  Sample: Sampled based on hash of name+id+num\n\n")
		sample.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\n")

		fmt.Fprintf(os.Stderr, "\n  Random: Emit a percentage of the input\n\n")
		random.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\n")
	}

	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}

	var (
		subcommand = os.Args[1]
		args       = os.Args[2:]
	)

	switch subcommand {
	case "transform":
		transform.Parse(args)
	case "prop":
		prop.Parse(args)
	case "on-orbit", "orbit":
		orbit.Parse(args)
	case "walk":
		walk.Parse(args)
	case "rename":
		rename.Parse(args)
	case "sample":
		sample.Parse(args)
	case "random":
		random.Parse(args)
	default:
		usage()
		os.Exit(1)
	}

	rand.Seed(*seed)

	in := os.Stdin

	defer in.Close()
	r := bufio.NewReader(in)

	es := make([]gpelements.Elements, 0, 1024)
	err := gpelements.Do(r, func(e gpelements.Elements) error {
		es = append(es, e)
		return nil
	})

	if err != nil {
		return err
	}

	switch subcommand {
	case "transform":
		csv := func() error {
			for i, e := range es {
				s, err := e.MarshalCSV()
				if err != nil {
					if !*tolerate {
						return err
					}
					log.Printf("set %d error: %s", i, err)
				} else {
					fmt.Printf("%s", s)
				}
			}
			return nil
		}

		switch *emit {
		case "csvh":
			fmt.Printf("%s\n", gpelements.CSVHeader)
			if err := csv(); err != nil {
				return err
			}

		case "csv":
			if err := csv(); err != nil {
				return err
			}

		case "xml":
			list := gpelements.ElementsList{
				Es: es,
			}
			bs, err := xml.Marshal(list)
			if err != nil {
				return err
			}
			fmt.Printf("%s\n", bs)

		case "json":
			for i, e := range es {
				bs, err := json.Marshal(e)
				if err != nil {
					if !*tolerate {
						return err
					}
					log.Printf("set %d error: %s", i, err)
				} else {
					fmt.Printf("%s\n", bs)
				}
			}
		case "jsonarray":
			bs, err := json.MarshalIndent(es, "", "  ")
			if err != nil {
				if !*tolerate {
					return err
				}
				log.Printf("error: %s", err)
			} else {
				fmt.Printf("%s\n", bs)
			}

		case "kvn":
			for i, e := range es {
				s, err := e.MarshalKVN()
				if err != nil {
					if !*tolerate {
						return err
					}
					log.Printf("set %d error: %s", i, err)
				} else {
					fmt.Printf("%s", s)
				}
			}
		case "tle":
			for i, e := range es {
				l0, l1, l2, err := e.MarshalTLE()
				if err != nil {
					if !*tolerate {
						return err
					}
					log.Printf("set %d error: %s", i, err)
				} else {
					fmt.Printf("%s\n%s\n%s\n", l0, l1, l2)
				}
			}
		default:
			return fmt.Errorf("unknown output representation '%s'", *emit)
		}
	case "prop":
		t0, err := time.Parse(time.RFC3339Nano, *propFrom)
		if err != nil {
			return err
		}
		t1, err := time.Parse(time.RFC3339Nano, *propTo)
		if err != nil {
			return err
		}
		gpelements.HigherPrecisionSGP4 = *propHigher
		for i, e := range es {
			err := Prop(&e, t0, t1, *propInterval, true)
			if err != nil {
				if !*tolerate {
					return err
				}
				log.Printf("set %d error: %s", i, err)
			}
		}

	case "sample":
		for _, e := range es {
			var (
				k = e.Name + "/" + e.Id + "/" + string(e.NoradCatId)
				h = Hash(k)
				r = h % uint64(*sampleMod)
			)
			if r != uint64(*sampleRem) {
				continue
			}

			bs, err := json.Marshal(e)
			if err != nil {
				return err
			}
			fmt.Printf("%s\n", bs)
		}

	case "random":
		percent := *randomPercentage
		if 1 < percent {
			percent /= 100
		}
		for _, e := range es {
			if percent < rand.Float64() {
				continue
			}
			bs, err := json.Marshal(e)
			if err != nil {
				return err
			}
			fmt.Printf("%s\n", bs)
		}

	case "orbit", "on-orbit":
		t0, err := time.Parse(time.RFC3339Nano, *orbitFrom)
		if err != nil {
			return err
		}
		t1, err := time.Parse(time.RFC3339Nano, *orbitTo)
		if err != nil {
			return err
		}
		for i, e := range es {
			err := Prop(&e, t0, t1, *orbitInterval, false)
			if err != nil {
				log.Printf("set %d propagation error: %s", i, err)
				continue
			}
			bs, err := json.Marshal(e)
			if err != nil {
				return err
			}
			fmt.Printf("%s\n", bs)
		}

	case "rename":
		state := *renameState
		for i, e := range es {
			var id string
			id, state, err = gpelements.NextAlpha5Num(state)
			if err != nil {
				return err
			}
			if *renameClear {
				e.Name = "#" + id
			} else {
				e.UpdateName(id)
			}
			e.NoradCatId = gpelements.NoradCatId(id)
			e.ElementSet = 0
			e.Epoch = gpelements.NewTime(time.Now().UTC())

			// Probably should emit in a high-precision format.
			l0, l1, l2, err := e.MarshalTLE()
			if err != nil {
				if !*tolerate {
					return err
				}
				log.Printf("set %d error: %s", i, err)
			} else {
				fmt.Printf("%s\n%s\n%s\n", l0, l1, l2)
			}
		}

	case "walk":
		for i, e := range es {
			if err := e.Walk(*minSteps, *maxSteps); err != nil {
				return err
			}

			if *incSet {
				if err := e.IncSetNum(); err != nil {
					return err
				}
			}

			if *resetEpoch {
				e.Epoch = gpelements.NewTime(time.Now().UTC())
			}

			// Probably should emit in a high-precision format.
			l0, l1, l2, err := e.MarshalTLE()
			if err != nil {
				if !*tolerate {
					return err
				}
				log.Printf("set %d error: %s", i, err)
			} else {
				fmt.Printf("%s\n%s\n%s\n", l0, l1, l2)
			}
		}

	}

	return nil
}

func Prop(e *gpelements.Elements, from, to time.Time, interval time.Duration, print bool) error {
	o, err := e.SGP4()
	if err != nil {
		return err
	}
	for t := from; t.Before(to); t = t.Add(interval) {
		s, err := gpelements.Prop(o, t)
		if err != nil {
			return err
		}

		lla, err := gpelements.ECIToLLA(t, s.ECI)
		if err != nil {
			return err
		}

		if !print {
			continue
		}

		m := map[string]interface{}{
			"Name":  e.Name,
			"Id":    e.Id,
			"Norad": e.NoradCatId,
			"At":    t,
			"State": s,
			"LLA":   lla,
			"Age":   t.Sub(time.Time(*e.Epoch)).Seconds(),
		}
		js, err := json.Marshal(&m)
		if err != nil {
			log.Fatalf("prop json.Marshal error %s on %#v", err, m)
		}
		fmt.Printf("%s\n", js)
	}

	return nil
}

func Hash(s string) uint64 {
	h := fnv.New64a()
	h.Write([]byte(s))
	return h.Sum64()
}
