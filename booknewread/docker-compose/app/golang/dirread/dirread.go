package dirread

import (
	"io/ioutil"
	"log"
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
//Dirtypeは[]filedata objectをバッファする必要あり
type Dirtype struct {
	path  string
	Data  []filedata
	Count []int
	Renew bool
}

func (t *Dirtype) Setup(s string) {
	t.path = s
	var tmp []filedata
	var tmp2 []int
	if (len(t.Data) == 0) || (t.Renew) {
		t.Data = tmp
		t.Count = tmp2
		t.Renew = false
	}
}

func (t *Dirtype) Read(s string) int {
	var tmp []filedata
	tmp = append(tmp, t.Data...)
	if t.path == "" {
		return -1
	}
	files, err := ioutil.ReadDir(t.path + s)
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {
		tmp2 := filedata{s + f.Name(), f.IsDir(), f.Size(), f.ModTime(), t.path + s}
		tmp = append(tmp, tmp2)
	}
	t.Data = tmp
	return 0

}
