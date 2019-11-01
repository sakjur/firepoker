package config

import (
	"io"

	"github.com/BurntSushi/toml"
	"github.com/sakjur/firepoker/internal/providers"
)

type Config struct {
	Providers Providers `toml:"providers"`
}

type Providers struct {
	Elks providers.Elks `toml:"elks"`
}

func Read(reader io.Reader) (Config, error) {
	c := &Config{}

	_, err := toml.DecodeReader(reader, c)
	if err != nil {
		return Config{}, err
	}

	return *c, nil
}
