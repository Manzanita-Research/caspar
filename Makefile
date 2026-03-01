.PHONY: build test clean install

build:
	go build -o ghostctl .

test:
	go test ./...

clean:
	rm -f ghostctl

install: build
	cp ghostctl $(GOPATH)/bin/ghostctl 2>/dev/null || cp ghostctl /usr/local/bin/ghostctl
