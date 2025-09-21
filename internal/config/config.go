package config

import (
	"strings"

	"github.com/knadh/koanf/parsers/dotenv"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

type (
	App struct {
		Production bool `koanf:"production"`
	}

	Database struct {
		Name     string `koanf:"name"`
		Host     string `koanf:"host"`
		Port     string `koanf:"port"`
		User     string `koanf:"user"`
		Password string `koanf:"password"`
	}

	Server struct {
		Host string `koanf:"host"`
		Port string `koanf:"port"`
	}

	Config struct {
		App      App
		Database Database
		Server   Server
	}
)

var k = koanf.New(".")

func New() (*Config, error) {
	config := new(Config)

	if err := k.Load(file.Provider("configs/config.yaml"), yaml.Parser()); err != nil {
		return nil, err
	}

	if err := k.Load(env.Provider("", ".", func(s string) string {
		return strings.ToLower(strings.ReplaceAll(s, "_", "."))
	}), nil); err != nil {
		return nil, err
	}

	k.Load(file.Provider(".env"), dotenv.ParserEnv("", ".", func(s string) string {
		return strings.ToLower(strings.ReplaceAll(s, "_", "."))
	}))

	if err := k.Unmarshal("", config); err != nil {
		return nil, err
	}

	return config, nil
}
