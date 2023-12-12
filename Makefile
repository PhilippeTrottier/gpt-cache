.PHONY: build
build:
	go build -o bin/server cmd/server/main.go

.PHONY: clean
clean:
	rm -r bin

.PHONY: run
run:
	go run cmd/server/main.go

.PHONY: test
test:
	go test -v -race -count=1 ./...