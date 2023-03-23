package loop

import (
	"booknewread/novel_chack"
	"fmt"
	"sort"
	"sync"
	"time"
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

func NobelLoop(urllists []string) {
	now := time.Now()
	ch := make(chan bool, len(urllists))
	listtemp = Listdata{data: []novel_chack.List{}, count: 1}
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
	endtime := time.Now()
	fmt.Println("read novel data end", (endtime.Sub(now)).Seconds(), "s")

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
	sort.Slice(ListData, func(i, j int) bool { return ListData[i].Lastdate.Unix() > ListData[j].Lastdate.Unix() })
}

func Count() int {
	listtemp.mu.Lock()
	count := listtemp.count
	listtemp.mu.Unlock()
	return count
}
