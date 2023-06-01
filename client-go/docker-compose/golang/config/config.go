package config

import "github.com/caarlos0/env/v6"

type SetupServer struct {
	Protocol string `env:"PROTOCOL" envDefault:"tcp"`
	Hostname string `env:"WEB_HOST" envDefault:""`
	Port     string `env:"WEB_PORT" envDefault:"8080"`
}

type SetupKube struct {
	KubePath     string `env:"KUBE_PATH" envDefault:"kubeconfig"`
	KubePathFlag bool   `env:"KUBE_PATH_FLAG" envDefault:"false"`
}

type Config struct {
	Server *SetupServer
	Kube   *SetupKube
}

// 環境設定
func Setup() (*Config, error) {
	serverCfg := &SetupServer{}
	if err := env.Parse(serverCfg); err != nil {
		return nil, err
	}
	kubeCfg := &SetupKube{}
	if err := env.Parse(kubeCfg); err != nil {
		return nil, err
	}
	return &Config{
		Server: serverCfg,
		Kube:   kubeCfg,
	}, nil

}
