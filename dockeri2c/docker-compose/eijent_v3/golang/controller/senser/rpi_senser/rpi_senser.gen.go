package rpisenser

import (
	"context"
	"errors"
	"log"
	"os"
	"strconv"
	"sync"
	"time"
)

var (
	CPU_TMP_PASS = "/sys/class/thermal/thermal_zone0/temp"
)

// 記録用
type rpi_data struct {
	temp []float64
	mu   sync.Mutex
}

// 出力用
type Rpi struct {
	Temp string
}

var datastore rpi_data
var shutdown chan bool
var done chan bool
var wait chan bool

func Init() error {
	if _, err := os.Stat(CPU_TMP_PASS); err != nil {
		return err
	}
	datastore.temp = []float64{}
	shutdown = make(chan bool, 1)
	done = make(chan bool, 1)
	wait = make(chan bool, 1)
	return nil
}

func Run(ctx context.Context) error {
	log.Println("info:", "start rpi read loop")
	if tmp, err := cpu_temp_read(); err != nil { //CPUの温度読み取り
		log.Println("error:", err)
	} else {
		datastore.temp_Add(tmp)
	}
loop:
	for {
		select {
		case <-ctx.Done():
			return errors.New("ctx done")
		case <-shutdown:
			done <- true
			break loop
		case <-wait:
			done <- true
		case <-time.After(500 * time.Millisecond): //500msあと
			if tmp, err := cpu_temp_read(); err != nil { //CPUの温度読み取り
				log.Println("error:", err)
			} else {
				datastore.temp_Add(tmp)
			}
		}
	}
	log.Println("info:", "stop rpi read loop")

	return nil
}

func Stop() error {
	if len(shutdown) > 0 {
		return nil
	}
	shutdown <- true
	select {
	case <-done:
		break
	case <-time.After(1 * time.Second):
		return errors.New("time over 1 sec")
	}
	return nil
}

func Wait() {
	wait <- true
	select {
	case <-done:
		break
	case <-time.After(1 * time.Second):
		log.Panic("time over 1 sec")
	}
}

func ReadNow() Rpi {
	datastore.mu.Lock()
	defer datastore.mu.Unlock()
	temp_f := datastore.temp[len(datastore.temp)-1]
	return Rpi{Temp: strconv.FormatFloat(temp_f, 'f', 3, 64)}
}

func ReadAve() Rpi {
	datastore.mu.Lock()
	temp_ary_f := datastore.temp
	datastore.mu.Unlock()
	var temp_f float64 = 0
	for i, temp := range temp_ary_f {
		temp_f = temp_f*float64(i) + temp
		temp_f = temp_f / float64(i+1)
	}
	return Rpi{Temp: strconv.FormatFloat(temp_f, 'f', 3, 64)}
}

func (t *rpi_data) temp_Add(f float64) {
	t.mu.Lock()
	defer t.mu.Unlock()
	tmp := t.temp
	tmp = append(tmp, f)
	if len(tmp) > 100 {
		tmp = tmp[1:]
	}
	t.temp = tmp
}
