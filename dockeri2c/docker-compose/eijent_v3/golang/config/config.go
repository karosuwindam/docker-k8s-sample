package config

import "github.com/caarlos0/env"

// Webサーバの設定
type WebConfig struct {
	Protocol   string `env:"WEB_PROTOCOL" envDefault:"tcp"`  //接続プロトコル
	Hostname   string `env:"WEB_HOST" envDefault:""`         //接続DNS名
	Port       string `env:"WEB_PORT" envDefault:"8080"`     //接続ポート
	StaticPage string `env:"WEB_FOLDER" envDefault:"./html"` //静的ページの参照先
}

var Web WebConfig

// 環境設定
func Init() error {
	if err := env.Parse(&Web); err != nil {
		return err
	}
	return nil
}
