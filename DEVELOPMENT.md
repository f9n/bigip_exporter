# Development

## Prerequisites

- Go 1.26 or later
- Make (optional)

## Building from Source

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

## Running Tests

```bash
go test ./...
```

## Local Release Build

```bash
# Install goreleaser
go install github.com/goreleaser/goreleaser@latest

# Create a snapshot release
goreleaser release --snapshot --clean

# Check dist/ directory
ls -lh dist/
```

## Project Structure

```
bigip_exporter/
├── cmd/
│   └── bigip_exporter/    # Main application
│       └── main.go
├── internal/
│   ├── collector/         # Prometheus collectors
│   └── config/            # Configuration
├── contrib/               # Example configs & deployment files
├── docs/
│   └── bigip-mixin/       # Grafana dashboard (Jsonnet mixin)
├── .build/                # Build artifacts
└── .github/               # CI/CD workflows
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request
