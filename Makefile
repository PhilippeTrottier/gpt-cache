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

.PHONY: generate-swagger
generate-swagger:
	mkdir -p pkg/api
	oapi-codegen -package api -generate types,chi-server,spec api/swagger.yaml > pkg/api/gptcacheapi.gen.go

.PHONY: fmt
fmt:
	go fmt ./...
