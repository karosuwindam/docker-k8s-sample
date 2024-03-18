package text

import "net/http"

type TXTData struct {
	Year  string   `json:"Year"`  // 年
	Quart string   `json:"Quart"` // 四半期
	Title []string `json:"Title"` // タイトル
}

func webTextRead(w http.ResponseWriter, r *http.Request) {

}
