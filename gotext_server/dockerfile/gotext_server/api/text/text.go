package text

import (
	"bufio"
	"encoding/json"
	"gocsvserver/config"
	"gocsvserver/dirread"
	"gocsvserver/webserver"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
)

type TXTData struct {
	Year  string   `json:"Year"`  // 年
	Quart string   `json:"Quart"` // 四半期
	Title []string `json:"Title"` // タイトル
}

var apiname string = "text"
var folder string
var TXTFolder *dirread.Dirtype

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

func webTextRead(w http.ResponseWriter, r *http.Request) {
	output := []TXTData{}
	TXTFolder, _ = dirread.Setup(folder)
	TXTFolder.Read("")
	for _, data := range TXTFolder.Data {
		if t := *readTxt(data.RootPath, data.Name); t.Year != "" {
			output = append(output, t)
		}
	}
	b, _ := json.Marshal(output)
	w.Write(b)
}

var route []webserver.WebConfig = []webserver.WebConfig{
	{Pass: "/" + apiname, Handler: webTextRead},
}

func Setup(cfg *config.Config) ([]webserver.WebConfig, error) {
	folder = cfg.TXT.RootPath
	if csv, err := dirread.Setup(folder); err != nil {
		return nil, err
	} else {
		if err := csv.Read(""); err != nil {
			return nil, err
		} else {
			TXTFolder = csv
		}
	}

	return route, nil
}
