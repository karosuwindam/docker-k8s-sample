package serialsenser_test

import (
	"eijent/controller/senser/serial_senser/co2sennser"
	"log"
	"testing"
	"time"
)

func TestUartSennser(t *testing.T) {
	api := co2sennser.NewAPI()
	if err := api.Init(); err != nil {
		t.Fatal(err)
	}
	api.Wait()
	go func() {
		if err := api.Run(); err != nil {
			t.Fatal(err)
		}
	}()
	time.Sleep(time.Microsecond)
	log.Println(api.Read())
	api.Stop()
}
