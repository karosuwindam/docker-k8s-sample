package text

import (
	"bufio"
	"context"
	"encoding/json"
	"gocsvserver/config"
	"log/slog"
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
	ctx := r.Context()
	ctx, span := config.TracerS(ctx, "webTextRead", "web read Text File")
	defer span.End()
	slog.DebugContext(ctx, r.Method+":"+r.URL.Path, "method", r.Method, "url", r.URL.Path)
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
		slog.DebugContext(ctx, "webTextRead "+tmppass+file.Name(), "file", file.Name(), "tmppass", tmppass)
		ctxfile := contextWriteFilename(ctx, tmppass, file.Name())
		if t := *readTxt(ctxfile); t.Year != "" {
			output = append(output, t)
		}
	}
	slog.DebugContext(ctx, "webTextRead out sort", "output", output)
	sort.Slice(output, func(i, j int) bool {
		return output[i].Year+output[i].Quart > output[j].Year+output[j].Quart
	})
	b, _ := json.Marshal(output)
	w.Write(b)
}

func readTxt(ctx context.Context) *TXTData {
	ctx, span := config.TracerS(ctx, "readTxt", "read Text File")
	defer span.End()
	filepass, filename, ok := contextReadFilename(ctx)
	if !ok {
		slog.WarnContext(ctx, "readTxt", "error", "filename not found")
		return &TXTData{}
	}
	var title []string
	re := regexp.MustCompile(`(\d{4})_(\d{1})Q.txt`)
	if !re.MatchString(filename) {
		slog.WarnContext(ctx, "readTxt", "error", "filename not match")
		return &TXTData{}
	}
	slog.DebugContext(ctx, "readText Open file "+filepass+filename, "filename", filename)
	//行ごとに読み込む
	if f, err := os.Open(filepass + filename); err != nil {
		slog.ErrorContext(ctx, "readTxt os.OpenError", "error", err)
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
