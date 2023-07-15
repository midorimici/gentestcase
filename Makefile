PKG=./...

run:
	@go run cmd/main/main.go

test:
	@go test -v $(PKG)
