package amedas

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"sync"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

const (
	LASTTIME_URL      = "https://www.jma.go.jp/bosai/amedas/data/latest_time.txt"     //2021-02-27T16:40:00+09:00
	MAP_DATA_URL      = "https://www.jma.go.jp/bosai/amedas/data/map/"                //+20210227164000.jsonで取得可能{lasttime}.json
	AMEDAS_TABLE_URL  = "https://www.jma.go.jp/bosai/amedas/const/amedastable.json"   //IDに対応した名称テーブル
	AMEDAS_ELEMNT_URL = "https://www.jma.go.jp/bosai/const/selectorinfos/amedas.json" //要素テーブル
	CLEAN_TIME        = 1 * time.Hour
)

type MapData struct {
	Temp             []float64 //気温
	Humidity         []int     //湿度
	Snow             []float64 //積雪
	Snow1h           []float64 //"積雪深
	Snow6h           []float64
	Snow12h          []float64
	Snow24h          []float64
	Sun10m           []int //日照時間
	Sun1h            []float64
	Visibility       []int     //視程
	Precipitation10m []float64 //降水量
	Precipitation1h  []float64
	Precipitation3h  []float64
	Precipitation24h []float64
	WindDirection    []int     //風向 16方位 0:北 1:北北東 2:北東 3:東北東 4:東 5:東南東 6:南東 7:南南東 8:南 9:南南西 10:南西 11:西南西 12:西 13:西北西 14:北西 15:北北西
	Wind             []float64 //風速
}

type TableData struct {
	Type   string //地点種別
	Elems  string //要素
	Lat    []float64
	Lon    []float64
	alt    int    //高度
	KjName string //地点名
	KnName string //地点名
	EnName string //地点名
}

type PrometesusData struct {
	Name     string    `json:"Name"`    //要素名
	EName    string    `json:"EName"`   //要素名
	Value    string    `json:"Value"`   //値
	Elems    string    `json:"Elems"`   //要素
	Type     string    `json:"Type"`    //地点種別
	NumberId string    `json:"Nuber"`   //地点ID
	KjName   string    `json:"KjName"`  //地点名
	KnName   string    `json:"KnName"`  //地点名
	EnName   string    `json:"EnName"`  //地点名
	AddTime  time.Time `json:"AddTime"` //追加時間
}

// 　保管用のPrometesusDataの構造体
type PrometesusDatas struct {
	Temp             map[string]PrometesusData //気温
	Humidity         map[string]PrometesusData //湿度
	Snow             map[string]PrometesusData //積雪
	Snow1h           map[string]PrometesusData //積雪深 1h
	Snow6h           map[string]PrometesusData //積雪深 6h
	Snow12h          map[string]PrometesusData //積雪深 12h
	Snow24h          map[string]PrometesusData //積雪深 24h
	Sun10m           map[string]PrometesusData //日照時間 10m
	Sun1h            map[string]PrometesusData //日照時間 1h
	Visibility       map[string]PrometesusData //視程
	Precipitation10m map[string]PrometesusData //降水量 10m
	Precipitation1h  map[string]PrometesusData //降水量 1h
	Precipitation3h  map[string]PrometesusData //降水量 3h
	Precipitation24h map[string]PrometesusData //降水量 24h
	WindDirection    map[string]PrometesusData //風向
	Wind             map[string]PrometesusData //風速
}

func getLastTime(ctx context.Context) (string, error) {

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, LASTTIME_URL, nil)
	// client := new(http.Client)
	client := http.Client{
		Transport: otelhttp.NewTransport(
			http.DefaultTransport,
			otelhttp.WithMessageEvents(otelhttp.ReadEvents, otelhttp.WriteEvents),
		),
	}

	resp, err := client.Do(req)

	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	t, _ := time.Parse("2006-01-02T15:04:05-07:00", string(b))

	return t.Format("20060102150405"), nil
}

// 全国のアメダステーブルを取得
func getJsonTable(ctx context.Context) (map[string]TableData, error) {

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, AMEDAS_TABLE_URL, nil)
	// client := new(http.Client)
	client := http.Client{
		Transport: otelhttp.NewTransport(
			http.DefaultTransport,
			otelhttp.WithMessageEvents(otelhttp.ReadEvents, otelhttp.WriteEvents),
		),
	}

	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	out := map[string]TableData{}
	json.Unmarshal(b, &out)
	return out, nil

}

// 全国のアメダスデータを取得
func getJsonDataMapData(ctx context.Context, lasttime string) (map[string]MapData, error) {

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, MAP_DATA_URL+lasttime+".json", nil)
	// client := new(http.Client)
	client := http.Client{
		Transport: otelhttp.NewTransport(
			http.DefaultTransport,
			otelhttp.WithMessageEvents(otelhttp.ReadEvents, otelhttp.WriteEvents),
		),
	}

	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	out := map[string]MapData{}
	json.Unmarshal(b, &out)
	return out, nil

}

func covertToePrometesusDatas(td map[string]TableData, md map[string]MapData) []PrometesusData {
	var out []PrometesusData
	wg := sync.WaitGroup{}
	mu := sync.Mutex{}
	waitch := make(chan struct{}, 10)
	if td == nil || md == nil {
		return out
	} else if len(td) == 0 {
		return out
	}
	for d, v := range md {
		if k, ok := td[d]; ok {
			wg.Add(1)
			waitch <- struct{}{}
			go func(d string, k TableData, v MapData) {
				defer func() {
					<-waitch
					wg.Done()
				}()
				tmp := covertToePrometesusData(d, k, v)
				mu.Lock()
				out = append(out, tmp...)
				mu.Unlock()
			}(d, k, v)
		} else {
			k = TableData{}
			wg.Add(1)
			waitch <- struct{}{}
			go func(d string, k TableData, v MapData) {
				defer func() {
					<-waitch
					wg.Done()
				}()
				tmp := covertToePrometesusData(d, k, v)
				mu.Lock()
				out = append(out, tmp...)
				mu.Unlock()
			}(d, k, v)
		}
	}
	wg.Wait()
	sort.Slice(out, func(i, j int) bool {
		return out[i].EName < out[j].EName
	})
	return out
}

func (out *PrometesusDatas) covertToePrometesusDataMaps(td map[string]TableData, md map[string]MapData) {
	wg := sync.WaitGroup{}
	waitch := make(chan struct{}, 10)
	mu := sync.Mutex{}
	if td == nil || md == nil {
		return
	} else if len(td) == 0 {
		return
	}
	for d, v := range md {
		if k, ok := td[d]; ok {
			wg.Add(1)
			waitch <- struct{}{}
			go func(d string, k TableData, v MapData) {
				defer func() {
					<-waitch
					wg.Done()
				}()
				tmp := covertToePrometesusDataMap(d, k, v)
				mu.Lock()
				out.Temp[d] = tmp.Temp[d]
				out.Humidity[d] = tmp.Humidity[d]
				out.Snow[d] = tmp.Snow[d]
				out.Snow1h[d] = tmp.Snow1h[d]
				out.Snow6h[d] = tmp.Snow6h[d]
				out.Snow12h[d] = tmp.Snow12h[d]
				out.Snow24h[d] = tmp.Snow24h[d]
				out.Sun10m[d] = tmp.Sun10m[d]
				out.Sun1h[d] = tmp.Sun1h[d]
				out.Visibility[d] = tmp.Visibility[d]
				out.Precipitation10m[d] = tmp.Precipitation10m[d]
				out.Precipitation1h[d] = tmp.Precipitation1h[d]
				out.Precipitation3h[d] = tmp.Precipitation3h[d]
				out.Precipitation24h[d] = tmp.Precipitation24h[d]
				out.WindDirection[d] = tmp.WindDirection[d]
				out.Wind[d] = tmp.Wind[d]
				mu.Unlock()
			}(d, k, v)
		} else {
			k = TableData{}
			wg.Add(1)
			waitch <- struct{}{}
			go func(d string, k TableData, v MapData) {
				defer func() {
					<-waitch
					wg.Done()
				}()
				tmp := covertToePrometesusDataMap(d, k, v)
				mu.Lock()
				out.Temp[d] = tmp.Temp[d]
				out.Humidity[d] = tmp.Humidity[d]
				out.Snow[d] = tmp.Snow[d]
				out.Snow1h[d] = tmp.Snow1h[d]
				out.Snow6h[d] = tmp.Snow6h[d]
				out.Snow12h[d] = tmp.Snow12h[d]
				out.Snow24h[d] = tmp.Snow24h[d]
				out.Sun10m[d] = tmp.Sun10m[d]
				out.Sun1h[d] = tmp.Sun1h[d]
				out.Visibility[d] = tmp.Visibility[d]
				out.Precipitation10m[d] = tmp.Precipitation10m[d]
				out.Precipitation1h[d] = tmp.Precipitation1h[d]
				out.Precipitation3h[d] = tmp.Precipitation3h[d]
				out.Precipitation24h[d] = tmp.Precipitation24h[d]
				out.WindDirection[d] = tmp.WindDirection[d]
				out.Wind[d] = tmp.Wind[d]
				mu.Unlock()
			}(d, k, v)
		}
	}
	wg.Wait()
	return
}

func covertToePrometesusData(d string, k TableData, v MapData) []PrometesusData {
	var out []PrometesusData
	timeNow := time.Now()
	if len(v.Temp) > 0 {
		out = append(out, PrometesusData{
			Name:     "気温",
			EName:    "Temp",
			Value:    fmt.Sprintf("%v", v.Temp[0]),
			Elems:    k.Elems,
			Type:     k.Type,
			NumberId: d,
			KjName:   k.KjName,
			KnName:   k.KnName,
			EnName:   k.EnName,
			AddTime:  timeNow,
		})
	}
	if len(v.Humidity) > 0 {
		out = append(out, PrometesusData{
			Name:     "湿度",
			EName:    "Humidity",
			Value:    fmt.Sprintf("%v", v.Humidity[0]),
			Elems:    k.Elems,
			Type:     k.Type,
			NumberId: d,
			KjName:   k.KjName,
			KnName:   k.KnName,
			EnName:   k.EnName,
			AddTime:  timeNow,
		})
	}
	if len(v.Snow) > 0 {
		out = append(out, PrometesusData{
			Name:     "積雪",
			EName:    "Snow",
			Value:    fmt.Sprintf("%v", v.Snow[0]),
			Elems:    k.Elems,
			Type:     k.Type,
			NumberId: d,
			KjName:   k.KjName,
			KnName:   k.KnName,
			EnName:   k.EnName,
			AddTime:  timeNow,
		})
	}
	if len(v.Snow1h) > 0 {
		out = append(out, PrometesusData{
			Name:     "積雪深 1h",
			EName:    "Snow1h",
			Value:    fmt.Sprintf("%v", v.Snow1h[0]),
			Elems:    k.Elems,
			Type:     k.Type,
			NumberId: d,
			KjName:   k.KjName,
			KnName:   k.KnName,
			EnName:   k.EnName,
			AddTime:  timeNow,
		})
	}
	if len(v.Snow6h) > 0 {
		out = append(out, PrometesusData{
			Name:     "積雪深 6h",
			EName:    "Snow6h",
			Value:    fmt.Sprintf("%v", v.Snow6h[0]),
			Elems:    k.Elems,
			Type:     k.Type,
			NumberId: d,
			KjName:   k.KjName,
			KnName:   k.KnName,
			EnName:   k.EnName,
			AddTime:  timeNow,
		})
	}
	if len(v.Snow12h) > 0 {
		out = append(out, PrometesusData{
			Name:     "積雪深 12h",
			EName:    "Snow12h",
			Value:    fmt.Sprintf("%v", v.Snow12h[0]),
			Elems:    k.Elems,
			Type:     k.Type,
			NumberId: d,
			KjName:   k.KjName,
			KnName:   k.KnName,
			EnName:   k.EnName,
			AddTime:  timeNow,
		})
	}
	if len(v.Snow24h) > 0 {
		out = append(out, PrometesusData{
			Name:     "積雪深 24h",
			EName:    "Snow24h",
			Value:    fmt.Sprintf("%v", v.Snow24h[0]),
			Elems:    k.Elems,
			Type:     k.Type,
			NumberId: d,
			KjName:   k.KjName,
			KnName:   k.KnName,
			EnName:   k.EnName,
			AddTime:  timeNow,
		})
	}
	if len(v.Sun10m) > 0 {
		out = append(out, PrometesusData{
			Name:     "日照時間 10m",
			EName:    "Sun10m",
			Value:    fmt.Sprintf("%v", v.Sun10m[0]),
			Elems:    k.Elems,
			Type:     k.Type,
			NumberId: d,
			KjName:   k.KjName,
			KnName:   k.KnName,
			EnName:   k.EnName,
			AddTime:  timeNow,
		})
	}
	if len(v.Sun1h) > 0 {
		out = append(out, PrometesusData{
			Name:     "日照時間 1h",
			EName:    "Sun1h",
			Value:    fmt.Sprintf("%v", v.Sun1h[0]),
			Elems:    k.Elems,
			Type:     k.Type,
			NumberId: d,
			KjName:   k.KjName,
			KnName:   k.KnName,
			EnName:   k.EnName,
			AddTime:  timeNow,
		})
	}
	if len(v.Visibility) > 0 {
		out = append(out, PrometesusData{
			Name:     "視程",
			EName:    "Visibility",
			Value:    fmt.Sprintf("%v", v.Visibility[0]),
			Elems:    k.Elems,
			Type:     k.Type,
			NumberId: d,
			KjName:   k.KjName,
			KnName:   k.KnName,
			EnName:   k.EnName,
			AddTime:  timeNow,
		})
	}
	if len(v.Precipitation10m) > 0 {
		out = append(out, PrometesusData{
			Name:     "降水量 10m",
			EName:    "Precipitation10m",
			Value:    fmt.Sprintf("%v", v.Precipitation10m[0]),
			Elems:    k.Elems,
			Type:     k.Type,
			NumberId: d,
			KjName:   k.KjName,
			KnName:   k.KnName,
			EnName:   k.EnName,
			AddTime:  timeNow,
		})
	}
	if len(v.Precipitation1h) > 0 {
		out = append(out, PrometesusData{
			Name:     "降水量 1h",
			EName:    "Precipitation1h",
			Value:    fmt.Sprintf("%v", v.Precipitation1h[0]),
			Elems:    k.Elems,
			Type:     k.Type,
			NumberId: d,
			KjName:   k.KjName,
			KnName:   k.KnName,
			EnName:   k.EnName,
			AddTime:  timeNow,
		})
	}
	if len(v.Precipitation3h) > 0 {
		out = append(out, PrometesusData{
			Name:     "降水量 3h",
			EName:    "Precipitation3h",
			Value:    fmt.Sprintf("%v", v.Precipitation3h[0]),
			Elems:    k.Elems,
			Type:     k.Type,
			NumberId: d,
			KjName:   k.KjName,
			KnName:   k.KnName,
			EnName:   k.EnName,
			AddTime:  timeNow,
		})
	}
	if len(v.Precipitation24h) > 0 {
		out = append(out, PrometesusData{
			Name:     "降水量 24h",
			EName:    "Precipitation24h",
			Value:    fmt.Sprintf("%v", v.Precipitation24h[0]),
			Elems:    k.Elems,
			Type:     k.Type,
			NumberId: d,
			KjName:   k.KjName,
			KnName:   k.KnName,
			EnName:   k.EnName,
			AddTime:  timeNow,
		})
	}
	if len(v.WindDirection) > 0 {
		out = append(out, PrometesusData{
			Name:     "風向",
			EName:    "WindDirection",
			Value:    fmt.Sprintf("%v", v.WindDirection[0]),
			Elems:    k.Elems,
			Type:     k.Type,
			NumberId: d,
			KjName:   k.KjName,
			KnName:   k.KnName,
			EnName:   k.EnName,
			AddTime:  timeNow,
		})
	}
	if len(v.Wind) > 0 {
		out = append(out, PrometesusData{
			Name:     "風速",
			EName:    "Wind",
			Value:    fmt.Sprintf("%v", v.Wind[0]),
			Elems:    k.Elems,
			Type:     k.Type,
			NumberId: d,
			KjName:   k.KjName,
			KnName:   k.KnName,
			EnName:   k.EnName,
			AddTime:  timeNow,
		})
	}

	return out
}
func covertToePrometesusDataMap(d string, k TableData, v MapData) PrometesusDatas {
	timeNow := time.Now()
	pd := PrometesusDatas{
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
	if len(v.Temp) > 0 {
		pd.Temp[d] = PrometesusData{
			Name:     "気温",
			EName:    "Temp",
			Value:    fmt.Sprintf("%v", v.Temp[0]),
			Elems:    k.Elems,
			Type:     k.Type,
			NumberId: d,
			KjName:   k.KjName,
			KnName:   k.KnName,
			EnName:   k.EnName,
			AddTime:  timeNow,
		}
	}
	if len(v.Humidity) > 0 {
		pd.Humidity[d] = PrometesusData{
			Name:     "湿度",
			EName:    "Humidity",
			Value:    fmt.Sprintf("%v", v.Humidity[0]),
			Elems:    k.Elems,
			Type:     k.Type,
			NumberId: d,
			KjName:   k.KjName,
			KnName:   k.KnName,
			EnName:   k.EnName,
			AddTime:  timeNow,
		}
	}
	if len(v.Snow) > 0 {
		pd.Snow[d] = PrometesusData{
			Name:     "積雪",
			EName:    "Snow",
			Value:    fmt.Sprintf("%v", v.Snow[0]),
			Elems:    k.Elems,
			Type:     k.Type,
			NumberId: d,
			KjName:   k.KjName,
			KnName:   k.KnName,
			EnName:   k.EnName,
			AddTime:  timeNow,
		}
	}
	if len(v.Snow1h) > 0 {
		pd.Snow1h[d] = PrometesusData{
			Name:     "積雪深 1h",
			EName:    "Snow1h",
			Value:    fmt.Sprintf("%v", v.Snow1h[0]),
			Elems:    k.Elems,
			Type:     k.Type,
			NumberId: d,
			KjName:   k.KjName,
			KnName:   k.KnName,
			EnName:   k.EnName,
			AddTime:  timeNow,
		}
	}
	if len(v.Snow6h) > 0 {
		pd.Snow6h[d] = PrometesusData{
			Name:     "積雪深 6h",
			EName:    "Snow6h",
			Value:    fmt.Sprintf("%v", v.Snow6h[0]),
			Elems:    k.Elems,
			Type:     k.Type,
			NumberId: d,
			KjName:   k.KjName,
			KnName:   k.KnName,
			EnName:   k.EnName,
			AddTime:  timeNow,
		}
	}
	if len(v.Snow12h) > 0 {
		pd.Snow12h[d] = PrometesusData{
			Name:     "積雪深 12h",
			EName:    "Snow12h",
			Value:    fmt.Sprintf("%v", v.Snow12h[0]),
			Elems:    k.Elems,
			Type:     k.Type,
			NumberId: d,
			KjName:   k.KjName,
			KnName:   k.KnName,
			EnName:   k.EnName,
			AddTime:  timeNow,
		}
	}
	if len(v.Snow24h) > 0 {
		pd.Snow24h[d] = PrometesusData{
			Name:     "積雪深 24h",
			EName:    "Snow24h",
			Value:    fmt.Sprintf("%v", v.Snow24h[0]),
			Elems:    k.Elems,
			Type:     k.Type,
			NumberId: d,
			KjName:   k.KjName,
			KnName:   k.KnName,
			EnName:   k.EnName,
			AddTime:  timeNow,
		}
	}
	if len(v.Sun10m) > 0 {
		pd.Sun10m[d] = PrometesusData{
			Name:     "日照時間 10m",
			EName:    "Sun10m",
			Value:    fmt.Sprintf("%v", v.Sun10m[0]),
			Elems:    k.Elems,
			Type:     k.Type,
			NumberId: d,
			KjName:   k.KjName,
			KnName:   k.KnName,
			EnName:   k.EnName,
			AddTime:  timeNow,
		}
	}
	if len(v.Sun1h) > 0 {
		pd.Sun1h[d] = PrometesusData{
			Name:     "日照時間 1h",
			EName:    "Sun1h",
			Value:    fmt.Sprintf("%v", v.Sun1h[0]),
			Elems:    k.Elems,
			Type:     k.Type,
			NumberId: d,
			KjName:   k.KjName,
			KnName:   k.KnName,
			EnName:   k.EnName,
			AddTime:  timeNow,
		}
	}
	if len(v.Visibility) > 0 {
		pd.Visibility[d] = PrometesusData{
			Name:     "視程",
			EName:    "Visibility",
			Value:    fmt.Sprintf("%v", v.Visibility[0]),
			Elems:    k.Elems,
			Type:     k.Type,
			NumberId: d,
			KjName:   k.KjName,
			KnName:   k.KnName,
			EnName:   k.EnName,
			AddTime:  timeNow,
		}
	}
	if len(v.Precipitation10m) > 0 {
		pd.Precipitation10m[d] = PrometesusData{
			Name:     "降水量 10m",
			EName:    "Precipitation10m",
			Value:    fmt.Sprintf("%v", v.Precipitation10m[0]),
			Elems:    k.Elems,
			Type:     k.Type,
			NumberId: d,
			KjName:   k.KjName,
			KnName:   k.KnName,
			EnName:   k.EnName,
			AddTime:  timeNow,
		}
	}
	if len(v.Precipitation1h) > 0 {
		pd.Precipitation1h[d] = PrometesusData{
			Name:     "降水量 1h",
			EName:    "Precipitation1h",
			Value:    fmt.Sprintf("%v", v.Precipitation1h[0]),
			Elems:    k.Elems,
			Type:     k.Type,
			NumberId: d,
			KjName:   k.KjName,
			KnName:   k.KnName,
			EnName:   k.EnName,
			AddTime:  timeNow,
		}
	}
	if len(v.Precipitation3h) > 0 {
		pd.Precipitation3h[d] = PrometesusData{
			Name:     "降水量 3h",
			EName:    "Precipitation3h",
			Value:    fmt.Sprintf("%v", v.Precipitation3h[0]),
			Elems:    k.Elems,
			Type:     k.Type,
			NumberId: d,
			KjName:   k.KjName,
			KnName:   k.KnName,
			EnName:   k.EnName,
			AddTime:  timeNow,
		}
	}
	if len(v.Precipitation24h) > 0 {
		pd.Precipitation24h[d] = PrometesusData{
			Name:     "降水量 24h",
			EName:    "Precipitation24h",
			Value:    fmt.Sprintf("%v", v.Precipitation24h[0]),
			Elems:    k.Elems,
			Type:     k.Type,
			NumberId: d,
			KjName:   k.KjName,
			KnName:   k.KnName,
			EnName:   k.EnName,
			AddTime:  timeNow,
		}
	}
	if len(v.WindDirection) > 0 {
		pd.WindDirection[d] = PrometesusData{
			Name:     "風向",
			EName:    "WindDirection",
			Value:    fmt.Sprintf("%v", v.WindDirection[0]),
			Elems:    k.Elems,
			Type:     k.Type,
			NumberId: d,
			KjName:   k.KjName,
			KnName:   k.KnName,
			EnName:   k.EnName,
			AddTime:  timeNow,
		}
	}
	if len(v.Wind) > 0 {
		pd.Wind[d] = PrometesusData{
			Name:     "風速",
			EName:    "Wind",
			Value:    fmt.Sprintf("%v", v.Wind[0]),
			Elems:    k.Elems,
			Type:     k.Type,
			NumberId: d,
			KjName:   k.KjName,
			KnName:   k.KnName,
			EnName:   k.EnName,
			AddTime:  timeNow,
		}
	}
	return pd
}

func (pd *PrometesusDatas) cleanData() PrometesusDatas {
	wg := sync.WaitGroup{}
	var out PrometesusDatas
	//気温
	wg.Add(1)
	go func() {
		defer wg.Done()
		tmp := map[string]PrometesusData{}
		for k, v := range pd.Temp {
			//vの時刻が1時間を超えている場合は削除
			if time.Now().Sub(v.AddTime) < CLEAN_TIME {
				tmp[k] = v
			}
		}
		out.Temp = tmp
	}()
	//湿度
	wg.Add(1)
	go func() {
		defer wg.Done()
		tmp := map[string]PrometesusData{}
		for k, v := range pd.Humidity {
			//vの時刻が1時間を超えている場合は削除
			if time.Now().Sub(v.AddTime) < CLEAN_TIME {
				tmp[k] = v
			}
		}
		out.Humidity = tmp
	}()
	//積雪
	wg.Add(1)
	go func() {
		defer wg.Done()
		tmp := map[string]PrometesusData{}
		for k, v := range pd.Snow {
			//vの時刻が1時間を超えている場合は削除
			if time.Now().Sub(v.AddTime) < CLEAN_TIME {
				tmp[k] = v
			}
		}
		out.Snow = tmp
	}()
	//積雪深 1h
	wg.Add(1)
	go func() {
		defer wg.Done()
		tmp := map[string]PrometesusData{}
		for k, v := range pd.Snow1h {
			//vの時刻が1時間を超えている場合は削除
			if time.Now().Sub(v.AddTime) < CLEAN_TIME {
				tmp[k] = v
			}
		}
		out.Snow1h = tmp
	}()
	//積雪深 6h
	wg.Add(1)
	go func() {
		defer wg.Done()
		tmp := map[string]PrometesusData{}
		for k, v := range pd.Snow6h {
			//vの時刻が1時間を超えている場合は削除
			if time.Now().Sub(v.AddTime) < CLEAN_TIME {
				tmp[k] = v
			}
		}
		out.Snow6h = tmp
	}()
	//積雪深 12h
	wg.Add(1)
	go func() {
		defer wg.Done()
		tmp := map[string]PrometesusData{}
		for k, v := range pd.Snow12h {
			//vの時刻が1時間を超えている場合は削除
			if time.Now().Sub(v.AddTime) < CLEAN_TIME {
				tmp[k] = v
			}
		}
		out.Snow12h = tmp
	}()
	//積雪深 24h
	wg.Add(1)
	go func() {
		defer wg.Done()
		tmp := map[string]PrometesusData{}
		for k, v := range pd.Snow24h {
			//vの時刻が1時間を超えている場合は削除
			if time.Now().Sub(v.AddTime) < CLEAN_TIME {
				tmp[k] = v
			}
		}
		out.Snow24h = tmp
	}()
	//日照時間 10m
	wg.Add(1)
	go func() {
		defer wg.Done()
		tmp := map[string]PrometesusData{}
		for k, v := range pd.Sun10m {
			//vの時刻が1時間を超えている場合は削除
			if time.Now().Sub(v.AddTime) < CLEAN_TIME {
				tmp[k] = v
			}
		}
		out.Sun10m = tmp
	}()
	//日照時間 1h
	wg.Add(1)
	go func() {
		defer wg.Done()
		tmp := map[string]PrometesusData{}
		for k, v := range pd.Sun1h {
			//vの時刻が1時間を超えている場合は削除
			if time.Now().Sub(v.AddTime) < CLEAN_TIME {
				tmp[k] = v
			}
		}
		out.Sun1h = tmp
	}()
	//視程
	wg.Add(1)
	go func() {
		defer wg.Done()
		tmp := map[string]PrometesusData{}
		for k, v := range pd.Visibility {
			//vの時刻が1時間を超えている場合は削除
			if time.Now().Sub(v.AddTime) < CLEAN_TIME {
				tmp[k] = v
			}
		}
		out.Visibility = tmp
	}()
	//降水量 10m
	wg.Add(1)
	go func() {
		defer wg.Done()
		tmp := map[string]PrometesusData{}
		for k, v := range pd.Precipitation10m {
			//vの時刻が1時間を超えている場合は削除
			if time.Now().Sub(v.AddTime) < CLEAN_TIME {
				tmp[k] = v
			}
		}
		out.Precipitation10m = tmp
	}()
	//降水量 1h
	wg.Add(1)
	go func() {
		defer wg.Done()
		tmp := map[string]PrometesusData{}
		for k, v := range pd.Precipitation1h {
			//vの時刻が1時間を超えている場合は削除
			if time.Now().Sub(v.AddTime) < CLEAN_TIME {
				tmp[k] = v
			}
		}
		out.Precipitation1h = tmp
	}()
	//降水量 3h
	wg.Add(1)
	go func() {
		defer wg.Done()
		tmp := map[string]PrometesusData{}
		for k, v := range pd.Precipitation3h {
			//vの時刻が1時間を超えている場合は削除
			if time.Now().Sub(v.AddTime) < CLEAN_TIME {
				tmp[k] = v
			}
		}
		out.Precipitation3h = tmp
	}()
	//降水量 24h
	wg.Add(1)
	go func() {
		defer wg.Done()
		tmp := map[string]PrometesusData{}
		for k, v := range pd.Precipitation24h {
			//vの時刻が1時間を超えている場合は削除
			if time.Now().Sub(v.AddTime) < CLEAN_TIME {
				tmp[k] = v
			}
		}
		out.Precipitation24h = tmp
	}()
	//風向
	wg.Add(1)
	go func() {
		defer wg.Done()
		tmp := map[string]PrometesusData{}
		for k, v := range pd.WindDirection {
			//vの時刻が1時間を超えている場合は削除
			if time.Now().Sub(v.AddTime) < CLEAN_TIME {
				tmp[k] = v
			}
		}
		out.WindDirection = tmp
	}()
	//風速
	wg.Add(1)
	go func() {
		defer wg.Done()
		tmp := map[string]PrometesusData{}
		for k, v := range pd.Wind {
			//vの時刻が1時間を超えている場合は削除
			if time.Now().Sub(v.AddTime) < CLEAN_TIME {
				tmp[k] = v
			}
		}
		out.Wind = tmp
	}()
	wg.Wait()
	return out
}

func (v *PrometesusData) ToJson() string {
	b, _ := json.Marshal(v)
	return string(b)
}

func (v *PrometesusData) ToPrometesus() string {
	return fmt.Sprintf("amedas_%s{Name=\"%s\",Elems=\"%s\",Type=\"%s\",NumberId=\"%s\",KjName=\"%s\",KnName=\"%s\",EnName=\"%s\"} %s", v.EName, v.Name, v.Elems, v.Type, v.NumberId, v.KjName, v.KnName, v.EnName, v.Value)
}
