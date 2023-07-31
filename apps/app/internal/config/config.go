package config

import (
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
)

var C config

type config struct {
	Host      string `yaml:"Host"`
	UserHost  string `yaml:"UserHost"`
	VideoHost string `yaml:"VideoHost"`
}

func Initial() {
	dir, _ := os.Getwd()
	yamlPath := filepath.Join(dir, "/etc/cfg.yaml")
	println(yamlPath)
	file, err := os.ReadFile(yamlPath)
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(file, &C)
	if err != nil {
		return
	}
}
