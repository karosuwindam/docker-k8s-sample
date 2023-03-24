package loop

import (
	"booknewread/novel_chack"
	"fmt"
	"sort"
	"strconv"
	"sync"
	"time"
)

var NListData []novel_chack.List

type Listdata struct {
	data     []novel_chack.List
	listdata []BListData
	status   Status
	count    int
	mu       sync.Mutex
}

var listtemp Listdata

func (t *Listdata) chackurl(url string) novel_chack.List {
	if tmp, err := novel_chack.ChackUrldata(url); err != nil {
		time.Sleep(time.Microsecond * 100)
		tmp, err = novel_chack.ChackUrldata(url)
		return tmp
	} else {
		return tmp
	}
}

func NobelLoop(urllists []string) {
	fmt.Println("start novel data count")
	dataStatusSet(NOBEL_SELECT, "Reload")
	now := time.Now()
	ch := make(chan bool, len(urllists))
	listtemp.count = 1
	novel_chack.Setup()

	for i, url := range urllists {
		go func(i int, url string) {
			tmp := listtemp.chackurl(url)
			listtemp.add(tmp)
			dataStatusSet(NOBEL_SELECT, "Reload:"+strconv.Itoa(Count()*100/len(urllists))+"%")
			ch <- true
		}(i, url)
	}
	for i := 0; i < len(urllists); i++ {
		<-ch
	}
	endtime := time.Now()
	timeStatusSet(NOBEL_SELECT, endtime)
	dataStatusSet(NOBEL_SELECT, "OK")
	fmt.Println("read novel data end", (endtime.Sub(now)).Seconds(), "s")

}

func (t *Listdata) add(data novel_chack.List) {
	t.mu.Lock()
	t.count++
	t.mu.Unlock()
	if data.Url == "" {
		return
	}
	t.mu.Lock()
	flag := true
	for i, tmp := range t.data {
		if data.Title == "" {
			flag = false
			break
		}
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
	Statusdata = listtemp.status
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
	listtemp.data = []novel_chack.List{}
	listtemp.listdata = make([]BListData, 3)
	listtemp.status = Status{
		BookNowTIme:     time.Time{},
		BookStatus:      "Reload",
		BookMarkNowTime: time.Time{},
		BookMarkStatus:  "Reload",
	}
}
