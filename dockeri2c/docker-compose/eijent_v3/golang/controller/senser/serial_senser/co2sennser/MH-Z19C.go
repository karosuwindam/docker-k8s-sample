package co2sennser

import (
	msgsenser "eijent/controller/senser/msg_senser"
	"errors"
	"log"
	"sync"
	"time"
)

func Init(i2cMu *sync.Mutex) error {
	memory = datastore{
		Flag:     false,
		StopFlag: false,
		msg:      msgsenser.Msg{},
		i2cMu:    i2cMu,
	}
	shudown = make(chan bool, 1)
	done = make(chan bool, 1)
	reset = make(chan bool, 1)
	wait = make(chan bool, 1)
	memory.msg.Create("MH-Z19C")
	for i := 0; i < 3; i++ {
		if Test() {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}
	if !memory.Flag {
		return errors.New("not Init Error for AM2320")
	}

	return nil
}

func Run() error {
	return nil
}

func Stop() error {
	shudown <- true
	memory.chageRunFlag(false)
	select {
	case <-done:
		break
	case <-time.After(1 * time.Second):
		msg := "shutdown time out"
		memory.changeMsg(msg)
		return errors.New(msg)
	}
	memory.changeMsg("shutdown")
	return nil
}

func Health() (bool, msgsenser.Msg) {
	return memory.readFlag(), memory.readMsg()
}

func Wait() {

	wait <- true
	select {
	case <-done:
		break
	case <-time.After(1 * time.Second):
		log.Println("error:", "time over 1 sec")
	}
}

func ReadValue() (interface{}, bool) {
	return nil, true
}

func ResetMessage() {
	if len(reset) > 0 {
		return
	}
	reset <- true
}

func Test() bool {
	flag := false
	return flag
}
