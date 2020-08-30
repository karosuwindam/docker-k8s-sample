package main

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"strconv"
	"strings"
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
func cmdOutput(name string, arg ...string) string {
	out, err := exec.Command(name, arg...).Output()
	if err != nil {
		return ""
	}
	if out[len(out)-1] == 10 {
		out = out[:len(out)-1]
	}
	return string(out)
}
func tmpOut(data string) string {
	tmp := strings.Split(data, "=")
	return tmp[1][:len(tmp[1])-2]
}
func clockOutput(data string) string {
	tmp := strings.Split(data, "=")
	return tmp[1]
}
func voltOutput(data string) string {
	tmp := strings.Split(data, "=")
	return tmp[1][:len(tmp[1])-1]
}
func memOut(data string) string {
	tmp := strings.Split(data, "=")
	return tmp[1][:len(tmp[1])-1]
}
