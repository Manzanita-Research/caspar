.PHONY: build test clean install

build:
	go build -o caspar .

test:
	go test ./...

clean:
	rm -f caspar

install: build
	cp caspar $(GOPATH)/bin/caspar 2>/dev/null || cp caspar /usr/local/bin/caspar
