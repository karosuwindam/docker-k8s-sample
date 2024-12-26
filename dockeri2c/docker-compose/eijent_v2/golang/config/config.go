package config

import "github.com/caarlos0/env/v6"

type SetupServer struct {
	Protocol   string `env:"PROTOCOL" envDefault:"tcp"`
	Hostname   string `env:"WEB_HOST" envDefault:""`
	Port       string `env:"WEB_PORT" envDefault:"8080"`
	StaticPage string `env:"WEB_FOLDER" envDefault:"./html"`
}

type TracerData struct {
	// GrpcURL        string `env:"TRACER_GRPC_URL" envDefault:"otel-grpc.bookserver.home:4317"`
	GrpcURL        string `env:"TRACER_GRPC_URL" envDefault:"localhost:4317"`
	ServiceName    string `env:"TRACER_SERVICE_URL" envDefault:"booknewRead-test"`
	TracerUse      bool   `env:"TRACER_ON" envDefault:"true"`
	ServiceVersion string `env:"TRACER_SERVICE_VERSION" envDefault:"0.26.2"`
}

var Web SetupServer
var TraData TracerData

func Init() error {
	Web = SetupServer{}
	if err := env.Parse(&Web); err != nil {
		return err
	}

	TraData = TracerData{}
	if err := env.Parse(&TraData); err != nil {
		return err
	}
	return nil
}
