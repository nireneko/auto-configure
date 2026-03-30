.PHONY: build test lint run clean

build:
	go build -o bin/so-install ./cmd/so-install

test:
	go test -v -race -count=1 ./...

lint:
	go vet ./...

clean:
	rm -rf bin/ coverage.out

run: build
	sudo ./bin/so-install
