# k8ctl

Enhanced kubectl with colored tables and built-in context/namespace management.

## Features

- **Colored Table Output**: Beautiful, color-coded tables with status indicators
- **Context & Namespace Management**: Built-in `ctx` and `ns` commands (replaces kubectx and kubens)
- **Enhanced Commands**: Improved `get`, `describe`, `logs`, and `watch` commands
- **Resource Search**: Fuzzy search across Kubernetes resources
- **Health Dashboard**: Comprehensive cluster and resource health overview
- **Simplified Port Forwarding**: Auto-discovery of ports and services
- **Resource Diff**: Compare resources across namespaces or contexts
- **Command Aliases**: Configurable shortcuts for common commands
- **Shell Completion**: Full completion support for bash, zsh, and fish

## Installation

### Manual Installation

```bash
# Download and install
curl -fsSL https://github.com/robertusnegoro/k8ctl/releases/latest/download/install.sh | bash

# Or download manually
wget https://github.com/robertusnegoro/k8ctl/releases/latest/download/k8ctl_linux_amd64.tar.gz
tar -xzf k8ctl_linux_amd64.tar.gz
sudo mv k8ctl /usr/local/bin/
```

### Homebrew (macOS/Linux)

```bash
brew tap robertusnegoro/k8ctl
brew install k8ctl
```

### APT (Debian/Ubuntu)

```bash
# Add repository
curl -fsSL https://robertusnegoro.github.io/k8ctl/deb/KEY.gpg | sudo gpg --dearmor -o /usr/share/keyrings/k8ctl-archive-keyring.gpg
echo "deb [signed-by=/usr/share/keyrings/k8ctl-archive-keyring.gpg] https://robertusnegoro.github.io/k8ctl/deb stable main" | sudo tee /etc/apt/sources.list.d/k8ctl.list

# Install
sudo apt update
sudo apt install k8ctl
```

### YUM/RPM (RHEL/CentOS/Fedora)

```bash
# Add repository
sudo tee /etc/yum.repos.d/k8ctl.repo <<EOF
[k8ctl]
name=k8ctl
baseurl=https://robertusnegoro.github.io/k8ctl/rpm/\$releasever/\$basearch
enabled=1
gpgcheck=1
gpgkey=https://robertusnegoro.github.io/k8ctl/rpm/KEY.gpg
EOF

# Install
sudo yum install k8ctl
# Or for Fedora/dnf
sudo dnf install k8ctl
```

## Quick Start

### Context Management

```bash
# List all contexts
k8ctl ctx

# Switch context
k8ctl ctx my-context
```

### Namespace Management

```bash
# List all namespaces
k8ctl ns

# Switch namespace
k8ctl ns production
```

### Enhanced Get Command

```bash
# Get pods with colored table output
k8ctl get pods

# Use shortcuts
k8ctl get po

# Get specific resource
k8ctl get pod my-pod

# Output as JSON or YAML
k8ctl get pods -o json
k8ctl get pods -o yaml
```

### Enhanced Describe

```bash
# Describe a pod with formatted output
k8ctl describe pod my-pod
```

### Enhanced Logs

```bash
# View logs with color-coded log levels
k8ctl logs my-pod

# Follow logs
k8ctl logs my-pod -f

# View logs from specific container
k8ctl logs my-pod -c my-container
```

### Watch Resources

```bash
# Watch pods for changes
k8ctl watch pods
```

### Resource Search

```bash
# Fuzzy search for resources
k8ctl search

# Search specific resource type
k8ctl search -t pods
```

### Health Dashboard

```bash
# Show cluster health
k8ctl health

# Show health for specific namespace
k8ctl health -n production
```

### Port Forwarding

```bash
# Port forward with auto-discovery
k8ctl port-forward my-pod

# Specify ports
k8ctl port-forward my-pod -l 8080 -r 80

# Port forward to service
k8ctl port-forward svc/my-service
```

### Resource Diff

```bash
# Compare pods across namespaces
k8ctl diff pod my-pod --namespace1 dev --namespace2 prod
```

### Command Aliases

Pre-configured aliases:
- `k8ctl g` → `k8ctl get`
- `k8ctl d` → `k8ctl describe`
- `k8ctl l` → `k8ctl logs`
- `k8ctl w` → `k8ctl watch`
- `k8ctl pf` → `k8ctl port-forward`
- `k8ctl h` → `k8ctl health`

### Shell Completion

```bash
# Bash
source <(k8ctl completion bash)
k8ctl completion bash > /etc/bash_completion.d/k8ctl

# Zsh
source <(k8ctl completion zsh)
k8ctl completion zsh > "${fpath[1]}/_k8ctl"

# Fish
k8ctl completion fish | source
k8ctl completion fish > ~/.config/fish/completions/k8ctl.fish
```

## Configuration

Configuration is stored in `~/.k8ctl/config.yaml`:

```yaml
current_context: my-context
current_namespace: production
aliases:
  g: get
  d: describe
  l: logs
colors:
  enabled: true
output:
  format: table
```

## Building from Source

```bash
# Clone repository
git clone https://github.com/robertusnegoro/k8ctl.git
cd k8ctl

# Build
make build

# Install
make install

# Run tests
make test

# Build for all platforms
make build-all
```

## Requirements

- Go 1.21 or later
- Kubernetes cluster access (via kubeconfig)
- kubectl (for reference, k8ctl is a replacement)

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

MIT License - see LICENSE file for details

## Acknowledgments

- Built with [Cobra](https://github.com/spf13/cobra)
- Uses [kubernetes/client-go](https://github.com/kubernetes/client-go)
- Inspired by kubectx and kubens
