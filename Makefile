build:
	go build -o bin/server cmd/server/main.go

test:
	go test -v -race ./...