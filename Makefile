.PHONY: build test clean install ghost-up ghost-down ghost-logs

build:
	go build -o caspar .

test:
	go test ./...

clean:
	rm -f caspar

install: build
	cp caspar $(GOPATH)/bin/caspar 2>/dev/null || cp caspar /usr/local/bin/caspar

ghost-up:
	docker compose up -d

ghost-down:
	docker compose down

ghost-logs:
	docker compose logs -f ghost
