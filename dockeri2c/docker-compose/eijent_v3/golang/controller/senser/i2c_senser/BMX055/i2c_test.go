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
	var ch chan bool = make(chan bool, 1)
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
		x := []float64{}
		y := []float64{}
		z := []float64{}
		for _, tm := range tmp {
			x = append(x, tm.X)
			y = append(y, tm.Y)
			z = append(z, tm.Z)
		}
		fmt.Println("ACC Read End")
		ch <- true
		fmt.Println("ACC-X:", x)
		fmt.Println("ACC ave:", average(x), "ACC med:", median(x))
		fmt.Println("ACC-Y:", y)
		fmt.Println("ACC ave:", average(y), "ACC med:", median(y))
		fmt.Println("ACC-Z:", z)
		fmt.Println("ACC ave:", average(z), "ACC med:", median(z))
		<-ch
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
		x := []float64{}
		y := []float64{}
		z := []float64{}
		for _, tm := range tmp {
			x = append(x, tm.X)
			y = append(y, tm.Y)
			z = append(z, tm.Z)
		}
		fmt.Println("GYRO Read End")
		ch <- true
		fmt.Println("GYRO-X:", x)
		fmt.Println("GYRO ave:", average(x), "GYRO med:", median(x))
		fmt.Println("GYRO-Y:", y)
		fmt.Println("GYRO ave:", average(y), "GYRO med:", median(y))
		fmt.Println("GYRO-Z:", z)
		fmt.Println("GYRO ave:", average(z), "GYRO med:", median(z))
		<-ch
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

		x := []int{}
		y := []int{}
		z := []int{}
		for _, tm := range tmp {
			x = append(x, tm.X)
			y = append(y, tm.Y)
			z = append(z, tm.Z)
		}
		fmt.Println("MAG Read End")
		ch <- true
		fmt.Println("MAG-X:", x)
		fmt.Println("MAG ave:", average(x), "MAG med:", median(x))
		fmt.Println("MAG-Y:", y)
		fmt.Println("MAG ave:", average(y), "MAG med:", median(y))
		fmt.Println("MAG-Z:", z)
		fmt.Println("MAG ave:", average(z), "MAG med:", median(z))
		<-ch
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
