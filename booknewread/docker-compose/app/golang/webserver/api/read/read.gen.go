package read

import (
	"book-newread/config"
	"net/http"
)

func Init(url string, mux *http.ServeMux) error {
	// mux.HandleFunc("GET "+url, ReadWeb)
	config.TraceHttpHandleFunc(mux, "GET "+url, ReadWeb)
	return nil
}
