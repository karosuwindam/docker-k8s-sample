package api

import (
	"book-newread/webserver/api/hello"
	"net/http"
)

func Init(mux *http.ServeMux) error {
	hello.Init("/hello", mux)
	return nil
}
