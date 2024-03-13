package hello

import "net/http"

func Init(url string, mux *http.ServeMux) error {
	mux.HandleFunc("GET "+url, HelloWeb)
	mux.HandleFunc("GET "+url+"/{id}", HelloId)
	return nil
}
