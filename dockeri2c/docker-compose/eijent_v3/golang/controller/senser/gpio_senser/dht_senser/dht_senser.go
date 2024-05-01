package dhtsenser

import (
	"eijent/config"
	msgsenser "eijent/controller/senser/msg_senser"
	"errors"
	"log"
	"sync"
	"time"

	"github.com/d2r2/go-dht"
)

type datastore struct {
	values   []dhtsenser
	typedata dhtsenser_type
	Flag     bool
	StopFlag bool
	msg      msgsenser.Msg

	mu sync.Mutex
}

type dhtsenser_type struct {
	sensertype dht.SensorType
	pin        int
}

type dhtsenser struct {
	Hum  float64
	Temp float64
}

var memory datastore

var shudown chan bool
var reset chan bool
var done chan bool
var wait chan bool

func Init() error {
	memory = datastore{
		values:   []dhtsenser{},
		Flag:     false,
		StopFlag: false,
		msg:      msgsenser.Msg{},
	}
	shudown = make(chan bool, 1)
	done = make(chan bool, 1)
	reset = make(chan bool, 1)
	wait = make(chan bool, 1)
	memory.msg.Create(config.Senser.DHT_senser_type)
	if err := memory.basecreate(); err != nil {
		return err
	}
	for i := 0; i < 3; i++ {
		if Test() {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}
	if !memory.Flag {
		return errors.New("not Init Error for TSL2561")
	}
	return nil
}

func Run() error {
	memory.chageRunFlag(true)
	log.Println("info:", "DHT loop start")
	var readone chan bool = make(chan bool, 1)
	readone <- true
loop:
	for {
		select {
		case <-reset:
			for i := 0; i < 3; i++ {
				if Test() {
					break
				}
				time.Sleep(100 * time.Millisecond)
			}
		case <-shudown:
			done <- true
			break loop
		case <-wait:
			done <- true
		case <-readone:
			if memory.readFlag() {
				readdate()
			}
		case <-time.After(time.Duration(config.Senser.DHT_senser_Count) * time.Millisecond):
			if memory.readFlag() {
				readdate()
			}
		}
	}
	memory.changeFlag(false)
	log.Println("info:", "DHT loop stop")

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

func Wait() {

	wait <- true
	select {
	case <-done:
		break
	case <-time.After(1 * time.Second):
		log.Println("error:", "time over 1 sec")
	}
}

func readdate() {
	flag := true
	for i := 0; i < 3; i++ {
		hum, temp := readSenser()
		if hum == -1 && temp == -1 {
			continue
		}
		memory.changeValue(dhtsenser{
			Hum:  hum,
			Temp: temp,
		})
		flag = false
		break
	}
	if flag {
		memory.changeFlag(false)
		memory.changeMsg("Stop DHT")

	}
}
func ResetMessage() {
	if len(reset) > 0 {
		return
	}
	reset <- true
}

func Health() (bool, msgsenser.Msg) {
	return memory.readFlag(), memory.readMsg()

}

func (v *datastore) basecreate() error {

	v.mu.Lock()
	defer v.mu.Unlock()
	if v.msg.Senser == "" {
		return errors.New("Senser Name enptiy")
	}
	v.typedata.pin = config.Senser.DHT_senser_pin
	switch v.msg.Senser {
	case "DHT11":
		v.typedata.sensertype = dht.DHT11
		break
	case "DHT12":
		v.typedata.sensertype = dht.DHT12
		break
	case "AM2302":
		v.typedata.sensertype = dht.AM2302
		break
	}
	return nil

}

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

func (v *datastore) changeValue(num dhtsenser) {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.values = append(v.values, num)
	if len(v.values) > 200 {
		v.values = v.values[1:]
	}
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

func (v *datastore) readSenserType() dhtsenser_type {
	v.mu.Lock()
	defer v.mu.Unlock()
	return v.typedata
}

func (v *datastore) readValue() []dhtsenser {
	v.mu.Lock()
	defer v.mu.Unlock()
	return v.values
}

// 記録したデータの平均を取得
func ReadValue() (dhtsenser, bool) {
	var out dhtsenser
	tmps := memory.readValue()
	hum := 0.0
	temp := 0.0
	for i, tmp := range tmps {
		hum = hum*float64(i) + tmp.Hum
		temp = temp*float64(i) + tmp.Temp
		hum = hum / float64(i+1)
		temp = temp / float64(i+1)
	}
	out = dhtsenser{
		Hum:  hum,
		Temp: temp,
	}
	return out, memory.readFlag()
}

func Test() bool {
	flag := true
	tmp := memory.readSenserType()
	_, _, _, err :=
		dht.ReadDHTxxWithRetry(tmp.sensertype, tmp.pin, false, 2)
	if err != nil {
		log.Println(err.Error())
		memory.changeMsg(err.Error())
		flag = false
	}
	memory.changeFlag(flag)
	return flag
}

func readSenser() (float64, float64) {

	tmp := memory.readSenserType()
	temperature, humidity, _, err :=
		dht.ReadDHTxxWithRetry(tmp.sensertype, tmp.pin, false, 2)
	if err != nil {
		log.Println(err.Error())
		memory.changeMsg(err.Error())
		return -1, -1
	}
	hum := float64(humidity)
	temp := float64(temperature)
	if hum < 5 || hum > 95 {
		return -1, -1
	}
	if temp < -20 {
		return -1, -1
	}
	memory.changeMsg("OK")
	return hum, temp
}
