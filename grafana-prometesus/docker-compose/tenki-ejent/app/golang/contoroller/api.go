package contoroller

import "tenkiej/contoroller/amedas"

type API struct {
}

func NewAPI() *API {
	return &API{}
}

func (a *API) GetAmedasMapData() []amedas.PrometesusData {
	return amedasAPI.PData
}
