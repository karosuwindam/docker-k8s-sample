package reset

import (
	"book-newread/loop"
	"fmt"
	"log/slog"
	"net/http"
)

func reset(w http.ResponseWriter, r *http.Request) {
	slog.DebugContext(r.Context(), "Reset")
	loop.Reset()
	fmt.Fprintf(w, "OK")
}
