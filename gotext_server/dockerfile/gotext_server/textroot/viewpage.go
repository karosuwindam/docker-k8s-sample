package textroot

import (
	"fmt"
	"gocsvserver/textroot/textread"
	"log"
	"net/http"
	"strings"
)

const (
	ROOTPATH = "./html"
)

// 静的HTMLのページを返す
func viewhtml(w http.ResponseWriter, r *http.Request) {
	textdata := []string{".html", ".htm", ".css", ".js"}
	upath := r.URL.Path
	tmp := map[string]string{}
	if !strings.HasPrefix(upath, "/") {
		upath = "/" + upath
		r.URL.Path = upath
	}
	fmt.Println(r.Method + ":" + r.URL.Path)
	if upath[len(upath)-1:] == "/" {
		upath += "index.html"
	}
	if !textread.Exists(ROOTPATH + upath) {
		w.WriteHeader(404)
		log.Printf("ERROR request:%v\n", r.URL.Path)
		return
	} else {
		for _, data := range textdata {
			if len(upath) > len(data) {
				if upath[len(upath)-len(data):] == data {
					fmt.Fprint(w, textread.ConvertData(textread.ReadHtml(ROOTPATH+upath), tmp))
					return
				}
			}
		}
		buffer := textread.ReadOther(ROOTPATH + upath)
		// bodyに書き込み
		w.Write(buffer)
	}
	return
}
