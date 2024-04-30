package health

import (
	"eijent/controller/senser"
	"encoding/json"
	"log"
	"net/http"
)

func GetReadHealth(w http.ResponseWriter, r *http.Request) {
	tmp := senser.Health()
	json, err := json.Marshal(&tmp)
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
