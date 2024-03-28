package config

import (
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	Host       string `yaml:"host"`
	Port       string `yaml:"port"`
	DbUrl      string `yaml:"DB_URL"`
	DbUser     string `yaml:"DB_USER"`
	DbPassword string `yaml:"DB_PASSWORD"`
	DbDatabase string `yaml:"DB_DATABASE"`
}

func LoadConfig() *Config {
	var err error

	cfg := Config{
		Host:       "",
		Port:       "",
		DbUrl:      "",
		DbUser:     "",
		DbPassword: "",
		DbDatabase: "",
	}
	file, err := os.ReadFile("./application.yml")
	if err != nil {
		logrus.Fatalf("Error read config file: %s", err.Error())
	}

	if err = yaml.Unmarshal(file, &cfg); err != nil {
		logrus.Panicln("LoadConfig - %w", err)
	}

	url, exists := os.LookupEnv("DB_URL")
	if exists {
		cfg.DbUrl = url
	}

	user, exists := os.LookupEnv("DB_USER")
	if exists {
		cfg.DbUser = user
	}

	password, exists := os.LookupEnv("DB_PASSWORD")
	if exists {
		cfg.DbPassword = password
	}

	database, exists := os.LookupEnv("DB_DATABASE")
	if exists {
		cfg.DbDatabase = database
	}

	return &cfg
}
