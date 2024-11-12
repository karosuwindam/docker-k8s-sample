package readdatastore

import (
	"book-newread/config"
	"net/http"
)

func Init(url string, mux *http.ServeMux) error {
	// mux.HandleFunc("GET "+url+"/book/{page}", readNewBook)
	// mux.HandleFunc("GET "+url+"/nobel", readNewNobel)
	config.TraceHttpHandleFunc(mux, "GET "+url+"/book/{page}", readNewBook)
	config.TraceHttpHandleFunc(mux, "GET "+url+"/nobel", readNewNobel)
	return nil
}
