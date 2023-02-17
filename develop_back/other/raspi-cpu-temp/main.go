package main

import (
	"time"
)

func main() {
	var web WebSetupData
	web.Data.Port = "9500"
	web.Tmp = cpuTmp()
	go func() {
		for {
			web.Tmp = cpuTmp()
			// fmt.Println(tmp)
			time.Sleep(1 * time.Second)
		}
	}()
	web.Webstart()
}
