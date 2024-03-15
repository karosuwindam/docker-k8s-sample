package readdatastore

import "net/http"

func Init(url string, mux *http.ServeMux) error {
	mux.HandleFunc("GET "+url+"/book/{page}", readNewBook)
	mux.HandleFunc("GET "+url+"/nobel", readNewNobel)
	return nil
}
