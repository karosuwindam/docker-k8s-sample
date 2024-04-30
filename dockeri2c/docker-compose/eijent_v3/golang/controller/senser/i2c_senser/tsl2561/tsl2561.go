package tsl2561

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

var (
	I2C_BUS = 1
)

type datastore struct {
	Lux      int
	Flag     bool
	StopFlag bool
	msg      msgsenser.Msg

	mu    sync.Mutex
	i2cMu *sync.Mutex
}

var memory datastore

var shudown chan bool
var reset chan bool
var done chan bool
var wait chan bool

func Init(i2cMu *sync.Mutex) error {
	memory = datastore{
		Lux:      -1,
		Flag:     false,
		StopFlag: false,
		msg:      msgsenser.Msg{},
		i2cMu:    i2cMu,
	}
	shudown = make(chan bool, 1)
	done = make(chan bool, 1)
	reset = make(chan bool, 1)
	wait = make(chan bool, 1)
	memory.msg.Create("TSL2561")
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

func Test() bool {
	if byte, err := readByte(IDDATA, 0x1); err != nil {
		log.Println("error:", err)
		memory.changeMsg("Read Error I2C")
		memory.changeFlag(false)
		return false
	} else {
		if byte[0] != 0x50 {
			msg := fmt.Sprintf("TSL2561 test header data 0x50 !=%x", byte[0])
			log.Println("error:", msg)
			memory.changeMsg(msg)
			memory.changeFlag(false)
			return false
		}
	}
	memory.changeMsg("OK")
	memory.changeFlag(true)
	return true
}

// xミリ秒ごとの読み取り開始
func Run() error {
	// if !memory.Flag {
	// 	return errors.New("errors Run")
	// }
	memory.chageRunFlag(true)
	log.Println("info:", "TSL2561 loop start")
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
		case <-time.After(time.Duration(config.Senser.Tsl2561_Count) * time.Millisecond):
			if memory.readFlag() {
				readdate()
			}
		}
	}
	memory.changeFlag(false)
	log.Println("info:", "TSL2561 loop stop")

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
		log.Panic("time over 1 sec")
	}
}

func readdate() {
	flag := true
	for i := 0; i < 3; i++ {
		num := readVisibleLux()
		if num > 0 {
			memory.changeValue(num)
			memory.changeMsg("OK")
			flag = false
			break
		} else {
			msg := fmt.Sprintf("Error Read Tsl2561 Count %v", i)
			memory.changeMsg(msg)
		}
		time.Sleep(time.Duration(config.Senser.Tsl2561_Count) * time.Millisecond)
	}
	if flag {
		memory.changeFlag(false)
		memory.changeMsg("Stop Tsl2561")

	}

}

func ResetMessage() {
	if len(reset) > 0 {
		return
	}
	reset < true
}

// 状態確認
// 有効とメッセージ情報
func Health() (bool, msgsenser.Msg) {
	return memory.readFlag(), memory.readMsg()
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

func (v *datastore) changeValue(num int) {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.Lux = num
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

// データストアの値を読み取る
func ReadValue() (int, bool) {
	memory.mu.Lock()
	defer memory.mu.Unlock()
	return memory.Lux, memory.readFlag()
}

func (v *datastore) readMsg() msgsenser.Msg {
	v.mu.Lock()
	defer v.mu.Unlock()
	return v.msg
}

// 現在の値を読み取る
func readVisibleLux() int {
	timing := 2
	gain := 0
	up()
	full, ir := readLux(gain, timing)
	if full < 500 && timing == 0 {
		timing = 1
		time.Sleep(time.Millisecond * 5)
		fmt.Println("dark 13.7ms to 101ms")
		full, ir = readLux(gain, timing)
	}
	if full < 500 && timing == 1 {
		timing = 2
		time.Sleep(time.Millisecond * 5)
		fmt.Println("dark 101ms to 402ms")
		full, ir = readLux(gain, timing)
	}
	if full < 500 && timing == 2 && gain == 0 {
		gain = 1
		time.Sleep(time.Millisecond * 5)
		fmt.Println("dark high gain")
		full, ir = readLux(gain, timing)
	}
	if (full < 20000 || ir > 20000) && timing == 2 && gain == 1 {
		gain = 0
		time.Sleep(time.Millisecond * 5)
		fmt.Println("light low gain")
		full, ir = readLux(gain, timing)
	}
	if (full < 20000 || ir > 20000) && timing == 2 {
		timing = 1
		time.Sleep(time.Millisecond * 5)
		fmt.Println("light 402ms to 101ms")
		full, ir = readLux(gain, timing)
	}
	if (full < 10000 || ir > 10000) && timing == 1 {
		timing = 0
		time.Sleep(time.Millisecond * 5)
		fmt.Println("light 101ms to 13.7ms")
		full, ir = readLux(gain, timing)
	}
	down()
	return calculateLux(gain, timing, full, ir)
}

func gainSelect(onoff, gain int) error {
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
	return writeByte(TIMING, data)
}

// Luxのデータを取得
func readLux(onoff, gain int) (int, int) {
	var full int = -1
	var ir int = -1
	sleep_time := time.Microsecond * 0
	if err := gainSelect(onoff, gain); err != nil {
		log.Println("error:", err)
		return full, ir
	}
	if gain == GAIN_INTEG_LOW {
		sleep_time = time.Millisecond*13 + time.Microsecond*700
	} else if gain == GAIN_INTEG_MID {
		sleep_time = time.Millisecond * 101
	} else {
		sleep_time = time.Millisecond * 402
	}
	time.Sleep(sleep_time + time.Millisecond*1)
	if tmp, err := readByte(DATA0LOW, 2); err != nil {
		log.Println("error:", err)
		return full, ir
	} else {
		full = (int(tmp[1])*256 + int(tmp[0]))
	}
	if tmp, err := readByte(DATA1LOW, 2); err != nil {
		log.Println("error:", err)
		return full, ir
	} else {
		ir = (int(tmp[1])*256 + int(tmp[0]))
	}
	return full, ir
}

func calculateLux(gain, timing, full, ir int) int {

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
	return lux
}

func writeByte(command, data byte) error {
	memory.i2cMu.Lock()
	defer memory.i2cMu.Unlock()
	i2c, err := i2c.New(TSL2561, I2C_BUS)
	if err != nil {
		return errors.Wrapf(err, "i2c.New(%v,%v)", TSL2561, I2C_BUS)
	}
	defer i2c.Close()
	_, err = i2c.Write([]byte{command, data})
	if err != nil {
		return errors.Wrapf(err, "i2c.Write(%v,%v)", command, data)
	}
	return nil
}

func readByte(command byte, size int) ([]byte, error) {
	buf := make([]byte, size)
	memory.i2cMu.Lock()
	defer memory.i2cMu.Unlock()
	i2c, err := i2c.New(TSL2561, I2C_BUS)
	if err != nil {
		return buf, errors.Wrapf(err, "i2c.New(%v,%v)", TSL2561, I2C_BUS)
	}
	defer i2c.Close()
	_, err = i2c.Write([]byte{command})
	if err != nil {
		return buf, errors.Wrapf(err, "i2c.Write(%v)", command)
	}
	_, err = i2c.Read(buf)
	if err != nil {
		return buf, errors.Wrap(err, "i2c.Read()")
	}
	return buf, nil
}

func up() error {
	return writeByte(CONTROL, CONTROL_ON)
}
func down() error {
	return writeByte(CONTROL, CONTROL_OFF)
}
