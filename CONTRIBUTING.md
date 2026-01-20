# Contributing to k8ctl

## Local Development Setup

### Prerequisites

- Go 1.21 or later
- Kubernetes cluster access (via kubeconfig)
- Make (optional, for convenience)

### Building Locally

```bash
# Clone the repository
git clone https://github.com/robertusnegoro/k8ctl.git
cd k8ctl

# Install dependencies
go mod download

# Build the binary
make build
# or
go build -o bin/k8ctl ./cmd/k8ctl

# Install to system (optional)
make install
# or
sudo cp bin/k8ctl /usr/local/bin/
```

### Running Tests

```bash
# Run all tests
make test
# or
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Run specific package tests
go test ./internal/config -v
go test ./internal/output -v
go test ./internal/errors -v
```

### Testing the CLI Locally

#### 1. Basic Functionality Test

```bash
# Build the binary
make build

# Test version command
./bin/k8ctl version

# Test help
./bin/k8ctl --help

# Test context listing (requires kubeconfig)
./bin/k8ctl ctx

# Test namespace listing (requires Kubernetes cluster)
./bin/k8ctl ns
```

#### 2. With a Kubernetes Cluster

**Prerequisites:**
- A running Kubernetes cluster (minikube, kind, or cloud cluster)
- Valid kubeconfig file at `~/.kube/config`

```bash
# List contexts
./bin/k8ctl ctx

# Switch context
./bin/k8ctl ctx <context-name>

# List namespaces
./bin/k8ctl ns

# Switch namespace
./bin/k8ctl ns <namespace>

# Get pods with colored output
./bin/k8ctl get pods

# Get pods with shortcuts
./bin/k8ctl get po

# Describe a pod
./bin/k8ctl describe pod <pod-name>

# View logs
./bin/k8ctl logs <pod-name>

# Health dashboard
./bin/k8ctl health

# Search resources
./bin/k8ctl search
```

#### 3. Testing Without a Cluster (Mock Mode)

For testing without a real cluster, you can:

1. **Use a test cluster** (recommended):
   ```bash
   # Start minikube
   minikube start
   
   # Or use kind
   kind create cluster
   ```

2. **Mock Kubernetes API** (for unit tests):
   - Use `k8s.io/client-go/testing` package
   - Create fake clients for testing

### Development Workflow

#### TDD Approach

1. **Write test first** (Red):
   ```bash
   # Create test file: internal/feature/feature_test.go
   # Write failing test
   ```

2. **Implement feature** (Green):
   ```bash
   # Write minimal code to pass test
   ```

3. **Refactor**:
   ```bash
   # Improve code while keeping tests green
   ```

4. **Run tests**:
   ```bash
   go test ./internal/feature -v
   ```

#### Adding New Commands

1. Create command file: `internal/commands/newcommand.go`
2. Write tests: `internal/commands/newcommand_test.go`
3. Register command in `internal/commands/root.go`
4. Test the command:
   ```bash
   make build
   ./bin/k8ctl newcommand --help
   ```

### Debugging

#### Enable Verbose Logging

```bash
# Run with debug output
./bin/k8ctl --help -v
```

#### Check Configuration

```bash
# View config file
cat ~/.k8ctl/config.yaml

# Edit config manually
vim ~/.k8ctl/config.yaml
```

### Common Issues

#### Issue: "no such host" or "connection refused"

**Solution**: Check your kubeconfig:
```bash
kubectl cluster-info
./bin/k8ctl ctx
```

#### Issue: "context not found"

**Solution**: List available contexts:
```bash
kubectl config get-contexts
./bin/k8ctl ctx
```

#### Issue: Tests failing

**Solution**: 
```bash
# Clean and rebuild
go clean -cache
go mod tidy
go test ./...
```

### Code Style

- Follow Go conventions: https://golang.org/doc/effective_go
- Run linters:
  ```bash
  make lint
  # or
  golangci-lint run
  ```

### Submitting Changes

1. Create a feature branch
2. Write tests (TDD approach)
3. Implement feature
4. Ensure all tests pass
5. Update documentation
6. Submit pull request

### Project Structure

```
k8ctl/
├── cmd/k8ctl/          # Main entry point
├── internal/
│   ├── commands/       # CLI commands
│   ├── config/         # Configuration management
│   ├── k8s/            # Kubernetes client
│   ├── output/         # Output formatting
│   ├── errors/         # Error handling
│   └── aliases/       # Command aliases
├── pkg/                # Public packages
├── packaging/          # Distribution files
└── .github/workflows/  # CI/CD
```

### Next Steps

1. **Complete remaining commands**: Add full implementations for get, describe, logs, etc.
2. **Add integration tests**: Test with real Kubernetes clusters
3. **Improve error handling**: Add more user-friendly error messages
4. **Add more resource types**: Extend get command to support all Kubernetes resources
5. **Performance optimization**: Profile and optimize slow operations
6. **Documentation**: Add more examples and use cases
