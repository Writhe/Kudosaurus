package config

import (
	"fmt"
	"os"

	"github.com/go-playground/validator"
	"gopkg.in/yaml.v2"
)

var validate *validator.Validate

// Config - app configuration
type Config struct {
	Slack struct {
		Token             string `yaml:"token" validate:"required"`
		SigningSecret     string `yaml:"signingSecret" validate:"required"`
		VerificationToken string `yaml:"verificationToken" validate:"required"`
	} `yaml:"slack" validate:"required"`
	Database struct {
		Path string `yaml:"path" validate:"required"`
	} `yaml:"db" validate:"required"`
	Server struct {
		Port int `yaml:"port" validate:"required"`
	} `yaml:"server" validate:"required"`
	Settings struct {
		MaxKudos int `yaml:"maxKudos" validate:"required"`
	} `yaml:"settings" validate:"required"`
}

// GetConfig - gets config from file
func GetConfig(path string) (Config, error) {
	validate = validator.New()
	f, err := os.Open(path)
	if err != nil {
		return Config{}, fmt.Errorf("could not open config file - '%s'", path)
	}
	defer f.Close()

	var cfg Config
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		return Config{}, fmt.Errorf("could not decode the config: %s", err.Error())
	}
	err = validate.Struct(cfg)
	if err != nil {
		return Config{}, fmt.Errorf("could not validate the config: %s", err.Error())
	}

	return cfg, nil
}
