package read

import (
	"book-newread/loop/datastore"
	"fmt"
	"net/http"
)

func ReadWeb(w http.ResponseWriter, r *http.Request) {
	status := datastore.Status{}
	if err := datastore.Read(&status); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(status)
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Read Web"))
}
