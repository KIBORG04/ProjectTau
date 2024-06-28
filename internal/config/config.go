package config

import (
	"os"
	"syscall"

	"gopkg.in/yaml.v3"
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
	AdminConfig struct {
		Login    string `yaml:"login"`
		Password string `yaml:"password"`
	} `yaml:"admin"`
	TelegramBot struct {
		Token          string `yaml:"token"`
		TrustedChatIDs []int  `yaml:"trustedChatIDs"`
	} `yaml:"telegram"`
	Proxy   string `yaml:"proxy"`
	BaseUrl string
}

func LoadConfigurations() {
	file, err := os.ReadFile("config/config.yaml")
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(file, &Config)
	if err != nil {
		panic(err)
	}
	baseUrl, found := syscall.Getenv("BASE_URL")
	if !found {
		baseUrl = ""
	}
	Config.BaseUrl = baseUrl
}
