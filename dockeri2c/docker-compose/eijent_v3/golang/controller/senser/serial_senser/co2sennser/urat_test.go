package co2sennser

import (
	"log"
	"sync"
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

// このテストはTxとRxを直接接続した状態で行う
func TestUart(t *testing.T) {
	if err := uartInit("/dev/ttyAMA0", 9600); err != nil {
		// if err := uartInit("/dev/ttyAMA0", 9600, time.Second); err != nil {
		t.Fatal(err)
	}
	if err := uartdata.open(); err != nil {
		t.Fatal(err)
	}
	var wg sync.WaitGroup
	defer uartdata.close()
	var shutdown chan bool = make(chan bool, 1)
	wg.Add(1)
	// スレッドによる読み込み処理
	go func() {
		defer wg.Done()
	loop:
		for {
			log.Println("info:", "read start")
			buf, err := uartdata.read()
			log.Println("info", buf, err)
			select {
			case <-shutdown:
				break loop
			case <-time.After(time.Millisecond * 10):
			}
		}
		log.Println("stop loop")
		return
	}()
	//定期的に書き込み処理
	log.Println("Write start")
	if err := uartdata.Write(INIT_DATA); err != nil {
		t.Fatal(err)
	}
	time.Sleep(100 * time.Millisecond)
	if err := uartdata.Write(READ_DATA); err != nil {
		t.Fatal(err)
	}
	log.Println("Write end")
	time.Sleep(time.Millisecond * 100)
	shutdown <- true
	wg.Wait()
	log.Println("main stop")
}
