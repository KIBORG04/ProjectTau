package config

import (
	"os"

	"github.com/go-yaml/yaml"
)

var Config GlobalConfig

type GlobalConfig struct {
	DatabaseConfig struct {
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		Dbname   string `yaml:"dbname"`
	} `yaml:"db"`
}

func LoadConfigurations() {
	file, err := os.ReadFile("../../config/config.yaml")
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(file, &Config)
	if err != nil {
		panic(err)
	}
}
