package senser

import (
	"fmt"
	"strconv"
)

var (
	I2C_BUS = 1
)

type Sennser struct {
	Bme280_data Bme280
}

type SenserValue struct {
	Bme280 Bme280_Vaule
}

var SennserData Sennser = Sennser{}

var SennserDataValue SenserValue = SenserValue{
	Bme280: Bme280_Vaule{},
}

func SennserSetup() {
	SennserData.Bme280_data = Bme280{}

	if !SennserData.Bme280_data.Init() {
		fmt.Println("I2C not for BME280")
	}

	return
}

func SenserRead() {
	press, temp, hum := SennserData.Bme280_data.ReadData()
	SennserDataValue.Bme280.Press = strconv.FormatFloat(press, 'f', 2, 64)
	SennserDataValue.Bme280.Temp = strconv.FormatFloat(temp, 'f', 2, 64)
	SennserDataValue.Bme280.Hum = strconv.FormatFloat(hum, 'f', 2, 64)
}
