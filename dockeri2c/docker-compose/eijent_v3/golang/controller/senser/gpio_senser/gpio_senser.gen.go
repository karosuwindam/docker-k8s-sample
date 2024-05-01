package gpiosenser

import (
	"log"
	"sync"
	"time"

	dhtsenser "eijent/controller/senser/gpio_senser/dht_senser"
	msgsenser "eijent/controller/senser/msg_senser"
)

type APIList interface {
	Init() error
	Stop() error
	Run() error
	Health() (bool, msgsenser.Msg)
	Wait()
	Read() ([]msgsenser.SenserData, bool)
	Reset()
	Name() string
}

var apilists map[string]APIList = make(map[string]APIList)

func AddApi(api APIList) {
	apilists[api.Name()] = api
}

func Init() error {
	AddApi(dhtsenser.NewAPI())
	var wg sync.WaitGroup
	wg.Add(len(apilists))
	for _, api := range apilists {
		go func() {
			defer wg.Done()
			if err := api.Init(); err != nil {
				log.Println("error:", err)
			}
		}()
	}
	wg.Wait()
	return nil
}

func Run() error {
	var wg sync.WaitGroup
	wg.Add(len(apilists))

	for _, api := range apilists {
		go func() {
			defer wg.Done()
			if err := api.Run(); err != nil {
				log.Println("error:", err)

			}
		}()
	}
	wg.Wait()
	return nil
}

func Stop() error {

	var wg sync.WaitGroup
	wg.Add(len(apilists))
	for _, api := range apilists {
		go func() {
			defer wg.Done()
			if err := api.Stop(); err != nil {
				log.Println("error:", err)

			}
		}()
	}
	wg.Wait()
	return nil
}

func Wait() {

	var wg sync.WaitGroup
	wg.Add(len(apilists))
	for _, api := range apilists {
		go func() {
			defer wg.Done()
			api.Wait()
		}()
	}
	wg.Wait()
}

func Read() []msgsenser.SenserData {
	var out []msgsenser.SenserData
	for _, api := range apilists {
		vs, ok := api.Read()
		if !ok {
			continue
		}
		for _, v := range vs {
			out = append(out, v)
		}
	}
	return out
}

func Health() []msgsenser.HealthData {
	var out []msgsenser.HealthData
	for _, api := range apilists {
		flag, msg := api.Health()
		if tmp, ok := addHealth(flag, msg); ok {
			out = append(out, tmp)
		}
	}
	return out
}

func Reset() {
	var wg sync.WaitGroup
	wg.Add(len(apilists))
	for _, api := range apilists {
		go func() {
			defer wg.Done()
			api.Reset()

		}()
	}
	wg.Wait()
}

func addHealth(f bool, msg msgsenser.Msg) (msgsenser.HealthData, bool) {
	var out msgsenser.HealthData
	out.Run = "Stop"
	flag := false
	if f || (time.Now().Sub(msg.MessageTime) < time.Hour) {
		out.Sennserdata = msg.Senser
		out.Message = msg.Message
		if f {
			out.Run = "Start"
		}
		flag = true
	}
	return out, flag
}
