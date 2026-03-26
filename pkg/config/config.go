package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	SiteTitle    string `yaml:"siteTitle"`
	AuthorName   string `yaml:"authorName"`
	AuthorRole   string `yaml:"authorRole"`
	AuthorBio    string `yaml:"authorBio"`
	DefaultTheme string `yaml:"defaultTheme"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
