package bmx055

import (
	"fmt"
	"log"
	"sort"
	"sync"
	"testing"
	"time"
)

func TestRead(t *testing.T) {
	var i2cMu sync.Mutex
	Init(&i2cMu)
	Test()
	accInit()
	var wg sync.WaitGroup
	wg.Add(3)
	go func() {
		tmp := []ACCAxis{}
		defer wg.Done()
		for i := 0; i < 1000; i++ {
			if d, err := getACCRAW(); err != nil {
				log.Println("error:", err)
			} else {
				num := axis_to_ACC(d)
				tmp = append(tmp, num)
			}
			time.Sleep(time.Millisecond * 10)
		}
		fmt.Println("ACC:", tmp)
		fmt.Println("ACC ave:", average(tmp), "ACC med:", median(tmp))
	}()
	go func() {
		tmp := []GyroAxis{}
		defer wg.Done()
		for i := 0; i < 1000; i++ {
			if d, err := getGyroRAW(); err != nil {
				log.Println("error:", err)
			} else {
				num := axis_to_Gyro(d)
				tmp = append(tmp, num)
			}
			time.Sleep(time.Millisecond * 10)
		}
		fmt.Println("GYRO:", tmp)
		fmt.Println("GYRO ave:", average(tmp), "GYRO med:", median(tmp))
	}()
	go func() {
		tmp := []Axis{}
		defer wg.Done()
		for i := 0; i < 1000; i++ {
			if d, err := getMag(); err != nil {
				log.Println("error:", err)
			} else {
				tmp = append(tmp, d)
			}
			time.Sleep(time.Millisecond * 10)
		}
		fmt.Println("MAG:", tmp)
		fmt.Println("MAG ave:", average(tmp), "MAG med:", median(tmp))
	}()
	wg.Wait()
}

func average(nums interface{}) float64 {
	out := 0.0
	switch nums.(type) {
	case []int:
		for i, num := range nums.([]int) {
			tmp := num
			out = (out*float64(i) + float64(tmp)) / float64(i+1)
		}
	case []float64:
		for i, num := range nums.([]float64) {
			tmp := num
			out = (out*float64(i) + float64(tmp)) / float64(i+1)
		}
	}
	return out
}

func median(nums interface{}) float64 {
	out := 0.0
	var tmps []float64
	switch nums.(type) {
	case []int:
		for _, num := range nums.([]int) {
			tmp := num
			tmps = append(tmps, float64(tmp))
		}
	case []float64:
		for _, num := range nums.([]float64) {
			tmp := num
			tmps = append(tmps, tmp)
		}
	}
	if len(tmps) == 0 {
		return 0
	} else {
		if len(tmps) == 1 {
			return tmps[0]
		}
		sort.Slice(tmps, func(i, j int) bool { return tmps[i] < tmps[j] })
		if (len(tmps)/2)*2 != len(tmps) {
			out = tmps[(len(tmps) / 2)]
		} else {
			out = (tmps[(len(tmps)/2)] + tmps[(len(tmps)/2+1)]) / 2

		}
	}
	return out
}
