package msgsenser

import "time"

type Msg struct {
	Senser      string
	Message     string
	MessageTime time.Time
}

func (t *Msg) Create(str string) {
	t.Senser = str
}

func (t *Msg) Write(str string) {
	t.Message = str
	t.MessageTime = time.Now()
}

type HealthData struct {
	Sennserdata string `json:sennserdata`
	Message     string `json:message`
}

type SenserData struct {
	Senser string `json:"Senser"`
	Type   string `json:"Type"`
	Data   string `json:"Data"`
}
