package main

import (
	"io/ioutil"

	yaml "gopkg.in/yaml.v3"
)

type Config struct {
	Features map[string]interface{} `yaml:"features"`
	Flows    map[string]Flow        `yaml:"flows"`
}

type Flow struct {
	Steps []Step `yaml:"steps"`
}

type Step struct {
	Name      string            `yaml:"name"`
	View      string            `yaml:"view"`
	Link      string            `yaml:"link"`
	EnabledIf string            `yaml:"enabled_if"`
	Next      map[string]string `yaml:"next"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	cfg := &Config{}
	err = yaml.Unmarshal(data, cfg)
	return cfg, err
}
