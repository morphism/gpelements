#!/bin/bash

# Try to compare transformations based on (canonicalized?) Celestrak
# formats.  Base data obtained via get-test-data.sh. Also see
# notes.md.

set -e pipefail

# Canonicalize base data (with the realization that we might mess up).
for F in tle csv xml kvn jsonarray; do
    cat data/test.$F | tletool transform -emit $F > data/canonical.$F
done

cp data/test.csv data/test.csvh

# Compare transformed data to the canonical data.
for OUT in tle csv csvh xml kvn jsonarray; do

    for IN in tle csv csvh xml kvn jsonarray; do

	      cat data/test.$IN | tletool transform -emit $OUT > data/check.$OUT

	      (diff -Z data/canonical.$OUT data/check.$OUT > data/check-$IN-$OUT.diff &&
	       echo "passed: $IN $OUT" data/check-$IN-$OUT.diff) ||
		  (echo "failed: $IN $OUT" data/check-$IN-$OUT.diff)

    done
    
done
