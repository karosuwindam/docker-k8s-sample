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

// BMX055 MAG REGISTER
const ()

const (
	SENSER_NAME string = "BMX055"
	BMX055_ID   byte   = 0xfa
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
