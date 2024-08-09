package config

import "github.com/caarlos0/env/v6"

type SetupServer struct {
	Protocol   string `env:"PROTOCOL" envDefault:"tcp"`
	Hostname   string `env:"WEB_HOST" envDefault:""`
	Port       string `env:"WEB_PORT" envDefault:"8080"`
	StaticPage string `env:"WEB_FOLDER" envDefault:"./html"`
}

type SetupKube struct {
	KubePath     string `env:"KUBE_PATH" envDefault:"kubeconfig"`
	KubePathFlag bool   `env:"KUBE_PATH_FLAG" envDefault:"true"`
}

type TracerData struct {
	GrpcURL     string `env:"TRACER_GRPC_URL" envDefault:"otel-grpc.bookserver.home:4317"`
	ServiceName string `env:"TRACER_SERVICE_URL" envDefault:"client-go-test"`
	TracerUse   bool   `env:"TRACER_ON" envDefault:"false"`
}

var Web SetupServer
var Kube SetupKube
var TraData TracerData

// 環境設定
func Init() error {
	Web = SetupServer{}
	if err := env.Parse(&Web); err != nil {
		return err
	}
	Kube = SetupKube{}
	if err := env.Parse(&Kube); err != nil {
		return err
	}
	TraData = TracerData{}
	if err := env.Parse(&TraData); err != nil {
		return err
	}
	return nil

}
