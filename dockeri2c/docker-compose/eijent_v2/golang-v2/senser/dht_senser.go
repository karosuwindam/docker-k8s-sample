package senser

import (
	"fmt"
	"os"
	"strconv"

	"github.com/d2r2/go-dht"
)

type DhtSenser struct {
	Flag        bool
	Name        string
	Message     string
	sennserType dht.SensorType
	pin         int
}

func (t *DhtSenser) Test() bool {
	_, _, _, err :=
		dht.ReadDHTxxWithRetry(t.sennserType, t.pin, false, 10)
	if err != nil {
		fmt.Println(err.Error())
		t.Message = err.Error()
		t.Flag = false
		return false
	}
	t.Message = "OK"
	t.Flag = true
	return true
}

func (t *DhtSenser) Init() {
	tmp := dht.DHT11
	if int(tmp) != 0 {
	}
	t.sennserType = dht.DHT11
	t.Name = "DHT11"
	if str := os.Getenv("SENSER_TYPE"); str != "" {
		t.Name = str
		switch str {
		case "DHT11":
			t.sennserType = dht.DHT11
			break
		case "DHT12":
			t.sennserType = dht.DHT12
			break
		case "AM2302":
			t.sennserType = dht.AM2302
			break
		}
	}
	t.pin = 12
	if str := os.Getenv("DHT11_PORT"); str != "" {
		tmp, err := strconv.Atoi(str)
		if err == nil {
			t.pin = tmp
		}
	}
}

func (t *DhtSenser) Read() (float64, float64) {
	temperature, humidity, _, err :=
		dht.ReadDHTxxWithRetry(t.sennserType, t.pin, false, 10)
	if err != nil {
		fmt.Println(err.Error())
		t.Message = err.Error()
		return -1, -1
	}
	hum := float64(humidity)
	tmp := float64(temperature)
	t.Message = "OK"
	return hum, tmp
}
