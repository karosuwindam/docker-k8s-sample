package main

import (
	"log"
	"time"

	"github.com/davecheney/i2c"
)

var (
	AM2320  = uint8(0x5c)
	I2C_BUS = 1
)

func ReadAM2320() (float32, float32) {
	i2c, err := i2c.New(AM2320, I2C_BUS)
	if err != nil {
		log.Fatal(err)
	}
	defer i2c.Close()
	i2c.Write([]byte{0x0})
	time.Sleep(time.Microsecond * 10)

	i2c.Write([]byte{0x3, 0x0, 0x4})
	time.Sleep(time.Microsecond * 15)
	buf := make([]byte, 6)
	i2c.Read(buf)
	hum := float32((uint(buf[2])<<8)|uint(buf[3])) / 10 //湿度
	tmp := float32((uint(buf[4])<<8)|uint(buf[5])) / 10 //温度
	return hum, tmp
}

