package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type SennserData struct {
	Am2320  Am2320
	Dht     DhtSenser
	Tsl2561 Tsl2561
}

type RaspberrypiData struct {
	cpu_tmp string
}
type DataType struct {
	Hum float64
	Tmp float64
	Lux int
	Rpi RaspberrypiData
}

type ServerData struct {
	Port    string
	Ip      string
	Sennser SennserData
	Data    DataType
}

type JsonSenser struct {
	Senser string `json:senser`
	Type   string `json:type`
	Data   string `json:data`
}

type HealthData struct {
	Sennserdata string `json:sennserdata`
	Message     string `json:message`
}

func (t *ServerData) jsonData(w http.ResponseWriter, r *http.Request) {
	var outdata []JsonSenser
	var tmp JsonSenser
	if t.Sennser.Am2320.Flag {
		tmp.Senser = t.Sennser.Am2320.Name
		tmp.Type = "tmp"
		tmp.Data = strconv.FormatFloat(t.Data.Tmp, 'f', 1, 64)
		outdata = append(outdata, tmp)
		tmp.Type = "hum"
		tmp.Data = strconv.FormatFloat(t.Data.Hum, 'f', 1, 64)
		outdata = append(outdata, tmp)
	} else if t.Sennser.Dht.Flag {
		tmp.Senser = t.Sennser.Dht.Name
		tmp.Type = "tmp"
		tmp.Data = strconv.FormatFloat(t.Data.Tmp, 'f', 1, 64)
		outdata = append(outdata, tmp)
		tmp.Type = "hum"
		tmp.Data = strconv.FormatFloat(t.Data.Hum, 'f', 1, 64)
		outdata = append(outdata, tmp)
	}
	if t.Sennser.Tsl2561.Flag {
		tmp.Senser = t.Sennser.Tsl2561.Name
		tmp.Type = "lux"
		tmp.Data = strconv.Itoa(t.Data.Lux)
		outdata = append(outdata, tmp)
	}
	tmp.Senser = "raspberrypi"
	tmp.Type = "cpu_tmp"
	tmp.Data = t.Data.Rpi.cpu_tmp
	outdata = append(outdata, tmp)

	output := ""
	tmpdata, _ := json.Marshal(outdata)
	output += string(tmpdata)
	fmt.Fprintf(w, "%s", output)
}
func (t *ServerData) health(w http.ResponseWriter, r *http.Request) {
	var outdata []HealthData
	var tmp HealthData
	code := 200
	if t.Sennser.Am2320.Flag {
		tmp.Sennserdata = t.Sennser.Am2320.Name
		tmp.Message = t.Sennser.Am2320.Message
		if tmp.Message != "OK" {
			code = 405
		}
		outdata = append(outdata, tmp)
	} else if t.Sennser.Dht.Flag {
		tmp.Sennserdata = t.Sennser.Dht.Name
		tmp.Message = t.Sennser.Dht.Message
		if tmp.Message != "OK" {
			code = 405
		}
		outdata = append(outdata, tmp)
	}
	if t.Sennser.Tsl2561.Flag {
		tmp.Sennserdata = t.Sennser.Tsl2561.Name
		tmp.Message = t.Sennser.Tsl2561.Message
		if tmp.Message != "OK" {
			code = 405
		}
		outdata = append(outdata, tmp)
	}
	if len(outdata) < 1 {
		tmp.Sennserdata = "Raspberrypi"
		tmp.Message = "OK"
		outdata = append(outdata, tmp)
	}

	output := ""
	tmpdata, _ := json.Marshal(outdata)
	output += string(tmpdata)
	w.WriteHeader(code)
	fmt.Fprintf(w, "%s", output)
}
func (t *ServerData) metrics(w http.ResponseWriter, r *http.Request) {
	output := ""
	if t.Sennser.Am2320.Flag {
		output += "senser_data{type=\"tmp\"} " + strconv.FormatFloat(t.Data.Tmp, 'f', 1, 64)
		output += "\n" + "senser_data{type=\"hum\"} " + strconv.FormatFloat(t.Data.Hum, 'f', 1, 64)
	} else if t.Sennser.Dht.Flag {
		output += "senser_data{type=\"tmp\"} " + strconv.FormatFloat(t.Data.Tmp, 'f', 1, 64)
		output += "\n" + "senser_data{type=\"hum\"} " + strconv.FormatFloat(t.Data.Hum, 'f', 1, 64)
	}
	if t.Sennser.Tsl2561.Flag {
		if output != "" {
			output += "\n"
		}
		output += "senser_data{type=\"lux\"} " + strconv.Itoa(t.Data.Lux)
	}
	if output != "" {
		output += "\n"
	}
	output += "senser_data{type=\"cpu_tmp\"} " + t.Data.Rpi.cpu_tmp
	fmt.Fprintf(w, "%s", output)

}
func (t *ServerData) rootdata(w http.ResponseWriter, r *http.Request) {
	output := "<html><body><a href=\"/metrics\">metrics</a></body></html>"
	fmt.Fprintf(w, "%s", output)

}

func ServerInt() ServerData {
	var tmp ServerData
	tmp.Port = "9140"
	tmp.Ip = ""
	return tmp
}

func (t *ServerData) ServerStart() {
	fmt.Println("start server " + t.Port)
	http.HandleFunc("/metrics", t.metrics)
	http.HandleFunc("/health", t.health)
	http.HandleFunc("/json", t.jsonData)
	http.HandleFunc("/", t.rootdata)
	http.ListenAndServe(t.Ip+":"+t.Port, nil)
}
