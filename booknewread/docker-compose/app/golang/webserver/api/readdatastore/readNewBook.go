package readdatastore

import (
	"book-newread/loop/datastore"
	"encoding/json"
	"net/http"
	"strconv"
)

// 新刊のデータを取得する
func readNewBook(w http.ResponseWriter, r *http.Request) {
	var data []datastore.BListData
	if err := datastore.Read(&data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	pageString := r.PathValue("page")
	page, err := strconv.Atoi(pageString)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if page < 0 || page >= len(data) {
		http.Error(w, "page out of range", http.StatusBadRequest)
		return
	}
	if err := json.NewEncoder(w).Encode(data[page]); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
