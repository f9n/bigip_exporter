# BIG-IP Exporter

[![Go Version](https://img.shields.io/badge/Go-1.26+-00ADD8?style=flat&logo=go)](https://go.dev)
[![Release](https://img.shields.io/github/v/release/f9n/bigip_exporter?style=flat)](https://github.com/f9n/bigip_exporter/releases)
[![License](https://img.shields.io/github/license/f9n/bigip_exporter?style=flat)](LICENSE)
[![Docker](https://img.shields.io/badge/Docker-ghcr.io-2496ED?style=flat&logo=docker)](https://ghcr.io/f9n/bigip_exporter)

Prometheus exporter for F5 BIG-IP statistics. Exports metrics from BIG-IP systems using the iControl REST API.

## 🚀 Quick Start

### Binary

```bash
# Download latest release
curl -LO https://github.com/f9n/bigip_exporter/releases/latest/download/bigip_exporter_Linux_x86_64.tar.gz
tar -xzf bigip_exporter_Linux_x86_64.tar.gz

# Run
./bigip_exporter run \
  --bigip.host=bigip.example.com \
  --bigip.username=admin \
  --bigip.password=admin
```

### Docker

```bash
docker run -d \
  -p 9142:9142 \
  -e BE_BIGIP_HOST=bigip.example.com \
  -e BE_BIGIP_USERNAME=admin \
  -e BE_BIGIP_PASSWORD=admin \
  ghcr.io/f9n/bigip_exporter:latest
```

### Docker Compose

```bash
cd contrib/
docker-compose up -d
```

### Kubernetes

```bash
kubectl apply -f contrib/kubernetes.yaml
```

Access metrics at `http://localhost:9142/metrics`

## 📦 Installation

### Download Binary

Download the latest release for your platform from [GitHub Releases](https://github.com/f9n/bigip_exporter/releases).

Available platforms:
- Linux (amd64, arm64)
- macOS (amd64, arm64)
- Windows (amd64, arm64)

### Docker Image

```bash
# Pull latest
docker pull ghcr.io/f9n/bigip_exporter:latest

# Pull specific version
docker pull ghcr.io/f9n/bigip_exporter:v1.0.0

# Pull specific architecture
docker pull ghcr.io/f9n/bigip_exporter:latest-arm64
```

### Systemd Service

```bash
# Install binary
sudo cp bigip_exporter /usr/local/bin/
sudo chmod +x /usr/local/bin/bigip_exporter

# Create user
sudo useradd -r -s /sbin/nologin bigip-exporter

# Install service
sudo cp contrib/bigip-exporter.service /etc/systemd/system/
sudo systemctl daemon-reload
sudo systemctl enable --now bigip-exporter
```

## ⚙️ Configuration

The exporter supports three configuration methods (in order of precedence):

1. **Command-line flags**
2. **Environment variables** (prefix: `BE_`)
3. **Configuration file**

### Command-line Flags

```bash
./bigip_exporter run \
  --bigip.host=bigip.example.com \
  --bigip.port=443 \
  --bigip.username=admin \
  --bigip.password=admin \
  --bigip.basic_auth=false \
  --exporter.bind_address=0.0.0.0 \
  --exporter.bind_port=9142 \
  --exporter.partitions="Common,Production" \
  --exporter.namespace=bigip \
  --exporter.log_level=info
```

### Environment Variables

```bash
export BE_BIGIP_HOST=bigip.example.com
export BE_BIGIP_PORT=443
export BE_BIGIP_USERNAME=admin
export BE_BIGIP_PASSWORD=admin
export BE_EXPORTER_BIND_PORT=9142
export BE_EXPORTER_LOG_LEVEL=info

./bigip_exporter run
```

### Configuration File

```bash
./bigip_exporter run --config=/etc/bigip_exporter/config.yaml
```

Example configuration: [contrib/config.example.yaml](contrib/config.example.yaml)

### Configuration Options

| Flag | Environment Variable | Description | Default |
|------|---------------------|-------------|---------|
| `--bigip.host` | `BE_BIGIP_HOST` | BIG-IP hostname/IP | `localhost` |
| `--bigip.port` | `BE_BIGIP_PORT` | BIG-IP management port | `443` |
| `--bigip.username` | `BE_BIGIP_USERNAME` | BIG-IP username | `user` |
| `--bigip.password` | `BE_BIGIP_PASSWORD` | BIG-IP password | `pass` |
| `--bigip.basic_auth` | `BE_BIGIP_BASIC_AUTH` | Use HTTP Basic auth | `false` |
| `--exporter.bind_address` | `BE_EXPORTER_BIND_ADDRESS` | Bind address | `localhost` |
| `--exporter.bind_port` | `BE_EXPORTER_BIND_PORT` | Bind port | `9142` |
| `--exporter.partitions` | `BE_EXPORTER_PARTITIONS` | Comma-separated partitions | `` (all) |
| `--exporter.namespace` | `BE_EXPORTER_NAMESPACE` | Prometheus namespace | `bigip` |
| `--exporter.log_level` | `BE_EXPORTER_LOG_LEVEL` | Log level (debug/info/warn/error) | `info` |
| `--config` | `BE_CONFIG` | Config file path | `` |

## 🎯 CLI Commands

The exporter uses a modern CLI structure with subcommands:

```bash
# Show help and available commands
./bigip_exporter --help

# Start the exporter (main command)
./bigip_exporter run [flags]

# Show version information
./bigip_exporter version

# Generate shell completion scripts
./bigip_exporter completion bash   # For Bash
./bigip_exporter completion zsh    # For Zsh
./bigip_exporter completion fish   # For Fish
```

### Shell Completion

Enable shell completion for a better CLI experience:

```bash
# Bash
./bigip_exporter completion bash > /etc/bash_completion.d/bigip_exporter

# Zsh
./bigip_exporter completion zsh > "${fpath[1]}/_bigip_exporter"

# Fish
./bigip_exporter completion fish > ~/.config/fish/completions/bigip_exporter.fish
```

## 📊 Exported Metrics

### Virtual Servers

- `bigip_vs_status` - Virtual server availability status
- `bigip_vs_enabled_state` - Virtual server enabled state
- `bigip_vs_clientside_bits_in` - Clientside bits in
- `bigip_vs_clientside_bits_out` - Clientside bits out
- `bigip_vs_clientside_cur_conns` - Current client connections
- And more...

### Pools

- `bigip_pool_status` - Pool availability status
- `bigip_pool_available_members` - Number of available pool members
- `bigip_pool_up_members` - Number of up pool members
- `bigip_pool_serverside_bits_in` - Serverside bits in
- And more...

### Nodes

- `bigip_node_status` - Node availability status
- `bigip_node_session_status` - Node session status
- `bigip_node_serverside_bits_in` - Serverside bits in
- And more...

### Rules (iRules)

- `bigip_rule_aborts` - Rule aborts
- `bigip_rule_avg_cycles` - Average cycles per execution
- `bigip_rule_executions` - Total executions
- `bigip_rule_failures` - Total failures
- And more...

All metrics include labels for `partition` and resource name.

## 🔧 Development

### Prerequisites

- Go 1.26 or later
- Make (optional)

### Building from Source

```bash
# Clone repository
git clone https://github.com/f9n/bigip_exporter.git
cd bigip_exporter

# Install dependencies
go mod download

# Build
go build -v ./cmd/bigip_exporter

# Run
./bigip_exporter run --bigip.host=bigip.example.com
```

### Running Tests

```bash
go test ./...
```

### Local Release Build

```bash
# Install goreleaser
go install github.com/goreleaser/goreleaser@latest

# Create a snapshot release
goreleaser release --snapshot --clean

# Check dist/ directory
ls -lh dist/
```

### Project Structure

```
bigip_exporter/
├── cmd/
│   └── bigip_exporter/    # Main application
│       └── main.go
├── internal/
│   ├── collector/         # Prometheus collectors
│   └── config/            # Configuration
├── contrib/               # Example configs & deployment files
├── .build/                # Build artifacts
└── .github/               # CI/CD workflows
```

## 📝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 🐛 Reporting Issues

If you find a bug or have a feature request, please create an issue on [GitHub Issues](https://github.com/f9n/bigip_exporter/issues).

## 📚 Resources

- [Prometheus](https://prometheus.io/)
- [F5 BIG-IP iControl REST API](https://clouddocs.f5.com/api/icontrol-rest/)
- [Example Configurations](contrib/)
- [Releases](https://github.com/f9n/bigip_exporter/releases)

## 🔐 Security

- Never commit credentials to version control
- Use Kubernetes Secrets, Docker Secrets, or environment variables for sensitive data
- Run containers as non-root user (automatically configured)
- Use TLS for BIG-IP connections in production

## 📄 License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Original project by [ExpressenAB](https://github.com/ExpressenAB/bigip_exporter)
- F5 for the BIG-IP platform
- Prometheus community

## 📞 Support

- 📖 Documentation: [README](README.md) and [contrib/](contrib/)
- 🐛 Bug Reports: [GitHub Issues](https://github.com/f9n/bigip_exporter/issues)
- 💬 Discussions: [GitHub Discussions](https://github.com/f9n/bigip_exporter/discussions)

---

**Made with ❤️ for the DevOps community**
