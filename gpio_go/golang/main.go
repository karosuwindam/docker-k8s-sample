package main

import (
	"fmt"
	"strconv"
	"time"
)

var (
	co2limit   = int(1000)
	avglimit   = int(60)
	co2chktime = time.Duration(1000) //ms
)

type Ckdata struct {
	data    int
	avg     int
	avgdate []int
	lux     int
}

func (t *Ckdata) check_data() bool {
	flag := false
	for _, gdata := range getdata() {
		// fmt.Println(data)

		if gdata.Type = "lux" {
			t.lux , _ := strconv.Atoi(gdata.Data)
		}
		if gdata.Type == "co2" {
			co2data, _ := strconv.Atoi(gdata.Data)
			t.data = co2data
			if len(t.avgdate) > avglimit {
				t.avgdate = t.avgdate[1:avglimit]
			}
			t.avgdate = append(t.avgdate, co2data)
			tmp := 0
			for _, data := range t.avgdate {
				tmp += data
			}
			t.avg = tmp / len(t.avgdate)
			if t.avg > co2limit {
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
	localzone := "Asia/Tokyo"
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
	co2data.avgdate = []int{}
	co2data.lux = -1
	go func() {
		count := 1
		for {
			loc, _ := time.LoadLocation(localzone)
			starttime := time.Now().In(loc)
			if co2data.lux == 0 {
				gpiodata.OutPin[0] = LED_OFF
				fmt.Println(starttime, "LED", LED_OFF, co2data.data, co2data.avg, co2data.lux)
			} else {
				if co2data.check_data() {
					gpiodata.OutPin[0] = LED_ON
					if count != 0 {
						fmt.Println(starttime, "LED", LED_ON, co2data.data, co2data.avg, co2data.lux)
					}
					count = 0
				} else if count > 5 {
					if (count % 10) == 0 {
						fmt.Println(starttime, "LED", LED_OFF, co2data.data, co2data.avg, co2data.lux)
					}
					gpiodata.OutPin[0] = LED_OFF
					count++
				} else {
					gpiodata.OutPin[0] = LED_TOGGLE
					count++
				}
			}
			time.Sleep(time.Millisecond * co2chktime)
		}
	}()
	for {
		time.Sleep(time.Millisecond * 10000)
	}
}
