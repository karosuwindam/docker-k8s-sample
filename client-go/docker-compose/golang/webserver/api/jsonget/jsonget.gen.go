package jsonget

import "net/http"

func Init(url string, mux *http.ServeMux) error {
	if url == "/" {
		url = ""
	}
	mux.HandleFunc(url, jsonget)
	return nil
}
