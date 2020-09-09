package main

import (
	"time"

	"github.com/stianeikeland/go-rpio"
)

func main() {
	err := rpio.Open()
	if err != nil {
		return
	}
	pin := rpio.Pin(21)
	pin.Output()
	pin.High()
	time.Sleep(time.Millisecond * 1000)
	pin.Low()
	rpio.Close()
}
