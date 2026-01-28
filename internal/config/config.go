package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	InputDir string `yaml:"input_dir"`

	HTTP HTTPConfig `yaml:"http"`
	FTP  FTPConfig  `yaml:"ftp"`
}

type HTTPConfig struct {
	TimeoutSeconds int          `yaml:"timeout_seconds"`
	Retries        int          `yaml:"retries"`
	Sources        []HTTPSource `yaml:"sources"`
}

type HTTPSource struct {
	URL      string `yaml:"url"`
	Filename string `yaml:"filename"`
}

type FTPConfig struct {
	Enabled        bool   `yaml:"enabled"`
	Host           string `yaml:"host"`
	Port           int    `yaml:"port"`
	User           string `yaml:"user"`
	Password       string `yaml:"password"`
	RemoteDir      string `yaml:"remote_dir"`
	TimeoutSeconds int    `yaml:"timeout_seconds"`
	Retries        int    `yaml:"retries"`
}

func Load(path string) (*Config, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(b, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
