package config

import "github.com/caarlos0/env"

type WebConfig struct {
	Protocol   string `env:"WEB_PROTOCOL" envDefault:"tcp"`
	Hostname   string `env:"WEB_HOSTNAME" envDefault:""`
	Port       string `env:"WEB_PORT" envDefault:"8080"`
	StaticPage string `env:"WEB_FOLDER" envDefault:"./html"`
}

type ReadConfig struct {
	FilePass string `env:"READ_FILEPASS" envDefault:"./txt-tmp"`
}

var Web WebConfig
var Read ReadConfig

func Init() error {
	if err := env.Parse(&Web); err != nil {
		return err
	}
	if err := env.Parse(&Read); err != nil {
		return err
	}
	return nil
}
