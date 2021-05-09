package config

import (
	"io/ioutil"

	"github.com/BurntSushi/toml"
)

type Server struct {
	Host    string
	Port    string
	Timeout string
}

type Storage struct {
	DBConnectionString string `toml:"db_connection_string"`
}

type Logger struct {
	File  string
	Level string
}

type Config struct {
	Server  Server
	Storage Storage
	Log     Logger
}

func ReadConfig(configPath string) (Config, error) {
	var cfg Config

	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return cfg, err
	}

	_, err = toml.Decode(string(data), &cfg)
	return cfg, err
}
