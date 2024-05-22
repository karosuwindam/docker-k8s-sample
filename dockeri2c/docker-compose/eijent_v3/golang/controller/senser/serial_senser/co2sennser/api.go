package co2sennser

import (
	msgsenser "eijent/controller/senser/msg_senser"
	"strconv"
)

type API struct{}

func NewAPI() *API {
	return &API{}
}

func (api *API) Run() error {
	return Run()
}

func (api *API) Init() error {
	return Init()
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
	if v.Co2 != -1 || v.Temp != -1 {
		out = append(out, msgsenser.SenserData{
			Senser: memory.readMsg().Senser,
			Type:   "co2",
			Data:   strconv.Itoa(v.Co2),
		})
		out = append(out, msgsenser.SenserData{
			Senser: memory.readMsg().Senser,
			Type:   "tmp",
			Data:   strconv.Itoa(v.Temp),
		})
	}
	return out, ok
}

func (api *API) Reset() {
	ResetMessage()
}
