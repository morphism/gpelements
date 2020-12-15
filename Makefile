all:
	cd cmd/tletool && go install

test:
	go test -v -bench=.
