package config

import (
	"gopkg.in/yaml.v2"
	"os"
)

type Config struct {
	Database databaseConfig `yaml:"database"`
}

type databaseConfig struct {
	Postgres postgresConfig `yaml:"postgres"`
}

type postgresConfig struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	DBName   string `yaml:"dbname"`
}

func Load() (*Config, error) {
	data, err := os.ReadFile("config/app.yaml")
	if err != nil {
		return nil, err
	}
	config := &Config{}
	err = yaml.Unmarshal(data, config)
	if err != nil {
		return nil, err
	}
	return config, nil
}
