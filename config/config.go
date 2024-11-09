package config

import (
	_ "github.com/jpfuentes2/go-env/autoload"
	"github.com/kelseyhightower/envconfig"
)

type AuthConfig struct {
	SecretKey string `envconfig:"secret_key"`
}

type HttpConfig struct {
	Host           string `envconfig:"host"`
	Port           string `envconfig:"port"`
	ContextTimeout int    `envconfig:"context_timeout"`
}

type DbConfig struct {
	Host     string `envconfig:"host"`
	Port     string `envconfig:"port"`
	Username string `envconfig:"username"`
	Password string `envconfig:"password"`
	Name     string `envconfig:"name"`
}

func (ac *AuthConfig) MustConfig() error {
	return envconfig.Process("", ac)
}

func (hc *HttpConfig) MustConfig() error {
	return envconfig.Process("", hc)
}

func (dc *DbConfig) MustConfig() error {
	return envconfig.Process("db", dc)
}
