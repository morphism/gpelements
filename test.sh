#!/bin/bash

# Try to compare transformations based on Celestrak formats.  See
# notes.md.

set -e pipefail

for OUT in tle csv xml kvn jsonarray; do

    for IN in tle csv xml kvn jsonarray; do

	      cat data/test.$IN | tletool transform -emit $OUT > data/check.$OUT

	      diff -Z data/test.$OUT data/check.$OUT > data/check-$IN-$OUT.diff ||
		  (echo "failed: $IN $OUT" data/check-$IN-$OUT.diff)

    done
    
done
