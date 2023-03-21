package senser

import (
	"fmt"
	"os"
	"time"

	"github.com/d2r2/go-dht"
)

type DhtSenser struct {
	Flag        bool
	Name        string
	Message     string
	sennserType dht.SensorType
	pin         int
}

type DhtSenser_Vaule struct {
	Hum  string
	Temp string
}

func (t *DhtSenser) Test() bool {
	_, _, _, err :=
		dht.ReadDHTxxWithRetry(t.sennserType, t.pin, false, 2)
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

func (t *DhtSenser) Init(pin int) bool {
	tmp := dht.DHT11
	flag := true
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
	t.pin = pin

	for i := 0; i < 3; i++ {
		flag = false
		if t.Test() {
			flag = true
			break
		}
		time.Sleep(500 * time.Millisecond)
	}
	return flag
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
