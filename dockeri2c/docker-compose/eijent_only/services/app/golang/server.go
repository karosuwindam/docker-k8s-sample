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
	// Co2senser Co2Sennser
	Co2senser MhZ19c
	Bme280    Bme280
}

type RaspberrypiData struct {
	cpu_tmp string
}
type Co2Data struct {
	Tmp int
	Co2 int
}
type MulData struct {
	Tmp   float64
	Hum   float64
	Press float64
}
type DataType struct {
	Hum  float64
	Tmp  float64
	Lux  int
	Co2  Co2Data
	MuDa MulData
	Rpi  RaspberrypiData
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
	if t.Sennser.Co2senser.Flag {
		tmp.Senser = t.Sennser.Co2senser.Name
		tmp.Type = "co2"
		tmp.Data = strconv.Itoa(t.Data.Co2.Co2)
		outdata = append(outdata, tmp)
		tmp.Type = "tmp"
		tmp.Data = strconv.Itoa(t.Data.Co2.Tmp)
		outdata = append(outdata, tmp)
	}
	if t.Sennser.Bme280.Flag {
		tmp.Senser = t.Sennser.Bme280.Name
		tmp.Type = "hum"
		tmp.Data = strconv.FormatFloat(t.Data.MuDa.Hum, 'f', 2, 64)
		outdata = append(outdata, tmp)
		tmp.Type = "tmp"
		tmp.Data = strconv.FormatFloat(t.Data.MuDa.Tmp, 'f', 2, 64)
		outdata = append(outdata, tmp)
		tmp.Type = "press"
		tmp.Data = strconv.FormatFloat(t.Data.MuDa.Press, 'f', 2, 64)
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
	} else {
		tmp.Sennserdata = "Temp Sensor"
		tmp.Message = "OFF"
		outdata = append(outdata, tmp)
	}
	if t.Sennser.Tsl2561.Flag {
		tmp.Sennserdata = t.Sennser.Tsl2561.Name
		tmp.Message = t.Sennser.Tsl2561.Message
		if tmp.Message != "OK" {
			code = 405
		}
		outdata = append(outdata, tmp)
	} else {
		tmp.Sennserdata = "Lux Sensor"
		tmp.Message = "OFF"
		outdata = append(outdata, tmp)
	}
	if t.Sennser.Co2senser.Flag {
		tmp.Sennserdata = t.Sennser.Co2senser.Name
		tmp.Message = t.Sennser.Co2senser.Message
		if tmp.Message != "OK" {
			code = 405
		}
		outdata = append(outdata, tmp)
	} else {
		tmp.Sennserdata = "CO2 sensor"
		tmp.Message = "OFF"
		outdata = append(outdata, tmp)
	}
	if t.Sennser.Bme280.Flag {
		tmp.Sennserdata = t.Sennser.Bme280.Name
		tmp.Message = t.Sennser.Bme280.Message
		if tmp.Message != "OK" {
			code = 405
		}
		outdata = append(outdata, tmp)

	} else {
		tmp.Sennserdata = "Bme Sensor"
		tmp.Message = "OFF"
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
		output += "senser_data{type=\"tmp\",sennser=\"AM2320\"} " + strconv.FormatFloat(t.Data.Tmp, 'f', 1, 64)
		output += "\n" + "senser_data{type=\"hum\",sennser=\"AM2320\"} " + strconv.FormatFloat(t.Data.Hum, 'f', 1, 64)
	} else if t.Sennser.Dht.Flag {
		output += "senser_data{type=\"tmp\",sennser=\"DHT11\"} " + strconv.FormatFloat(t.Data.Tmp, 'f', 1, 64)
		output += "\n" + "senser_data{type=\"hum\",sennser=\"DHT11\"} " + strconv.FormatFloat(t.Data.Hum, 'f', 1, 64)
	}
	if t.Sennser.Tsl2561.Flag {
		if output != "" {
			output += "\n"
		}
		output += "senser_data{type=\"lux\",sennser=\"TSL-2561\"} " + strconv.Itoa(t.Data.Lux)
	}
	if t.Sennser.Co2senser.Flag {
		if output != "" {
			output += "\n"
		}
		output += "senser_data{type=\"co2\"} " + strconv.Itoa(t.Data.Co2.Co2)
		output += "\n" + "senser_data{type=\"tmp\",sennser=\"co2\"}" + strconv.Itoa(t.Data.Co2.Tmp)
	}
	if t.Sennser.Bme280.Flag {
		if output != "" {
			output += "\n"
		}
		output += "senser_data{type=\"tmp\",sennser=\"BME280\"} " + strconv.FormatFloat(t.Data.MuDa.Tmp, 'f', 2, 64)
		output += "\n" + "senser_data{type=\"hum\",sennser=\"BME280\"} " + strconv.FormatFloat(t.Data.MuDa.Hum, 'f', 2, 64)
		output += "\n" + "senser_data{type=\"press\",sennser=\"BME280\"} " + strconv.FormatFloat(t.Data.MuDa.Press, 'f', 2, 64)
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
