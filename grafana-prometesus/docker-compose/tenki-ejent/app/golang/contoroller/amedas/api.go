package amedas

import (
	"context"
	"sync"
)

type API struct {
	mpdata    map[string]MapData
	tbdata    map[string]TableData
	PData     []PrometesusData
	PDMapData PrometesusDatas
	lttime    string
	mux       sync.Mutex
}

func NewAmedasAPI() *API {
	tmp := PrometesusDatas{
		Temp:             map[string]PrometesusData{},
		Humidity:         map[string]PrometesusData{},
		Snow:             map[string]PrometesusData{},
		Snow1h:           map[string]PrometesusData{},
		Snow6h:           map[string]PrometesusData{},
		Snow12h:          map[string]PrometesusData{},
		Snow24h:          map[string]PrometesusData{},
		Sun10m:           map[string]PrometesusData{},
		Sun1h:            map[string]PrometesusData{},
		Visibility:       map[string]PrometesusData{},
		Precipitation10m: map[string]PrometesusData{},
		Precipitation1h:  map[string]PrometesusData{},
		Precipitation3h:  map[string]PrometesusData{},
		Precipitation24h: map[string]PrometesusData{},
		WindDirection:    map[string]PrometesusData{},
		Wind:             map[string]PrometesusData{},
	}

	return &API{PDMapData: tmp}
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
	tmp := covertToePrometesusDatas(tmpTb, tmpMp)
	a.mux.Lock()
	a.PData = tmp
	a.mux.Unlock()
	return a.PData
}

func (a *API) GetPrometesusDatas() PrometesusDatas {
	a.mux.Lock()
	tmpTb := a.tbdata
	tmpMp := a.mpdata
	tmpMap := a.PDMapData
	a.mux.Unlock()
	tmpMap.covertToePrometesusDataMaps(tmpTb, tmpMp)
	a.mux.Lock()
	a.PDMapData = tmpMap
	a.mux.Unlock()
	return a.PDMapData
}

func (a *API) CleanData() {
	a.mux.Lock()
	tmpMap := a.PDMapData
	a.mux.Unlock()
	tmpMap = tmpMap.cleanData()
	a.mux.Lock()
	a.PDMapData = tmpMap
	a.mux.Unlock()
	return
}

func (a *API) CountSumMap() int {
	a.mux.Lock()
	tmpMap := a.PDMapData
	a.mux.Unlock()
	count := 0
	count += len(tmpMap.Temp)
	count += len(tmpMap.Humidity)
	count += len(tmpMap.Snow)
	count += len(tmpMap.Snow1h)
	count += len(tmpMap.Snow6h)
	count += len(tmpMap.Snow12h)
	count += len(tmpMap.Snow24h)
	count += len(tmpMap.Sun10m)
	count += len(tmpMap.Sun1h)
	count += len(tmpMap.Visibility)
	count += len(tmpMap.Precipitation10m)
	count += len(tmpMap.Precipitation1h)
	count += len(tmpMap.Precipitation3h)
	count += len(tmpMap.Precipitation24h)
	count += len(tmpMap.WindDirection)
	count += len(tmpMap.Wind)

	return count
}
