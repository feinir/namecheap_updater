package common

import (
	"gopkg.in/yaml.v3"
	"os"
	"time"
)

type Config struct {
	DomainHosts    []DomainHost `yaml:"domain_host"` //频率，秒
	ConfigFile     string
	CommonGetIpUrl string `yaml:"common_get_ip_url"`
}

type DomainHost struct {
	UpdateFrequency float64 `yaml:"frequency"` //频率，秒
	Host            string  `yaml:"host"`
	DomainName      string  `yaml:"domain"`
	Password        string  `yaml:"password"`
	IP              string  `yaml:"ip"`
	GetIpUrl        string  `yaml:"get_ip_url"`
	LastUpdateTime  time.Time
	LastStatus      bool
	LastUpIP        string
}

func ReadConfig(path string) (*Config, error) {
	configBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	result := &Config{}

	if err = yaml.Unmarshal(configBytes, result); err != nil {
		return nil, err
	}

	return result, nil
}
