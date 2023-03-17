package senser

import (
	"app/webserver"
	"encoding/json"
	"fmt"
	"net/http"
)

type JsonSenser struct {
	Senser string `json:senser`
	Type   string `json:type`
	Data   string `json:data`
}

type HealthData struct {
	Sennserdata string `json:sennserdata`
	Message     string `json:message`
}

var Route []webserver.WebConfig = []webserver.WebConfig{
	{Pass: "/metrics", Handler: metrics},
	{Pass: "/health", Handler: health},
	{Pass: "/json", Handler: jsonData},
	{Pass: "/", Handler: rootdate},
}

func createAddJsonData(name, typev, value string) JsonSenser {
	var output JsonSenser
	output.Senser = name
	output.Type = typev
	output.Data = value
	return output
}

// JSONDataのHTTP出力
func jsonData(w http.ResponseWriter, r *http.Request) {
	var outdata []JsonSenser
	if SennserData.Bme280_data.Flag {
		outdata = append(outdata, createAddJsonData(SennserData.Bme280_data.Name, "hum", SennserDataValue.Bme280.Hum))
		outdata = append(outdata, createAddJsonData(SennserData.Bme280_data.Name, "press", SennserDataValue.Bme280.Press))
		outdata = append(outdata, createAddJsonData(SennserData.Bme280_data.Name, "tmp", SennserDataValue.Bme280.Temp))
	}
	tmpdata, _ := json.Marshal(outdata)
	fmt.Fprintf(w, "%s", string(tmpdata))
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

func createLineMetrics(name, types, value string) string {
	return "senser_data{type=\"" + types + "\",sennser=\"" + name + "\"} " + value
}

// metricsの結果確認
func metrics(w http.ResponseWriter, r *http.Request) {
	var output []string
	if SennserData.Bme280_data.Flag {
		output = append(output, createLineMetrics(SennserData.Bme280_data.Name, "hum", SennserDataValue.Bme280.Hum))
		output = append(output, createLineMetrics(SennserData.Bme280_data.Name, "press", SennserDataValue.Bme280.Press))
		output = append(output, createLineMetrics(SennserData.Bme280_data.Name, "tmp", SennserDataValue.Bme280.Temp))
	}
	for _, line := range output {
		fmt.Fprintf(w, "%s\n", line)
	}
}

// ルートフォルダの結果
func rootdate(w http.ResponseWriter, r *http.Request) {
	output := "<html><body><a href=\"/metrics\">metrics</a></body></html>"
	fmt.Fprintf(w, "%s", output)
}
