package status

import (
	"book-newread/config"
	"net/http"
)

func Init(url string, mux *http.ServeMux) error {
	// mux.HandleFunc("GET "+url, status)
	config.TraceHttpHandleFunc(mux, "GET "+url, status)
	return nil
}
