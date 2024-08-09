package jsonget

import (
	"fmt"
	"ingresslist/controller"
	"net/http"
)

// httpハンドラ
func jsonget(w http.ResponseWriter, r *http.Request) {
	api := controller.NewApi()
	jdata := api.GetJson()
	// レスポンスを返す
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, jdata)
}
