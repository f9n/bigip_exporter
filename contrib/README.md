# BIG-IP Exporter - Contrib

This directory contains example configurations and deployment files for running the BIG-IP Exporter.

## Files

### Configuration

- **[config.example.yaml](config.example.yaml)** - Example configuration file with all available options

### Deployment Examples

- **[docker-compose.yaml](docker-compose.yaml)** - Docker Compose deployment example
- **[bigip-exporter.service](bigip-exporter.service)** - Systemd service unit file
- **[kubernetes.yaml](kubernetes.yaml)** - Kubernetes deployment with ConfigMap, Secret, Service, and ServiceMonitor

## Usage

### Docker Compose

```bash
cd contrib/
cp config.example.yaml config.yaml
# Edit config.yaml with your BIG-IP credentials
docker-compose up -d
```

### Systemd

```bash
# Install binary
sudo cp bigip_exporter /usr/local/bin/
sudo chmod +x /usr/local/bin/bigip_exporter

# Create user
sudo useradd -r -s /sbin/nologin bigip-exporter

# Install service
sudo cp contrib/bigip-exporter.service /etc/systemd/system/
sudo systemctl daemon-reload
sudo systemctl enable bigip-exporter
sudo systemctl start bigip-exporter
sudo systemctl status bigip-exporter
```

### Kubernetes

```bash
# Create namespace
kubectl create namespace monitoring

# Update kubernetes.yaml with your credentials
kubectl apply -f contrib/kubernetes.yaml

# Check status
kubectl -n monitoring get pods -l app=bigip-exporter
kubectl -n monitoring logs -l app=bigip-exporter
```

## Configuration Methods

The exporter supports three configuration methods (in order of precedence):

1. **Environment Variables** (prefix: `BE_`)
   - Example: `BE_BIGIP_HOST=bigip.example.com`

2. **Configuration File**
   - Flag: `--config=/path/to/config.yaml`

3. **Command-line Flags**
   - Example: `bigip_exporter run --bigip.host=bigip.example.com`

## Environment Variables

All configuration options can be set via environment variables:

| Environment Variable | Description | Default |
|---------------------|-------------|---------|
| `BE_BIGIP_HOST` | BIG-IP hostname/IP | `localhost` |
| `BE_BIGIP_PORT` | BIG-IP management port | `443` |
| `BE_BIGIP_USERNAME` | BIG-IP username | `user` |
| `BE_BIGIP_PASSWORD` | BIG-IP password | `pass` |
| `BE_BIGIP_BASIC_AUTH` | Use HTTP Basic auth | `false` |
| `BE_EXPORTER_BIND_ADDRESS` | Exporter bind address | `localhost` |
| `BE_EXPORTER_BIND_PORT` | Exporter bind port | `9142` |
| `BE_EXPORTER_PARTITIONS` | Partitions to monitor | `` (all) |
| `BE_EXPORTER_NAMESPACE` | Prometheus namespace | `bigip` |
| `BE_EXPORTER_LOG_LEVEL` | Log level | `info` |

## Security Considerations

- **Never commit credentials** to version control
- Use **secrets management** (Kubernetes Secrets, Docker Secrets, Vault)
- Run as **non-root user**
- Use **read-only file systems** where possible
- Enable **TLS** for BIG-IP connections in production
- Rotate credentials regularly
