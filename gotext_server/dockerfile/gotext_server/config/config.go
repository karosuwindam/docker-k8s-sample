package config

import "github.com/caarlos0/env/v6"

type SetupServer struct {
	Protocol string `env:"PROTOCOL" envDefault:"tcp"`
	Hostname string `env:"WEB_HOST" envDefault:""`
	Port     string `env:"WEB_PORT" envDefault:"8080"`
}

type TXTFolder struct {
	RootPath string `env:"TXT_ROOT_PATH" envDefault:"./txt/"`
}

type Config struct {
	Server *SetupServer
	TXT    *TXTFolder
}

// 環境設定
func Setup() (*Config, error) {
	serverCfg := &SetupServer{}
	if err := env.Parse(serverCfg); err != nil {
		return nil, err
	}
	txtCfg := &TXTFolder{}
	if err := env.Parse(txtCfg); err != nil {
		return nil, err
	}
	return &Config{
		Server: serverCfg,
		TXT:    txtCfg,
	}, nil

}

//フォルダ作成
