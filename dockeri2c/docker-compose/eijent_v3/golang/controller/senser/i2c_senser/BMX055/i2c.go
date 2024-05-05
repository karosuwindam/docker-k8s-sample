package bmx055

import (
	"github.com/davecheney/i2c"
	"github.com/pkg/errors"
)

var (
	I2C_BUS = 1
)

func writeByte(i2c_addr uint8, command, data byte) error {
	i2c, err := i2c.New(i2c_addr, I2C_BUS)
	if err != nil {
		return errors.Wrapf(err, "i2c.New(%v,%v)", i2c_addr, I2C_BUS)
	}
	defer i2c.Close()
	_, err = i2c.Write([]byte{command, data})
	if err != nil {
		return errors.Wrapf(err, "i2c.Write(%v,%v)", command, data)
	}
	return nil
}

func readByte(i2c_addr uint8, command byte, size int) ([]byte, error) {
	buf := make([]byte, size)
	i2c, err := i2c.New(i2c_addr, I2C_BUS)
	if err != nil {
		return buf, errors.Wrapf(err, "i2c.New(%v,%v)", i2c_addr, I2C_BUS)
	}
	defer i2c.Close()
	_, err = i2c.Write([]byte{command})
	if err != nil {
		return buf, errors.Wrapf(err, "i2c.Write(%v)", command)
	}
	_, err = i2c.Read(buf)
	if err != nil {
		return buf, errors.Wrap(err, "i2c.Read()")
	}
	return buf, nil

}
