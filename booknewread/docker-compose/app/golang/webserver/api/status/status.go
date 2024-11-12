package status

import (
	"book-newread/loop/datastore"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
)

func status(w http.ResponseWriter, r *http.Request) {
	status := datastore.Status{}
	if err := datastore.Read(&status); err != nil {
		slog.ErrorContext(r.Context(), "status Read", "error", err)
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusServiceUnavailable)
		fmt.Fprint(w, err.Error())
		return
	}
	if jsondata, err := json.Marshal(status); err != nil {
		slog.ErrorContext(r.Context(), "status Marshal", "error", err)
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, err.Error())
	} else {
		slog.DebugContext(r.Context(), "status", "status", status)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "%s", jsondata)

	}
}
