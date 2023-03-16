package senser

import (
	"app/webserver"
	"encoding/json"
	"fmt"
	"net/http"
)

type HealthData struct {
	Sennserdata string `json:sennserdata`
	Message     string `json:message`
}

var Route []webserver.WebConfig = []webserver.WebConfig{
	{Pass: "/metrics", Handler: metrics},
	{Pass: "/health", Handler: health},
	// {Pass: "/json", Handler: json},
	{Pass: "/", Handler: rootdate},
}

// JSONDataのHTTP出力
func jsonData(w http.ResponseWriter, r *http.Request) {

}

//Healthチェックのデータ追加のフラグ
func healthadd(flag bool, name string, message string) (HealthData, int) {
	var tmp HealthData
	code := 200
	if flag {
		tmp.Sennserdata = name
		tmp.Message = message
		if tmp.Message != "OK" {
			code = 405
		}
		return tmp, code
	}
	return tmp, 404
}

// Health Checkの結果確認
func health(w http.ResponseWriter, r *http.Request) {
	code := 200
	var outdata []HealthData
	//BME280のチェック
	if tmp, tcode := healthadd(SennserData.Bme280_data.Flag, SennserData.Bme280_data.Name, SennserData.Bme280_data.Message); tcode != 404 {
		outdata = append(outdata, tmp)
	}
	tmpdata, _ := json.Marshal(outdata)
	w.WriteHeader(code)
	fmt.Fprintf(w, "%s", string(tmpdata))
}

// metricsの結果確認
func metrics(w http.ResponseWriter, r *http.Request) {

}

// ルートフォルダの結果
func rootdate(w http.ResponseWriter, r *http.Request) {
	output := "<html><body><a href=\"/metrics\">metrics</a></body></html>"
	fmt.Fprintf(w, "%s", output)
}
