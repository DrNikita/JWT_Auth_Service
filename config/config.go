package config

import (
	_ "github.com/jpfuentes2/go-env/autoload"
	"github.com/kelseyhightower/envconfig"
)

type HttpConfig struct {
	Host string `envconfig:"host"`
	Port string `envconfig:"port"`
}

type DbConfig struct {
	Host     string `envconfig:"host"`
	Port     string `envconfig:"port"`
	Username string `envconfig:"username"`
	Password string `envconfig:"password"`
	Name     string `envconfig:"name"`
}

func (hc *HttpConfig) MustConfig() error {
	return envconfig.Process("", hc)
}

func (dc *DbConfig) MustConfig() error {
	return envconfig.Process("db", dc)
}
