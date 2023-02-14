package main

import (
	"fmt"
	"time"
)

func main() {
	var co2 Co2Sennser
	if !co2.Init("/dev/serial0") {
		fmt.Println(co2.InitData)
		return
	}
	defer co2.Close()
	fmt.Println(co2.InitData)
	for i := 0; i < 5; i++ {
		co2ppm, temp := co2.Read()
		fmt.Println(co2.ReadData)
		fmt.Println(co2ppm, temp)
		time.Sleep(1500 * time.Millisecond)
	}
}
