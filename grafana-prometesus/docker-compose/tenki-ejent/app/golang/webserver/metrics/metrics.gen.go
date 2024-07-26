package metrics

import "net/http"

func Init(url string, mux *http.ServeMux) error {
	if url[len(url)-1:] == "/" {
		url = url[:len(url)-1]
	}
	mux.HandleFunc("GET "+url, getMetrics)
	mux.HandleFunc("GET "+url+"/v2", getMetrics)
	return nil
}
