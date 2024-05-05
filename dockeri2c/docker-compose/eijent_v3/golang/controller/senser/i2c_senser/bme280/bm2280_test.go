package bme280_test

import (
	"eijent/controller/senser/i2c_senser/bme280"
	"fmt"
	"sync"
	"testing"
)

func TestRun(t *testing.T) {
	var i2cMu sync.Mutex
	if err := bme280.Init(&i2cMu); err != nil {
		t.Fatal(err)
	}
	go func() {
		bme280.Run()
	}()
	bme280.Wait()
	fmt.Println(bme280.ReadValue())
	api := bme280.NewAPI()
	fmt.Println(api.Read())

}
