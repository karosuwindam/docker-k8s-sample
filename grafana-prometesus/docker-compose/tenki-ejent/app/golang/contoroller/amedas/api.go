package amedas

import (
	"context"
	"sync"
)

type API struct {
	mpdata map[string]MapData
	tbdata map[string]TableData
	PData  []PrometesusData
	lttime string
	mux    sync.Mutex
}

func NewAmedasAPI() *API {
	return &API{}
}

func (a *API) GetAmedasMapData(ctx context.Context) error {
	ltime, err := getLastTime(ctx)
	if err != nil {
		return err
	}
	if a.lttime == ltime {
		return nil
	}
	tmp, err := getJsonDataMapData(ctx, ltime)
	if err != nil {
		return err
	}
	a.mux.Lock()
	a.lttime = ltime
	a.mpdata = tmp
	a.mux.Unlock()
	return nil
}

func (a *API) GetAmedasTableData(ctx context.Context) error {
	tmp, err := getJsonTable(ctx)
	if err != nil {
		return err
	}
	a.mux.Lock()
	a.tbdata = tmp
	a.mux.Unlock()
	return nil
}

func (a *API) GetPrometesusData() []PrometesusData {
	a.mux.Lock()
	tmpTb := a.tbdata
	tmpMp := a.mpdata
	a.mux.Unlock()
	a.PData = covertToePrometesusDatas(tmpTb, tmpMp)
	return a.PData
}
