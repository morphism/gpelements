package main

import (
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

		generic  = flag.NewFlagSet("generic", flag.ContinueOnError)
		bufSize  = generic.Int("buf-size", 4096, "Buffer size")
		tolerate = generic.Bool("tolerate", false, "Log errors instead of stopping")

		transform = flag.NewFlagSet("transform", flag.ExitOnError)
		emit      = transform.String("emit", "csv", "Output represention: csv|csvh|json|jsonarray|tle|kvn|xml")

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

		sample    = flag.NewFlagSet("sample", flag.ExitOnError)
		sampleMod = sample.Int("mod", 10, "Sampling hash modulus")
		sampleRem = sample.Int("rem", 0, "Sampling hash remainder")

		random           = flag.NewFlagSet("random", flag.ExitOnError)
		randomPercentage = random.Float64("percent", 0, "Approximate percent of lines to emit")
	)

	usage := func() {
		fmt.Fprintf(os.Stderr, `Usage: %s transform|prop|on-orbit|walk|rename|sample|random ...

Subcommands:

  transform:

`, os.Args[0])
		generic.PrintDefaults()

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

	generic.Parse(os.Args[1:])

	var (
		remaining  = generic.NArg()
		next       = len(os.Args) - remaining
		subcommand = os.Args[next]
		args       = os.Args[next+1:]
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

	if *bufSize < gpelements.MinBufferSize {
		*bufSize = gpelements.MinBufferSize
	}

	t0, err := time.Parse(time.RFC3339Nano, *propFrom)
	if err != nil {
		return err
	}
	t1, err := time.Parse(time.RFC3339Nano, *propTo)
	if err != nil {
		return err
	}
	gpelements.HigherPrecisionSGP4 = *propHigher

	var (
		i  = 0
		es = make([]gpelements.Elements, 0, 1024)
	)

	in := os.Stdin
	defer in.Close()

	err = gpelements.Do(in, *bufSize, func(e gpelements.Elements) error {
		var (
			s   string
			err error
			bs  []byte
		)

		switch subcommand {
		case "transform":
			switch *emit {
			case "csv":
				s, err = e.MarshalCSV()
			case "csvh":
				if s, err = e.MarshalCSV(); err == nil {
					if i == 0 {
						s = gpelements.CSVHeader + "\n" + s
					}
				}
			case "json":
				if bs, err = json.Marshal(e); err == nil {
					s = string(bs)
				}
			case "kvn":
				s, err = e.MarshalKVN()
			case "tle":
				var l0, l1, l2 string
				if l0, l1, l2, err = e.MarshalTLE(); err == nil {
					s = fmt.Sprintf("%s\n%s\n%s\n", l0, l1, l2)
				}
			case "xml":
				es = append(es, e)
			case "jsonarray":
				es = append(es, e)
			default:
				return fmt.Errorf("unknown output representation '%s'", *emit)
			}
		case "prop":
			err = Prop(&e, t0, t1, *propInterval, true)
		case "sample":
			var (
				k = e.Name + "/" + e.Id + "/" + string(e.NoradCatId)
				h = Hash(k)
				r = h % uint64(*sampleMod)
			)
			if r == uint64(*sampleRem) {
				bs, err := json.Marshal(e)
				if err == nil {
					s = string(bs)
				}
			}
		case "random":
			percent := *randomPercentage
			if 1 < percent {
				percent /= 100
			}
			if percent > rand.Float64() {
				bs, err = json.Marshal(e)
				if bs == nil {
					s = string(bs)
				}
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
			err = Prop(&e, t0, t1, *orbitInterval, false)
			if err == nil {
				bs, err = json.Marshal(e)
				if err == nil {
					s = string(bs)
				}
			}

		case "rename":
			state := *renameState
			var id string
			id, state, err = gpelements.NextAlpha5Num(state)
			if err != nil {
				break
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
			var l0, l1, l2 string
			if l0, l1, l2, err = e.MarshalTLE(); err == nil {
				s = fmt.Sprintf("%s\n%s\n%s\n", l0, l1, l2)
			}
		case "walk":
			if err = e.Walk(*minSteps, *maxSteps); err == nil {
				if *incSet {
					if err = e.IncSetNum(); err == nil {
						if *resetEpoch {
							e.Epoch = gpelements.NewTime(time.Now().UTC())
						}

						// Probably should emit in a high-precision format.
						var l0, l1, l2 string
						l0, l1, l2, err = e.MarshalTLE()
						if err == nil {
							s = fmt.Sprintf("%s\n%s\n%s\n", l0, l1, l2)
						}
					}
				}
			}
		}

		if err != nil {
			if *tolerate {
				log.Printf("at %d %v", i, err)
				err = nil
			}
		}

		i++

		if 0 < len(s) {
			fmt.Println(s)
		}

		return err
	})

	if err == nil {
		var bs []byte
		switch subcommand {
		case "jsonarray":
			bs, err = json.MarshalIndent(es, "", "  ")
			if err == nil {
				fmt.Printf("%s\n", bs)
			}

		case "xml":
			list := gpelements.ElementsList{
				Es: es,
			}
			bs, err := xml.Marshal(list)
			if err == nil {
				fmt.Printf("%s\n", bs)
			}
		}
	}

	return err
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
