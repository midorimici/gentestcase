PKG=./...

run:
	@go run cmd/gentestcase/main.go

test:
	@go test -v $(PKG)

build:
	@go build -o gentestcase cmd/gentestcase/main.go
