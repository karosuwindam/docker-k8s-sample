package co2sennser

import (
	"log"
	"testing"
	"time"

	"github.com/tarm/serial"
)

func TestConfig(t *testing.T) {
	c, err := uartInitConfig("a", 1, byte(0x2), time.Nanosecond, serial.ParityMark, serial.Stop1)
	if err != nil {
		t.Fatal(err)
	}
	log.Println(c)
	c, err = uartInitConfig("a", 1, 0.0)
	if err == nil {
		t.Fatal(err)
	}
	log.Println(c, err)
}
