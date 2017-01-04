GO_FMT = gofmt -s -w -l .

all: deps compile

compile:
	go build ./...

deps:
	go get

format:
	$(GO_FMT)
