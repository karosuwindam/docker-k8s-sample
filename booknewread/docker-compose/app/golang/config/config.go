package config

import "github.com/caarlos0/env/v6"

type SetupServer struct {
	Protocol   string `env:"PROTOCOL" envDefault:"tcp"`
	Hostname   string `env:"WEB_HOST" envDefault:""`
	Port       string `env:"WEB_PORT" envDefault:"8080"`
	StaticPage string `env:"WEB_FOLDER" envDefault:"./html"`
}

type SetupLoop struct {
	MaxAccess  int    `env:"DEF_ACCESS" envDefault:"2"`
	MaxBackDay int    `env:"DEF_BACKDAY" envDefault:"5"`
	LoopTIme   int    `env:"DEF_LOOPTIME" envDefault:"3600"` //待ち時間 デフォルト3600s = 1h
	BookmarkF  string `env:"DEF_BOOKMARK" envDefault:"./bookmark"`
}

type NobelChack struct {
	MaxNarouAPI    int `env:"NOBEL_MAX_NAROU_API" envDefault:"2"`
	MaxKakuyomuAPI int `env:"NOBEL_MAX_NAROU_API" envDefault:"2"`
	MaxNarou18API  int `env:"NOBEL_MAX_NAROU_18_API" envDefault:"2"`
	MaxAlphaAPI    int `env:"NOBEL_MAX_ALPHA_API" envDefault:"2"`
	Timeout        int `env:"NOBEL_TIMEOUT" envDefault:"2000"` //ms
}

type TracerData struct {
	GrpcURL        string `env:"TRACER_GRPC_URL" envDefault:"otel-grpc.bookserver.home:4317"`
	ServiceName    string `env:"TRACER_SERVICE_URL" envDefault:"booknewRead-test"`
	TracerUse      bool   `env:"TRACER_ON" envDefault:"true"`
	ServiceVersion string `env:"TRACER_SERVICE_VERSION" envDefault:"0.26.2"`
}

var Web SetupServer
var Loop SetupLoop
var Nobel NobelChack
var TraData TracerData

func Init() error {
	Web = SetupServer{}
	if err := env.Parse(&Web); err != nil {
		return err
	}
	Loop = SetupLoop{}
	if err := env.Parse(&Loop); err != nil {
		return err
	}
	Nobel = NobelChack{}
	if err := env.Parse(&Nobel); err != nil {
		return err
	}
	TraData = TracerData{}

	if err := env.Parse(&TraData); err != nil {
		return err
	}
	// if err := logConfig(); err != nil {
	// 	return err
	// }
	return nil
}
