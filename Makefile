PKG=./...

run:
	@go run main.go run $(ARGS)

test:
	@go test -v $(PKG)

build:
	@go build -o gentestcase main.go

schema:
	@go run main.go schema $(ARGS)
