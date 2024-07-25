package jsons

import (
	"encoding/json"
	"net/http"
	"tenkiej/contoroller"
)

func getJsons(w http.ResponseWriter, r *http.Request) {

	api := contoroller.NewAPI()
	datas := api.GetAmedasMapData()
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	b, _ := json.Marshal(datas)
	w.Write(b)

}
