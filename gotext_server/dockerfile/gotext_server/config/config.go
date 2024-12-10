package config

import "github.com/caarlos0/env"

type WebConfig struct {
	Protocol   string `env:"WEB_PROTOCOL" envDefault:"tcp"`
	Hostname   string `env:"WEB_HOSTNAME" envDefault:""`
	Port       string `env:"WEB_PORT" envDefault:"8080"`
	StaticPage string `env:"WEB_FOLDER" envDefault:"./html"`
}

type ReadConfig struct {
	FilePass string `env:"READ_FILEPASS" envDefault:"./txt"`
	// FilePass string `env:"READ_FILEPASS" envDefault:"./txt-tmp"`
}

type TracerData struct {
	GrpcURL        string `env:"TRACER_GRPC_URL" envDefault:"otel-grpc.bookserver.home:4317"`
	ServiceName    string `env:"TRACER_SERVICE_URL" envDefault:"gotext-server-test"`
	TracerUse      bool   `env:"TRACER_ON" envDefault:"true"`
	ServiceVersion string `env:"TRACER_SERVICE_VERSION" envDefault:"0.1.0"`
}

var Web WebConfig
var Read ReadConfig
var TraData TracerData

func Init() error {
	if err := env.Parse(&Web); err != nil {
		return err
	}
	if err := env.Parse(&Read); err != nil {
		return err
	}
	if err := env.Parse(&TraData); err != nil {
		return err
	}
	return nil
}
