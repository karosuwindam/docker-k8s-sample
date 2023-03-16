package senser

import (
	"fmt"
	"time"

	"github.com/davecheney/i2c"
)

var (
	AM2320 = uint8(0x5c)
)

type Am2320 struct {
	Flag    bool
	Name    string
	Message string
}

func (t *Am2320) Init() {
	t.Name = "AM2320"
}

func (t *Am2320) Test() bool {
	t.Flag = false
	i2c, err := i2c.New(AM2320, I2C_BUS)
	if err != nil {
		fmt.Println(err.Error())
		t.Message = err.Error()
		return false
	}
	defer i2c.Close()
	_, err = i2c.Write([]byte{0x0})
	if err != nil {
		fmt.Println(err.Error())
		t.Message = err.Error()
		return false
	}
	time.Sleep(time.Microsecond * 10)

	_, err = i2c.Write([]byte{0x3, 0x0, 0x4})
	if err != nil {
		fmt.Println(err.Error())
		t.Message = err.Error()
		return false
	}
	time.Sleep(time.Microsecond * 15)
	buf := make([]byte, 6)
	_, err = i2c.Read(buf)
	if err != nil {
		fmt.Println(err.Error())
		t.Message = err.Error()
		return false
	}
	t.Flag = true
	t.Message = "OK"
	return true
}

func (t *Am2320) Read() (float32, float32) {
	i2c, err := i2c.New(AM2320, I2C_BUS)
	if err != nil {
		t.Message = err.Error()
		return -1, -1
	}
	defer i2c.Close()
	_, err = i2c.Write([]byte{0x0})
	if err != nil {
		t.Message = err.Error()
		// return -1, -1
	}
	time.Sleep(time.Microsecond * 10)

	_, err = i2c.Write([]byte{0x3, 0x0, 0x4})
	if err != nil {
		t.Message = err.Error()
		// return -1, -1
	}
	time.Sleep(time.Microsecond * 15)
	buf := make([]byte, 6)
	_, err = i2c.Read(buf)
	if err != nil {
		t.Message = err.Error()
		// return -1, -1
	}
	hum := float32((uint(buf[2])<<8)|uint(buf[3])) / 10 //湿度
	tmp := float32((uint(buf[4])<<8)|uint(buf[5])) / 10 //温度
	t.Message = "OK"
	return hum, tmp
}
