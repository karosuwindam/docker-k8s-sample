package bmx055

import (
	msgsenser "eijent/controller/senser/msg_senser"
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
	_, ok := ReadValue()
	// out = append(out, msgsenser.SenserData{
	// 	Senser: memory.readMsg().Senser,
	// 	Type:   "lux",
	// 	Data:   strconv.Itoa(v),
	// })
	return out, ok
}

func (api *API) Reset() {
	ResetMessage()
}
