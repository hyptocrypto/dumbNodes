package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// ServerConfig represents the server configuration.
type Config struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Protocol string `yaml:"protocol"`
}

func (c *Config) FullHost() string {
	return c.Host + c.Port
}

func NewConfig() *Config {
	file, err := os.Open("internal/config/config.yaml")
	if err != nil {
		fmt.Printf("Error opening config.yaml: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	var config Config
	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		fmt.Printf("Error decoding YAML: %v\n", err)
		os.Exit(1)
	}
	return &config
}
