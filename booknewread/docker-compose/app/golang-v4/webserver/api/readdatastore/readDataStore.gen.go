package readdatastore

import "net/http"

func Init(url string, mux *http.ServeMux) error {
	mux.HandleFunc("GET "+url, readNewBook)
	return nil
}
