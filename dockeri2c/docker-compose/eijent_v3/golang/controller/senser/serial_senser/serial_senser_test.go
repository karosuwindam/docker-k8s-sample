package serialsenser_test

import (
	"eijent/config"
	"eijent/controller/senser/serial_senser/co2sennser"
	"log"
	"testing"
	"time"
)

func TestUartSennser(t *testing.T) {
	config.Init()
	api := co2sennser.NewAPI()
	if err := api.Init(); err != nil {
		t.Fatal(err)
	}
	go func() {
		if err := api.Run(); err != nil {
			t.Fatal(err)
		}
	}()
	api.Wait()
	time.Sleep(time.Second)
	log.Println(api.Read())
	api.Stop()
}
