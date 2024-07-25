package metrics

import (
	"fmt"
	"net/http"
	"tenkiej/contoroller"
	"tenkiej/contoroller/amedas"
	"tenkiej/logger"
)

const (
	HeaderText1 string = "（気象庁ホームページ https://www.jma.go.jp/ のアメダスのデータをもとに作成"
)

type Elems struct {
	Name  string
	Kname string
	Type  string
}

var metas []Elems = []Elems{
	{"Temp", "気温", "℃"},
	{"Humidity", "湿度", "%"},
	{"Snow", "積雪", "cm"},
	{"Snow1h", "積雪深 1h", "cm"},
	{"Snow6h", "積雪深 6h", "cm"},
	{"Snow12h", "積雪深 12h", "cm"},
	{"Snow24h", "積雪深 24h", "cm"},
	{"Sun10m", "日照時間 10m", "分"},
	{"Sun1h", "日照時間 1h", "時"},
	{"Visibility", "視程", ""},
	{"Precipitation10m", "降水量 10m", "mm"},
	{"Precipitation1h", "降水量 1h", "mm"},
	{"Precipitation3h", "降水量 3h", "mm"},
	{"Precipitation24h", "降水量 24h", "mm"},
	{"WindDirection", "風向", "0:北 1:北北東 2:北東 3:東北東 4:東 5:東南東 6:南東 7:南南東 8:南 9:南南西 10:南西 11:西南西 12:西 13:西北西 14:北西 15:北北西"},
	{"Wind", "風速", "m/s"},
}

func getMetrics(w http.ResponseWriter, r *http.Request) {
	var output []string
	api := contoroller.NewAPI()
	datas := api.GetAmedasMapData()
	output = convertData(datas)
	w.Header().Set("Content-Type", "text/plan; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	for _, d := range output {
		fmt.Fprintln(w, d)

	}
}

func convertData(pd []amedas.PrometesusData) []string {
	tmp := pd
	var out []string
	for _, d := range metas {
		var tmps []string
		var tmp_tmp []amedas.PrometesusData
		for i := 0; i < len(tmp); i++ {
			if tmp[i].EName == d.Name {
				tmps = append(tmps, tmp[i].ToPrometesus())
			} else {
				tmp_tmp = append(tmp_tmp, tmp[i])
			}
		}
		if len(tmps) > 0 {
			header := []string{
				fmt.Sprintf("# HELP amedas_%v %v (%v) %v", d.Name, d.Kname, d.Type, HeaderText1),
			}
			tmps = append(header, tmps...)
		}
		if len(tmps) != 0 {
			out = append(out, tmps...)
		}
		tmp = tmp_tmp
	}
	if len(tmp) != 0 {
		logger.Error("Unexpected Data", "data", tmp)
	}
	return out
}
