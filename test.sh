#!/bin/bash

# Try to compare transformations based on Celestrak formats.

set -e pipefail

for OUT in tle csv xml kvn jsonarray; do

    # Celestrak CSV doesn't quote strings.
    #
    # Celestrak exponent flag 'E' is upper case.
    #
    # Celestrak XML has OMM entries on separate lines.
    #
    # Celestrak KVN floating point formatting ... needs further study.

    
    for IN in tle csv xml kvn jsonarray; do

	      cat data/test.$IN | tletool transform -emit $OUT > data/check.$OUT

	      diff -Z data/test.$OUT data/check.$OUT > data/check-$IN-$OUT.diff ||
		  (echo "failed: $IN $OUT" data/check-$IN-$OUT.diff)

    done
    
done
