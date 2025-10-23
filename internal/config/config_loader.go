package config

import (
	"io/ioutil"

	"github.com/iabdulzahid/product-update-service/internal/domain"
	"gopkg.in/yaml.v3"
)

func LoadConfig(path string) (*domain.Config, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg domain.Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
