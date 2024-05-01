package config

import (
	"github.com/caarlos0/env"
	"github.com/pkg/errors"
)

// Webサーバの設定
type WebConfig struct {
	Protocol   string `env:"WEB_PROTOCOL" envDefault:"tcp"`  //接続プロトコル
	Hostname   string `env:"WEB_HOST" envDefault:""`         //接続DNS名
	Port       string `env:"WEB_PORT" envDefault:"8080"`     //接続ポート
	StaticPage string `env:"WEB_FOLDER" envDefault:"./html"` //静的ページの参照先
}

type SenserConfig struct {
	Tsl2561_Count    int    `env:"TSL2561_INTERVAL" envDefault:"100"`  //Tsl2561のセンサーのインターバル ms
	DHT_senser_type  string `env:"DHT_SENSER_TYPE" envDefault:"DHT11"` //DHT系のセンサーの種類
	DHT_senser_pin   int    `env:"DHT_SENSER_PIN" envDefault:"583"`    //DHT系のセンサーのピン番号
	DHT_senser_Count int    `env:"DHT_SENSER_PIN" envDefault:"200"`    //DHT系のセンサーのインターバル ms
}

type LogConfig struct {
	Colors bool `env:"LOG_COLORS" envDefault:"true"`
}

var Web WebConfig
var Log LogConfig
var Senser SenserConfig

// 環境設定
func Init() error {
	if err := env.Parse(&Web); err != nil {
		return errors.Wrap(err, "env.Parse(Web)")
	}
	if err := env.Parse(&Log); err != nil {
		return errors.Wrap(err, "env.Parse(Log)")
	}
	if err := env.Parse(&Senser); err != nil {
		return errors.Wrap(err, "env.Parse(Senser)")
	}
	return nil
}
