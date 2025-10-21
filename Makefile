.PHONY: build test lint clean run run-light

# Build the application
build:
	CGO_CXXFLAGS="-Wno-vla-cxx-extension" go build -o midi-viewer ./cmd/midi-viewer

# Run tests
test:
	go test ./... -v

# Run tests with coverage
test-coverage:
	go test ./... -cover

# Run linter
lint:
	golint ./...

# Clean build artifacts
clean:
	rm -f midi-viewer

# Run with dark theme (default)
run: build
	./midi-viewer

# Run with light theme
run-light: build
	./midi-viewer -theme light

# Install dependencies
deps:
	go mod download
	go mod tidy

# Format code
fmt:
	go fmt ./...

# Run all checks (fmt, lint, test)
check: fmt lint test

# Install the binary to GOPATH/bin
install:
	CGO_CXXFLAGS="-Wno-vla-cxx-extension" go install ./cmd/midi-viewer
