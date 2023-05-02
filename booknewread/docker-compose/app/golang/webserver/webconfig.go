package webserver

import (
	"log"
	"net/http"
)

type WebConfig struct {
	Pass    string
	Handler func(http.ResponseWriter, *http.Request)
}

// Config(cfg *SetupServer, wconfs []WebConfig) = error
//
// 複数のhttp上に定義されたハンドラ関数を紐づける
//
// cfg(*SetupServer) : Webサーバの設定
// wconfs([]WebConfig) : 登録する定義
// root(string): ホームルートからのパス
func Config(cfg *SetupServer, wconfs []WebConfig, pass string) error {
	tmp := pass
	if tmp == "" {
	} else if tmp[len(tmp)-1:] == "/" {
		tmp = tmp[:len(tmp)-1]
	}
	for _, wconf := range wconfs {
		if err := cfg.Add(tmp+wconf.Pass, wconf.Handler); err != nil {
			log.Panicln(err)
		}
	}
	return nil
}
