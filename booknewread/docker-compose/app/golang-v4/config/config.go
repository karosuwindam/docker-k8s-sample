package config

import "github.com/caarlos0/env/v6"

type SetupServer struct {
	Protocol string `env:"PROTOCOL" envDefault:"tcp"`
	Hostname string `env:"WEB_HOST" envDefault:""`
	Port     string `env:"WEB_PORT" envDefault:"8080"`
}

type SetupLoop struct {
	MaxAccess  int    `env:"DEF_ACCESS" envDefault:"2"`
	MaxBackDay int    `env:"DEF_BACKDAY" envDefault:"5"`
	LoopTIme   int    `env:"DEF_LOOPTIME" envDefault:"3600"` //待ち時間 デフォルト3600s = 1h
	BookmarkF  string `env:"DEF_BOOKMARK" envDefault:"./bookmark"`
}

type Config struct {
	Server *SetupServer
	Loop   *SetupLoop
}

func EnvRead() (*Config, error) {
	serverCfg := &SetupServer{}
	if err := env.Parse(serverCfg); err != nil {
		return nil, err
	}
	loopCfg := &SetupLoop{}
	if err := env.Parse(loopCfg); err != nil {
		return nil, err
	}
	return &Config{
		Server: serverCfg,
		Loop:   loopCfg,
	}, nil

}
