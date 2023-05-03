package text

import (
	"bufio"
	"gocsvserver/config"
	"gocsvserver/dirread"
	"log"
	"os"
	"regexp"
	"strings"
)

type TXTData struct {
	Year  string   `json:"year"`  // 年
	Quart string   `json:"quart"` // 四半期
	Title []string `json:"title"` // タイトル
}

var apiname string = "text"
var folder string
var TXTFolder *dirread.Dirtype

func readTxt(filepass, filename string) *TXTData {
	var title []string
	re := regexp.MustCompile(`(\d{4})_(\d{1})Q.txt`)
	if !re.MatchString(filename) {
		return nil
	}
	//行ごとに読み込む
	if f, err := os.Open(filepass + filename); err != nil {
		log.Println(err)
		return nil
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

func Setup(cfg *config.Config) error {
	folder = cfg.TXT.RootPath
	if csv, err := dirread.Setup(folder); err != nil {
		return err
	} else {
		if err := csv.Read(""); err != nil {
			return err
		} else {
			TXTFolder = csv
		}
	}

	return nil
}
