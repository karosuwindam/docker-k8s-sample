package senser

import (
	"fmt"
	"time"

	"github.com/davecheney/i2c"
)

const (
	TSL2561   uint8 = 0x29
	CONTROL   byte  = 0x80
	TIMING    byte  = 0x81
	INTERRUPT byte  = 0x86
	IDDATA    byte  = 0x8A
	DATA0LOW  byte  = 0x8C
	DATA0HIGH byte  = 0x8D
	DATA1LOW  byte  = 0x8E
	DATA1HIGH byte  = 0x8F
)

const (
	LUX_SCALE     = 14     // scale by 2^14
	RATIO_SCALE   = 9      // scale ratio by 2^9
	CH_SCALE      = 10     // scale channel values by 2^10
	CHSCALE_TINT0 = 0x7517 // 322/11 * 2^CH_SCALE
	CHSCALE_TINT1 = 0x0fe7 // 322/81 * 2^CH_SCALE
	K1T           = 0x0040 // 0.125 * 2^RATIO_SCALE
	B1T           = 0x01f2 // 0.0304 * 2^LUX_SCALE
	M1T           = 0x01be // 0.0272 * 2^LUX_SCALE
	K2T           = 0x0080 // 0.250 * 2^RATIO_SCA
	B2T           = 0x0214 // 0.0325 * 2^LUX_SCALE
	M2T           = 0x02d1 // 0.0440 * 2^LUX_SCALE
	K3T           = 0x00c0 // 0.375 * 2^RATIO_SCALE
	B3T           = 0x023f // 0.0351 * 2^LUX_SCALE
	M3T           = 0x037b // 0.0544 * 2^LUX_SCALE
	K4T           = 0x0100 // 0.50 * 2^RATIO_SCALE
	B4T           = 0x0270 // 0.0381 * 2^LUX_SCALE
	M4T           = 0x03fe // 0.0624 * 2^LUX_SCALE
	K5T           = 0x0138 // 0.61 * 2^RATIO_SCALE
	B5T           = 0x016f // 0.0224 * 2^LUX_SCALE
	M5T           = 0x01fc // 0.0310 * 2^LUX_SCALE
	K6T           = 0x019a // 0.80 * 2^RATIO_SCALE
	B6T           = 0x00d2 // 0.0128 * 2^LUX_SCALE
	M6T           = 0x00fb // 0.0153 * 2^LUX_SCALE
	K7T           = 0x029a // 1.3 * 2^RATIO_SCALE
	B7T           = 0x0018 // 0.00146 * 2^LUX_SCALE
	M7T           = 0x0012 // 0.00112 * 2^LUX_SCALE
	K8T           = 0x029a // 1.3 * 2^RATIO_SCALE
	B8T           = 0x0000 // 0.000 * 2^LUX_SCALE
	M8T           = 0x0000 // 0.000 * 2^LUX_SCALE

	K1C = 0x0043 // 0.130 * 2^RATIO_SCALE
	B1C = 0x0204 // 0.0315 * 2^LUX_SCALE
	M1C = 0x01ad // 0.0262 * 2^LUX_SCALE
	K2C = 0x0085 // 0.260 * 2^RATIO_SCALE
	B2C = 0x0228 // 0.0337 * 2^LUX_SCALE
	M2C = 0x02c1 // 0.0430 * 2^LUX_SCALE
	K3C = 0x00c8 // 0.390 * 2^RATIO_SCALE
	B3C = 0x0253 // 0.0363 * 2^LUX_SCALE
	M3C = 0x0363 // 0.0529 * 2^LUX_SCALE
	K4C = 0x010a // 0.520 * 2^RATIO_SCALE
	B4C = 0x0282 // 0.0392 * 2^LUX_SCALE
	M4C = 0x03df // 0.0605 * 2^LUX_SCALE
	K5C = 0x014d // 0.65 * 2^RATIO_SCALE
	B5C = 0x0177 // 0.0229 * 2^LUX_SCALE
	M5C = 0x01dd // 0.0291 * 2^LUX_SCALE
	K6C = 0x019a // 0.80 * 2^RATIO_SCALE
	B6C = 0x0101 // 0.0157 * 2^LUX_SCALE
	M6C = 0x0127 // 0.0180 * 2^LUX_SCALE
	K7C = 0x029a // 1.3 * 2^RATIO_SCALE
	B7C = 0x0037 // 0.00338 * 2^LUX_SCALE
	M7C = 0x002b // 0.00260 * 2^LUX_SCALE
	K8C = 0x029a // 1.3 * 2^RATIO_SCALE
	B8C = 0x0000 // 0.000 * 2^LUX_SCALE
	M8C = 0x0000 // 0.000 * 2^LUX_SCALE
)

const (
	CONTROL_ON          byte = 0x03
	CONTROL_OFF         byte = 0x00
	TIMING_GAIN_OFF     byte = 0x00
	TIMING_GAIN_ON      byte = 0x10
	TIMING_MANUAL_OFF   byte = 0x00
	TIMING_MANUAL_ON    byte = 0x08
	TIMING_INTEG_13_7MS byte = 0x00
	TIMING_INTEG_101MS  byte = 0x01
	TIMING_INTEG_402MS  byte = 0x02
	GAIN_ON             int  = 0
	GAIN_OFF            int  = 1
	GAIN_INTEG_LOW      int  = 0
	GAIN_INTEG_MID      int  = 1
	GAIN_INTEG_HIGH     int  = 2
)

var packageType int = 0

type Tsl2561 struct {
	Flag    bool
	Name    string
	Message string
}

func (t *Tsl2561) Init() {
	t.Name = "tsl256"
	t.Up()
	t.WriteByte(INTERRUPT, 0)
	t.Down()
}

func (t *Tsl2561) WriteByte(command, data byte) {
	i2c, err := i2c.New(TSL2561, I2C_BUS)
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
func (t *Tsl2561) ReadByte(command byte, size int) []byte {
	i2c, err := i2c.New(TSL2561, I2C_BUS)
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

func (t *Tsl2561) Test() bool {
	t.Flag = false
	i2c, err := i2c.New(TSL2561, I2C_BUS)
	if err != nil {
		fmt.Println(err.Error())
		t.Message = err.Error()
		return false
	}
	defer i2c.Close()
	tmp := t.ReadByte(IDDATA, 1)
	if tmp[0] != 0x50 {
		t.Message = "NG"
		return false
	}
	t.Message = "OK"
	t.Flag = true
	return true
}
func (t *Tsl2561) Up() {
	t.WriteByte(CONTROL, CONTROL_ON)
}
func (t *Tsl2561) GainSelect(onoff, gain int) {
	var data byte
	if onoff == GAIN_ON {
		data = TIMING_GAIN_ON
	} else {
		data = TIMING_GAIN_OFF
	}
	if gain == GAIN_INTEG_LOW {
		data |= TIMING_INTEG_13_7MS
	} else if gain == GAIN_INTEG_MID {
		data |= TIMING_INTEG_101MS
	} else {
		data |= TIMING_INTEG_402MS
	}
	t.WriteByte(TIMING, data)

}
func (t *Tsl2561) Down() {
	t.WriteByte(CONTROL, CONTROL_OFF)
}
func (t *Tsl2561) ReadLux(onoff, gain int) (int, int) {
	sleep_time := time.Microsecond * 0
	t.GainSelect(onoff, gain)
	if gain == GAIN_INTEG_LOW {
		sleep_time = time.Millisecond*13 + time.Microsecond*700
	} else if gain == GAIN_INTEG_MID {
		sleep_time = time.Millisecond * 101
	} else {
		sleep_time = time.Millisecond * 402
	}
	time.Sleep(sleep_time + time.Millisecond*1)
	ch0_data := t.ReadByte(DATA0LOW, 2)
	ch1_data := t.ReadByte(DATA1LOW, 2)
	full := (int(ch0_data[1])*256 + int(ch0_data[0]))
	ir := (int(ch1_data[1])*256 + int(ch1_data[0]))
	fmt.Println("ch0", full, "ch1", ir)
	return full, ir
}

func (t *Tsl2561) CalculateLux(gain, timing, full, ir int) int {
	chScale := 0
	ratio := 0
	b := 0
	m := 0
	if timing == GAIN_INTEG_LOW {
		chScale = CHSCALE_TINT0
	} else if timing == GAIN_INTEG_MID {
		chScale = CHSCALE_TINT1
	} else {
		chScale = 1 << CH_SCALE
	}
	if gain == GAIN_OFF {
		chScale = chScale << 4
	}
	schannel0 := (full * chScale) >> CH_SCALE
	schannel1 := (ir * chScale) >> CH_SCALE
	if schannel0 != 0 {
		ratio = (schannel1 << (RATIO_SCALE + 1)) / schannel0
	}
	ratio = (ratio + 1) >> 1
	if packageType == 0 { // T package
		if (ratio >= 0) && (ratio <= K1T) {
			b = B1T
			m = M1T
		} else if ratio <= K2T {
			b = B2T
			m = M2T
		} else if ratio <= K3T {
			b = B3T
			m = M3T
		} else if ratio <= K4T {
			b = B4T
			m = M4T
		} else if ratio <= K5T {
			b = B5T
			m = M5T
		} else if ratio <= K6T {
			b = B6T
			m = M6T
		} else if ratio <= K7T {
			b = B7T
			m = M7T
		} else if ratio > K8T {
			b = B8T
			m = M8T
		}

	} else if packageType == 1 { // CS package
		if (ratio >= 0) && (ratio <= K1C) {
			b = B1C
			m = M1C
		} else if ratio <= K2C {
			b = B2C
			m = M2C
		} else if ratio <= K3C {
			b = B3C
			m = M3C
		} else if ratio <= K4C {
			b = B4C
			m = M4C
		} else if ratio <= K5C {
			b = B5C
			m = M5C
		} else if ratio <= K6C {
			b = B6C
			m = M6C
		} else if ratio <= K7C {
			b = B7C
			m = M7C
		}
	}
	temp := ((schannel0 * b) - (schannel1 * m))
	if temp < 0 {
		temp = 0
	}

	temp += (1 << (LUX_SCALE - 1))

	lux := temp >> LUX_SCALE
	fmt.Println(ratio, schannel0, schannel1, b, m, lux)
	return lux
}

func (t *Tsl2561) ReadVisibleLux() int {
	timing := 2
	gain := 0
	t.Up()
	full, ir := t.ReadLux(gain, timing)
	if full < 500 && timing == 0 {
		timing = 1
		time.Sleep(time.Millisecond * 5)
		fmt.Println("dark 13.7ms to 101ms")
		full, ir = t.ReadLux(gain, timing)
	}
	if full < 500 && timing == 1 {
		timing = 2
		time.Sleep(time.Millisecond * 5)
		fmt.Println("dark 101ms to 402ms")
		full, ir = t.ReadLux(gain, timing)
	}
	if full < 500 && timing == 2 && gain == 0 {
		gain = 1
		time.Sleep(time.Millisecond * 5)
		fmt.Println("dark high gain")
		full, ir = t.ReadLux(gain, timing)
	}
	if (full < 20000 || ir > 20000) && timing == 2 && gain == 1 {
		gain = 0
		time.Sleep(time.Millisecond * 5)
		fmt.Println("light low gain")
		full, ir = t.ReadLux(gain, timing)
	}
	if (full < 20000 || ir > 20000) && timing == 2 {
		timing = 1
		time.Sleep(time.Millisecond * 5)
		fmt.Println("light 402ms to 101ms")
		full, ir = t.ReadLux(gain, timing)
	}
	if (full < 10000 || ir > 10000) && timing == 1 {
		timing = 0
		time.Sleep(time.Millisecond * 5)
		fmt.Println("light 101ms to 13.7ms")
		full, ir = t.ReadLux(gain, timing)
	}
	t.Down()
	return t.CalculateLux(gain, timing, full, ir)
}
func (t *Tsl2561) ReadData(onoff, gain int) {
	t.Up()
	t.GainSelect(onoff, gain)
	if gain == GAIN_INTEG_LOW {
		time.Sleep(time.Millisecond*14 + time.Microsecond*700)
	} else if gain == GAIN_INTEG_MID {
		time.Sleep(time.Millisecond*102 + time.Microsecond*0)
	} else {
		time.Sleep(time.Millisecond*403 + time.Microsecond*0)
	}
	ch0_data := t.ReadByte(DATA0LOW, 2)
	ch1_data := t.ReadByte(DATA1LOW, 2)
	full := float64(int(ch0_data[1])*256 + int(ch0_data[0]))
	ir := float64(int(ch1_data[1])*256 + int(ch1_data[0]))
	fmt.Println(ch0_data, ch1_data)
	fmt.Println(full, ir, float64(full-ir)*0.03)
	t.Down()

}
