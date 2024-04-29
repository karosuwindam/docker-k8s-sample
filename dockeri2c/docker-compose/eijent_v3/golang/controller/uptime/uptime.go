package uptime

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
)

const (
	uptimefile string = "/proc/uptime"
)

// Read() float64
//
// 電源起動後の時刻(s)
func Read() float64 {
	file, err := os.Open(uptimefile)
	if err != nil {
		log.Panic(err)
		return -1
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		tmp := strings.Split(line, " ")
		if lat, err := strconv.ParseFloat(tmp[0], 64); err == nil {
			return lat
		} else {
			log.Panic(err)
		}

	}
	return -1

}
