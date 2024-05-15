package mma8452q

import (
	"eijent/config"
	msgsenser "eijent/controller/senser/msg_senser"
	"sync"
	"time"
)

type datastore struct {
	values   interface{}
	readtime time.Time
	Flag     bool
	StopFlag bool
	msg      msgsenser.Msg

	mu    sync.Mutex
	i2cMu *sync.Mutex
}

var memory datastore

var shudown chan bool
var reset chan bool
var done chan bool
var wait chan bool
var busy chan bool

func (v *datastore) changeFlag(flag bool) {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.Flag = flag

}

func (v *datastore) chageRunFlag(flag bool) {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.StopFlag = flag
}

func (v *datastore) changeValue(num interface{}) {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.readtime = time.Now()
	v.values = num
}

func (v *datastore) changeMsg(str string) {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.msg.Write(str)
}

func (v *datastore) readFlag() bool {
	v.mu.Lock()
	defer v.mu.Unlock()
	return v.Flag
}

func (v *datastore) readRunFlag() bool {
	v.mu.Lock()
	defer v.mu.Unlock()
	return v.StopFlag
}

func (v *datastore) readMsg() msgsenser.Msg {
	v.mu.Lock()
	defer v.mu.Unlock()
	return v.msg
}

func (v *datastore) readValue() interface{} {
	v.mu.Lock()
	defer v.mu.Unlock()
	return v.values
}

func (v *datastore) chackEnable() bool {
	v.mu.Lock()
	defer v.mu.Unlock()
	if !v.Flag {
		return false
	}
	return time.Now().Sub(v.readtime) < time.Duration(config.Senser.HorldTime)*time.Minute
}
