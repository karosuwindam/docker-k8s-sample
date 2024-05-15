package mma8452q

import (
	"eijent/config"
	msgsenser "eijent/controller/senser/msg_senser"
	"log"
	"math"
	"sort"
	"sync"
	"time"

	"github.com/pkg/errors"
)

const (
	MMA8452Q uint8 = 0x1d

	MMA8452Q_ADDR_STATUS       byte = 0x00
	MMA8452Q_ADDR_OUT_X_MSB    byte = 0x01
	MMA8452Q_ADDR_OUT_X_LSB    byte = 0x02
	MMA8452Q_ADDR_OUT_Y_MSB    byte = 0x03
	MMA8452Q_ADDR_OUT_Y_LSB    byte = 0x04
	MMA8452Q_ADDR_OUT_Z_MSB    byte = 0x05
	MMA8452Q_ADDR_OUT_Z_LSB    byte = 0x06
	MMA8452Q_WHO_AM_I          byte = 0x0D
	MMA8452Q_ADDR_XYZ_DATA_CFG byte = 0x0E
	MMA8452Q_ADDR_CTRL_REG1    byte = 0x2A
	MMA8452Q_ADDR_CTRL_REG2    byte = 0x2B
	MMA8452Q_ADDR_CTRL_REG3    byte = 0x2C
	MMA8452Q_ADDR_CTRL_REG4    byte = 0x2D
	MMA8452Q_ADDR_CTRL_REG5    byte = 0x2E
)

const (
	MMA8452Q_VAL_XYZ_DATA_CFG_FS_2G byte = 0
	MMA8452Q_VAL_XYZ_DATA_CFG_FS_4G byte = 1
	MMA8452Q_VAL_XYZ_DATA_CFG_FS_8G byte = 2

	MMA8452Q_DIVECE_ID byte = 0x2A

	MMA8452Q_VAL_CTRL_REG1_STANBY byte = 0
	MMA8452Q_VAL_CTRL_REG1_ACTIVE byte = 1
)

type Mma8452q_Vaule struct {
	Acc    []AccValue
	Zero_X float64
	Zero_Y float64
	Zero_Z float64
}

type Mma8452q_Raw struct {
	x int
	y int
	z int
}

type AccValue struct {
	X float64
	Y float64
	Z float64
}

const (
	SENSER_NAME string = "MMA8452Q"
	MAX_HOLED   int    = 1000
	AVE_COUNT   int    = 200
)

func Init(i2cMu *sync.Mutex) error {
	memory = datastore{
		values:   []Mma8452q_Raw{},
		Flag:     false,
		StopFlag: false,
		msg:      msgsenser.Msg{},
		i2cMu:    i2cMu,
	}
	shudown = make(chan bool, 1)
	done = make(chan bool, 1)
	reset = make(chan bool, 1)
	wait = make(chan bool, 1)
	busy = make(chan bool, 1)
	memory.msg.Create(SENSER_NAME)
	for i := 0; i < 3; i++ {
		if Test() {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}
	if !memory.Flag {
		return errors.New("not Init Error for " + SENSER_NAME)
	}

	return nil
}

func Run() error {
	memory.chageRunFlag(true)
	log.Println("info:", SENSER_NAME+" loop start")
	var readone chan bool = make(chan bool, 1)
	if memory.readFlag() {
		readone <- true
		up()
	}
loop:
	for {
		select {
		case <-reset:
			down()
			for i := 0; i < 3; i++ {
				up()
				if Test() {
					break
				}
				down()
				time.Sleep(100 * time.Millisecond)
			}
		case <-shudown:
			done <- true
			break loop
		case <-wait:
			tmp, _ := ReadValue()
			if tmp.Zero_X == 0 && tmp.Zero_Y == 0 && tmp.Zero_Z == 0 {
				readdate()
				busy <- true
			} else {
				done <- true
			}

		case <-readone:
			if memory.readFlag() {
				readdate()
			}
		case <-time.After(time.Duration(config.Senser.MMA8452Q_Count) * time.Microsecond):
			if memory.readFlag() {
				readdate()
			}
		}
	}
	down()
	memory.changeFlag(false)
	log.Println("info:", SENSER_NAME+" loop stop")

	return nil
}

func Stop() error {
	shudown <- true
	memory.chageRunFlag(false)
	select {
	case <-done:
		break
	case <-time.After(1 * time.Second):
		msg := "shutdown time out"
		memory.changeMsg(msg)
		return errors.New(msg)
	}
	memory.changeMsg("shutdown")
	return nil
}

func Health() (bool, msgsenser.Msg) {
	return memory.readFlag(), memory.readMsg()
}

func Wait() {

	wait <- true
	select {
	case <-busy:
		time.Sleep(10 * time.Millisecond)
		wait <- true
	case <-done:
		break
	case <-time.After(1 * time.Second):
		log.Println("error:", "time over 1 sec")
	}
}

func ReadValue() (Mma8452q_Vaule, bool) {
	tmp, ok := memory.readValue().(Mma8452q_Vaule)
	if !ok {
		tmp = Mma8452q_Vaule{
			Acc:    []AccValue{},
			Zero_X: 0,
			Zero_Y: 0,
			Zero_Z: 0,
		}
	}
	return tmp, memory.readFlag()
}

func ResetMessage() {
	if len(reset) > 0 {
		return
	}
	reset <- true
}

func readdate() {
	tmp := readSenserData()
	if tmp.x == 0 && tmp.y == 0 && tmp.z == 0 {
		return
	}
	vv := tmp.ChangeToAccell()
	v, ok := memory.readValue().(Mma8452q_Vaule)
	if ok {
		v.Acc = append(v.Acc, vv)
		if len(v.Acc) > MAX_HOLED {
			v.Acc = v.Acc[1:]
		}
	} else {
		v = Mma8452q_Vaule{
			Acc:    []AccValue{vv},
			Zero_X: 0,
			Zero_Y: 0,
			Zero_Z: 0,
		}
	}
	if len(v.Acc) > AVE_COUNT {

		avg := average_zero(v.Acc[len(v.Acc)-200:])
		mid := median_zero(v.Acc[len(v.Acc)-200:])
		flag := true
		for i := 0; i < len(avg); i++ {
			if math.Abs(avg[i])*0.95 > math.Abs(mid[i]) || math.Abs(mid[i]) > math.Abs(avg[i])*1.05 {
				flag = false
				break
			}
		}
		if flag {
			v.Zero_X = avg[0]
			v.Zero_Y = avg[1]
			v.Zero_Z = avg[2]
		}
	}
	memory.changeValue(v)

}

func Test() bool {
	flag := false
	if err := up(); err != nil {
		log.Println("error:", err)
	}
	tmp := readSenserData()
	if tmp.x == 0 && tmp.y == 0 && tmp.z == 0 {
		log.Println("error:", SENSER_NAME+" Not Fontd")
	} else {
		flag = true
	}
	if flag {
		for i := 0; i < 20; i++ { //初回読み込み
			readdate()
		}
	}
	down()
	return flag
}

func up() error {
	if err := writeByte(MMA8452Q_ADDR_XYZ_DATA_CFG, MMA8452Q_VAL_XYZ_DATA_CFG_FS_2G); err != nil {
		return errors.Wrapf(err, "writeByte(%x,%x)", MMA8452Q_ADDR_XYZ_DATA_CFG, MMA8452Q_VAL_XYZ_DATA_CFG_FS_2G)
	}
	if err := writeByte(MMA8452Q_ADDR_CTRL_REG1, MMA8452Q_VAL_CTRL_REG1_ACTIVE); err != nil {
		return errors.Wrapf(err, "writeByte(%x,%x)", MMA8452Q_ADDR_CTRL_REG1, MMA8452Q_VAL_CTRL_REG1_ACTIVE)
	}
	return nil
}

func down() error {
	if err := writeByte(MMA8452Q_ADDR_CTRL_REG1, MMA8452Q_VAL_CTRL_REG1_STANBY); err != nil {
		return errors.Wrapf(err, "writeByte(%x,%x)", MMA8452Q_ADDR_CTRL_REG1, MMA8452Q_VAL_CTRL_REG1_STANBY)
	}
	return nil
}

func readSenserData() Mma8452q_Raw {
	out := Mma8452q_Raw{0, 0, 0}
	if b, err := readByte(0x00, 7); err != nil {
		log.Println("error:", err)
		memory.changeMsg(err.Error())
	} else {
		memory.changeMsg("OK")
		tmp_data := []int{0, 0, 0}
		tmp_data[0] = int(b[MMA8452Q_ADDR_OUT_X_MSB])<<4 | int(b[MMA8452Q_ADDR_OUT_X_LSB]&0xf0)>>4
		tmp_data[1] = int(b[MMA8452Q_ADDR_OUT_Y_MSB])<<4 | int(b[MMA8452Q_ADDR_OUT_Y_LSB]&0xf0)>>4
		tmp_data[2] = int(b[MMA8452Q_ADDR_OUT_Z_MSB])<<4 | int(b[MMA8452Q_ADDR_OUT_Z_LSB]&0xf0)>>4
		for i := 0; i < len(tmp_data); i++ {
			if tmp_data[i] > 2048 {
				tmp_data[i] -= 4096
			}
		}
		out = Mma8452q_Raw{tmp_data[0], tmp_data[1], tmp_data[2]}
	}
	return out
}

// 読み取ったx,y,zの値を加速度へ変換する
func (t *Mma8452q_Raw) ChangeToAccell() AccValue {
	tmp := []float64{float64(t.x), float64(t.y), float64(t.z)}
	return AccValue{tmp[0] / 512.0, tmp[1] / 512.0, tmp[2] / 512.0}
}

func average_zero(nums []AccValue) []float64 {
	out := []float64{0, 0, 0}
	for i, num := range nums {
		out[0] = (out[0]*float64(i) + num.X) / float64(i+1)
		out[1] = (out[1]*float64(i) + num.Y) / float64(i+1)
		out[2] = (out[2]*float64(i) + num.Z) / float64(i+1)
	}
	return out
}

func median_zero(nums []AccValue) []float64 {
	out := []float64{0, 0, 0}
	var tmps_x []float64
	var tmps_y []float64
	var tmps_z []float64
	for _, num := range nums {
		tmps_x = append(tmps_x, float64(num.X))
		tmps_y = append(tmps_y, float64(num.Y))
		tmps_z = append(tmps_z, float64(num.Z))

	}

	if len(tmps_x) == 0 {
		return out
	} else {
		if len(tmps_x) == 1 {
			out[0] = tmps_x[0]
			out[1] = tmps_y[0]
			out[2] = tmps_z[0]
			return out
		}
		sort.Slice(tmps_x, func(i, j int) bool { return tmps_x[i] < tmps_x[j] })
		sort.Slice(tmps_y, func(i, j int) bool { return tmps_y[i] < tmps_y[j] })
		sort.Slice(tmps_z, func(i, j int) bool { return tmps_z[i] < tmps_z[j] })
		if (len(tmps_x)/2)*2 != len(tmps_x) {
			out[0] = tmps_x[(len(tmps_x) / 2)]
			out[1] = tmps_y[(len(tmps_y) / 2)]
			out[2] = tmps_z[(len(tmps_z) / 2)]
		} else {
			out[0] = (tmps_x[(len(tmps_x)/2)] + tmps_x[(len(tmps_x)/2+1)]) / 2
			out[1] = (tmps_y[(len(tmps_y)/2)] + tmps_y[(len(tmps_y)/2+1)]) / 2
			out[2] = (tmps_z[(len(tmps_z)/2)] + tmps_z[(len(tmps_z)/2+1)]) / 2
		}
	}
	return out
}
