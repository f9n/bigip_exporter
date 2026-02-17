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

### Virtual Servers (`bigip_vs_*`)

Labels: `partition`, `vs`

| Metric | Type | Description |
|--------|------|-------------|
| `bigip_vs_status_availability_state` | Gauge | Availability (1=available, 0=unavailable) |
| `bigip_vs_status_enabled_state` | Gauge | Enabled state (1=enabled, 0=disabled) |
| `bigip_vs_clientside_bytes_in` | Counter | Clientside bytes in |
| `bigip_vs_clientside_bytes_out` | Counter | Clientside bytes out |
| `bigip_vs_clientside_cur_conns` | Gauge | Current clientside connections |
| `bigip_vs_clientside_tot_conns` | Counter | Total clientside connections |
| `bigip_vs_clientside_max_conns` | Counter | Max clientside connections |
| `bigip_vs_clientside_pkts_in` | Counter | Clientside packets in |
| `bigip_vs_clientside_pkts_out` | Counter | Clientside packets out |
| `bigip_vs_clientside_evicted_conns` | Counter | Clientside evicted connections |
| `bigip_vs_clientside_slow_killed` | Counter | Clientside slow killed |
| `bigip_vs_ephemeral_bytes_in` | Counter | Ephemeral bytes in |
| `bigip_vs_ephemeral_bytes_out` | Counter | Ephemeral bytes out |
| `bigip_vs_ephemeral_cur_conns` | Gauge | Current ephemeral connections |
| `bigip_vs_ephemeral_tot_conns` | Counter | Total ephemeral connections |
| `bigip_vs_ephemeral_max_conns` | Counter | Max ephemeral connections |
| `bigip_vs_ephemeral_pkts_in` | Counter | Ephemeral packets in |
| `bigip_vs_ephemeral_pkts_out` | Counter | Ephemeral packets out |
| `bigip_vs_ephemeral_evicted_conns` | Counter | Ephemeral evicted connections |
| `bigip_vs_ephemeral_slow_killed` | Counter | Ephemeral slow killed |
| `bigip_vs_tot_requests` | Counter | Total requests |
| `bigip_vs_syncookie_accepts` | Counter | SYN cookie accepts |
| `bigip_vs_syncookie_rejects` | Counter | SYN cookie rejects |
| `bigip_vs_syncookie_syncookies` | Counter | SYN cookies created |
| `bigip_vs_syncookie_syncache_curr` | Gauge | Current SYN cache entries |
| `bigip_vs_syncookie_syncache_over` | Counter | SYN cache overflow |
| `bigip_vs_syncookie_hw_accepts` | Counter | Hardware SYN cookie accepts |
| `bigip_vs_syncookie_hw_syncookies` | Counter | Hardware SYN cookies |
| `bigip_vs_syncookie_hwsyncookie_instance` | Counter | Hardware SYN cookie instances |
| `bigip_vs_syncookie_swsyncookie_instance` | Counter | Software SYN cookie instances |
| `bigip_vs_cs_min_conn_dur` | Gauge | Min clientside connection duration |
| `bigip_vs_cs_mean_conn_dur` | Gauge | Mean clientside connection duration |
| `bigip_vs_cs_max_conn_dur` | Counter | Max clientside connection duration |
| `bigip_vs_five_sec_avg_usage_ratio` | Gauge | 5-second average usage ratio |
| `bigip_vs_one_min_avg_usage_ratio` | Gauge | 1-minute average usage ratio |
| `bigip_vs_five_min_avg_usage_ratio` | Gauge | 5-minute average usage ratio |

### Pools (`bigip_pool_*`)

Labels: `partition`, `pool`

| Metric | Type | Description |
|--------|------|-------------|
| `bigip_pool_status_availability_state` | Gauge | Availability (1=available, 0=unavailable) |
| `bigip_pool_status_enabled_state` | Gauge | Enabled state (1=enabled, 0=disabled) |
| `bigip_pool_serverside_bytes_in` | Counter | Serverside bytes in |
| `bigip_pool_serverside_bytes_out` | Counter | Serverside bytes out |
| `bigip_pool_serverside_cur_conns` | Gauge | Current serverside connections |
| `bigip_pool_serverside_tot_conns` | Counter | Total serverside connections |
| `bigip_pool_serverside_max_conns` | Counter | Max serverside connections |
| `bigip_pool_serverside_pkts_in` | Counter | Serverside packets in |
| `bigip_pool_serverside_pkts_out` | Counter | Serverside packets out |
| `bigip_pool_tot_requests` | Counter | Total requests |
| `bigip_pool_cur_sessions` | Gauge | Current sessions |
| `bigip_pool_active_member_cnt` | Gauge | Active member count |
| `bigip_pool_member_total_cnt` | Gauge | Total member count |
| `bigip_pool_min_active_members` | Gauge | Minimum active members |
| `bigip_pool_connq_depth` | Gauge | Connection queue depth |
| `bigip_pool_connq_serviced` | Counter | Connection queue serviced |
| `bigip_pool_connq_age_head` | Gauge | Connection queue age (head) |
| `bigip_pool_connq_age_max` | Counter | Connection queue age (max) |
| `bigip_pool_connq_age_ema` | Gauge | Connection queue age (EMA) |
| `bigip_pool_connq_age_edm` | Gauge | Connection queue age (EDM) |
| `bigip_pool_connq_all_depth` | Gauge | Connection queue all depth |
| `bigip_pool_connq_all_serviced` | Counter | Connection queue all serviced |
| `bigip_pool_connq_all_age_head` | Gauge | Connection queue all age (head) |
| `bigip_pool_connq_all_age_max` | Counter | Connection queue all age (max) |
| `bigip_pool_connq_all_age_ema` | Gauge | Connection queue all age (EMA) |
| `bigip_pool_connq_all_age_edm` | Gauge | Connection queue all age (EDM) |

### Pool Members (`bigip_pool_member_*`)

Labels: `partition`, `pool`, `member`

| Metric | Type | Description |
|--------|------|-------------|
| `bigip_pool_member_status_availability_state` | Gauge | Availability (1=available, 0=unavailable) |
| `bigip_pool_member_status_enabled_state` | Gauge | Enabled state (1=enabled, 0=disabled) |
| `bigip_pool_member_serverside_bytes_in` | Counter | Serverside bytes in |
| `bigip_pool_member_serverside_bytes_out` | Counter | Serverside bytes out |
| `bigip_pool_member_serverside_cur_conns` | Gauge | Current serverside connections |
| `bigip_pool_member_serverside_tot_conns` | Counter | Total serverside connections |
| `bigip_pool_member_serverside_max_conns` | Counter | Max serverside connections |
| `bigip_pool_member_serverside_pkts_in` | Counter | Serverside packets in |
| `bigip_pool_member_serverside_pkts_out` | Counter | Serverside packets out |
| `bigip_pool_member_tot_requests` | Counter | Total requests |
| `bigip_pool_member_cur_sessions` | Gauge | Current sessions |

### Nodes (`bigip_node_*`)

Labels: `partition`, `node`

| Metric | Type | Description |
|--------|------|-------------|
| `bigip_node_status_availability_state` | Gauge | Availability (1=available, 0=unavailable) |
| `bigip_node_status_enabled_state` | Gauge | Enabled state (1=enabled, 0=disabled) |
| `bigip_node_monitor_status` | Gauge | Monitor status (1=up, 0=down) |
| `bigip_node_session_status` | Gauge | Session status (1=enabled, 0=disabled) |
| `bigip_node_serverside_bytes_in` | Counter | Serverside bytes in |
| `bigip_node_serverside_bytes_out` | Counter | Serverside bytes out |
| `bigip_node_serverside_cur_conns` | Gauge | Current serverside connections |
| `bigip_node_serverside_tot_conns` | Counter | Total serverside connections |
| `bigip_node_serverside_max_conns` | Counter | Max serverside connections |
| `bigip_node_serverside_pkts_in` | Counter | Serverside packets in |
| `bigip_node_serverside_pkts_out` | Counter | Serverside packets out |
| `bigip_node_tot_requests` | Counter | Total requests |
| `bigip_node_cur_sessions` | Gauge | Current sessions |

### iRules (`bigip_rule_*`)

Labels: `partition`, `rule`, `event`

| Metric | Type | Description |
|--------|------|-------------|
| `bigip_rule_priority` | Gauge | Rule priority |
| `bigip_rule_total_executions` | Counter | Total executions |
| `bigip_rule_failures` | Counter | Rule failures |
| `bigip_rule_aborts` | Counter | Rule aborts |
| `bigip_rule_min_cycles` | Gauge | Min CPU cycles per execution |
| `bigip_rule_max_cycles` | Counter | Max CPU cycles per execution |
| `bigip_rule_avg_cycles` | Gauge | Average CPU cycles per execution |

### Exporter Info

| Metric | Type | Description |
|--------|------|-------------|
| `bigip_exporter_build_info` | Gauge | Build info (constant 1). Labels: `version`, `commit`, `build_date`, `go_version` |

### Collector Metadata

| Metric | Type | Description |
|--------|------|-------------|
| `bigip_collector_scrape_status` | Gauge | Scrape success (1=ok, 0=fail). Labels: `collector` |
| `bigip_collector_scrape_duration` | Summary | Scrape duration in seconds. Labels: `collector` |
| `bigip_total_scrape_duration` | Summary | Total scrape duration across all collectors |

## 📈 Grafana Dashboard

A pre-built Grafana dashboard is available via the [Jsonnet mixin](docs/bigip-mixin/).

To generate the dashboard JSON:

```bash
cd docs/bigip-mixin
make generated
```

Then import `docs/bigip-mixin/generated/bigip.json` into Grafana.

The dashboard includes panels for virtual servers, pools, pool members, nodes, and iRules with template variables for instance and partition filtering.

## 🔧 Development

See [DEVELOPMENT.md](DEVELOPMENT.md) for building from source, running tests, project structure, and contributing guidelines.

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
