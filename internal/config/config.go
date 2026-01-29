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
	Type     string `yaml:"type"`
}

type FTPConfig struct {
	Addr           string      `yaml:"addr"`
	User           string      `yaml:"user"`
	Password       string      `yaml:"password"`
	TimeoutSeconds int         `yaml:"timeout_seconds"`
	Sources        []FTPSource `yaml:"sources"`
}

type FTPSource struct {
	Remote   string `yaml:"remote"`
	Filename string `yaml:"filename"`
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
