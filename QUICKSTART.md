# Quick Start Guide

## Build and Run Locally

### 1. Build the Binary

```bash
# Using Makefile
make build

# Or directly with go
go build -o bin/k8ctl ./cmd/k8ctl
```

### 2. Test Basic Commands (No Cluster Required)

```bash
# Check version
./bin/k8ctl version

# View help
./bin/k8ctl --help

# Test shell completion
./bin/k8ctl completion bash
```

### 3. Test with Kubernetes Cluster

**Prerequisites:**
- A running Kubernetes cluster
- Valid kubeconfig at `~/.kube/config`

```bash
# List contexts
./bin/k8ctl ctx

# Switch context
./bin/k8ctl ctx <context-name>

# List namespaces
./bin/k8ctl ns

# Switch namespace
./bin/k8ctl ns <namespace>

# Get pods (with colored output!)
./bin/k8ctl get pods
./bin/k8ctl get po  # shortcut

# Describe a pod
./bin/k8ctl describe pod <pod-name>

# View logs
./bin/k8ctl logs <pod-name>

# Health dashboard
./bin/k8ctl health

# Search resources
./bin/k8ctl search
```

### 4. Install System-Wide (Optional)

```bash
# Install to /usr/local/bin
make install

# Or manually
sudo cp bin/k8ctl /usr/local/bin/

# Verify installation
k8ctl version
```

## Running Tests

```bash
# Run all tests
make test

# Run with coverage
go test -cover ./...

# Run specific package tests
go test ./internal/config -v
go test ./internal/output -v
go test ./internal/errors -v
```

## Development Workflow

```bash
# 1. Make changes to code

# 2. Run tests
make test

# 3. Build
make build

# 4. Test manually
./bin/k8ctl <command>

# 5. Repeat
```

## Next Steps

1. **Extend Get Command**: Add support for more Kubernetes resources
2. **Add Integration Tests**: Test with real Kubernetes clusters
3. **Improve Error Handling**: Add more user-friendly messages
4. **Performance**: Profile and optimize slow operations

## Troubleshooting

### Build Issues
```bash
go clean -cache
go mod tidy
make build
```

### Test Issues
```bash
go test -v ./...
```

### Runtime Issues
```bash
# Check kubeconfig
kubectl config view

# Verify cluster access
kubectl cluster-info
```
