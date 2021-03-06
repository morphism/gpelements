_Warning: This code is experimental; please see the [LICENSE](LICENSE)_.

A tool to convert to and from various general perturbations orbital
data formats.  See [this
discussion](https://celestrak.com/NORAD/documentation/gp-data-formats.php)
for some background.

Supported formats:

1. TLE
1. [OMM](https://public.ccsds.org/Pubs/502x0b2c1e2.pdf) XML
1. [OMM](https://public.ccsds.org/Pubs/502x0b2c1e2.pdf) KVN
1. JSON
1. CSV

Currently not all round trips starting from Celestrak data are
perfect.  (Floating point number formatting is only one class of
challenge there.)

This tool can also perform SGP4 propagation (using [this
implementation](https://github.com/morphism/sgp4go)), object renaming,
and random walks.

## Usage

```
Usage: tletool transform|prop|on-orbit|walk|rename|sample|random ...

Subcommands:

  transform:

  -buf-size int
    	Buffer size (default 4096)
  -help
    	Just get help
  -tolerate
    	Log errors instead of stopping
  -emit string
    	Output represention: csv|csvh|json|jsonarray|tle|kvn|xml (default "csv")


  csvh emits CSV output with a header line.

  json emits each element set as a single line of JSON.

  jsonarray emits an array of element sets as one big blob of JSON.


  prop: Propagate

  -duration duration
    	Duration of propagation (instead of -to)
  -from string
    	Propagation start time (default "2020-12-18T16:42:22.644764256Z")
  -higher-precision
    	Higher-precision (as able) (default true)
  -interval duration
    	Propagation end time (default 10m0s)
  -to string
    	Propagation end time (default "2020-12-18T17:42:22.644764256Z")

  on-orbit: Filter for on-orbit

  -from string
    	Propagation start time (default "2020-12-18T16:42:22.644764256Z")
  -interval duration
    	Propagation end time (default 10m0s)
  -to string
    	Propagation end time (default "2020-12-18T17:42:22.644764256Z")

  walk: Random walk

  -inc-set
    	Increment element set number (default true)
  -max-steps int
    	Maximum number of steps (default 3)
  -min-steps int
    	Minimum number of steps (default 1)
  -reset-epoch
    	Set Epoch to now (default true)
  -seed int
    	RNG seed (defaults to current time in ns) (default 1608309742644776587)

  rename: Update name, catalog number

  -clear
    	Remove original name (suffix)
  -state int
    	Next catalog number in Alpha-5 A range


  Sample: Sampled based on hash of name+id+num

  -mod int
    	Sampling hash modulus (default 10)
  -rem int
    	Sampling hash remainder


  Random: Emit a percentage of the input

  -percent float
    	Approximate percent of lines to emit

```

(The default timestamps are acutally the current time.)

## Examples

With [this
data](https://www.space-track.org/basicspacedata/query/class/gp/EPOCH/%3Enow-30/NORAD_CAT_ID/270000--339999/orderby/NORAD_CAT_ID/format/json)
from [`space-track.org`](https://www.space-track.org/):

```Shell
# We'll needlessly exercise our transformations.
cat tmp/now.json |
    tletool transform -emit xml |
    tletool transform -emit kvn |
    tletool transform -emit csv |
    tletool transform -emit json |
    tletool prop |
    tail -4 |
    jq -r -c '{"NORAD":.Norad,"LLA":.LLA}'

```

gives

```
{"NORAD":270288,"LLA":{"Lat":65.84531,"Lon":-14.565547,"Alt":1292.6202}}
{"NORAD":270288,"LLA":{"Lat":33.46682,"Lon":-17.518644,"Alt":1285.534}}
{"NORAD":270288,"LLA":{"Lat":0.9417589,"Lon":-20.208746,"Alt":1283.9185}}
{"NORAD":270288,"LLA":{"Lat":-31.554588,"Lon":-22.894846,"Alt":1292.9807}}
```

With
[data](https://celestrak.com/NORAD/elements/gp.php?GROUP=STATIONS&FORMAT=KVN)
from [Celestrak](https://celestrak.com):

```Shell
if [ ! -f data/test.kvn ]; then
	curl 'https://celestrak.com/NORAD/elements/gp.php?GROUP=STATIONS&FORMAT=KVN' > data/test.kvn
fi

# We'll needlessly exercise our transformations, and we'll keep 
# the intermediate output.

cat data/test.kvn |
  tletool transform -emit json | 
  tee data.json | 
  tletool transform -emit xml | 
  tee data.xml | 
  tletool transform -emit csv | 
  tee data.csv | 
  tletool transform -emit kvn | 
  tee data.kvn | 
  tletool on-orbit | 
  tee on-orbit.json |
  tletool prop | 
  tee prop.json |
  head -1 | 
  jq -r .
```

results in something that looks like

```JSON
{
  "Age": 32705.934428877,
  "At": "2020-12-14T21:29:18.874844877Z",
  "Id": "1998-067A",
  "LLA": {
    "Lat": 17.104858,
    "Lon": 87.40023,
    "Alt": 420.42767
  },
  "Name": "ISS (ZARYA)",
  "Norad": 25544,
  "State": {
    "V": {
      "X": -4.775169,
      "Y": -2.1982052,
      "Z": -5.5731597
    },
    "ECI": {
      "X": -4493.6265,
      "Y": 4695.992,
      "Z": 1987.5846
    }
  }
}
```

in addition to several intermediate output files.

The example above uses [jq](https://stedolan.github.io/jq/).

## References

1. [Consultative Committee for Space Data Systems (CCSDS)](https://public.ccsds.org/default.aspx)

1. [Recommended Standard CCSDS 502.0-B-2](https://public.ccsds.org/Pubs/502x0b2c1e2.pdf)

1. [Celestrak's discussion of element set formats](https://celestrak.com/NORAD/documentation/gp-data-formats.php)

1. The [SGP4 implementation used here](https://github.com/morphism/sgp4go)
