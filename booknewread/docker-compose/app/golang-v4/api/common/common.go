package common

import (
	"book-newread/config"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

// URLの解析
func UrlAnalysis(url string) []string {
	tmp := []string{}
	for _, str := range strings.Split(url[1:], "/") {
		tmp = append(tmp, str)
	}
	return tmp
}

type Result struct {
	Name   string      `json:"Name"`
	Url    string      `json:"Url"`
	Code   int         `json:"Code"`
	Option string      `json:"Option"`
	Date   time.Time   `json:"Date"`
	Result interface{} `json:"Result"`
}

// 共通の出力設定
func CommonBack(msg Result, w http.ResponseWriter) {
	jsondata, _ := json.Marshal(msg)
	w.WriteHeader(msg.Code)
	fmt.Fprintf(w, "%v\n", string(jsondata))
}

func Setup(cfg *config.Config) error {
	return nil
}

// ファイルの存在確認
func Exists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}
