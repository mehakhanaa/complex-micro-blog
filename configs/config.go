package configs

import (
	"os"

	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/pelletier/go-toml/v2"
)

type Config struct {
	Server struct {
		Host string `toml:"host"`

		Port int `toml:"port"`
	} `toml:"server"`

	Database struct {
		Host string `toml:"host"`

		Port int `toml:"port"`

		User string `toml:"user"`

		Password string `toml:"password"`

		DBName string `toml:"db_name"`
	} `toml:"database"`

	Redis struct {
		Host string `toml:"host"`

		Port int `toml:"port"`

		Username string `toml:"username"`

		Password string `toml:"password"`

		DB int `toml:"db"`
	} `toml:"redis"`

	SearchService struct {
		Host string `toml:"host"`
		Port int    `toml:"port"`
	} `toml:"search_service"`

	Compress struct {
		Level compress.Level `toml:"level"`
	} `toml:"compress"`

	Env struct {
		Type string `toml:"type"`
	} `toml:"env"`
}

func NewConfig() (*Config, error) {

	file, err := os.ReadFile("./configuration.toml")
	if err != nil {
		return nil, err
	}

	config := new(Config)
	err = toml.Unmarshal(file, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
