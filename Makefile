PKG=./...

run:
	@go run cmd/gentestcase/main.go $(ARGS)

test:
	@go test -v $(PKG)

build:
	@go build -o gentestcase cmd/gentestcase/main.go

schema:
	@go run cmd/schema/main.go
