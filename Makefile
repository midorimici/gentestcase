PKG=./...

run:
	@go run cmd/integtest/integtest.go

test:
	@go test -v $(PKG)
