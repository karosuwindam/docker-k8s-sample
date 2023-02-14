package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/d2r2/go-dht"

	logger "github.com/d2r2/go-logger"
)

var lg = logger.NewPackageLogger("main",
	logger.DebugLevel,
	// logger.InfoLevel,
)
var data_hum, data_tmp float64

// var cpu_tmp string

func metrics(w http.ResponseWriter, r *http.Request) {
	output := ""
	if data_hum != 0 || data_tmp != 0 {
		output += "senser_data{type=\"tmp\"} " + strconv.FormatFloat(data_tmp, 'f', 1, 64)
		output += "\n" + "senser_data{type=\"hum\"} " + strconv.FormatFloat(data_hum, 'f', 1, 64)
	}
	// output += "\n" + "senser_data{type=\"cpu_tmp\"} " + cpu_tmp
	fmt.Fprintf(w, "%s", output)

}
func data(w http.ResponseWriter, r *http.Request) {
	output := "<html><body><a href=\"/metrics\">metrics</a></body></html>"
	fmt.Fprintf(w, "%s", output)

}

func main() {
	defer logger.FinalizeLogger()

	lg.Notify("***************************************************************************************************")
	lg.Notify("*** You can change verbosity of output, to modify logging level of module \"dht\"")
	lg.Notify("*** Uncomment/comment corresponding lines with call to ChangePackageLogLevel(...)")
	lg.Notify("***************************************************************************************************")
	// Uncomment/comment next line to suppress/increase verbosity of output
	// logger.ChangePackageLogLevel("dht", logger.InfoLevel)

	sensorType := dht.DHT11
	if str := os.Getenv("SENSER_TYPE"); str != "" {
		switch str {
		case "DHT11":
			sensorType = dht.DHT11
			break
		case "DHT12":
			sensorType = dht.DHT12
			break
		case "AM2302":
			sensorType = dht.AM2302
			break
		}
	}
	// Read DHT11 sensor data from specific pin, retrying 10 times in case of failure.
	pin := 12
	if str := os.Getenv("DHT11_PORT"); str != "" {
		tmp, err := strconv.Atoi(str)
		if err == nil {
			pin = tmp
		}
	}
	temperature, humidity, retried, err :=
		dht.ReadDHTxxWithRetry(sensorType, pin, false, 10)
	if err != nil {
		lg.Fatal(err)
	}
	data_hum = float64(humidity)
	data_tmp = float64(temperature)
	lg.Infof("Sensor = %v: Temperature = %v*C, Humidity = %v%% (retried %d times)",
		sensorType, temperature, humidity, retried)
	go func() {
		for {
			nowtime := time.Now()

			temperature, humidity, retried, err :=
				dht.ReadDHTxxWithRetry(sensorType, pin, false, 10)
			if err != nil {
				lg.Fatal(err)
			}
			data_hum = float64(humidity)
			data_tmp = float64(temperature)
			// print temperature and humidity
			lg.Infof("Sensor = %v: Temperature = %v*C, Humidity = %v%% (retried %d times)",
				sensorType, temperature, humidity, retried)
			duration_time := (time.Now()).Sub(nowtime)
			if 5*time.Second > duration_time {
				time.Sleep(10 * time.Second)
			} else {
				time.Sleep(5*time.Second - duration_time)

			}

		}
	}()
	port := "9141"
	fmt.Println("start server " + port)
	http.HandleFunc("/metrics", metrics)
	http.HandleFunc("/", data)
	http.ListenAndServe(":"+port, nil)
	fmt.Println("stop server")
}
