package dirread

import (
	"errors"
	"io/ioutil"
	"os"
	"time"
)

type filedata struct {
	Name     string
	Folder   bool
	Size     int64
	Time     time.Time
	RootPath string
}

// Dirtype implements buffering for an []filedata object.
// Dirtypeは[]filedata objectをバッファする必要あり
type Dirtype struct {
	path  string     //rootパス
	Data  []filedata //フォルダデータ
	Count []int      //
	Renew bool       //
}

// フォルダパスの確認
func dirpasscheck(s string) bool {
	if f, err := os.Stat(s); os.IsNotExist(err) || !f.IsDir() {
		return false
	}
	return true
}

// セットアップ
func Setup(s string) (*Dirtype, error) {
	t := &Dirtype{}
	if !dirpasscheck(s) {
		return t, errors.New("folder pass err :" + s)
	}
	if s[len(s)-1] == "/"[0] {
		t.path = s
	} else {
		t.path = s + "/"
	}
	var tmp []filedata
	var tmp2 []int
	if (len(t.Data) == 0) || (t.Renew) {
		t.Data = tmp
		t.Count = tmp2
		t.Renew = false
	}
	return t, nil
}

// 読み取り
func (t *Dirtype) Read(s string) error {
	var tmp []filedata
	tmp = append(tmp, t.Data...)
	if t.path == "" {
		return errors.New("error root pass")
	}
	files, err := ioutil.ReadDir(t.path + s)
	if err != nil {
		return err
	}
	for _, f := range files {
		tmp2 := filedata{f.Name(), f.IsDir(), f.Size(), f.ModTime(), t.path + s}
		tmp = append(tmp, tmp2)
	}
	t.Data = tmp
	return nil

}
