package senser

import (
	"fmt"
	"io/ioutil"
	"strconv"
)

var (
	CPU_TMP_PASS = "/sys/class/thermal/thermal_zone0/temp"
)

func useIoutilReadFile(fileName string) string {
	bytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		panic(err)
	}

	return string(bytes)
}

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
func cpuTmp() string {
	return cpuTmpOutput(useIoutilReadFile(CPU_TMP_PASS))
}
