package senser

import (
	"context"
	rpisenser "eijent/controller/senser/rpi_senser"
	"log"
	"sync"
)

func Init() error {
	if err := rpisenser.Init(); err != nil {
		return err
	}
	return nil
}

func Run(ctx context.Context) error {
	var wg sync.WaitGroup
	wg.Add(1)
	go func(ctx context.Context) {
		defer wg.Done()
		if err := rpisenser.Run(ctx); err != nil {
			log.Println("error:", err)
		}
	}(ctx)
	wg.Wait()
	return nil
}

func Stop(ctx context.Context) error {
	var wg sync.WaitGroup
	wg.Add(1)
	go func(ctx context.Context) {
		if err := rpisenser.Stop(); err != nil {
			log.Println("errors:", err)
		}
	}(ctx)
	wg.Wait()
	return nil
}

type OuteSenserData struct { //各有効センサーの読み取り結果
	CpuTmp *string
}

func ReadValue() OuteSenserData { //各有効センサの読み取り
	var out OuteSenserData
	//Raspberry Piのセンサー読み取り
	tmpRpi := rpisenser.ReadNow()
	if tmpRpi.Temp != "" {
		out.CpuTmp = new(string)
		*out.CpuTmp = tmpRpi.Temp
	}
	return out

}
