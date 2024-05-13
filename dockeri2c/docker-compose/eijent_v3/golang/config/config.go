package config

import (
	"github.com/caarlos0/env"
	"github.com/pkg/errors"
)

// Webサーバの設定
type WebConfig struct {
	Protocol   string `env:"WEB_PROTOCOL" envDefault:"tcp"`  //接続プロトコル
	Hostname   string `env:"WEB_HOST" envDefault:""`         //接続DNS名
	Port       string `env:"WEB_PORT" envDefault:"9140"`     //接続ポート
	StaticPage string `env:"WEB_FOLDER" envDefault:"./html"` //静的ページの参照先
}

type SenserConfig struct {
	GPIO_ON          bool   `env:"GPIO_ON" envDefault:"true"`            //GPIO読み取りの有効
	I2C_ON           bool   `env:"I2C_ON" envDefault:"true"`             //I2C読み取りの有効
	UART_ON          bool   `env:"UART_ON" envDefault:"true"`            //UART読み取りの有効
	HorldTime        int    `env:"HOLDTIME" envDefault:"30"`             //センサー読み取りの有効読み込み保持時間 (分)
	Am2320_Count     int    `env:"AM2320_INTERVAL" envDefault:"100"`     //Am2320のセンサーのインターバル ms
	Tsl2561_Count    int    `env:"TSL2561_INTERVAL" envDefault:"100"`    //Tsl2561のセンサーのインターバル ms
	BME280_Count     int    `env:"BME280_INTERVAL" envDefault:"100"`     //BME280のセンサーのインターバル ms
	MMA8452Q_Count   int    `env:"MMA8452Q_INTERVAL" envDefault:"500"`   //MMA8452Qのセンサーのインターバル µs
	DHT_senser_type  string `env:"DHT_SENSER_TYPE" envDefault:"DHT11"`   //DHT系のセンサーの種類
	DHT_senser_pin   int    `env:"DHT_SENSER_PIN" envDefault:"12"`       //DHT系のセンサーのピン番号
	DHT_senser_Count int    `env:"DHT_SENSER_INTERVAL" envDefault:"200"` //DHT系のセンサーのインターバル ms
	CO2_SENSER_Count int    `env:"CO2_SENSER_INTERVAL" envDefault:"200"` //CO2のセンサーのインターバル ms
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
