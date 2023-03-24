package loop

import (
	"booknewread/novel_chack"
	"fmt"
	"sort"
	"sync"
	"time"
)

var NListData []novel_chack.List

type Listdata struct {
	data     []novel_chack.List
	listdata []BListData
	count    int
	mu       sync.Mutex
}

var listtemp Listdata

func (t *Listdata) chackurl(url string) novel_chack.List {
	return novel_chack.ChackUrldata(url)
}

func NobelLoop(urllists []string) {
	fmt.Println("start novel data count")
	now := time.Now()
	ch := make(chan bool, len(urllists))
	listtemp.data = []novel_chack.List{}
	listtemp.count = 1
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
	NListData = listtemp.data
	BookListData = listtemp.listdata
	listtemp.mu.Unlock()
	sort.Slice(NListData, func(i, j int) bool { return NListData[i].Lastdate.Unix() > NListData[j].Lastdate.Unix() })
}

func Count() int {
	listtemp.mu.Lock()
	count := listtemp.count
	listtemp.mu.Unlock()
	return count
}

func Setup() {
	EnvData = SetupEnv()
	BookListData = make([]BListData, 3)
	listtemp.listdata = make([]BListData, 3)
}
