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

// このテストはTxとRxを直接接続した状態で行う
func TestUart(t *testing.T) {
	if err := uartInit("/dev/ttyS0", 9600); err != nil {
		t.Fatal(err)
	}
	if err := uartdata.open(); err != nil {
		t.Fatal(err)
	}
	defer uartdata.close()
	//スレッドによる読み込み処理
	go func() {
		for {
			buf, err := uartdata.read()
			log.Println("info", buf, err)
			time.Sleep(time.Second / 10)
		}
	}()
	//定期的に書き込み処理
	time.Sleep(time.Second)
}
