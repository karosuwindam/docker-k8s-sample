package main

import (
	"fmt"
	"time"

	"github.com/davecheney/i2c"
)

type Ame280 struct {
	Flag    bool
	Name    string
	Message string
	calib   ame280_cal
}

type ame280_cal struct {
	press    []int
	temp     []int
	hum      []int
	timefine float64
}

var (
	AME280  = uint8(0x77)
	I2C_BUS = 1
)

const (
	AME280_ID         byte = 0xD0
	AME280_CTRL_HUM   byte = 0xF2
	AME280_STATUS     byte = 0xF3
	AME280_CTRL_MEAS  byte = 0xF4
	AME280_CONFIG     byte = 0xF5
	AME280_PRESS_MSB  byte = 0xF7
	AME280_PRESS_LSB  byte = 0xF8
	AME280_PRESS_XLSB byte = 0xF9
	AME280_TEMP_MSB   byte = 0xFA
	AME280_TEMP_LSB   byte = 0xFB
	AME280_TEMP_XLSB  byte = 0xFC
	AME280_HUM_MSB    byte = 0xFD
	AME280_HUM_LSB    byte = 0xFE

	AME280_CALIB00 byte = 0x88
	AME280_CALIB25 byte = 0xA1
	AME280_CALIB26 byte = 0xE1
	AME280_CALIB41 byte = 0xF0
)

const (
	AME280_ID_DATA byte = 0x60
)

func (t *Ame280) Init() bool {
	t.Name = "AME280"
	t.up()
	if !t.Test() {
		t.Close()
	} else {
		t.ReadCalib()
	}
	return t.Flag
}

func (t *Ame280) Close() {
	t.down()
}
func (t *Ame280) up() {
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

	t.WriteByte(AME280_CTRL_HUM, ctrl_hum_reg)
	t.WriteByte(AME280_CTRL_MEAS, ctrl_meas_reg)
	t.WriteByte(AME280_CONFIG, config_reg)
}

func (t *Ame280) down() {
	t.WriteByte(AME280_CTRL_HUM, 0x0)
	t.WriteByte(AME280_CTRL_MEAS, 0x0)
	t.WriteByte(AME280_CONFIG, 0x0)
}

func (t *Ame280) WriteByte(command, data byte) {
	i2c, err := i2c.New(AME280, I2C_BUS)
	if err != nil {
		t.Message = err.Error()
		return
	}
	defer i2c.Close()
	_, err = i2c.Write([]byte{command, data})
	if err != nil {
		t.Message = err.Error()
		return
	}
	t.Message = "OK"
}

func (t *Ame280) ReadCalib() {
	var calib []int
	buf := t.ReadByte(AME280_CALIB00, int(AME280_CALIB25-AME280_CALIB00+1))
	for _, b := range buf {
		calib = append(calib, int(b))
	}
	buf = t.ReadByte(AME280_CALIB26, int(AME280_CALIB41-AME280_CALIB26+1))
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
		t.calib.temp = append(t.calib.temp, num)
	}
	for i := 0; i < 9; i++ { //6-23
		num := (calib[7+i*2] << 8) | calib[6+i*2]
		if i != 0 {
			if (num & 0x8000) != 0 {
				num = (-num ^ 0xffff) + 1
			}
		}
		t.calib.press = append(t.calib.press, num)
	}
	//24-31
	t.calib.hum = append(t.calib.hum, calib[24])
	t.calib.hum = append(t.calib.hum, (calib[26]<<8)|calib[25])
	t.calib.hum = append(t.calib.hum, calib[27])
	t.calib.hum = append(t.calib.hum, (calib[28]<<4)|(0x0F&calib[29]))
	t.calib.hum = append(t.calib.hum, (calib[30]<<4)|((calib[29]>>4)&0x0F))
	t.calib.hum = append(t.calib.hum, calib[31])
	for i := 0; i < len(t.calib.hum); i++ {
		if (i != 0) || (i != 2) {
			if (t.calib.hum[i] & 0x8000) != 0 {
				t.calib.hum[i] = (-t.calib.hum[i] ^ 0xffff) + 1
			}
		}
	}

}

func (t *Ame280) ReadByte(command byte, size int) []byte {
	i2c, err := i2c.New(AME280, I2C_BUS)
	if err != nil {
		t.Message = err.Error()
		return []byte{}
	}
	defer i2c.Close()
	i2c.Write([]byte{command})
	if err != nil {
		t.Message = err.Error()
		return []byte{}
	}
	buf := make([]byte, size)
	i2c.Read(buf)
	if err != nil {
		t.Message = err.Error()
		return []byte{}
	}
	t.Message = "OK"
	return buf
}
func (t *Ame280) Test() bool {
	t.Flag = false
	i2c, err := i2c.New(AME280, I2C_BUS)
	if err != nil {
		fmt.Println(err.Error())
		t.Message = err.Error()
		return false
	}
	i2c.Close()
	fmt.Println(t.ReadByte(AME280_CTRL_MEAS, 1))
	for i := 0; i < 3; i++ {
		buf := t.ReadByte(AME280_ID, 1)
		if AME280_ID_DATA == buf[0] {
			t.Flag = true
			return true
		}
		fmt.Println(buf)
		t.down()
		time.Sleep(time.Microsecond * 1000)
		t.up()
	}
	fmt.Println(t.Message)
	return false
}
func (t *Ame280) Read() (int, int, int) {
	buf := t.ReadByte(AME280_PRESS_MSB, 8)
	press := int(buf[0])<<12 | int(buf[1])<<4 | int(buf[2])>>4
	temp := int(buf[3])<<12 | int(buf[4])<<4 | int(buf[5])>>4
	hum := int(buf[6])<<8 | int(buf[7])
	fmt.Println(buf)
	return press, temp, hum
}

func (t *Ame280) ReadData() {

}
