package read

import (
	"book-newread/loop/datastore"
	"log/slog"
	"net/http"
)

func ReadWeb(w http.ResponseWriter, r *http.Request) {
	status := datastore.Status{}
	if err := datastore.Read(&status); err != nil {

		slog.Error("ReadWeb Read", "error", err)
		return
	}
	slog.Debug("ReadWeb", "status", status)
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Read Web"))
}
