package loop

import (
	"booknewread/novel_chack"
	"sync"
)

var ListData []novel_chack.List

type Listdata struct {
	data  []novel_chack.List
	count int
	mu    sync.Mutex
}

var listtemp Listdata

func (t *Listdata) chackurl(url string) novel_chack.List {
	return novel_chack.ChackUrldata(url)
}

func Loop(urllists []string) {
	ch := make(chan bool, len(urllists))
	listtemp = Listdata{data: []novel_chack.List{}, count: 0}
	novel_chack.Setup()

	for i, url := range urllists {
		go func(i int, url string) {
			tmp := listtemp.chackurl(url)
			listtemp.add(tmp)
			ch <- true
		}(i, url)
	}
	for i := 0; i < len(urllists); i++ {
		<-ch
	}
}

func (t *Listdata) add(data novel_chack.List) {
	t.count++
	if data.Url == "" {
		return
	}
	t.mu.Lock()
	flag := true
	for i, tmp := range t.data {
		if tmp.Url == data.Url {
			t.data[i] = data
			flag = false
			break
		}
	}
	if flag {
		t.data = append(t.data, data)
	}
	t.mu.Unlock()
}

func Read() {
	listtemp.mu.Lock()
	ListData = listtemp.data
	listtemp.mu.Unlock()
}

func Count() int {
	listtemp.mu.Lock()
	count := listtemp.count
	listtemp.mu.Unlock()
	return count

}
