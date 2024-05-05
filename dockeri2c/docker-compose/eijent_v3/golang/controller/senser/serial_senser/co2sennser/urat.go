package co2sennser

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/tarm/serial"
)

// uartInitConfig(v ...interface{}) (*serial.Config, error)
//
// input type:
// string, int,time.Duration,byte,serial.Parity,serial.StopBits
//
// Parity: ParityNone, ParityOdd, ParityEven, ParityMark, ParitySpace
//
// StopBits: Stop1 Stop1Half Stop2
func uartInitConfig(v ...interface{}) (*serial.Config, error) {
	out := serial.Config{}
	for _, d := range v {
		switch d.(type) {
		case string:
			out.Name = d.(string)
		case int:
			out.Baud = d.(int)
		case time.Duration:
			out.ReadTimeout = d.(time.Duration)
		case byte:
			out.Size = d.(byte)
		case serial.Parity:
			out.Parity = d.(serial.Parity)
		case serial.StopBits:
			out.StopBits = d.(serial.StopBits)
		default:
			msg := fmt.Sprintf("data = %T", d)
			return &out, errors.New(msg)
		}
	}
	return &out, nil
}

type UartSet struct {
	config   *serial.Config
	openflag bool
	port     *serial.Port
}

var uartdata UartSet

func uartInit(v ...interface{}) error {
	c, err := uartInitConfig(v)
	if err != nil {
		return errors.Wrap(err, "uartInitConfig")
	}
	uartdata = UartSet{
		config: c,
	}
	return nil
}

func (t *UartSet) open() error {
	var err error
	t.openflag = true
	t.port, err = serial.OpenPort(t.config)
	return err
}

func (t *UartSet) close() {
	if !t.openflag {
		return
	}
	t.openflag = false
	t.port.Close()
}

func (t *UartSet) read()  {}
func (t *UartSet) Write() {}
