package readdatastore

import (
	"book-newread/loop/datastore"
	"book-newread/loop/novelchack"
	"encoding/json"
	"net/http"
)

// データストアからWeb小説のデータを読み取る
func readNewNobel(w http.ResponseWriter, r *http.Request) {
	var data []novelchack.List
	if err := datastore.Read(&data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
