package bmx055

import (
	"eijent/config"
	msgsenser "eijent/controller/senser/msg_senser"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"
)

var (
	BMX055_ACC = uint8(0x19)
	// BMX055_ACC = uint8(0x18)
	BMX055_GYRO = uint8(0x69)
	// BMX055_GYRO = uint8(0x68)
	BMX055_MAG = uint8(0x13)
	// BMX055_MAG = uint8(0x12)
	// BMX055_MAG = uint8(0x11)
	// BMX055_MAG = uint8(0x10)
)

// BMX055 ACC REGISTER
const (
	ACC_BGW_CHIPID    byte = 0x00
	ACC_IS_RESERVED_0 byte = iota
	ACC_ACCD_X_LSB
	ACC_ACCD_X_MSB
	ACC_ACCD_Y_LSB
	ACC_ACCD_Y_MSB
	ACC_ACCD_Z_LSB
	ACC_ACCD_Z_MSB
	ACC_ACCD_TEMP
	ACC_INT_STATUS_0
	ACC_INT_STATUS_1
	ACC_INT_STATUS_2
	ACC_INT_STATUS_3
	ACC_IS_RESERVED_1
	ACC_FIFO_STATUS
	ACC_PMU_RANGE
	ACC_PMU_BW
	ACC_PMU_LPW
	ACC_PMU_LOW_POWER
	ACC_ACC_HRW
	ACC_BGW_SOFTRESET
	ACC_IS_RESERVED_2
	ACC_INT_EN_0
	ACC_INT_EN_1
	ACC_INT_EN_2
	ACC_INT_MAP_0
	ACC_INT_MAP_1
	ACC_INT_MAP_2
	ACC_IS_RESERVED_3
	ACC_IS_RESERVED_4
	ACC_INT_SRC
	ACC_IS_RESERVED_5
	ACC_INT_OUT_CTRL
	ACC_INT_RST_LATCH
	ACC_INT_0
	ACC_INT_1
	ACC_INT_2
	ACC_INT_3
	ACC_INT_4
	ACC_INT_5
	ACC_INT_6
	ACC_INT_7
	ACC_INT_8
	ACC_INT_9
	ACC_INT_A
	ACC_INT_B
	ACC_INT_C
	ACC_INT_D
	ACC_FIFO_CONFIG_0
	ACC_IS_RESERVED_6
	ACC_PMU_SELF_TEST
	ACC_TRIM_NVM_CTRL
	ACC_BGW_SPI3_WDT
	ACC_IS_RESERVED_7
	ACC_OFC_CTRL
	ACC_OFC_SETTING
	ACC_OFC_OFFSET_X
	ACC_OFC_OFFSET_Y
	ACC_OFC_OFFSET_Z
	ACC_TRIM_GP0
	ACC_TRIM_GP1
	ACC_IS_RESERVED_8
	ACC_FIFO_CONFIG_1
	ACC_FIFO_DATA
)

const ( //BMX055 ACC SetUP
	ACC_RANGE_2G  byte = 0x1<<1 | 0x1    //2g
	ACC_RANGE_4G  byte = 0x1<<2 | 0x1    //4g
	ACC_RANGE_8G  byte = 0x1 << 3        //8g
	ACC_RANGE_16G byte = 0x1<<3 | 0x1<<2 //16g

	ACC_BW_7_81_HZ  byte = 0x1 << 3                       //7.81Hz
	ACC_BW_15_63_HZ byte = 0x1<<3 | 0x1                   //15.63Hz
	ACC_BW_31_25_HZ byte = 0x1<<3 | 0x1<<1                //31.25Hz
	ACC_BW_62_5_HZ  byte = 0x1<<3 | 0x1<<1 | 0x1          //62.5Hz
	ACC_BW_125_HZ   byte = 0x1<<3 | 0x1<<2                //125Hz
	ACC_BW_250_HZ   byte = 0x1<<3 | 0x1<<2 | 0x01         //250Hz
	ACC_BW_500_HZ   byte = 0x1<<3 | 0x1<<2 | 0x1<<1       //500Hz
	ACC_BW_1000_HZ  byte = 0x1<<3 | 0x1<<2 | 0x1<<1 | 0x1 //1000Hz

	ACC_LPW_0_5_MS  byte = 0x1<<2 | 0x1                   //0.5ms
	ACC_LPW_1_MS    byte = 0x1<<2 | 0x1<<1                //1ms
	ACC_LPW_2_MS    byte = 0x1<<2 | 0x1<<1 | 0x1          //2ms
	ACC_LPW_4_MS    byte = 0x1 << 3                       //4ms
	ACC_LPW_6_MS    byte = 0x1<<3 | 0x1                   //6ms
	ACC_LPW_10_MS   byte = 0x1<<3 | 0x1<<1                //10ms
	ACC_LPW_25_MS   byte = 0x1<<3 | 0x1<<1 | 0x1          //25ms
	ACC_LPW_50_MS   byte = 0x1<<3 | 0x1<<2                //50ms
	ACC_LPW_100_MS  byte = 0x1<<3 | 0x1<<2 | 0x1          //100ms
	ACC_LPW_500_MS  byte = 0x1<<3 | 0x1<<2 | 0x1<<1       //500ms
	ACC_LPW_1000_MS byte = 0x1<<3 | 0x1<<2 | 0x1<<1 | 0x1 //1000ms

	ACC_LPW_NOMAL        byte = 0x0<<7 | 0x0<<6 | 0x0<<5 //Nomal mode
	ACC_LPW_DEEP_SUSPEND byte = 0x0<<7 | 0x0<<6 | 0x1<<5 //deep suspend mode
	ACC_LPW_LOW_POWER    byte = 0x0<<7 | 0x1<<6 | 0x0<<5 //Low power mode
	ACC_LPW_SUSPEND      byte = 0x1<<7 | 0x0<<6 | 0x0<<5 //Suspend mode
)

// BMX055 GYRO REGISTER
const (
	GYR_CHIP_ID       byte = 0x00
	GYR_IS_RESERVED_0 byte = iota
	GYR_RATE_X_LSB
	GYR_RATE_X_MSB
	GYR_RATE_Y_LSB
	GYR_RATE_Y_MSB
	GYR_RATE_Z_LSB
	GYR_RATE_Z_MSB
	GYR_RESERVED
	GYR_INT_STATUS_0
	GYR_INT_STATUS_1
	GYR_INT_STATUS_2
	GYR_INT_STATUS_3
	GYR_IS_RESERVED_1
	GYR_FIFO_STATUS
	GYR_RANGE
	GYR_BW
	GYR_LPM1
	GYR_LPM2
	GYR_RATE_HBW
	GYR_BGW_SOFTRESET
	GYR_INT_EN_0
	GYR_INT_EN_1
	GYR_INT_MAP_0
	GYR_INT_MAP_1
	GYR_INT_MAP_2
	GYR_NONE_1
	GYR_NONE_2
	GYR_NONE_3
	GYR_IS_RESERVED_2
	GYR_NONE_4
	GYR_AND_RESERVED_1
	GYR_AND_RESERVED_2
	GYR_INT_RST_LATCH
	GYR_HIGH_TH_X
	GYR_HIGH_DUR_X
	GYR_HIGH_TH_Y
	GYR_HIGH_DUR_Y
	GYR_HIGH_TH_Z
	GYR_HIGH_DUR_Z
	GYR_AND_RESERVED_3
	GYR_AND_RESERVED_4
	GYR_AND_RESERVED_5
	GYR_AND_RESERVED_6
	GYR_AND_RESERVED_7
	GYR_AND_RESERVED_8
	GYR_AND_RESERVED_9
	GYR_AND_RESERVED_10
	GYR_AND_RESERVED_11
	GYR_SOC
	GYR_A_FOC
	GYR_TRIM_NVM_CTRL
	GYR_BGW_SPI3_WDT
	GYR_IS_RESERVED_3
	GYR_OFC1
	GYR_OFC2
	GYR_OFC3
	GYR_OFC4
	GYR_TRIM_GP0
	GYR_TRIM_GP1
	GYR_BIST
	GYR_FIFO_CONFIG_0
	GYR_FIFO_CONFIG_1
)

const ( //BMX055 GYR SetUP
	GYR_RANGE_2000 byte = 0x0<<2 | 0x0<<1 | 0x0 //Full Scale 2000 1/s
	GYR_RANGE_1000 byte = 0x0<<2 | 0x0<<1 | 0x1 //Full Scale 1000 1/s
	GYR_RANGE_500  byte = 0x0<<2 | 0x1<<1 | 0x0 //Full Scale 500 1/s
	GYR_RANGE_250  byte = 0x0<<2 | 0x1<<1 | 0x1 //Full Scale 250 1/s
	GYR_RANGE_125  byte = 0x1<<2 | 0x0<<1 | 0x0 //Full Scale 125 1/s

	GYR_BW_32_HZ  byte = 0x0<<3 | 0x1<<2 | 0x1<<1 | 0x1 //BandWidth 32Hz
	GYR_BW_64_HZ  byte = 0x0<<3 | 0x1<<2 | 0x1<<1 | 0x0 //BandWidth 64Hz
	GYR_BW_12_HZ  byte = 0x0<<3 | 0x1<<2 | 0x0<<1 | 0x1 //BandWidth 12Hz
	GYR_BW_23_HZ  byte = 0x0<<3 | 0x1<<2 | 0x0<<1 | 0x0 //BandWidth 23Hz
	GYR_BW_47_HZ  byte = 0x0<<3 | 0x0<<2 | 0x1<<1 | 0x1 //BandWidth 47Hz
	GYR_BW_116_HZ byte = 0x0<<3 | 0x0<<2 | 0x1<<1 | 0x0 //BandWidth 116Hz
	GYR_BW_230_HZ byte = 0x0<<3 | 0x0<<2 | 0x0<<1 | 0x0 //BandWidth 230Hz

	GYR_LPM1_SLEEP_DUR_2MS  byte = 0x0<<3 | 0x0<<2 | 0x0<<1 //Sleep Duration time 2ms
	GYR_LPM1_SLEEP_DUR_4MS  byte = 0x0<<3 | 0x0<<2 | 0x1<<1 //Sleep Duration time 4ms
	GYR_LPM1_SLEEP_DUR_5MS  byte = 0x0<<3 | 0x1<<2 | 0x0<<1 //Sleep Duration time 5ms
	GYR_LPM1_SLEEP_DUR_8MS  byte = 0x0<<3 | 0x1<<2 | 0x1<<1 //Sleep Duration time 8ms
	GYR_LPM1_SLEEP_DUR_10MS byte = 0x1<<3 | 0x0<<2 | 0x0<<1 //Sleep Duration time 10ms
	GYR_LPM1_SLEEP_DUR_15MS byte = 0x1<<3 | 0x0<<2 | 0x1<<1 //Sleep Duration time 15ms
	GYR_LPM1_SLEEP_DUR_18MS byte = 0x1<<3 | 0x1<<2 | 0x0<<1 //Sleep Duration time 18ms
	GYR_LPM1_SLEEP_DUR_20MS byte = 0x1<<3 | 0x1<<2 | 0x1<<1 //Sleep Duration time 20ms

	GYR_LPM1_NORMAL       byte = 0x0<<7 | 0x0<<5 //Normal mode
	GYR_LPM1_DEEP_SUSPEND byte = 0x0<<7 | 0x1<<5 //Deep suspend mode
	GYR_LPM1_SUSPEND      byte = 0x1<<7 | 0x0<<5 //Suspen mode
)

// BMX055 MAG REGISTER
const (
	MAG_CHIP_ID    byte = 0x40
	MAG_DATA_X_LSB byte = iota + 0x41
	MAG_DATA_X_MSB
	MAG_DATA_Y_LSB
	MAG_DATA_Y_MSB
	MAG_DATA_Z_LSB
	MAG_DATA_Z_MSB
	MAG_RHALL_LSB
	MAG_RHALL_MSB
	MAG_CTL_INTERRUPT_1
	MAG_CTL_POWER_RESET_SPI
	MAG_CTL_OPT_OUT_RATE_SELFTEST
	MAG_CTL_INTERPUT_2
	MAG_CTL_INTERPUT_3 //Interrupt settings and axes enable bits control register
	MAG_LOW_THRESHOLD
	MAG_HIGH_THRESHOLD
	MAG_REPXY
	MAG_REPZ
)

const (
	MAG_CTL_POEWER_RESET_ON       byte = 0x1<<7 | 0x1<<1 //One of the soft reset trigger bits
	MAG_CTL_POEWER_RESET_OFF      byte = 0x0<<7 | 0x0<<1 //One of the soft reset trigger bits
	MAG_CTL_POWER_CONTROL_NORMAL  byte = 0x1
	MAG_CTL_POWER_CONTROL_SUSPEND byte = 0x0
	MAG_CTL_SPI3EN_ON             byte = 0x1 << 2
	MAG_CTL_SPI3EN_OFF            byte = 0x0 << 2

	MAG_CTL_OPT_DATA_RATE_10_HZ byte = 0x0<<5 | 0x0<<4 | 0x0<<3
	MAG_CTL_OPT_DATA_RATE_2_HZ  byte = 0x0<<5 | 0x0<<4 | 0x1<<3
	MAG_CTL_OPT_DATA_RATE_6_HZ  byte = 0x0<<5 | 0x1<<4 | 0x0<<3
	MAG_CTL_OPT_DATA_RATE_8_HZ  byte = 0x0<<5 | 0x1<<4 | 0x1<<3
	MAG_CTL_OPT_DATA_RATE_15_HZ byte = 0x1<<5 | 0x0<<4 | 0x0<<3
	MAG_CTL_OPT_DATA_RATE_20_HZ byte = 0x1<<5 | 0x0<<4 | 0x1<<3
	MAG_CTL_OPT_DATA_RATE_25_HZ byte = 0x1<<5 | 0x1<<4 | 0x0<<3
	MAG_CTL_OPT_DATA_RATE_30_HZ byte = 0x1<<5 | 0x1<<4 | 0x1<<3
	MAG_CTL_OPT_ADV_NORMAL      byte = 0x0<<2 | 0x0<<1
	MAG_CTL_OPT_ADV_FORCED      byte = 0x0<<2 | 0x1<<1
	MAG_CTL_OPT_ADV_SLEEPbyte   byte = 0x1<<2 | 0x1<<1

	MAG_CTL_DATA_READY_PIN_ON       byte = 0x1 << 7
	MAG_CTL_DATA_READY_PIN_OFF      byte = 0x0 << 7
	MAG_CTL_INTERRUPT_PIN_ON        byte = 0x1 << 6
	MAG_CTL_INTERRUPT_PIN_OFF       byte = 0x0 << 6
	MAG_CTL_CHANNEL_Z_ON            byte = 0x0 << 5
	MAG_CTL_CHANNEL_Z_OFF           byte = 0x1 << 5
	MAG_CTL_CHANNEL_Y_ON            byte = 0x0 << 4
	MAG_CTL_CHANNEL_Y_OFF           byte = 0x1 << 4
	MAG_CTL_CHANNEL_X_ON            byte = 0x0 << 3
	MAG_CTL_CHANNEL_X_OFF           byte = 0x1 << 3
	MAG_CTL_DR_POLARITY_HIGH        byte = 0x1 << 2
	MAG_CTL_DR_POLARITY_LOW         byte = 0x0 << 2
	MAG_CTL_INTERRUPT_LATCH_MEANS   byte = 0x1 << 1
	MAG_CTL_INTERRUPT_LATCH_NON     byte = 0x0 << 1
	MAG_CTL_INTERRUPT_POLARITY_HIGH byte = 0x1 << 0
	MAG_CTL_INTERRUPT_POLARITY_LOW  byte = 0x0 << 0
)
const (
	SENSER_NAME    string  = "BMX055"
	BMX055_ID      byte    = 0xfa
	ACC_THRESHOLD  float64 = 0.0196
	GYRO_THRESHOLD float64 = 0.0038
)

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
	for i := 0; i < 3; i++ {
		if Test() {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}
	if !memory.Flag {
		return errors.New("not Init Error for " + SENSER_NAME)
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
		case <-time.After(time.Duration(config.Senser.MMA8452Q_Count) * time.Microsecond):
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

func ReadValue() (interface{}, bool) {
	return nil, true
}

func ResetMessage() {
	if len(reset) > 0 {
		return
	}
	reset <- true
}

func readdate() {
}

func Test() bool {
	flag := false
	msg := ""
	if buf, err := readByte(BMX055_ACC, ACC_BGW_CHIPID, 1); err != nil {
		msg = fmt.Sprintf("%v Test Read Error Addr %x", SENSER_NAME, BMX055_ACC)
	} else if buf[0] != BMX055_ID {
		msg = fmt.Sprintf("%v Test test header data %x !=%x", SENSER_NAME, BMX055_ID, buf[0])
	} else {
		msg = "OK"
		flag = true
	}

	memory.changeFlag(flag)
	memory.changeMsg(msg)
	return flag
}

func accInit() {
	writeByte(BMX055_ACC, ACC_PMU_RANGE, ACC_RANGE_2G)               //範囲を+-2G設定
	writeByte(BMX055_ACC, ACC_PMU_BW, ACC_BW_7_81_HZ)                //帯域幅7.81Hz
	writeByte(BMX055_ACC, ACC_PMU_LPW, ACC_LPW_NOMAL|ACC_LPW_0_5_MS) //電力モードと低電力スリープ期間 0.5msでNomal mode
	time.Sleep(100 * time.Millisecond)

	writeByte(BMX055_GYRO, GYR_RANGE, GYR_RANGE_125)                         //範囲125/s
	writeByte(BMX055_GYRO, GYR_BW, GYR_BW_32_HZ)                             //帯域 32Hz
	writeByte(BMX055_GYRO, GYR_LPM1, GYR_LPM1_NORMAL|GYR_LPM1_SLEEP_DUR_2MS) //電力モードと低電力スリープ期間 2msでNomal mode
	time.Sleep(100 * time.Millisecond)

	writeByte(BMX055_MAG, MAG_CTL_POWER_RESET_SPI, MAG_CTL_POEWER_RESET_ON|MAG_CTL_POWER_CONTROL_NORMAL) //リセットとNomal mode
	time.Sleep(100 * time.Millisecond)
	writeByte(BMX055_MAG, MAG_CTL_POWER_RESET_SPI, MAG_CTL_POEWER_RESET_OFF|MAG_CTL_POWER_CONTROL_NORMAL) //Nomal mode
	time.Sleep(100 * time.Millisecond)
	writeByte(BMX055_MAG, MAG_CTL_OPT_OUT_RATE_SELFTEST, MAG_CTL_OPT_ADV_NORMAL|MAG_CTL_OPT_DATA_RATE_10_HZ) //拡張出力を Nomalかつ10Hz
	writeByte(BMX055_MAG, MAG_CTL_INTERPUT_3, MAG_CTL_DATA_READY_PIN_ON|MAG_CTL_DR_POLARITY_HIGH)
	writeByte(BMX055_MAG, MAG_REPXY, 0x04) //1+2*4=9
	writeByte(BMX055_MAG, MAG_REPZ, 0x16)  //1+22=23
}

type Axis struct {
	X int
	Y int
	Z int
}

type ACCAxis struct {
	X float64
	Y float64
	Z float64
}

type GyroAxis struct {
	X float64
	Y float64
	Z float64
}

func getACCRAW() (Axis, error) {
	var out Axis
	if buf, err := readByte(BMX055_ACC, ACC_ACCD_X_LSB, 6); err != nil {
		return out, err
	} else {
		tmp := make([]int, 3)
		tmp[0] = int(uint16(buf[1])*256|uint16(buf[0]&0xf0)) >> 4
		tmp[1] = int(uint16(buf[3])*256|uint16(buf[2]&0xf0)) >> 4
		tmp[2] = int(uint16(buf[5])*256|uint16(buf[4]&0xf0)) >> 4
		for i := 0; i < len(tmp); i++ {
			if tmp[i] > 2047 {
				tmp[i] -= 4096
			}
		}
		out = Axis{tmp[0], tmp[1], tmp[2]}
	}

	return out, nil
}

func axis_to_ACC(num Axis) ACCAxis {
	return ACCAxis{
		float64(num.X) * ACC_THRESHOLD,
		float64(num.Y) * ACC_THRESHOLD,
		float64(num.Z) * ACC_THRESHOLD,
	}
}

func getGyroRAW() (Axis, error) {
	var out Axis
	if buf, err := readByte(BMX055_GYRO, GYR_RATE_X_LSB, 6); err != nil {
		return out, err
	} else {
		tmp := make([]int, 3)
		tmp[0] = int(uint16(buf[1])<<8 | uint16(buf[0]))
		tmp[1] = int(uint16(buf[3])<<8 | uint16(buf[2]))
		tmp[2] = int(uint16(buf[5])<<8 | uint16(buf[4]))
		for i := 0; i < len(tmp); i++ {
			if tmp[i] > 32767 {
				tmp[i] -= 65536
			}
		}
		out = Axis{tmp[0], tmp[1], tmp[2]}
	}

	return out, nil
}

func axis_to_Gyro(num Axis) GyroAxis {
	return GyroAxis{
		float64(num.X) * GYRO_THRESHOLD,
		float64(num.Y) * GYRO_THRESHOLD,
		float64(num.Z) * GYRO_THRESHOLD,
	}
}
func getMag() (Axis, error) {
	var out Axis
	if buf, err := readByte(BMX055_MAG, MAG_DATA_X_LSB, 6); err != nil {
		return out, err
	} else {
		tmp := make([]int, 3)
		tmp[0] = int(uint16(buf[1])<<8|uint16(buf[0]&0xf8)) / 8
		tmp[1] = int(uint16(buf[3])<<8|uint16(buf[2]&0xf8)) / 8
		tmp[2] = int(uint16(buf[5])<<8|uint16(buf[4]&0xf8)) / 8
		for i := 0; i < len(tmp); i++ {
			if tmp[i] > 4095 {
				tmp[i] -= 8192
			}
		}
		out = Axis{tmp[0], tmp[1], tmp[2]}
	}

	return out, nil
}
