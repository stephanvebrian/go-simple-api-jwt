test:
	@echo ">  Running tests..."
	go test -cover -race -v ./...

build:
	go build -o ./bin/go-simple-api-jwt ./cmd

run: build
	./bin/go-simple-api-jwt