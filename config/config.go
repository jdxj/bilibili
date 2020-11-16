package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

var (
	cfg *Config
)

func init() {
	err := readConfig("./config.yaml")
	if err != nil {
		panic(err)
	}
}

type Config struct {
	Email *Email  `yaml:"email"`
	Users []*User `yaml:"users"`
}

type Email struct {
	User  string `yaml:"user"`
	Token string `yaml:"token"`
}

type User struct {
	Email  string `yaml:"email"`
	Cookie string `yaml:"cookie"`
}

func readConfig(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	cfg = new(Config)
	decoder := yaml.NewDecoder(file)
	return decoder.Decode(cfg)
}

func GetEmail() *Email {
	return cfg.Email
}

func GetUsers() []*User {
	return cfg.Users
}
