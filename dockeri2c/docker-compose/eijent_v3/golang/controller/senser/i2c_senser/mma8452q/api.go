package mma8452q

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
	if len(v.Acc) != 0 {
		out = append(out, msgsenser.SenserData{
			Senser: memory.readMsg().Senser,
			Type:   "ax",
			Data:   strconv.FormatFloat(v.Acc[len(v.Acc)-1].X, 'f', 1, 64),
		})
		out = append(out, msgsenser.SenserData{
			Senser: memory.readMsg().Senser,
			Type:   "ay",
			Data:   strconv.FormatFloat(v.Acc[len(v.Acc)-1].Y, 'f', 1, 64),
		})
		out = append(out, msgsenser.SenserData{
			Senser: memory.readMsg().Senser,
			Type:   "az",
			Data:   strconv.FormatFloat(v.Acc[len(v.Acc)-1].Z, 'f', 1, 64),
		})
		out = append(out, msgsenser.SenserData{
			Senser: memory.readMsg().Senser,
			Type:   "zero_x",
			Data:   strconv.FormatFloat(v.Zero_X, 'f', 1, 64),
		})
		out = append(out, msgsenser.SenserData{
			Senser: memory.readMsg().Senser,
			Type:   "zero_y",
			Data:   strconv.FormatFloat(v.Zero_Y, 'f', 1, 64),
		})
		out = append(out, msgsenser.SenserData{
			Senser: memory.readMsg().Senser,
			Type:   "zero_z",
			Data:   strconv.FormatFloat(v.Zero_Z, 'f', 1, 64),
		})
	}
	return out, ok
}

func (api *API) Reset() {
	ResetMessage()
}
