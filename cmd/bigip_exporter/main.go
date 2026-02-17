package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/f9n/bigip_exporter/internal/collector"
	"github.com/f9n/bigip_exporter/internal/config"
	"github.com/pr8kerl/f5er/f5"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	logger  = slog.New(slog.NewJSONHandler(os.Stdout, nil))
	cfgFile string
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

// rootCmd represents the base command
var rootCmd = &cobra.Command{
	Use:   "bigip_exporter",
	Short: "Prometheus exporter for F5 BIG-IP statistics",
	Long: `BIG-IP Exporter - Prometheus exporter for F5 BIG-IP systems

Exports metrics from BIG-IP using the iControl REST API.
Supports Virtual Servers, Pools, Nodes, and iRules.

Use "bigip_exporter run" to start the exporter.`,
}

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Start the BIG-IP exporter server",
	Long: `Start the BIG-IP exporter HTTP server and begin collecting metrics.

The exporter will connect to the specified BIG-IP system and expose
Prometheus metrics on the configured port (default: 9142).`,
	Run: func(cmd *cobra.Command, args []string) {
		runExporter()
	},
}

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("bigip_exporter %s\n", version)
		fmt.Printf("  commit: %s\n", commit)
		fmt.Printf("  built:  %s\n", date)
	},
}

func init() {
	cobra.OnInitialize(initConfig)

	// Global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default searches for config.yaml in current directory)")

	// BIG-IP flags (for run command)
	runCmd.Flags().String("bigip.host", "localhost", "BIG-IP hostname or IP address")
	runCmd.Flags().Int("bigip.port", 443, "BIG-IP management port")
	runCmd.Flags().String("bigip.username", "user", "BIG-IP username")
	runCmd.Flags().String("bigip.password", "pass", "BIG-IP password")
	runCmd.Flags().Bool("bigip.basic_auth", false, "Use HTTP Basic authentication instead of token-based")

	// Exporter flags (for run command)
	runCmd.Flags().String("exporter.bind_address", "localhost", "Address to bind the exporter")
	runCmd.Flags().Int("exporter.bind_port", 9142, "Port to expose metrics")
	runCmd.Flags().String("exporter.partitions", "", "Comma-separated list of partitions to monitor (empty = all)")
	runCmd.Flags().String("exporter.namespace", "bigip", "Prometheus metrics namespace")
	runCmd.Flags().String("exporter.log_level", "info", "Log level (debug, info, warn, error)")

	// Bind flags to viper
	for _, pair := range []struct {
		key  string
		flag string
	}{
		{"bigip.host", "bigip.host"},
		{"bigip.port", "bigip.port"},
		{"bigip.username", "bigip.username"},
		{"bigip.password", "bigip.password"},
		{"bigip.basic_auth", "bigip.basic_auth"},
		{"exporter.bind_address", "exporter.bind_address"},
		{"exporter.bind_port", "exporter.bind_port"},
		{"exporter.partitions", "exporter.partitions"},
		{"exporter.namespace", "exporter.namespace"},
		{"exporter.log_level", "exporter.log_level"},
	} {
		if err := viper.BindPFlag(pair.key, runCmd.Flags().Lookup(pair.flag)); err != nil {
			panic(fmt.Sprintf("failed to bind flag %s: %v", pair.flag, err))
		}
	}

	// Environment variables
	viper.SetEnvPrefix("BE")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// Add commands
	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(versionCmd)
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag
		viper.SetConfigFile(cfgFile)
	} else {
		// Search for config in current directory
		viper.AddConfigPath(".")
		viper.AddConfigPath("/etc/bigip_exporter/")
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
	}

	// Read config file if it exists
	if err := viper.ReadInConfig(); err == nil {
		logger.Info("Using config file", "file", viper.ConfigFileUsed())
	}

	// Configure logger level
	configureLogger()
}

func configureLogger() {
	logLevel := strings.ToLower(viper.GetString("exporter.log_level"))
	var level slog.Level

	switch logLevel {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn", "warning":
		level = slog.LevelWarn
	case "error", "critical":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
		logger.Warn("Invalid log level, using info", "provided", logLevel)
	}

	logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: level}))
}

func runExporter() {
	cfg := config.GetConfig()
	logger.Debug("Config loaded", "config", cfg)

	// Configure BIG-IP connection
	bigipEndpoint := cfg.Bigip.Host + ":" + strconv.Itoa(cfg.Bigip.Port)

	var exporterPartitionsList []string
	if cfg.Exporter.Partitions != "" {
		exporterPartitionsList = strings.Split(cfg.Exporter.Partitions, ",")
	}

	authMethod := f5.TOKEN
	if cfg.Bigip.BasicAuth {
		authMethod = f5.BASIC_AUTH
	}

	logger.Info("Connecting to BIG-IP",
		"host", cfg.Bigip.Host,
		"port", cfg.Bigip.Port,
		"username", cfg.Bigip.Username,
		"auth_method", authMethod,
	)

	bigip := f5.New(bigipEndpoint, cfg.Bigip.Username, cfg.Bigip.Password, authMethod)

	// Register Prometheus collector
	bigipCollector, err := collector.NewBigipCollector(bigip, cfg.Exporter.Namespace, exporterPartitionsList)
	if err != nil {
		logger.Error("Failed to create collector", "error", err)
		os.Exit(1)
	}

	prometheus.MustRegister(bigipCollector)

	// Start HTTP server
	listen(cfg.Exporter.BindAddress, cfg.Exporter.BindPort)
}

func listen(exporterBindAddress string, exporterBindPort int) {
	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`<html>
			<head><title>BIG-IP Exporter</title></head>
			<body>
			<h1>BIG-IP Exporter</h1>
			<p><a href="/metrics">Metrics</a></p>
			</body>
			</html>`))
	})

	exporterBind := exporterBindAddress + ":" + strconv.Itoa(exporterBindPort)
	logger.Info("Starting HTTP server", "address", exporterBind)

	if err := http.ListenAndServe(exporterBind, nil); err != nil {
		logger.Error("HTTP server failed", "error", err)
		os.Exit(1)
	}
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
