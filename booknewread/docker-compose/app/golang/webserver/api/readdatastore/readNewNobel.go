package readdatastore

import (
	"book-newread/loop/datastore"
	"book-newread/loop/novelchack"
	"encoding/json"
	"net/http"
	"sort"
)

// データストアからWeb小説のデータを読み取る
func readNewNobel(w http.ResponseWriter, r *http.Request) {
	var data []novelchack.List
	if err := datastore.Read(&data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	sort.Slice(data, func(i, j int) bool {
		return data[i].Lastdate.Unix() > data[j].Lastdate.Unix()
	})
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
