package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"
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
		if hum == tmp && hum == -1 {

		} else {
			server.Data.Hum = float64(hum)
			server.Data.Tmp = float64(tmp)
		}
	} else if server.Sennser.Dht.Flag {
		hum, tmp := server.Sennser.Dht.Read()
		if hum == tmp && hum == -1 {

		} else {
			server.Data.Hum = float64(hum)
			server.Data.Tmp = float64(tmp)
		}
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
	if server.Sennser.Mma8452.Flag {
		x, y, z := server.Sennser.Mma8452.ReadData()
		ax, ay, az := server.Sennser.Mma8452.Ch_data_to_accel(x, y, z)
		server.Data.Acc = AccelData{Accele_x: ax, Accele_y: ay, Accele_z: az}
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
	server.Sennser.Mma8452.Init()

	senserDataCk(&server)
	ctx := context.Background()
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer stop()
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		go func() {
			for {
				senserDataCk(&server)

				time.Sleep(15 * time.Second)
			}
		}()
		server.ServerStart()
		return nil
	})
	<-ctx.Done()
	if err := eg.Wait(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
	fmt.Println("end")
}
