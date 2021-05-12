package main

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/tarm/serial"
)

type Co2Sennser struct {
	Flag     bool
	Name     string
	Message  string
	InitData []byte
	ReadData []byte
	port     *serial.Port
}

var (
	INIT_DATA = []byte{0xff, 0x87, 0x87, 0x00, 0x00, 0x00, 0x00, 0x00, 0xf2}
	READ_DATA = []byte{0xff, 0x01, 0x86, 0x00, 0x00, 0x00, 0x00, 0x00, 0x79}
)

const (
	INIT       = 0
	READ       = 1
	BAUDRATE   = 9600
	CO2SLEEP   = 10                    //10us	Time out count interval
	CO2TIMEOUT = 500 * 1000 / CO2SLEEP //500ms Time Out
)

func (t *Co2Sennser) Init(name string) bool {
	var err error
	t.Name = "Grove - CO2 Sensor"

	t.Flag = false
	c := &serial.Config{Name: name, Baud: BAUDRATE}
	t.port, err = serial.OpenPort(c)
	if err != nil {
		fmt.Println(err.Error())
		t.Message = err.Error()
		return false
	}
	t.InitData, err = t.write(INIT_DATA)
	if err != nil {
		fmt.Println(err.Error())
		t.Message = err.Error()
		return false
	}
	if !t.checkdata(INIT) {
		fmt.Println("checkdata is Error")
		t.Message = err.Error()
		return false
	}
	t.Message = "OK"
	t.Flag = true
	return true
}
func (t *Co2Sennser) Close() {
	t.port.Close()
	t.Flag = false
	t.InitData = []byte{}
	t.ReadData = []byte{}
}
func (t *Co2Sennser) Read() (int, int) {
	var err error
	if !t.Flag {
		return -1, -1
	}
	t.ReadData, err = t.write(READ_DATA)
	if err != nil {
		fmt.Println(err.Error())
		t.Message = err.Error()
		return -1, -1
	}
	if !t.checkdata(READ) {
		fmt.Println("checkdata is Error")
		t.Message = "checkdata is Error"
		return -1, -1
	}
	t.Message = "OK"
	return t.output()
}

func (t *Co2Sennser) checkdata(sel int) bool {
	tmp := t.ReadData
	if sel == INIT {
		tmp = t.InitData
	}
	var num byte
	var num_s byte
	var num_e byte
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
	// fmt.Println(((num_s ^ num) + 1), num_e)
	return ((num_s ^ num) + 1) == num_e
}

func (t *Co2Sennser) output() (int, int) {
	data := t.ReadData
	co2ppm := int(data[2])*256 + int(data[3])
	temp := int(data[4]) - 40
	return co2ppm, temp
}

func (t *Co2Sennser) write(data []byte) ([]byte, error) {
	var werr error
	s := t.port
	output := []byte{}
	m := 0
	go func() {
		_, werr = s.Write(data)
		if werr != nil {
			log.Fatal(werr)
			return
		}
		for {
			buf := make([]byte, 32)
			n, err := s.Read(buf)
			if err != nil {
				log.Fatal(err)
				break
			}
			for _, v := range buf[:n] {
				output = append(output, v)
				if output[0] != 0xff{
					output = []byte{}
				}
			}
			m += n
			if m > 8 {
				break
			}
		}
	}()
	i := 0
	for {
		if werr != nil{
			return output, werr
		}
		if i > CO2TIMEOUT {
			return output, errors.New("Srial Read Time Out")
		}
		if len(output) > 8 {
			break
		}
		time.Sleep(CO2SLEEP * time.Microsecond) //10us
		i++
	}
	return output, nil
}
