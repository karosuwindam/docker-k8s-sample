package indexpage

import (
	"book-newread/config"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"strings"
)

var baseurl string

func Init(url string) func(w http.ResponseWriter, r *http.Request) {
	baseurl = url

	return index
}

func index(w http.ResponseWriter, r *http.Request) {
	slog.DebugContext(r.Context(), r.Method+":"+r.URL.Path, "method", r.Method, "url", r.URL.Path)
	url := r.URL.Path[len(baseurl):]
	pass := config.Web.StaticPage
	if pass[len(pass)-1:] != "/" {
		pass += "/"
	}
	if url == "" || url == "index.html" || url == "index.htm" {
		filepath := pass + "index.html"
		_, err := os.Stat(filepath)
		if err == nil {
			tmp := make(map[string]string)
			tmp["base_title"] = "新刊取得"
			title := os.Getenv("WEB_TITLE")
			if title != "" {
				tmp["base_title"] = title
			}
			tpl := template.Must(template.ParseFiles(filepath))
			tpl.Execute(w, tmp)
			return
		}

	} else {
		filepath := pass + url
		_, err := os.Stat(filepath)
		if err == nil {
			if strings.Index(filepath, ".js") >= 0 || strings.Index(filepath, ".css") >= 0 {

				fp, _ := os.Open(filepath)
				defer fp.Close()
				buf := make([]byte, 1024)
				var buffer []byte
				for {
					n, err := fp.Read(buf)
					if err != nil {
						break
					}
					if n == 0 {
						break
					}
					buffer = append(buffer, buf[:n]...)
				}
				w.Write(buffer)
				return

			} else {
				tpl := template.Must(template.ParseFiles(filepath))
				tpl.Execute(w, nil)
				return

			}
		}
	}
	slog.WarnContext(r.Context(), "404 Not Found", "url", r.URL.Path)
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("404 Not Found"))

}
