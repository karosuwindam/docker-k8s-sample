package reset

import (
	"eijent/controller/senser"
	"fmt"
	"net/http"
)

func PostReset(w http.ResponseWriter, r *http.Request) {
	senser.Reset()
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	output := "<html><body><a href=\"/\">index</a></body></html>"
	fmt.Fprintf(w, "%s", output)
}

func GetReset(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	output := "<html><body><a href=\"/metrics\">metrics</a></body></html>"
	fmt.Fprintf(w, "%s", output)
}
