package text

import (
	"gocsvserver/config"
	"net/http"
	"os"
)

func Init(url string, mux *http.ServeMux) error {
	// config.Read.FilePassによるフォルダ存在確認
	if f, err := os.Stat(config.Read.FilePass); os.IsNotExist(err) || !f.IsDir() {
		return err
	}
	mux.HandleFunc("GET "+url, webTextRead)
	return nil
}
