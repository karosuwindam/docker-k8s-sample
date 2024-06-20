package readdatastore

import (
	"book-newread/config"
	"book-newread/loop/datastore"
	"book-newread/loop/novelchack"
	"encoding/json"
	"net/http"
	"sort"
)

// データストアからWeb小説のデータを読み取る
func readNewNobel(w http.ResponseWriter, r *http.Request) {
	var data []novelchack.List
	var errch chan error = make(chan error, 1)
	ctx := r.Context()
	ctx, handerTracer := config.TracerS(ctx, "readNewNobel", r.URL.Path)
	defer func() {
		handerTracer.End()
	}()
	go func() {
		_, span := config.TracerS(ctx, "datastore.Read", "datastore read")
		defer func() {
			span.End()
		}()
		errch <- datastore.Read(&data)
	}()
	if err := <-errch; err != nil {
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
