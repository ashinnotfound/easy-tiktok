package config

import (
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
)

var C config

type config struct {
	NetworkType        string `yaml:"NetworkType"`
	Host               string `yaml:"Host"`
	Etcd               string `yaml:"Etcd"`
	UserHost           string `yaml:"UserHost"`
	UserService        string `yaml:"UserService"`
	VideoHost          string `yaml:"VideoHost"`
	VideoService       string `yaml:"VideoService"`
	InteractionHost    string `yaml:"InteractionHost"`
	InteractionService string `yaml:"InteractionService"`
	SocialHost         string `yaml:"SocialHost"`
	SocialService      string `yaml:"SocialService"`
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
