#!/bin/bash

set -e

mkdir -p data

curl 'https://celestrak.com/NORAD/elements/gp.php?GROUP=STATIONS&FORMAT=TLE' > data/test.tle
sleep 3
curl 'https://celestrak.com/NORAD/elements/gp.php?GROUP=STATIONS&FORMAT=KVN' > data/test.kvn
sleep 3
curl 'https://celestrak.com/NORAD/elements/gp.php?GROUP=STATIONS&FORMAT=XML' > data/test.xml
sleep 3
curl 'https://celestrak.com/NORAD/elements/gp.php?GROUP=STATIONS&FORMAT=JSON' > data/test.jsonarray
sleep 3
curl 'https://celestrak.com/NORAD/elements/gp.php?GROUP=STATIONS&FORMAT=CSV' > data/test.csv

