package config

import (
	"github.com/spf13/viper"
)

type bigipConfig struct {
	Username  string `yaml:"username"`
	Password  string `yaml:"password"`
	BasicAuth bool   `yaml:"basic_auth"`
	Host      string `yaml:"host"`
	Port      int    `yaml:"port"`
}

type exporterConfig struct {
	BindAddress string `yaml:"bind_address"`
	BindPort    int    `yaml:"bind_port"`
	Partitions  string `yaml:"partitions"`
	Config      string `yaml:"config"`
	Namespace   string `yaml:"namespace"`
	LogLevel    string `yaml:"log_level"`
}

// Config is a container for settings modifiable by the user
type Config struct {
	Bigip    bigipConfig    `yaml:"bigip"`
	Exporter exporterConfig `yaml:"exporter"`
}

// GetConfig returns an instance of Config containing the resulting parameters
// to the program. Configuration is managed by Viper and can come from:
// 1. Command-line flags (highest priority)
// 2. Environment variables (prefix: BE_)
// 3. Configuration file
func GetConfig() *Config {
	return &Config{
		Bigip: bigipConfig{
			Username:  viper.GetString("bigip.username"),
			Password:  viper.GetString("bigip.password"),
			BasicAuth: viper.GetBool("bigip.basic_auth"),
			Host:      viper.GetString("bigip.host"),
			Port:      viper.GetInt("bigip.port"),
		},
		Exporter: exporterConfig{
			BindAddress: viper.GetString("exporter.bind_address"),
			BindPort:    viper.GetInt("exporter.bind_port"),
			Partitions:  viper.GetString("exporter.partitions"),
			Config:      viper.GetString("exporter.config"),
			Namespace:   viper.GetString("exporter.namespace"),
			LogLevel:    viper.GetString("exporter.log_level"),
		},
	}
}
