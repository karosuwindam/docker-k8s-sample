package senser

import (
	"fmt"
	"time"

	"github.com/davecheney/i2c"
)

type Bme280 struct {
	Flag    bool
	Name    string
	Message string
	calib   bme280_cal
}

type Bme280_Vaule struct {
	Press string
	Temp  string
	Hum   string
}

type bme280_cal struct {
	press    []int
	temp     []int
	hum      []int
	timefine float64
}

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

func (t *Bme280) Init() bool {
	t.Name = "BME280"
	t.up()
	if !t.Test() {
		t.Close()
	} else {
		t.ReadCalib()
	}
	return t.Flag
}

func (t *Bme280) Close() {
	t.Flag = false
	t.down()
	t.Message = "Close"
}
func (t *Bme280) up() {
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

	t.WriteByte(BME280_CTRL_HUM, ctrl_hum_reg)
	t.WriteByte(BME280_CTRL_MEAS, ctrl_meas_reg)
	t.WriteByte(BME280_CONFIG, config_reg)
}

func (t *Bme280) down() {
	t.WriteByte(BME280_CTRL_HUM, 0x0)
	t.WriteByte(BME280_CTRL_MEAS, 0x0)
	t.WriteByte(BME280_CONFIG, 0x0)
}

func (t *Bme280) WriteByte(command, data byte) {
	i2c, err := i2c.New(BME280, I2C_BUS)
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

func (t *Bme280) ReadCalib() {
	var calib []int
	buf := t.ReadByte(BME280_CALIB00, int(BME280_CALIB25-BME280_CALIB00))
	for _, b := range buf {
		calib = append(calib, int(b))
	}
	buf = t.ReadByte(BME280_CALIB26, int(BME280_CALIB41-BME280_CALIB26+1))
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

func (t *Bme280) ReadByte(command byte, size int) []byte {
	i2c, err := i2c.New(BME280, I2C_BUS)
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
func (t *Bme280) Test() bool {
	t.Flag = false
	i2c, err := i2c.New(BME280, I2C_BUS)
	if err != nil {
		fmt.Println(err.Error())
		t.Message = err.Error()
		return false
	}
	i2c.Close()
	// fmt.Println(t.ReadByte(BME280_CTRL_MEAS, 1))
	for i := 0; i < 3; i++ {
		buf := t.ReadByte(BME280_ID, 1)
		if BME280_ID_DATA == buf[0] {
			t.Flag = true
			return true
		}
		t.down()
		time.Sleep(time.Microsecond * 1000)
		t.up()
	}
	fmt.Println(t.Message)
	return false
}
func (t *Bme280) Read() (int, int, int) {
	buf := t.ReadByte(BME280_PRESS_MSB, 8)
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

func (t *Bme280) Calib_Temp(b_temp int) (float64, float64) {
	tmp := float64(b_temp)
	var calib []float64
	for _, flt := range t.calib.temp {
		calib = append(calib, float64(flt))
	}
	v1 := (tmp/16384.0 - calib[0]/1024.0) * calib[1]
	v2 := (tmp/131072.0 - calib[0]/8192.0) * (tmp/131072.0 - calib[0]/8192.0) * calib[2]
	t_fine := v1 + v2
	temperature := t_fine / 5120.0
	return temperature, t_fine
	// #print "temp : %-6.2f ℃" % (temperature)
	// return "%.2f" % (temperature)
}
func (t *Bme280) Calib_Press(b_press int, t_fine float64) float64 {
	tmp := float64(b_press)
	var calib []float64
	for _, flt := range t.calib.press {
		calib = append(calib, float64(flt))
	}

	v1 := (t_fine / 2.0) - 64000.0
	v2 := (((v1 / 4.0) * (v1 / 4.0)) / 2048) * calib[5]
	v2 = v2 + ((v1 * calib[4]) * 2.0)
	v2 = (v2 / 4.0) + (calib[3] * 65536.0)
	v1 = (((calib[2] * (((v1 / 4.0) * (v1 / 4.0)) / 8192)) / 8) + ((calib[1] * v1) / 2.0)) / 262144
	v1 = ((32768 + v1) * calib[0]) / 32768
	if v1 == 0 {
		return 0
	} else {
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

	// #print "pressure : %7.2f hPa" % (pressure/100)
	// return "%7.2f" % (pressure / 100)

}
func (t *Bme280) Calib_Hum(b_hum int, t_fine float64) float64 {
	tmp := float64(b_hum)
	var calib []float64
	for _, flt := range t.calib.hum {
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
	// #print "hum : %6.2f ％" % (var_h)
	// return "%.2f" % (var_h)

}

func (t *Bme280) ReadData() (float64, float64, float64) {
	press, temp, hum := t.Read()
	c_tmp, t_fine := t.Calib_Temp(temp)
	c_press := t.Calib_Press(press, t_fine)
	c_hum := t.Calib_Hum(hum, t_fine)
	if c_hum <= 0 {
		return -1, -1, -1
	}
	return c_press, c_tmp, c_hum
}
