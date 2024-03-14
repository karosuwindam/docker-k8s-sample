package api

import (
	"book-newread/webserver/api/hello"
	"book-newread/webserver/api/reset"
	"book-newread/webserver/api/status"
	"net/http"
)

func Init(mux *http.ServeMux) error {
	hello.Init("/hello", mux)
	status.Init("/status", mux)
	reset.Init("/reset", mux)
	return nil
}
