.PHONY: all
all: lint build test coverage

.PHONY: lint
lint:
	gometalinter $(go list ./... | grep -v /vendor/)

.PHONY: build
build:
	go build -v ./...

.PHONY: test
test:
	go test -v ./...

.PHONY: coverage
coverage:

.PHONY: runsvc
runsvc: build
	./hoardsvc

.PHONY: runsvctest
runsvctest: build
	./hoardsvc -test

.PHONY: clean
clean:
	rm -f hoardsvc
