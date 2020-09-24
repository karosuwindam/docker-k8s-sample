package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type AM2320Data struct {
	Id        int64     `json:"id"`
	Tmp       float32   `json:"tmp"`
	Hum       float32   `json:"hum"`
	CreatedAt time.Time `json:"createdAt"`
}

var data_hum, data_tmp float64
var cpu_tmp string

func metrics(w http.ResponseWriter, r *http.Request) {
	output := ""
	if data_hum != 0 || data_tmp != 0 {
		output += "senser_data{type=\"tmp\"} " + strconv.FormatFloat(data_tmp, 'f', 1, 64)
		output += "\n" + "senser_data{type=\"hum\"} " + strconv.FormatFloat(data_hum, 'f', 1, 64)
	}
	output += "\n" + "senser_data{type=\"cpu_tmp\"} " + cpu_tmp
	fmt.Fprintf(w, "%s", output)

}
func data(w http.ResponseWriter, r *http.Request) {
	output := "<html><body><a href=\"/metrics\">metrics</a></body></html>"
	fmt.Fprintf(w, "%s", output)

}

func main() {
	hum, tmp := ReadAM2320()
	data_hum = float64(hum)
	data_tmp = float64(tmp)
	cpu_tmp = cpuTmp()
	go func() {
		for {
			hum, tmp := ReadAM2320()
			data_hum = float64(hum)
			data_tmp = float64(tmp)
			cpu_tmp = cpuTmp()
			// sample1 := AM2320Data{Tmp: tmp, Hum: hum}
			// fmt.Printf("%v,%v\n", sample1.Tmp, sample1.Hum)
			time.Sleep(15 * time.Second)
		}
	}()

	port := "9140"
	fmt.Println("start server " + port)
	http.HandleFunc("/metrics", metrics)
	http.HandleFunc("/", data)
	http.ListenAndServe(":"+port, nil)
	fmt.Println("stop server")

}
