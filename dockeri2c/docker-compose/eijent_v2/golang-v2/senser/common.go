package senser

import (
	"fmt"
	"math"
	"strconv"
	"sync"
	"time"
)

var (
	I2C_BUS         = 1
	SRIAL_PORT      = "/dev/ttyS0"
	DHT_GPIO        = 12
	GYRO_SLEEP_TIME = time.Microsecond * 500
	GYRO_COUNT_MAX  = 360
)

type ResetFlag struct {
	reset bool
	mu    sync.Mutex
}

type Sennser struct {
	Bme280_data    Bme280
	Am2320_data    Am2320
	CO2Sensor_data MhZ19c
	Tsl2561_data   Tsl2561
	DhtSenser_data DhtSenser
	Mma8452q_data  Mma8452q
}

type SenserValue struct {
	Mu        sync.Mutex
	Bme280    Bme280_Vaule
	Am2320    Am2320_Vaule
	CO2       MhZ19c_Vaule
	Tsl2561   Tsl2561_Vaule
	DhtSenser DhtSenser_Vaule
	Mma8452q  Mma8452q_Vaule
	CpuTmp    string
}

type ZeroMma8452q struct {
	Mu sync.Mutex
	X  float64
	Y  float64
	Z  float64
}

type TempMma8452q struct {
	Mu   sync.Mutex
	Data []ZeroMma8452q
}

var zeromma8452qdata ZeroMma8452q = ZeroMma8452q{
	X: 0, Y: 0, Z: 0,
}

var SennserData Sennser = Sennser{}

var SennserDataValue SenserValue = SenserValue{
	Bme280:    Bme280_Vaule{},
	Am2320:    Am2320_Vaule{},
	CO2:       MhZ19c_Vaule{},
	Tsl2561:   Tsl2561_Vaule{},
	DhtSenser: DhtSenser_Vaule{},
	Mma8452q:  Mma8452q_Vaule{},
	CpuTmp:    "",
}

var i2cmu sync.Mutex

var resetFlagData ResetFlag = ResetFlag{}

func SennserSetup() {
	ch1 := make(chan bool)
	ch2 := make(chan bool)
	ch3 := make(chan bool)
	SennserData.Bme280_data = Bme280{}
	SennserData.Am2320_data = Am2320{}
	SennserData.CO2Sensor_data = MhZ19c{}
	SennserData.Tsl2561_data = Tsl2561{}
	SennserData.DhtSenser_data = DhtSenser{}
	SennserData.Mma8452q_data = Mma8452q{}

	//sriale
	go func() {
		for i := 0; i < 3; i++ {
			if !SennserData.CO2Sensor_data.Init(SRIAL_PORT) {
				fmt.Printf("Count:%v Srical Error for MH-Z19C\n", i+1)
				if "open "+SRIAL_PORT+": permission denied" == SennserData.CO2Sensor_data.Message {
					break
				}
			} else {
				break
			}
			time.Sleep(500 * time.Millisecond)
		}
		ch1 <- true
		return
	}()
	//i2c
	go func() {
		mma := make(chan bool)
		go func() {
			if !SennserData.Mma8452q_data.Init() {
				fmt.Println("I2C not for MMA8452q")
			} else {
				i2cmu.Lock()
				tmpzero := []ZeroMma8452q{}
				for i := 0; i < 20; i++ {
					x, y, z := SennserData.Mma8452q_data.ReadData()
					ax, ay, az := SennserData.Mma8452q_data.Ch_data_to_accel(x, y, z)
					tmp := ZeroMma8452q{X: ax, Y: ay, Z: az}
					tmpzero = append(tmpzero, tmp)
					time.Sleep(GYRO_SLEEP_TIME)
				}
				zeromma8452qdata.Mu.Lock()
				zeromma8452qdata.X = 0
				zeromma8452qdata.Y = 0
				zeromma8452qdata.Z = 0
				for i, tmp := range tmpzero {
					zeromma8452qdata.X = (zeromma8452qdata.X*float64(i) + tmp.X) / (float64(i) + 1)
					zeromma8452qdata.Y = (zeromma8452qdata.Y*float64(i) + tmp.Y) / (float64(i) + 1)
					zeromma8452qdata.Z = (zeromma8452qdata.Z*float64(i) + tmp.Z) / (float64(i) + 1)
				}
				zeromma8452qdata.Mu.Unlock()
				i2cmu.Unlock()
			}
			mma <- true
		}()
		i2cmu.Lock()
		if !SennserData.Bme280_data.Init() {
			fmt.Println("I2C not for BME280")
		}
		i2cmu.Unlock()
		i2cmu.Lock()
		if !SennserData.Am2320_data.Init() {
			fmt.Println("I2C not for AM2320")
		}
		i2cmu.Unlock()
		i2cmu.Lock()
		if !SennserData.Tsl2561_data.Init() {
			fmt.Println("I2C not for TSL2561")
		}
		i2cmu.Unlock()
		<-mma
		ch2 <- true
	}()
	//GPIO
	go func() {
		if !SennserData.DhtSenser_data.Init(DHT_GPIO) {
			fmt.Println("GPIO not for DHT Senser")
		}
		ch3 <- true
	}()
	<-ch1
	<-ch2
	<-ch3

	SennserResetSet(false)
	return
}

func SennserResetSet(flag bool) {
	resetFlagData.mu.Lock()
	resetFlagData.reset = flag
	resetFlagData.mu.Unlock()
}

func SennserResetRead() bool {
	resetFlagData.mu.Lock()
	flag := resetFlagData.reset
	resetFlagData.mu.Unlock()
	return flag
}

var tempMma8452q TempMma8452q = TempMma8452q{
	Data: []ZeroMma8452q{},
}

//センサーの短時間読み取り
func SenserMoveRead() {
	if SennserData.Mma8452q_data.Flag {
		tmp := ZeroMma8452q{}
		i2cmu.Lock()
		x, y, z := SennserData.Mma8452q_data.ReadData()
		i2cmu.Unlock()
		tmp.X, tmp.Y, tmp.Z = SennserData.Mma8452q_data.Ch_data_to_accel(x, y, z)
		tempMma8452q.Mu.Lock()
		arytmp := tempMma8452q.Data
		temparymma := []ZeroMma8452q{tmp}
		for i := 0; i < GYRO_COUNT_MAX; i++ {
			if i < len(arytmp) {
				temparymma = append(temparymma, arytmp[i])
			} else {
				break
			}
		}
		tempMma8452q.Data = temparymma
		tempMma8452q.Mu.Unlock()
	}
}

func SennserMoveReadData() {
	tmp := ZeroMma8452q{}
	tempMma8452q.Mu.Lock()
	tempdata := tempMma8452q.Data
	tempMma8452q.Mu.Unlock()
	zeromma8452qdata.Mu.Lock()
	tmpzero := zeromma8452qdata
	zeromma8452qdata.Mu.Unlock()
	//ゼロ位置調整を実施してその差分から絶対値が大きいものを抜き出す
	for i, tempary := range tempdata {
		if i == 0 {
			tmp.X = tempary.X
			tmp.Y = tempary.Y
			tmp.Z = tempary.Z
		} else {
			if math.Abs(tmp.X-tmpzero.X) < math.Abs(tempary.X-tmpzero.X) {
				tmp.X = tempary.X
			}
			if math.Abs(tmp.Y-tmpzero.Y) < math.Abs(tempary.Y-tmpzero.Y) {
				tmp.Y = tempary.Y
			}
			if math.Abs(tmp.Z-tmpzero.Z) < math.Abs(tempary.Z-tmpzero.Z) {
				tmp.Z = tempary.Z
			}
		}
	}
	//出力準備
	SennserDataValue.Mu.Lock()
	SennserDataValue.Mma8452q.X = strconv.FormatFloat(tmp.X, 'f', 2, 64)
	SennserDataValue.Mma8452q.Y = strconv.FormatFloat(tmp.Y, 'f', 2, 64)
	SennserDataValue.Mma8452q.Z = strconv.FormatFloat(tmp.Z, 'f', 2, 64)
	SennserDataValue.Mma8452q.Zero_X = strconv.FormatFloat(tmpzero.X, 'f', 2, 64)
	SennserDataValue.Mma8452q.Zero_Y = strconv.FormatFloat(tmpzero.Y, 'f', 2, 64)
	SennserDataValue.Mma8452q.Zero_Z = strconv.FormatFloat(tmpzero.Z, 'f', 2, 64)
	SennserDataValue.Mu.Unlock()
}

//通常センサーの読み取り
func SenserRead() {
	if SennserResetRead() {
		fmt.Println("Reset Sennser Set Up.")
		SennserSetup()
	}
	if SennserData.Bme280_data.Flag {
		i2cmu.Lock()
		press, temp, hum := SennserData.Bme280_data.ReadData()
		i2cmu.Unlock()
		SennserDataValue.Mu.Lock()
		SennserDataValue.Bme280.Press = strconv.FormatFloat(press, 'f', 2, 64)
		SennserDataValue.Bme280.Temp = strconv.FormatFloat(temp, 'f', 2, 64)
		SennserDataValue.Bme280.Hum = strconv.FormatFloat(hum, 'f', 2, 64)
		SennserDataValue.Mu.Unlock()
	}
	if SennserData.Am2320_data.Flag {
		i2cmu.Lock()
		hum, temp := SennserData.Am2320_data.Read()
		i2cmu.Unlock()
		if hum != -1 && temp != -1 {
			SennserDataValue.Mu.Lock()
			SennserDataValue.Am2320.Temp = strconv.FormatFloat(temp, 'f', 1, 64)
			SennserDataValue.Am2320.Hum = strconv.FormatFloat(hum, 'f', 1, 64)
			SennserDataValue.Mu.Unlock()
		}
	}
	if SennserData.CO2Sensor_data.Flag {
		co2, temp := SennserData.CO2Sensor_data.Read()
		SennserDataValue.Mu.Lock()
		SennserDataValue.CO2.Co2 = strconv.Itoa(co2)
		SennserDataValue.CO2.Temp = strconv.Itoa(temp)
		SennserDataValue.Mu.Unlock()
	}
	if SennserData.Tsl2561_data.Flag {
		i2cmu.Lock()
		lux := SennserData.Tsl2561_data.ReadVisibleLux()
		i2cmu.Unlock()
		SennserDataValue.Mu.Lock()
		SennserDataValue.Tsl2561.Lux = strconv.Itoa(lux)
		SennserDataValue.Mu.Unlock()

	}
	if SennserData.DhtSenser_data.Flag {
		hum, temp := SennserData.DhtSenser_data.Read()
		if hum != -1 && temp != -1 {
			SennserDataValue.Mu.Lock()
			SennserDataValue.DhtSenser.Hum = strconv.FormatFloat(hum, 'f', 1, 64)
			SennserDataValue.DhtSenser.Temp = strconv.FormatFloat(temp, 'f', 1, 64)
			SennserDataValue.Mu.Unlock()
		}
	}
	if SennserData.Mma8452q_data.Flag {
		SennserMoveReadData()
	}
	SennserDataValue.Mu.Lock()
	SennserDataValue.CpuTmp = cpuTmp()
	SennserDataValue.Mu.Unlock()

}
