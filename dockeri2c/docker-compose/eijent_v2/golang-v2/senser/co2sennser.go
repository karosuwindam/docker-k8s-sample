package senser

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
type MhZ19c struct {
	Flag     bool
	Name     string
	com      string
	Message  string
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
	CO2_MAX           = 5000
	CO2_MIN           = 400
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
	t.com = name
	c := &serial.Config{Name: t.com, Baud: BAUDRATE}
	t.port, err = serial.OpenPort(c)
	if err != nil {
		fmt.Println(err.Error())
		t.Message = err.Error()
		return false
	}
	for i := 0; i < 4; i++ {
		if t.ReadChack() {
			t.Message = "OK"
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
func (t *MhZ19c) Read() (int, int) {
	var err error
	if !t.Flag {
		return -1, -1
	}
	c := &serial.Config{Name: t.com, Baud: BAUDRATE}
	t.port, err = serial.OpenPort(c)
	if err != nil {
		fmt.Println(err.Error())
		t.Message = err.Error()
		return -1, -1
	}
	for i := 0; i < 3; i++ {
		if t.ReadChack() {
			t.Message = "OK"
			break
		} else {
			t.Message = "Read CO2 Error"
		}
		// fmt.Println("count:", i+1)
		t.port.Close()
		time.Sleep(500 * time.Microsecond)
		t.port, _ = serial.OpenPort(c)
	}
	t.port.Close()
	return t.Output()

}
func (t *MhZ19c) ReadChack() bool {
	s := t.port
	flag := false
	tmp := []byte{}
	go func() {
		n, err := s.Write(READ_DATA)
		// log.Printf("WriteData %v:%q", len(READ_DATA), READ_DATA)
		if err != nil {
			log.Printf(err.Error())
		}
		for {
			buf := make([]byte, 128)
			n, err = s.Read(buf)
			if err != nil {
				if n != 0 {
					log.Println(n, err.Error())
				}
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
	// log.Printf("ReadData %v:%q,flag:%v", len(tmp), tmp, flag)
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
	if (co2ppm >= CO2_MIN) && (co2ppm <= CO2_MAX) {
		return co2ppm, temp
	} else {
		return -1, -1
	}
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
				if output[0] != 0xff {
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
		if werr != nil {
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
