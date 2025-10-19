build:
	@go build -o bin/kriminoDB

run: build
	@./bin/kriminoDB

test:
	@go test ./... -v
