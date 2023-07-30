package config

import (
	"gopkg.in/yaml.v3"
	"os"
)

var C config

type config struct {
	Host string `yaml:"Host"`
	RPC  struct {
		Host string `yaml:"Host"`
	}
}

func init() {
	file, err := os.ReadFile("..../etc/config.yaml")
	if err != nil {
		return
	}
	err = yaml.Unmarshal(file, &C)
	if err != nil {
		return
	}
}
