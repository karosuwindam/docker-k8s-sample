package text

import (
	"bufio"
	"encoding/json"
	"gocsvserver/config"
	"log"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strings"
)

type TXTData struct {
	Year  string   `json:"Year"`  // 年
	Quart string   `json:"Quart"` // 四半期
	Title []string `json:"Title"` // タイトル
}

func webTextRead(w http.ResponseWriter, r *http.Request) {
	//config.Read.FilePassによるフォルダ指定からファイルリスト取得
	output := []TXTData{}
	tmppass := config.Read.FilePass
	if tmppass[len(tmppass)-1] != '/' {
		tmppass += "/"
	}
	files, err := os.ReadDir(tmppass)
	if err != nil {
		w.Write([]byte("Error: " + err.Error()))
		return
	}
	for _, file := range files {
		if t := *readTxt(tmppass, file.Name()); t.Year != "" {
			output = append(output, t)
		}
	}
	sort.Slice(output, func(i, j int) bool {
		return output[i].Year+output[i].Quart > output[j].Year+output[j].Quart
	})
	b, _ := json.Marshal(output)
	w.Write(b)
}

func readTxt(filepass, filename string) *TXTData {
	var title []string
	re := regexp.MustCompile(`(\d{4})_(\d{1})Q.txt`)
	if !re.MatchString(filename) {
		return &TXTData{}
	}
	//行ごとに読み込む
	if f, err := os.Open(filepass + filename); err != nil {
		log.Println(err)
		return &TXTData{}
	} else {
		defer f.Close()
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			title = append(title, scanner.Text())
		}
	}
	s := strings.Split(filename[:len(filename)-4], "_")
	tmp := &TXTData{
		Year:  s[0],
		Quart: s[1],
		Title: title,
	}
	return tmp

}
