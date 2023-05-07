test:
	@echo "Running tests..."
	go test -v ./... --count=1

check:
	@echo "Formatting and Checking code..."
	go mod tidy
	go fmt ./...
	go vet ./...
	golangci-lint run

update:
	@echo "Update dependencies..."
	go get -u -t
	go mod tidy

prereqs:
	@echo "Installing prerequisites..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/vektra/mockery/v2@v2.20.0
