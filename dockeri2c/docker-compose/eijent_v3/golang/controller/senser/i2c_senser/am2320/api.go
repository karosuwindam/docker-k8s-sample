package am2320

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
	out = append(out, msgsenser.SenserData{
		Senser: memory.readMsg().Senser,
		Type:   "temp",
		Data:   strconv.FormatFloat(v.Temp, 'f', 1, 64),
	})

	out = append(out, msgsenser.SenserData{
		Senser: memory.readMsg().Senser,
		Type:   "hum",
		Data:   strconv.FormatFloat(v.Hum, 'f', 1, 64),
	})
	return out, ok
}

func (api *API) Reset() {
	ResetMessage()
}