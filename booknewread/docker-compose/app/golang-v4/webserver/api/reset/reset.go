package reset

import (
	"book-newread/loop"
	"fmt"
	"net/http"
)

func reset(w http.ResponseWriter, r *http.Request) {
	loop.Reset()
	fmt.Fprintf(w, "OK")
}
