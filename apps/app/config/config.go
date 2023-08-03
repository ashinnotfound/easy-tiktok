package config

import (
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
)

var C config

type config struct {
	NetworkType     string `yaml:"NetworkType"`
	Host            string `yaml:"Host"`
	UserHost        string `yaml:"UserHost"`
	VideoHost       string `yaml:"VideoHost"`
	InteractionHost string `yaml:"InteractionHost"`
	SocialHost      string `yaml:"SocialHost"`
}

func Initial() {
	dir, _ := os.Getwd()
	yamlPath := filepath.Join(dir, "apps/app/etc/cfg.yaml")
	file, err := os.ReadFile(yamlPath)
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(file, &C)
	if err != nil {
		return
	}
}
