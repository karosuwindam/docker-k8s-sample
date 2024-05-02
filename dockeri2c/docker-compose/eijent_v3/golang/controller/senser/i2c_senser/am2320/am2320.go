package am2320

import (
	"eijent/config"
	msgsenser "eijent/controller/senser/msg_senser"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/davecheney/i2c"
	"github.com/pkg/errors"
)

const (
	AM2320      uint8  = 0x5c
	SENSER_NAME string = "AM2320"
	I2C_BUS            = 1
)

type Am2320_Vaule struct {
	Hum  float64
	Temp float64
}

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
	memory.msg.Create(SENSER_NAME)
	i2c_addr = AM2320
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
	memory.chageRunFlag(true)
	log.Println("info:", SENSER_NAME+" loop start")
	var readone chan bool = make(chan bool, 1)
	if memory.readFlag() {
		readone <- true
	}
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
		case <-time.After(time.Duration(config.Senser.Am2320_Count) * time.Millisecond):
			if memory.readFlag() {
				readdate()
			}
		}
	}
	memory.changeFlag(false)
	log.Println("info:", SENSER_NAME+" loop stop")
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
func readdate() {
	hum, temp := readSenserData()
	if hum == -1 && temp == -1 {
		return
	}
	memory.changeValue(Am2320_Vaule{hum, temp})
}
func ReadValue() (Am2320_Vaule, bool) {
	num, ok := memory.readValue().(Am2320_Vaule)
	if !ok {
		return Am2320_Vaule{-1, -1}, false
	}
	return num, memory.readFlag()
}

func ResetMessage() {
	if len(reset) > 0 {
		return
	}
	reset <- true
}

func Test() bool {
	flag := false
	if b, err := readbyte(); len(b) == 0 {
		log.Println("error:", err)
		memory.changeMsg("Read Error I2C")
		memory.changeFlag(false)
		return flag
	} else if err != nil {
		log.Println("error:", err)
		msg := fmt.Sprintf("%v Test Read Error Addr %x", SENSER_NAME, i2c_addr)
		memory.changeMsg(msg)
		memory.changeFlag(false)
		return flag

	}
	flag = true
	memory.changeMsg("OK")
	memory.changeFlag(flag)
	return flag
}

func readSenserData() (float64, float64) {
	var hum float64
	var temp float64
	b, err := readbyte()
	if err != nil {
		memory.changeMsg(err.Error())
		if !memory.chackEnable() {
			memory.changeFlag(false)
			memory.changeMsg("Stop " + SENSER_NAME)
		}
		return -1, -1
	}
	hum = float64((uint(b[2])<<8)|uint(b[3])) / 10  //湿度
	temp = float64((uint(b[4])<<8)|uint(b[5])) / 10 //温度
	if hum == 0 || temp < -40 || temp > 80 {
		msg := fmt.Sprintf("Error Value Hum=%v,Tmp=%v", hum, temp)
		memory.changeMsg(msg)
		return -1, -1
	}
	memory.changeMsg("OK")
	return hum, temp
}

var i2c_addr uint8

func readbyte() ([]byte, error) {
	buf := []byte{}

	i2c, err := i2c.New(i2c_addr, I2C_BUS)
	if err != nil {
		return buf, errors.Wrapf(err, "i2c.New(%v,%v)", i2c_addr, I2C_BUS)
	}
	defer i2c.Close()
	cmd := []byte{0x0}
	buf = make([]byte, 6)
	_, err = i2c.Write(cmd)
	if err != nil {
		return buf, errors.Wrapf(err, "i2c.Write(%v)", cmd)
	}
	time.Sleep(time.Microsecond * 10)
	cmd = []byte{0x3, 0x0, 0x4}
	_, err = i2c.Write(cmd)
	if err != nil {
		return buf, errors.Wrapf(err, "i2c.Write(%v)", cmd)
	}
	buf = make([]byte, 6)
	_, err = i2c.Read(buf)
	if err != nil {
		return buf, errors.Wrap(err, "i2c.Read()")
	}

	return buf, nil
}
