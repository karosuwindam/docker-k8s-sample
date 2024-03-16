package hello_test

import (
	"book-newread/config"
	"book-newread/webserver/api/hello"
	"context"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestHello(t *testing.T) {
	config.Init()
	mux := http.NewServeMux()
	if err := hello.Init("/", mux); err != nil {
		t.Fatal(err)
	}
	srv := &http.Server{
		Addr:    config.Web.Hostname + ":" + config.Web.Port,
		Handler: mux,
	}
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer srv.Shutdown(ctx)
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			t.Fatal(err)
		}
	}()
	// Httpサーバの起動待ち
	for {
		_, err := http.Get("http://localhost:8080/")
		if err == nil {
			break
		}
	}
	if getHello("") != "Hello Web" {
		t.Fatal("getHello() != Hello Web")
	}
	if getHello("1") != "Hello Id: 1" {
		t.Fatal("getHello(1) != Hello Id: 1")
	}
	cancel()
}

func getHello(url string) string {
	res, err := http.Get("http://localhost:8080/" + url)
	if err != nil {
		return ""
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return ""
	}
	return string(body)
}
