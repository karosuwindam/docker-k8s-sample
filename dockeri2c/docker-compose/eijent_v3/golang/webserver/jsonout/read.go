package jsonout

import (
	"eijent/controller/senser"
	"encoding/json"
	"log"
	"net/http"
)

func GetReadJson(w http.ResponseWriter, r *http.Request) {
	tmp := senser.ReadValue()
	json, err := json.Marshal(&tmp.Data)
	if err != nil {
		log.Println("error:", err)

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("error"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(json)
}
