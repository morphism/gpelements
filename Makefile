all:
	cd cmd/tletool && go install

install: all


test:
	go test -v -bench=.
