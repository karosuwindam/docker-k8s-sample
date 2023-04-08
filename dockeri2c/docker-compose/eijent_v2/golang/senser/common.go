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
	GYRO_COUNT_MAX  = 1000
	ERROR_COUNT_MAX = 1000
)

const (
	ERR_BME280    = 1 << 0
	ERR_AM2320    = 1 << 1
	ERR_CO2SENSOR = 1 << 2
	ERR_TSL2561   = 1 << 3
	ERR_DHTSENSER = 1 << 4
	ERR_MM8452Q   = 1 << 5
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

type senserErrorCount struct {
	Mu        sync.Mutex
	Bme280    int
	Am2320    int
	CO2       int
	Tsl2561   int
	DhtSenser int
	Mma8452q  int
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
var gpiomu sync.Mutex
var serialmu sync.Mutex

var resetFlagData ResetFlag = ResetFlag{}

var errorCount senserErrorCount = senserErrorCount{
	Bme280:    0,
	Am2320:    0,
	CO2:       0,
	Tsl2561:   0,
	DhtSenser: 0,
	Mma8452q:  0,
}

func errorCountUP(count int) int {
	output := 0
	errorCount.Mu.Lock()
	switch count {
	case ERR_AM2320:
		errorCount.Am2320++
		output = errorCount.Am2320
	case ERR_BME280:
		errorCount.Bme280++
		output = errorCount.Bme280
	case ERR_CO2SENSOR:
		errorCount.CO2++
		output = errorCount.CO2
	case ERR_TSL2561:
		errorCount.Tsl2561++
		output = errorCount.Tsl2561
	case ERR_DHTSENSER:
		errorCount.DhtSenser++
		output = errorCount.DhtSenser
	case ERR_MM8452Q:
		errorCount.Mma8452q++
		output = errorCount.Mma8452q
	}
	errorCount.Mu.Unlock()
	return output
}

func errorCountReset(count int) {
	errorCount.Mu.Lock()
	switch count {
	case ERR_AM2320:
		errorCount.Am2320 = 0
	case ERR_BME280:
		errorCount.Bme280 = 0
	case ERR_CO2SENSOR:
		errorCount.CO2 = 0
	case ERR_TSL2561:
		errorCount.Tsl2561 = 0
	case ERR_DHTSENSER:
		errorCount.DhtSenser = 0
	case ERR_MM8452Q:
		errorCount.Mma8452q = 0
	}
	errorCount.Mu.Unlock()
}

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
			serialmu.Lock()
			if !SennserData.CO2Sensor_data.Init(SRIAL_PORT) {
				serialmu.Unlock()
				fmt.Printf("Count:%v Srical Error for MH-Z19C\n", i+1)
				if "open "+SRIAL_PORT+": permission denied" == SennserData.CO2Sensor_data.Message {
					break
				}
			} else {
				serialmu.Unlock()
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
				i2cmu.Unlock()
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
		gpiomu.Lock()
		if !SennserData.DhtSenser_data.Init(DHT_GPIO) {
			fmt.Println("GPIO not for DHT Senser")
		}
		gpiomu.Unlock()
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

var zeroCount int = 0

// センサーの短時間読み取り
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
		if zeroCount > GYRO_COUNT_MAX {
			tmpx := 0.0
			tmpy := 0.0
			tmpz := 0.0
			for _, tmpdata := range temparymma {
				tmpx += tmpdata.X
				tmpy += tmpdata.Y
				tmpz += tmpdata.Z
			}
			zeromma8452qdata.Mu.Lock()
			zeromma8452qdata.X = tmpx / float64(len(temparymma))
			zeromma8452qdata.Y = tmpy / float64(len(temparymma))
			zeromma8452qdata.Z = tmpz / float64(len(temparymma))
			zeromma8452qdata.Mu.Unlock()
			zeroCount = 0
		} else {
			zeroCount++
		}
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

// 通常センサーの読み取り
func SenserRead() {
	ch := make(chan bool, 7)
	if SennserResetRead() {
		fmt.Println("Reset Sennser Set Up.")
		SennserSetup()
	}
	go func() { //ch1
		if SennserData.Bme280_data.Flag {
			for i := 0; i < 2; i++ {
				i2cmu.Lock()
				press, temp, hum := SennserData.Bme280_data.ReadData()
				i2cmu.Unlock()
				if press == temp && temp == hum && hum == -1 {
					SennserData.Bme280_data.Message = "Read Error Bm280"
					if errorCountUP(ERR_BME280) > ERROR_COUNT_MAX {
						SennserData.Bme280_data.Flag = false
						break
					}
					time.Sleep(time.Microsecond * 100)
				} else {
					SennserDataValue.Mu.Lock()
					SennserDataValue.Bme280.Press = strconv.FormatFloat(press, 'f', 2, 64)
					SennserDataValue.Bme280.Temp = strconv.FormatFloat(temp, 'f', 2, 64)
					SennserDataValue.Bme280.Hum = strconv.FormatFloat(hum, 'f', 2, 64)
					SennserDataValue.Mu.Unlock()
					errorCountReset(ERR_BME280)
					break
				}
			}

		}
		ch <- true
	}()
	go func() { //ch2
		if SennserData.Am2320_data.Flag {
			for i := 0; i < 2; i++ {
				i2cmu.Lock()
				hum, temp := SennserData.Am2320_data.Read()
				i2cmu.Unlock()
				if hum != -1 && temp != -1 {
					SennserDataValue.Mu.Lock()
					SennserDataValue.Am2320.Temp = strconv.FormatFloat(temp, 'f', 1, 64)
					SennserDataValue.Am2320.Hum = strconv.FormatFloat(hum, 'f', 1, 64)
					SennserDataValue.Mu.Unlock()
					errorCountReset(ERR_AM2320)
					break
				} else {
					fmt.Println("reread AM2320")
					if errorCountUP(ERR_AM2320) > ERROR_COUNT_MAX {
						SennserData.Am2320_data.Flag = false
						break
					}
					time.Sleep(time.Microsecond * 100)
				}
			}

		}
		ch <- true
	}()
	go func() { //ch3
		if SennserData.CO2Sensor_data.Flag {
			for i := 0; i < 2; i++ {
				serialmu.Lock()
				co2, temp := SennserData.CO2Sensor_data.Read()
				serialmu.Unlock()
				if co2 > 0 {
					SennserDataValue.Mu.Lock()
					SennserDataValue.CO2.Co2 = strconv.Itoa(co2)
					SennserDataValue.CO2.Temp = strconv.Itoa(temp)
					SennserDataValue.Mu.Unlock()
					errorCountReset(ERR_CO2SENSOR)
					break
				} else {
					fmt.Println("reread Co2Sensor")
					fmt.Println("reread AM2320")
					if errorCountUP(ERR_CO2SENSOR) > ERROR_COUNT_MAX {
						SennserData.CO2Sensor_data.Flag = false
						break
					}
					time.Sleep(time.Microsecond * 100)

				}
			}
		}
		ch <- true
	}()
	go func() { //ch4
		if SennserData.Tsl2561_data.Flag {
			i2cmu.Lock()
			lux := SennserData.Tsl2561_data.ReadVisibleLux()
			i2cmu.Unlock()
			SennserDataValue.Mu.Lock()
			SennserDataValue.Tsl2561.Lux = strconv.Itoa(lux)
			SennserDataValue.Mu.Unlock()
			errorCountReset(ERR_TSL2561)
		}
		ch <- true
	}()
	go func() { //ch5
		if SennserData.DhtSenser_data.Flag {
			for i := 0; i < 2; i++ {
				gpiomu.Lock()
				hum, temp := SennserData.DhtSenser_data.Read()
				gpiomu.Unlock()
				if hum != -1 && temp != -1 {
					SennserDataValue.Mu.Lock()
					SennserDataValue.DhtSenser.Hum = strconv.FormatFloat(hum, 'f', 1, 64)
					SennserDataValue.DhtSenser.Temp = strconv.FormatFloat(temp, 'f', 1, 64)
					SennserDataValue.Mu.Unlock()
					errorCountReset(ERR_DHTSENSER)
					break
				} else {
					fmt.Println("reread DhtSenser")
					if errorCountUP(ERR_DHTSENSER) > ERROR_COUNT_MAX {
						SennserData.DhtSenser_data.Flag = false
						break
					}
					time.Sleep(time.Microsecond * 100)
				}
			}
		}
		ch <- true
	}()
	go func() { //ch6
		if SennserData.Mma8452q_data.Flag {
			SennserMoveReadData()
		}
		ch <- true
	}()

	go func() { //ch7
		SennserDataValue.Mu.Lock()
		SennserDataValue.CpuTmp = cpuTmp()
		SennserDataValue.Mu.Unlock()
		ch <- true
	}()
	for {
		if len(ch) == 7 {
			break
		}
		time.Sleep(time.Microsecond * 10)
	}

}

func Close() {
	if SennserData.Bme280_data.Flag {
		fmt.Println("Close BME280")
		SennserData.Bme280_data.Close()
	}
	if SennserData.Am2320_data.Flag {
		fmt.Println("Close AM2320")
		SennserData.Am2320_data.Close()
	}
	if SennserData.CO2Sensor_data.Flag {
		fmt.Println("Close CO2Sensor")
		SennserData.CO2Sensor_data.Close()
	}
	if SennserData.Tsl2561_data.Flag {
		fmt.Println("Close TSL2561")
		SennserData.Tsl2561_data.Close()
	}
	if SennserData.Mma8452q_data.Flag {
		fmt.Println("Close MMA8452Q")
		SennserData.Mma8452q_data.Close()
	}
	if SennserData.DhtSenser_data.Flag {
		fmt.Println("Close DhtSenser")
		SennserData.DhtSenser_data.Close()
	}

}
