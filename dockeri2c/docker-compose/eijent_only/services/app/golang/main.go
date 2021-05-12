package main

import (
	"time"
)

type AM2320Data struct {
	Id        int64     `json:"id"`
	Tmp       float32   `json:"tmp"`
	Hum       float32   `json:"hum"`
	CreatedAt time.Time `json:"createdAt"`
}

var (
	I2C_BUS = 1
)

func senserDataCk(server *ServerData) {
	if server.Sennser.Am2320.Flag {
		hum, tmp := server.Sennser.Am2320.Read()
		server.Data.Hum = float64(hum)
		server.Data.Tmp = float64(tmp)
	} else if server.Sennser.Dht.Flag {
		hum, tmp := server.Sennser.Dht.Read()
		server.Data.Hum = float64(hum)
		server.Data.Tmp = float64(tmp)
	}
	if server.Sennser.Tsl2561.Flag {
		lux := server.Sennser.Tsl2561.ReadVisibleLux()
		server.Data.Lux = lux
	}
	if server.Sennser.Co2senser.Flag {
		co2ppm, temp := server.Sennser.Co2senser.Read()
		if co2ppm > 0 {
			server.Data.Co2.Co2 = co2ppm
			server.Data.Co2.Tmp = temp	
		}
	}
	if server.Sennser.Bme280.Flag {
		press, temp, hum := server.Sennser.Bme280.ReadData()
		server.Data.MuDa = MulData{Tmp: temp, Hum: hum, Press: press}
	}
	server.Data.Rpi.cpu_tmp = cpuTmp()
}

func main() {
	server := ServerInt()
	server.Sennser.Am2320.Init()
	for i := 0; i < 3; i++ {
		if server.Sennser.Am2320.Test() {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}
	if !server.Sennser.Am2320.Flag {
		server.Sennser.Dht.Init()
		for i := 0; i < 1; i++ {
			if server.Sennser.Dht.Test() {
				break
			}
			time.Sleep(100 * time.Millisecond)
		}
	}
	server.Sennser.Tsl2561.Init()
	server.Sennser.Co2senser.Init("/dev/ttyS0")

	for i := 0; i < 3; i++ {
		if server.Sennser.Tsl2561.Test() {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}

	server.Sennser.Bme280.Init()

	senserDataCk(&server)

	go func() {
		for {
			senserDataCk(&server)

			time.Sleep(15 * time.Second)
		}
	}()
	server.ServerStart()

}
