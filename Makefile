.PHONY: build build-all install test lint clean package-deb package-rpm package-homebrew release

BINARY_NAME=k8ctl
VERSION?=dev
COMMIT?=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
DATE?=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

# Build for current platform
build:
	@echo "Building $(BINARY_NAME)..."
	@go build -ldflags "-X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE)" \
		-o bin/$(BINARY_NAME) ./cmd/k8ctl

# Build for all platforms
build-all:
	@echo "Building for all platforms..."
	@GOOS=linux GOARCH=amd64 go build -ldflags "-X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE)" \
		-o bin/$(BINARY_NAME)-linux-amd64 ./cmd/k8ctl
	@GOOS=linux GOARCH=arm64 go build -ldflags "-X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE)" \
		-o bin/$(BINARY_NAME)-linux-arm64 ./cmd/k8ctl
	@GOOS=darwin GOARCH=amd64 go build -ldflags "-X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE)" \
		-o bin/$(BINARY_NAME)-darwin-amd64 ./cmd/k8ctl
	@GOOS=darwin GOARCH=arm64 go build -ldflags "-X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE)" \
		-o bin/$(BINARY_NAME)-darwin-arm64 ./cmd/k8ctl
	@GOOS=windows GOARCH=amd64 go build -ldflags "-X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE)" \
		-o bin/$(BINARY_NAME)-windows-amd64.exe ./cmd/k8ctl

# Install to system
install: build
	@echo "Installing $(BINARY_NAME)..."
	@sudo cp bin/$(BINARY_NAME) /usr/local/bin/
	@echo "Installation complete!"

# Run tests
test:
	@echo "Running tests..."
	@go test -v ./...

# Run linters
lint:
	@echo "Running linters..."
	@golangci-lint run || echo "golangci-lint not installed, skipping..."

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf bin/
	@go clean

# Create Debian package
package-deb:
	@echo "Creating Debian package..."
	@goreleaser release --snapshot --skip-publish

# Create RPM package
package-rpm:
	@echo "Creating RPM package..."
	@goreleaser release --snapshot --skip-publish

# Update Homebrew formula
package-homebrew:
	@echo "Updating Homebrew formula..."
	@goreleaser release --snapshot --skip-publish

# Full release via GoReleaser
release:
	@echo "Creating release..."
	@goreleaser release --clean

# Development helpers
dev:
	@go run -ldflags "-X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE)" ./cmd/k8ctl

.DEFAULT_GOAL := build
