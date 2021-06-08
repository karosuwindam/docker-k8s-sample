package main

import (
	"errors"
	"fmt"
	"time"

	"github.com/stianeikeland/go-rpio"
)

type GpioLed struct {
	PinNb  []int
	PinOut []int
	InOut  []int
	pin    []rpio.Pin
	// Pin
}

const (
	GPIO_LOW  = 0
	GPIO_HIGH = 1
	GPIO_OUT  = 1
	GPIO_IN   = 0
)
const (
	LED_OFF    = 0
	LED_ON     = 1
	LED_TOGGLE = 2
)

type GpioSys struct {
	OutPin []int
}

//gpioの設定初期化
func gpioInt(data []int) GpioSys {
	var tmp GpioSys
	for i := 0; i < len(data); i++ {
		tmp.OutPin = append(tmp.OutPin, LED_OFF)
	}
	return tmp
}

// LEDの制御
func (t *GpioSys) gpiostart(data *GpioLed) {
	for i := 0; i < len(data.PinNb); i++ {
		if data.InOut[i] == GPIO_OUT {
			switch t.OutPin[i] {
			case LED_OFF:
				data.Low(i)
			case LED_ON:
				data.High(i)
			case LED_TOGGLE:
				//点滅処理
				if data.PinOut[i] == GPIO_LOW {
					data.High(i)
				} else {
					data.Low(i)
				}
			}
		}
	}
	time.Sleep(time.Millisecond * 100)
}

func (t *GpioLed) Open() error {
	err := rpio.Open()
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	if len(t.PinNb) != len(t.InOut) {
		return errors.New("not count Pin Nunber Inout count")
	}
	for i, num := range t.PinNb {
		if num == 0 {
			return errors.New("not setup Pin Nuber")
		}
		// fmt.Println(i)
		pin := rpio.Pin(num)
		if t.InOut[i] == GPIO_OUT {
			pin.Output()
		} else {
			pin.Input()
		}
		t.PinOut = append(t.PinOut, GPIO_LOW)
		t.pin = append(t.pin, pin)
	}
	return nil
}

func (t *GpioLed) High(num int) {
	if num < len(t.pin) {
		t.pin[num].High()
		t.PinOut[num] = GPIO_HIGH
	} else {
		fmt.Println("error input number")
	}
}
func (t *GpioLed) Low(num int) {
	if num < len(t.pin) {
		t.pin[num].Low()
		t.PinOut[num] = GPIO_LOW
	} else {
		fmt.Println("error input number")
	}
}

func (t *GpioLed) Close() {
	for _, pin := range t.pin {
		pin.Low()
	}
	rpio.Close()
}
