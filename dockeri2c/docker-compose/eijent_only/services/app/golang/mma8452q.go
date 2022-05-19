package main

import (
	"fmt"

	"github.com/davecheney/i2c"
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

type Mma8452q struct {
	Flag    bool
	Name    string
	Message string
}

func (t *Mma8452q) Init() {
	t.Name = "MMA8452Q"
	if t.up() != nil {
		t.Flag = false
	}
	t.Flag = true
}

func (t *Mma8452q) Close() {
	t.down()
}

func (t *Mma8452q) WriteByte(command, data byte) error {
	i2c, err := i2c.New(MMA8452Q, I2C_BUS)
	if err != nil {
		t.Message = err.Error()
		return err
	}
	defer i2c.Close()
	_, err = i2c.Write([]byte{command, data})
	if err != nil {
		t.Message = err.Error()
		return err
	}
	t.Message = "OK"
	return err
}

func (t *Mma8452q) ReadByte(command byte, size int) []byte {
	i2c, err := i2c.New(MMA8452Q, I2C_BUS)
	if err != nil {
		t.Message = err.Error()
		return []byte{}
	}
	defer i2c.Close()
	i2c.Write([]byte{command})
	if err != nil {
		t.Message = err.Error()
		return []byte{}
	}
	buf := make([]byte, size)
	i2c.Read(buf)
	if err != nil {
		t.Message = err.Error()
		return []byte{}
	}
	t.Message = "OK"
	return buf
}

func (t *Mma8452q) up() error {
	err := t.WriteByte(MMA8452Q_ADDR_XYZ_DATA_CFG, MMA8452Q_VAL_XYZ_DATA_CFG_FS_2G)
	if err != nil {
		return err
	}
	err = t.WriteByte(MMA8452Q_ADDR_CTRL_REG1, MMA8452Q_VAL_CTRL_REG1_ACTIVE)
	return err
}

func (t *Mma8452q) down() {
	t.WriteByte(MMA8452Q_ADDR_CTRL_REG1, MMA8452Q_VAL_CTRL_REG1_STANBY)
}

func (t *Mma8452q) ReadData() (int, int, int) {
	tmp := t.ReadByte(0x00, 7)
	tmp_data := []int{0, 0, 0}
	tmp_data[0] = int(tmp[MMA8452Q_ADDR_OUT_X_MSB])<<4 | int(tmp[MMA8452Q_ADDR_OUT_X_LSB]&0xf0)>>4
	tmp_data[1] = int(tmp[MMA8452Q_ADDR_OUT_Y_MSB])<<4 | int(tmp[MMA8452Q_ADDR_OUT_Y_LSB]&0xf0)>>4
	tmp_data[2] = int(tmp[MMA8452Q_ADDR_OUT_Z_MSB])<<4 | int(tmp[MMA8452Q_ADDR_OUT_Z_LSB]&0xf0)>>4
	for i := 0; i < len(tmp_data); i++ {
		if tmp_data[i] > 2048 {
			tmp_data[i] -= 4096
		}
	}
	x := tmp_data[0]
	y := tmp_data[1]
	z := tmp_data[2]

	return x, y, z
}

func (t *Mma8452q) Ch_data_to_accel(x, y, z int) (float64, float64, float64) {
	tmp := []float64{float64(x), float64(y), float64(z)}
	return tmp[0] / 512.0, tmp[1] / 512.0, tmp[2] / 512.0
}

func (t *Mma8452q) Test() { //未完成
	str := t.ReadByte(MMA8452Q_WHO_AM_I, 40)
	fmt.Println(str)

}
