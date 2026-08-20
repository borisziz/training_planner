package config

import (
	"os"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Ports struct {
		Http string `yaml:"http"`
		Grpc string `yaml:"grpc"`
	} `yaml:"ports"`
	DBConnectURL string `yaml:"db_connect_url"`
}

func NewConfig() (*Config, error) {
	rawYAML, err := os.ReadFile("config.yml")
	if err != nil {
		return nil, errors.WithMessage(err, "reading config file")
	}
	var config Config

	err = yaml.Unmarshal(rawYAML, &config)
	if err != nil {
		return nil, errors.WithMessage(err, "parsing yaml")
	}

	return &config, nil
}
