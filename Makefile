fmt:
	go fmt ./...
	go mod tidy

test:
	go test -race -cover ./...

lint:
	 golangci-lint run --enable-all --disable gomnd --disable goerr113