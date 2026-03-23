package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type DBConfig struct {
	Host         string `yaml:"host"`
	Port         int    `yaml:"port"`
	User         string `yaml:"user"`
	Password     string `yaml:"password"`
	DBName       string `yaml:"dbname"`
	SSHHost      string `yaml:"ssh_host"`
	SSHPort      int    `yaml:"ssh_port"`
	SSHUser      string `yaml:"ssh_user"`
	SSHPublicKey string `yaml:"ssh_public_key"`
}

type OnlineStatsConfig struct {
	RemoteDB   DBConfig `yaml:"remote_db"`
	OutputPath string   `yaml:"output_path"`
	StatePath  string   `yaml:"state_path"`
}

func LoadOnlineStatsConfig(path string) (*OnlineStatsConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg OnlineStatsConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return &cfg, nil
}
