package bme280

import (
	"eijent/config"
	msgsenser "eijent/controller/senser/msg_senser"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"
)

const (
	SENSER_NAME string = "BME280"
)

var (
	BME280 = uint8(0x76)
	// BME280 = uint8(0x77)
)

const (
	BME280_ID         byte = 0xD0
	BME280_CTRL_HUM   byte = 0xF2
	BME280_STATUS     byte = 0xF3
	BME280_CTRL_MEAS  byte = 0xF4
	BME280_CONFIG     byte = 0xF5
	BME280_PRESS_MSB  byte = 0xF7
	BME280_PRESS_LSB  byte = 0xF8
	BME280_PRESS_XLSB byte = 0xF9
	BME280_TEMP_MSB   byte = 0xFA
	BME280_TEMP_LSB   byte = 0xFB
	BME280_TEMP_XLSB  byte = 0xFC
	BME280_HUM_MSB    byte = 0xFD
	BME280_HUM_LSB    byte = 0xFE

	BME280_CALIB00 byte = 0x88
	BME280_CALIB25 byte = 0xA1
	BME280_CALIB26 byte = 0xE1
	BME280_CALIB41 byte = 0xF0
)

const (
	BME280_ID_DATA byte = 0x60
)

type Bme280_Vaule struct {
	Press float64
	Temp  float64
	Hum   float64
}

type bme280_cal struct {
	press    []int
	temp     []int
	hum      []int
	timefine float64
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
	i2c_addr = BME280
	up()
	for i := 0; i < 3; i++ {
		if Test() {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}
	down()
	if !memory.Flag {
		return errors.New("not Init Error for AM2320")
	}

	return nil
}

func Run() error {
	memory.chageRunFlag(true)
	log.Println("info:", SENSER_NAME+" loop start")
	var readone chan bool = make(chan bool, 1)
	var calib bme280_cal
	if memory.readFlag() {
		readone <- true
		up()
	}
loop:
	for {
		select {
		case <-reset:
			down()
			up()
			for i := 0; i < 3; i++ {
				if Test() {
					calib = calibRead()
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
			calib = calibRead()
			if memory.readFlag() {
				readdate(calib)
			}
		case <-time.After(time.Duration(config.Senser.BME280_Count) * time.Millisecond):
			if memory.readFlag() {
				readdate(calib)
			}
		}
	}
	down()
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

func readdate(calib bme280_cal) {
	press, tmp, hum := readSenserData(calib)
	if press == -1 && tmp == -1 && hum == -1 {
		return
	}
	memory.changeValue(Bme280_Vaule{press, tmp, hum})
}

func ReadValue() (Bme280_Vaule, bool) {
	num, ok := memory.readValue().(Bme280_Vaule)
	if !ok {
		return Bme280_Vaule{-1, -1, -1}, false
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
	msg := ""
	if buf, err := readByte(BME280_ID, 1); err != nil {
		msg = fmt.Sprintf("%v Test Read Error Addr %x", SENSER_NAME, i2c_addr)
	} else if buf[0] != BME280_ID_DATA {
		msg = fmt.Sprintf("%v Test test header data %x !=%x", SENSER_NAME, BME280_ID_DATA, buf[0])
	} else {
		msg = "OK"
		flag = true
	}
	memory.changeFlag(flag)
	memory.changeMsg(msg)
	return flag
}

func up() {
	osrs_t := 1
	osrs_p := 1
	osrs_h := 1
	mode := 3
	t_sb := 5
	filter := 0
	spi3w_en := 0
	ctrl_meas_reg := byte((osrs_t << 5) | (osrs_p << 2) | mode)
	config_reg := byte((t_sb << 5) | (filter << 2) | spi3w_en)
	ctrl_hum_reg := byte(osrs_h)

	writeByte(BME280_CTRL_HUM, ctrl_hum_reg)
	writeByte(BME280_CTRL_MEAS, ctrl_meas_reg)
	writeByte(BME280_CONFIG, config_reg)

}
func down() {
	writeByte(BME280_CTRL_HUM, 0x0)
	writeByte(BME280_CTRL_MEAS, 0x0)
	writeByte(BME280_CONFIG, 0x0)
}

// rawRead()
//
// press, temp, hum
func rawRead() (int, int, int) {
	buf, err := readByte(BME280_PRESS_MSB, 8)
	if err != nil {
		log.Println("error:", err)
		memory.changeMsg(err.Error())
	}
	memory.changeMsg("OK")
	tmp := 0
	for _, bt := range buf {
		tmp += int(bt)
	}
	if tmp == 0 {
		return -1, -1, -1
	}
	press := int(buf[0])<<12 | int(buf[1])<<4 | int(buf[2])>>4
	temp := int(buf[3])<<12 | int(buf[4])<<4 | int(buf[5])>>4
	hum := int(buf[6])<<8 | int(buf[7])
	return press, temp, hum
}

func calibRead() bme280_cal {
	var out bme280_cal

	var calib []int
	buf, _ := readByte(BME280_CALIB00, int(BME280_CALIB25-BME280_CALIB00))
	for _, b := range buf {
		calib = append(calib, int(b))
	}
	buf, _ = readByte(BME280_CALIB26, int(BME280_CALIB41-BME280_CALIB26+1))
	for _, b := range buf {
		calib = append(calib, int(b))
	}
	for i := 0; i < 3; i++ { //0-5
		num := (calib[1+i*2] << 8) | calib[0+i*2]
		if i != 0 {
			if (num & 0x8000) != 0 {
				num = (-num ^ 0xffff) + 1
			}
		}
		out.temp = append(out.temp, num)
	}
	for i := 0; i < 9; i++ { //6-23
		num := (calib[7+i*2] << 8) | calib[6+i*2]
		if i != 0 {
			if (num & 0x8000) != 0 {
				num = (-num ^ 0xffff) + 1
			}
		}
		out.press = append(out.press, num)
	}
	//24-31
	out.hum = append(out.hum, calib[24])
	out.hum = append(out.hum, (calib[26]<<8)|calib[25])
	out.hum = append(out.hum, calib[27])
	out.hum = append(out.hum, (calib[28]<<4)|(0x0F&calib[29]))
	out.hum = append(out.hum, (calib[30]<<4)|((calib[29]>>4)&0x0F))
	out.hum = append(out.hum, calib[31])
	for i := 0; i < len(out.hum); i++ {
		if (i != 0) && (i != 2) {
			if (out.hum[i] & 0x8000) != 0 {
				out.hum[i] = (-out.hum[i] ^ 0xffff) + 1
			}
		}
	}
	return out
}

func (t *bme280_cal) CalibTemp(b_temp int) (float64, float64) {
	tmp := float64(b_temp)
	var calib []float64
	for _, flt := range t.temp {
		calib = append(calib, float64(flt))
	}
	v1 := (tmp/16384.0 - calib[0]/1024.0) * calib[1]
	v2 := (tmp/131072.0 - calib[0]/8192.0) * (tmp/131072.0 - calib[0]/8192.0) * calib[2]
	t_fine := v1 + v2
	temperature := t_fine / 5120.0
	return temperature, t_fine

}

func (t *bme280_cal) CalibPress(b_press int, t_fine float64) float64 {
	tmp := float64(b_press)
	var calib []float64
	for _, flt := range t.press {
		calib = append(calib, float64(flt))
	}

	v1 := (t_fine / 2.0) - 64000.0
	v2 := (((v1 / 4.0) * (v1 / 4.0)) / 2048) * calib[5]
	v2 = v2 + ((v1 * calib[4]) * 2.0)
	v2 = (v2 / 4.0) + (calib[3] * 65536.0)
	v1 = (((calib[2] * (((v1 / 4.0) * (v1 / 4.0)) / 8192)) / 8) + ((calib[1] * v1) / 2.0)) / 262144
	v1 = ((32768 + v1) * calib[0]) / 32768
	if v1 != 0 {
		pressure := ((1048576 - tmp) - (v2 / 4096)) * 3125
		if pressure < 0 {
			pressure = (pressure * 2.0) / v1
		} else {
			pressure = (pressure / v1) * 2
		}
		v1 = (calib[8] * (((pressure / 8.0) * (pressure / 8.0)) / 8192.0)) / 4096
		v2 = ((pressure / 4.0) * calib[7]) / 8192.0
		return pressure / 100
	}
	return 0
}

func (t *bme280_cal) CalibHum(b_hum int, t_fine float64) float64 {
	tmp := float64(b_hum)
	var calib []float64
	for _, flt := range t.hum {
		calib = append(calib, float64(flt))
	}
	var_h := t_fine - 76800.0
	if var_h != 0 {
		var_h = (tmp - (calib[3]*64.0 + calib[4]/16384.0*var_h)) * (calib[1] / 65536.0 * (1.0 + calib[5]/67108864.0*var_h*(1.0+calib[2]/67108864.0*var_h)))
	} else {
		return 0
	}
	var_h = var_h * (1.0 - calib[0]*var_h/524288.0)
	if var_h > 100 {
		var_h = 100
	} else if var_h < 0 {
		var_h = 0
	}
	return var_h
}

// readSenserData(calib bme280_cal) (float64, float64, float64)
//
// c_press, c_tmp, c_hum
func readSenserData(calib bme280_cal) (float64, float64, float64) {
	press, temp, hum := rawRead()
	c_tmp, t_fine := calib.CalibTemp(temp)
	c_press := calib.CalibPress(press, t_fine)
	c_hum := calib.CalibHum(hum, t_fine)
	if c_hum <= 0 || c_tmp < -40 || c_tmp > 85 {
		return -1, -1, -1
	}
	return c_press, c_tmp, c_hum

}
