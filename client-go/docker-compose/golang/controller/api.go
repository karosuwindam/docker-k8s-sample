package controller

import (
	"bytes"
	"encoding/json"
	"log/slog"
)

type api struct{}

func NewApi() *api {
	return &api{}
}

func (a *api) GetJson() string {
	d, err := getDomainData()
	if err != nil {
		slog.Error("GetJson", "error", err)
		return ""
	}
	tmpjson, _ := json.Marshal(d)
	var buf bytes.Buffer
	json.Indent(&buf, tmpjson, "", " ")
	return string(buf.Bytes())
}
