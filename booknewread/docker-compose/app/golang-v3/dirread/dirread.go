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
	if len(files) == 0 {
		var tmp2 []filedata
		for _, f := range tmp {
			if f.RootPath == t.path+s {

			} else {
				tmp2 = append(tmp2, f)
			}
		}
		t.Data = tmp2
		return 0
	}
	for _, f := range files {
		flag := true
		tmp2 := filedata{s + f.Name(), f.IsDir(), f.Size(), f.ModTime(), t.path + s}
		for _, ff := range tmp {
			if (ff.Name == tmp2.Name) && (ff.RootPath == tmp2.RootPath) {
				flag = false
				break
			}
		}
		if flag {
			tmp = append(tmp, tmp2)
		}
	}
	t.Data = tmp
	return 0

}
