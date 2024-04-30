package reset

import "net/http"

func Init(url string, mux *http.ServeMux) error {
	mux.HandleFunc("POST "+url, PostReset)
	mux.HandleFunc("GET "+url, GetReset)
	return nil
}
