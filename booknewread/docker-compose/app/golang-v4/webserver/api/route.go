package api

import (
	"book-newread/webserver/api/hello"
	"book-newread/webserver/api/read"
	"net/http"
)

func Init(mux *http.ServeMux) error {
	hello.Init("/hello", mux)
	read.Init("/read", mux)
	return nil
}
