package bmx055

import (
	msgsenser "eijent/controller/senser/msg_senser"
	"strconv"
	"sync"
)

type API struct{}

func NewAPI() *API {
	return &API{}
}

func (api *API) Run() error {
	return Run()
}

func (api *API) Init(i2cMu *sync.Mutex) error {
	return Init(i2cMu)
}

func (api *API) Stop() error {
	return Stop()
}

func (api *API) Health() (bool, msgsenser.Msg) {
	return Health()
}

func (api *API) Wait() {
	Wait()
}

func (api *API) Name() string {
	msg := memory.readMsg()
	return msg.Senser
}

func (api *API) Read() ([]msgsenser.SenserData, bool) {
	var out []msgsenser.SenserData
	v, ok := ReadValue()
	if len(v.acc) != 0 {
		out = append(out, msgsenser.SenserData{
			Senser: memory.readMsg().Senser,
			Type:   "ax",
			Data:   strconv.FormatFloat(v.acc[len(v.acc)-1].X, 'f', 1, 64),
		})
		out = append(out, msgsenser.SenserData{
			Senser: memory.readMsg().Senser,
			Type:   "ay",
			Data:   strconv.FormatFloat(v.acc[len(v.acc)-1].Y, 'f', 1, 64),
		})
		out = append(out, msgsenser.SenserData{
			Senser: memory.readMsg().Senser,
			Type:   "az",
			Data:   strconv.FormatFloat(v.acc[len(v.acc)-1].Z, 'f', 1, 64),
		})
		out = append(out, msgsenser.SenserData{
			Senser: memory.readMsg().Senser,
			Type:   "zero_x",
			Data:   strconv.FormatFloat(v.acc_zero.X, 'f', 1, 64),
		})
		out = append(out, msgsenser.SenserData{
			Senser: memory.readMsg().Senser,
			Type:   "zero_y",
			Data:   strconv.FormatFloat(v.acc_zero.Y, 'f', 1, 64),
		})
		out = append(out, msgsenser.SenserData{
			Senser: memory.readMsg().Senser,
			Type:   "zero_z",
			Data:   strconv.FormatFloat(v.acc_zero.Z, 'f', 1, 64),
		})
		out = append(out, msgsenser.SenserData{
			Senser: memory.readMsg().Senser,
			Type:   "gx",
			Data:   strconv.FormatFloat(v.gyro[len(v.gyro)-1].X, 'f', 1, 64),
		})
		out = append(out, msgsenser.SenserData{
			Senser: memory.readMsg().Senser,
			Type:   "gy",
			Data:   strconv.FormatFloat(v.gyro[len(v.gyro)-1].Y, 'f', 1, 64),
		})
		out = append(out, msgsenser.SenserData{
			Senser: memory.readMsg().Senser,
			Type:   "gz",
			Data:   strconv.FormatFloat(v.gyro[len(v.gyro)-1].Z, 'f', 1, 64),
		})
		out = append(out, msgsenser.SenserData{
			Senser: memory.readMsg().Senser,
			Type:   "zero_gx",
			Data:   strconv.FormatFloat(v.gyro_zero.X, 'f', 1, 64),
		})
		out = append(out, msgsenser.SenserData{
			Senser: memory.readMsg().Senser,
			Type:   "zero_gy",
			Data:   strconv.FormatFloat(v.gyro_zero.Y, 'f', 1, 64),
		})
		out = append(out, msgsenser.SenserData{
			Senser: memory.readMsg().Senser,
			Type:   "zero_gz",
			Data:   strconv.FormatFloat(v.gyro_zero.Z, 'f', 1, 64),
		})
		out = append(out, msgsenser.SenserData{
			Senser: memory.readMsg().Senser,
			Type:   "mx",
			Data:   strconv.Itoa(v.mag[len(v.mag)-1].X),
		})
		out = append(out, msgsenser.SenserData{
			Senser: memory.readMsg().Senser,
			Type:   "my",
			Data:   strconv.Itoa(v.mag[len(v.mag)-1].Y),
		})
		out = append(out, msgsenser.SenserData{
			Senser: memory.readMsg().Senser,
			Type:   "mz",
			Data:   strconv.Itoa(v.mag[len(v.mag)-1].Z),
		})
	}
	return out, ok
}

func (api *API) Reset() {
	ResetMessage()
}
