package main

import (
	"fmt"
	"strconv"
	"time"
)

var (
	co2limit   = int(1000)
	co2chktime = time.Duration(1000) //ms
)

type Ckdata struct {
	data int
}

func (t *Ckdata) check_data() bool {
	flag := false
	for _, data := range getdata() {
		// fmt.Println(data)

		if data.Type == "co2" {
			co2data, _ := strconv.Atoi(data.Data)
			t.data = co2data
			if co2data > co2limit {
				flag = true
			}
		}
	}
	return flag
}

func main() {
	//GPIOの初期化処理
	var data GpioLed
	data.PinNb = []int{5}
	data.InOut = []int{GPIO_OUT}
	err := data.Open()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	_, err = ConfJsonUrlRead()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	gpiodata := gpioInt(data.PinNb)
	defer data.Close()
	go func() {
		for {
			gpiodata.gpiostart(&data)
		}
	}()
	gpiodata.OutPin[0] = LED_OFF
	var co2data Ckdata
	go func() {
		count := 1
		for {
			if co2data.check_data() {
				gpiodata.OutPin[0] = LED_ON
				if count != 0 {
					fmt.Println("LED", LED_ON, co2data.data)
				}
				count = 0
			} else if count > 5 {
				if (count % 10) == 0 {
					fmt.Println("LED", LED_OFF, co2data.data)
				}
				gpiodata.OutPin[0] = LED_OFF
				count++
			} else {
				gpiodata.OutPin[0] = LED_TOGGLE
				count++
			}
			time.Sleep(time.Millisecond * co2chktime)
		}
	}()
	for {
		time.Sleep(time.Millisecond * 10000)
	}
}
