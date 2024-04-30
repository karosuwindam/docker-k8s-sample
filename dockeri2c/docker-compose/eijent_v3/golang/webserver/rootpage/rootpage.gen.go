package rootpage

import "net/http"

func Init(url string, mux *http.ServeMux) error {
	mux.HandleFunc(url, GetIndexPage)
	return nil
}
