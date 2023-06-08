package common

import (
	"fmt"
	"io/ioutil"
	"strconv"
)

var (
	CPU_TMP_PASS = "/sys/class/thermal/thermal_zone0/temp"
)

// RaspberryPiのCPU温度をファイルから読み込む
func useIoutilReadFile(fileName string) string {
	bytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		panic(err)
	}

	return string(bytes)
}

// RaspberryPiのCPU温度に変換する
func cpuTmpOutput(cputmp string) string {
	tmp := cputmp
	if tmp[len(tmp)-1] == 10 {
		tmp = tmp[:len(tmp)-1]
	}
	f, err := strconv.ParseFloat(tmp, 64)
	if err != nil {
		fmt.Println(err.Error())
		return cputmp
	}
	f = f / 1000
	return strconv.FormatFloat(f, 'f', 3, 64)
}

// RaspberryPiのCPU温度を取得する
func CpuTmp() string {
	return cpuTmpOutput(useIoutilReadFile(CPU_TMP_PASS))
}
