package reset

import (
	"book-newread/config"
	"net/http"
)

func Init(url string, mux *http.ServeMux) error {
	// mux.HandleFunc("POST "+url, reset)
	config.TraceHttpHandleFunc(mux, "POST "+url, reset)
	return nil
}
