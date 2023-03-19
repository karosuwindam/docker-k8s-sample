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
	{Pass: "/reset", Handler: reset},
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
	tmpdata := SennserData
	SennserDataValue.Mu.Lock()
	tmpvaule := SennserDataValue
	SennserDataValue.Mu.Unlock()
	if tmpdata.Bme280_data.Flag {
		outdata = append(outdata, createAddJsonData(tmpdata.Bme280_data.Name, "hum", tmpvaule.Bme280.Hum))
		outdata = append(outdata, createAddJsonData(tmpdata.Bme280_data.Name, "press", tmpvaule.Bme280.Press))
		outdata = append(outdata, createAddJsonData(tmpdata.Bme280_data.Name, "tmp", tmpvaule.Bme280.Temp))
	}
	if tmpdata.Am2320_data.Flag {
		outdata = append(outdata, createAddJsonData(tmpdata.Am2320_data.Name, "hum", tmpvaule.Am2320.Hum))
		outdata = append(outdata, createAddJsonData(tmpdata.Am2320_data.Name, "tmp", tmpvaule.Am2320.Temp))
	}
	if tmpdata.Tsl2561_data.Flag {
		outdata = append(outdata, createAddJsonData(tmpdata.Tsl2561_data.Name, "lux", tmpvaule.Tsl2561.Lux))
	}
	if tmpdata.CO2Sensor_data.Flag {
		outdata = append(outdata, createAddJsonData(tmpdata.CO2Sensor_data.Name, "co2", tmpvaule.CO2.Co2))
		outdata = append(outdata, createAddJsonData(tmpdata.CO2Sensor_data.Name, "tmp", tmpvaule.CO2.Temp))
	}
	if tmpdata.DhtSenser_data.Flag {
		outdata = append(outdata, createAddJsonData(tmpdata.DhtSenser_data.Name, "hum", tmpvaule.DhtSenser.Hum))
		outdata = append(outdata, createAddJsonData(tmpdata.DhtSenser_data.Name, "tmp", tmpvaule.DhtSenser.Temp))
	}
	if tmpdata.Mma8452q_data.Flag {
		outdata = append(outdata, createAddJsonData(tmpdata.Mma8452q_data.Name, "ax", tmpvaule.Mma8452q.X))
		outdata = append(outdata, createAddJsonData(tmpdata.Mma8452q_data.Name, "ay", tmpvaule.Mma8452q.Y))
		outdata = append(outdata, createAddJsonData(tmpdata.Mma8452q_data.Name, "az", tmpvaule.Mma8452q.Z))
		outdata = append(outdata, createAddJsonData(tmpdata.Mma8452q_data.Name, "zero_x", tmpvaule.Mma8452q.Zero_X))
		outdata = append(outdata, createAddJsonData(tmpdata.Mma8452q_data.Name, "zero_y", tmpvaule.Mma8452q.Zero_Y))
		outdata = append(outdata, createAddJsonData(tmpdata.Mma8452q_data.Name, "zero_z", tmpvaule.Mma8452q.Zero_Z))
	}
	outdata = append(outdata, createAddJsonData("localhost", "cpu_tmp", tmpvaule.CpuTmp))
	output, _ := json.Marshal(outdata)
	fmt.Fprintf(w, "%s", string(output))
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
	//AM2320のチェック
	if tmp, tcode := healthadd(SennserData.Am2320_data.Flag, SennserData.Am2320_data.Name, SennserData.Am2320_data.Message); tcode != 404 {
		outdata = append(outdata, tmp)
	}
	//TSL2561のチェック
	if tmp, tcode := healthadd(SennserData.Tsl2561_data.Flag, SennserData.Tsl2561_data.Name, SennserData.Tsl2561_data.Message); tcode != 404 {
		outdata = append(outdata, tmp)
	}
	//CO2Sensorのチェック
	if tmp, tcode := healthadd(SennserData.CO2Sensor_data.Flag, SennserData.CO2Sensor_data.Name, SennserData.CO2Sensor_data.Message); tcode != 404 {
		outdata = append(outdata, tmp)
	}
	//DHTSenserのチェック
	if tmp, tcode := healthadd(SennserData.DhtSenser_data.Flag, SennserData.DhtSenser_data.Name, SennserData.DhtSenser_data.Message); tcode != 404 {
		outdata = append(outdata, tmp)
	}
	//MMA8452Qのチェック
	if tmp, tcode := healthadd(SennserData.Mma8452q_data.Flag, SennserData.Mma8452q_data.Name, SennserData.Mma8452q_data.Message); tcode != 404 {
		outdata = append(outdata, tmp)
	}
	tmpdata, _ := json.Marshal(outdata)
	w.WriteHeader(code)
	fmt.Fprintf(w, "%s", string(tmpdata))
}

func createLineMetrics(name, types, value string) string {
	if name == "" {
		return "senser_data{type=\"" + types + "\"} " + value

	}
	return "senser_data{type=\"" + types + "\",sennser=\"" + name + "\"} " + value
}

// metricsの結果確認
func metrics(w http.ResponseWriter, r *http.Request) {
	var output []string
	tmpdata := SennserData
	SennserDataValue.Mu.Lock()
	tmpvaule := SennserDataValue
	SennserDataValue.Mu.Unlock()

	if tmpdata.Bme280_data.Flag {
		output = append(output, createLineMetrics(tmpdata.Bme280_data.Name, "hum", tmpvaule.Bme280.Hum))
		output = append(output, createLineMetrics(tmpdata.Bme280_data.Name, "press", tmpvaule.Bme280.Press))
		output = append(output, createLineMetrics(tmpdata.Bme280_data.Name, "tmp", tmpvaule.Bme280.Temp))
	}
	if tmpdata.Am2320_data.Flag {
		output = append(output, createLineMetrics(tmpdata.Am2320_data.Name, "hum", tmpvaule.Am2320.Hum))
		output = append(output, createLineMetrics(tmpdata.Am2320_data.Name, "tmp", tmpvaule.Am2320.Temp))
	}
	if tmpdata.Tsl2561_data.Flag {
		output = append(output, createLineMetrics(tmpdata.Tsl2561_data.Name, "lux", tmpvaule.Tsl2561.Lux))
	}
	if tmpdata.CO2Sensor_data.Flag {
		output = append(output, createLineMetrics(tmpdata.CO2Sensor_data.Name, "co2", tmpvaule.CO2.Co2))
		output = append(output, createLineMetrics(tmpdata.CO2Sensor_data.Name, "tmp", tmpvaule.CO2.Temp))
	}
	if tmpdata.DhtSenser_data.Flag {
		output = append(output, createLineMetrics(tmpdata.DhtSenser_data.Name, "hum", tmpvaule.DhtSenser.Hum))
		output = append(output, createLineMetrics(tmpdata.DhtSenser_data.Name, "tmp", tmpvaule.DhtSenser.Temp))
	}
	if tmpdata.Mma8452q_data.Flag {
		output = append(output, createLineMetrics(tmpdata.Mma8452q_data.Name, "ax", tmpvaule.Mma8452q.X))
		output = append(output, createLineMetrics(tmpdata.Mma8452q_data.Name, "ay", tmpvaule.Mma8452q.Y))
		output = append(output, createLineMetrics(tmpdata.Mma8452q_data.Name, "az", tmpvaule.Mma8452q.Z))
		output = append(output, createLineMetrics(tmpdata.Mma8452q_data.Name, "zero_x", tmpvaule.Mma8452q.Zero_X))
		output = append(output, createLineMetrics(tmpdata.Mma8452q_data.Name, "zero_y", tmpvaule.Mma8452q.Zero_Y))
		output = append(output, createLineMetrics(tmpdata.Mma8452q_data.Name, "zero_z", tmpvaule.Mma8452q.Zero_Z))
	}
	output = append(output, createLineMetrics("", "cpu_tmp", tmpvaule.CpuTmp))
	for _, line := range output {
		fmt.Fprintf(w, "%s\n", line)
	}
}

// ルートフォルダの結果
func rootdate(w http.ResponseWriter, r *http.Request) {
	output := "<html><body><a href=\"/metrics\">metrics</a></body></html>"
	fmt.Fprintf(w, "%s", output)
}

func reset(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		SennserResetSet(true)
		output := "<html><body><a href=\"/\">index</a></body></html>"
		fmt.Fprintf(w, "%s", output)

	} else {
		output := "<html><body><a href=\"/metrics\">metrics</a></body></html>"
		fmt.Fprintf(w, "%s", output)

	}
}
