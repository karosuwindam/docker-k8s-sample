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
	InitData []byte
	ReadData []byte
	port     *serial.Port
}
type MhZ19c struct {
	Flag     bool
	Name     string
	ReadData []byte
	port     *serial.Port
}

var (
	INIT_DATA = []byte{0xff, 0x87, 0x87, 0x00, 0x00, 0x00, 0x00, 0x00, 0xf2}
	READ_DATA = []byte{0xff, 0x01, 0x86, 0x00, 0x00, 0x00, 0x00, 0x00, 0x79}
)

const (
	GROVENAME  string = "Grove - CO2 Sensor"
	MHZ19CNAME string = "MH-Z19C"
)
const (
	INIT       = 0
	READ       = 1
	BAUDRATE   = 9600
	CO2SLEEP   = 10                    //10us	Time out count interval
	CO2TIMEOUT = 500 * 1000 / CO2SLEEP //500ms Time Out
)

func (t *MhZ19c) Init(name string) bool {
	var err error
	t.Name = MHZ19CNAME
	t.Flag = false
	c := &serial.Config{Name: name, Baud: BAUDRATE}
	t.port, err = serial.OpenPort(c)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	for i := 0; i < 4; i++ {
		if t.Read() {
			t.Flag = true
			break
		}
		fmt.Println("count:", i+1)
		t.port.Close()
		time.Sleep(500 * time.Microsecond)
		t.port, _ = serial.OpenPort(c)
	}
	t.port.Close()
	return t.Flag
}

func (t *MhZ19c) Read() bool {
	s := t.port
	flag := false
	tmp := []byte{}
	go func() {
		n, err := s.Write(READ_DATA)
		log.Printf("WriteData %v:%q", len(READ_DATA), READ_DATA)
		if err != nil {
			log.Printf(err.Error())
		}
		for {
			buf := make([]byte, 128)
			n, err = s.Read(buf)
			if err != nil {
				log.Printf(err.Error())
				break
			}
			if n > 0 {
				for _, v := range buf[:n] {
					tmp = append(tmp, v)
				}
			}
			if len(tmp) > 8 {
				break
			}
		}
	}()
	for i := 0; i < 100; i++ {
		if len(tmp) > 8 {
			flag = t.checkdata(tmp)
			if flag {
				t.ReadData = tmp
			}
			break
		}
		time.Sleep(time.Millisecond * 10)
	}
	log.Printf("ReadData %v:%q,flag:%v", len(tmp), tmp, flag)
	return flag

}

func (t *MhZ19c) checkdata(tmp []byte) bool {
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
	// fmt.Println(((num_s ^ num) + 1), num_e)
	return ((num_s ^ num) + 1) == num_e
}

func (t *MhZ19c) Output() (int, int) {
	data := t.ReadData
	co2ppm := int(data[2])*256 + int(data[3])
	temp := int(data[4]) - 40
	return co2ppm, temp
}

func (t *MhZ19c) Close() {
	t.port.Close()
	t.Flag = false
	t.ReadData = []byte{}
}

func (t *Co2Sennser) Init(name string) bool {
	var err error
	t.Name = GROVENAME

	t.Flag = false
	c := &serial.Config{Name: name, Baud: BAUDRATE}
	t.port, err = serial.OpenPort(c)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	t.InitData, err = t.write(INIT_DATA)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	if !t.checkdata(INIT) {
		fmt.Println("checkdata is Error")
		return false
	}
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
		return -1, -1
	}
	if !t.checkdata(READ) {
		fmt.Println("checkdata is Error")
		return -1, -1
	}
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
	s := t.port
	_, err := s.Write(data)
	if err != nil {
		return []byte{}, err
	}
	output := []byte{}
	m := 0
	go func() {
		for {
			buf := make([]byte, 32)
			n, err := s.Read(buf)
			if err != nil {
				log.Fatal(err)
				break
			}
			for _, v := range buf[:n] {
				output = append(output, v)
			}
			m += n
			if m > 8 {
				break
			}
		}
	}()
	i := 0
	for {
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
