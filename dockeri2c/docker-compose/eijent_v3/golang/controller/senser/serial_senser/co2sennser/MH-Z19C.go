package co2sennser

import (
	"eijent/config"
	msgsenser "eijent/controller/senser/msg_senser"
	"eijent/controller/uptime"
	"errors"
	"io"
	"log"
	"time"
)

var (
	INIT_DATA = []byte{0xff, 0x87, 0x87, 0x00, 0x00, 0x00, 0x00, 0x00, 0xf2}
	READ_DATA = []byte{0xff, 0x01, 0x86, 0x00, 0x00, 0x00, 0x00, 0x00, 0x79}
)

const (
	GROVENAME  string  = "Grove - CO2 Sensor"
	MHZ19CNAME string  = "MH-Z19C"
	CO2_MAX            = 5000
	CO2_MIN            = 400
	WAITTIME   float64 = 7 * 60 //wait time 7 min(420sec)
)
const (
	INIT              = 0
	READ              = 1
	BAUDRATE   int    = 9600
	UART_DEV   string = "/dev/ttyAMA0"
	CO2SLEEP          = 10                    //10us	Time out count interval
	CO2TIMEOUT        = 500 * 1000 / CO2SLEEP //500ms Time Out
)

type MhZ19c_Vaule struct {
	Co2  int
	Temp int
}

func Init() error {
	memory = datastore{
		Flag:     false,
		StopFlag: false,
		msg:      msgsenser.Msg{},
	}
	shudown = make(chan bool, 1)
	done = make(chan bool, 1)
	reset = make(chan bool, 1)
	wait = make(chan bool, 1)
	memory.msg.Create("MH-Z19C")
	//起動してから7経過していることを
	waittime := int((WAITTIME - uptime.Read()) * 1000)
	if waittime > 0 {
		log.Println("info:", "Serial sleep", waittime, "ms")
		time.Sleep(time.Millisecond * time.Duration(waittime))
	}
	if err := uartInit(UART_DEV, BAUDRATE); err != nil {
		return err
	}
	for i := 0; i < 3; i++ {
		if Test() {
			memory.changeMsg("OK")
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
	log.Println("info:", MHZ19CNAME+" loop start")
	var readone chan bool = make(chan bool, 1)
	var stopread chan bool = make(chan bool, 1)
	var readCash chan []byte = make(chan []byte, 10)

	if memory.readFlag() {
		if err := uartdata.open(); err != nil {
			memory.changeFlag(false)
			memory.changeMsg(err.Error())
		} else {
			readone <- true

		}
	}
	go func() { //永続読み取り
		for {
			if memory.readFlag() {
				if buf, err := uartdata.read(); err == nil && len(readCash) >= 10 {
					readCash <- buf
				} else if err != io.EOF {
					log.Println("error", err)
					memory.changeMsg(err.Error())
				}
			}
			select {
			case <-stopread:
				return
			case <-time.After(time.Millisecond * 10):
			}
		}
	}()
loop:
	for {
		select {
		case <-reset:
			memory.changeFlag(false)
			uartdata.close()
			if err := uartdata.open(); err != nil {
				memory.changeMsg(err.Error())
				log.Println("error:", err)
			} else {
				for i := 0; i < 3; i++ {
					if Test() {
						memory.changeMsg("OK")
						break
					}
					time.Sleep(100 * time.Millisecond)
				}
			}
		case data := <-readCash:
			if checkCO2Data(data) {
				readdate(data)
			}
		case <-shudown:
			stopread <- true
			done <- true
			uartdata.close()
			break loop
		case <-wait:
			done <- true
		case <-readone:
			if memory.readFlag() {
				if err := uartdata.Write(READ_DATA); err != nil {
					memory.changeMsg(err.Error())
				}
			}
		case <-time.After(time.Duration(config.Senser.CO2_SENSER_Count) * time.Millisecond):
			if memory.readFlag() {
				if err := uartdata.Write(READ_DATA); err != nil {
					memory.changeMsg(err.Error())
				}
			}

		}
	}
	memory.changeFlag(false)
	log.Println("info:", MHZ19CNAME+" loop stop")
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

func ReadValue() (MhZ19c_Vaule, bool) {
	tmp, ok := memory.readValue().(MhZ19c_Vaule)
	if !ok {
		tmp = MhZ19c_Vaule{
			Co2:  -1,
			Temp: -1,
		}
	}
	return tmp, memory.readFlag()
}

func readdate(b []byte) {
	mh := changeData(b)
	if mh.Co2 == -1 && mh.Temp == -1 {
		return
	}
	memory.changeValue(mh)
}

func ResetMessage() {
	if len(reset) > 0 {
		return
	}
	reset <- true
}

func Test() bool {
	flag := false
	if err := uartdata.open(); err != nil {
		log.Println("error:", err)
		return flag
	}
	defer uartdata.close()
	var readdata chan []byte = make(chan []byte, 3)
	var goFuncStop chan bool = make(chan bool, 1)
	var goFuncStopDone chan bool = make(chan bool, 1)
	go func() {
		for {
			buf, err := uartdata.read()
			if err == nil && len(readdata) != 3 {
				readdata <- buf
			} else if err != io.EOF {
				log.Println("error:", err)
				memory.changeMsg(err.Error())
			}
			select {
			case <-goFuncStop:
				goFuncStopDone <- true
				return
			case <-time.After(time.Millisecond * 10):
			}
		}
	}()
	if err := uartdata.Write(READ_DATA); err != nil {
		log.Println("error:", err)
		memory.changeMsg(err.Error())
		return flag
	}
	flag = checkCO2Data(<-readdata)
	if !flag {
		for i := 0; i <= len(readdata); i++ {
			if checkCO2Data(<-readdata) {
				flag = true
				break
			}
		}
	}
	goFuncStop <- true
	<-goFuncStopDone
	if flag {
		memory.changeFlag(flag)
		memory.changeMsg("OK")
	}
	return flag
}

func checkCO2Data(tmp []byte) bool {

	var num byte
	var num_s byte
	var num_e byte
	if len(tmp) == 0 {
		return false
	}
	if tmp[0] != 0xff {
		return false
	}
	i := 0
	for _, v := range tmp {
		if i == 0 {
			num_s = v
		} else if i == 8 {
			num_e = v
		} else {
			num += v
		}
		i++
	}
	return ((num_s ^ num) + 1) == num_e
}

func changeData(data []byte) MhZ19c_Vaule {
	var output MhZ19c_Vaule
	co2ppm := int(data[2])*256 + int(data[3])
	temp := int(data[4]) - 40
	if (co2ppm >= CO2_MIN) && (co2ppm <= CO2_MAX) {
		output.Co2 = co2ppm
		output.Temp = temp
	} else {
		output = MhZ19c_Vaule{-1, -1}
	}
	return output

}
