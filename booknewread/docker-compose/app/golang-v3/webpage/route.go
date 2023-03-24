package webpage

import (
	"booknewread/loop"
	"booknewread/textread"
	"booknewread/webserver"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

var Route []webserver.WebConfig = []webserver.WebConfig{
	{Pass: "/status", Handler: status},
	{Pass: "/jsonb", Handler: getlocaljson},
	{Pass: "/jsonnobel", Handler: getnowdata},
	{Pass: "/restart", Handler: restart},
	{Pass: "/", Handler: viewhtml},
}

const (
	ROOTPATH = "./html"
)

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "hello World")
}

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
	if upath == "/" {
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

// ステータスを取得する
func status(w http.ResponseWriter, r *http.Request) {
	loop.Read()
	jsondata, err := json.Marshal(loop.Statusdata)
	if err != nil {
		fmt.Fprint(w, err.Error())
	} else {
		fmt.Fprintf(w, "%s", jsondata)
	}
}

func getlocaljson(w http.ResponseWriter, r *http.Request) {
	form_data := ""
	// timedata := time.Now().Sub(GrobalStatus.BookNowTIme).Seconds()
	// if timedata > 300 {
	// 	Reloadflag.BookMarkFlag = true
	// 	Reloadflag.BookFlag = true
	// }
	r.ParseForm()
	for cnt, strs := range r.Form {
		form_data += " " + cnt + ":"
		for i, str := range strs {
			if i == 0 {
				form_data += str
			} else {
				form_data += "," + str
			}
		}
	}
	fmt.Println(r.Method + ":" + r.URL.Path + " " + form_data)
	loop.Read()
	if r.Method == "GET" {
		jsondata, err := json.Marshal(loop.BookListData[0])
		if err != nil {
			fmt.Fprint(w, err.Error())
		} else {
			fmt.Fprintf(w, "%s", jsondata)
		}
	} else if r.Method == "POST" {
		mf := r.MultipartForm

		// 通常のリクエスト
		if mf != nil {
			for k, v := range mf.Value {
				fmt.Printf("%v : %v", k, v)
			}
		}
		page := r.FormValue("page")
		tmp_page, _ := strconv.Atoi(page)
		if tmp_page >= len(loop.BookListData) {
			tmp_page = 0
		}
		jsondata, err := json.Marshal(loop.BookListData[tmp_page])
		if err != nil {
			fmt.Fprint(w, err.Error())
		} else {
			fmt.Fprintf(w, "%s", jsondata)
		}
	}
}

func getnowdata(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Method + ":" + r.URL.Path)
	loop.Read()

	if r.Method == "GET" {
		jsondata, err := json.Marshal(loop.NListData)
		if err != nil {
			fmt.Fprint(w, err.Error())
		} else {
			fmt.Fprintf(w, "%s", jsondata)
		}
	}
}

func restart(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Method + ":" + r.URL.Path)

	if r.Method == "POST" {
		loop.Reset_ON(loop.RESET_DATA)
		fmt.Fprintf(w, "OK")
	} else {
		fmt.Fprintf(w, "NG")

	}
}
