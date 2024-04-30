package senser

import (
	"context"
	i2csenser "eijent/controller/senser/i2c_senser"
	msgsenser "eijent/controller/senser/msg_senser"
	rpisenser "eijent/controller/senser/rpi_senser"
	"log"
	"sync"
)

func Init() error {
	if err := rpisenser.Init(); err != nil {
		return err
	}
	if err := i2csenser.Init(); err != nil {
		log.Println("error:", err)
	}
	return nil
}

func Run(ctx context.Context) error {
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		if err := i2csenser.Run(); err != nil {
			log.Println("error:", err)
		}
	}()
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
	wg.Add(2)
	go func() {
		defer wg.Done()
		if err := i2csenser.Stop(); err != nil {

			log.Println("errors:", err)
		}
	}()
	go func(ctx context.Context) {
		defer wg.Done()
		if err := rpisenser.Stop(); err != nil {
			log.Println("errors:", err)
		}
	}(ctx)
	wg.Wait()
	return nil
}

func Wait() {
	rpisenser.Wait()
	i2csenser.Wait()
}

type SenserData struct {
	Senser string `json:"Senser"`
	Type   string `json:"Type"`
	Data   string `json:"Data"`
}

// 各有効センサーの読み取り結果
type OuteSenserData struct {
	Data []*SenserData
}

// 各有効センサの読み取り
func ReadValue() OuteSenserData {
	var out OuteSenserData
	//i2cに関連したセンサーの読み取り
	for _, msg := range i2csenser.Read() {
		tmp := new(SenserData)
		tmp.Senser = msg.Senser
		tmp.Type = msg.Type
		tmp.Data = msg.Data
		out.Data = append(out.Data, tmp)
	}
	//Raspberry Piのセンサー読み取り
	tmpRpi := rpisenser.ReadNow()
	if tmpRpi.Temp != "" {
		tmp := new(SenserData)
		tmp.Senser = "localhost"
		tmp.Type = "cpu_tmp"
		tmp.Data = tmpRpi.Temp
		out.Data = append(out.Data, tmp)
	}
	return out

}

type HealthData struct {
	msgsenser.HealthData
}

// health状態確認
func Health() []HealthData {
	var out []HealthData
	for _, tmp := range i2csenser.Health() {
		tt := HealthData{tmp}
		out = append(out, tt)
	}
	return out
}

// senserの再確認動作
func Reset() {
	i2csenser.Reset()
	return
}
