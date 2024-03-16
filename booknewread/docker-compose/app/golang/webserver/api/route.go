package api

import (
	"book-newread/webserver/api/hello"
	"book-newread/webserver/api/readdatastore"
	"book-newread/webserver/api/reset"
	"book-newread/webserver/api/status"
	"net/http"
)

type api struct {
	Router string
	Func   func(string, *http.ServeMux) error
}

var v1apis = []api{
	{"/hello", hello.Init},
	{"/status", status.Init},
	{"/reset", reset.Init},
	{"/json", readdatastore.Init},
}

func Init(mux *http.ServeMux) error {
	if err := v1apisetup(mux, "/v1"); err != nil {
		return err
	}
	return nil
}

func v1apisetup(mux *http.ServeMux, router string) error {
	if router == "/" {
		router = ""
	}
	if router[len(router)-1] == '/' {
		router = router[:len(router)-1]

	}
	for _, v := range v1apis {
		if err := v.Func(router+v.Router, mux); err != nil {
			return err
		}
	}
	return nil
}
