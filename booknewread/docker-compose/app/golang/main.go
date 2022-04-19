package main

import (
	"fmt"
	"log"
	"runtime"
	"sort"
	"strconv"
	"time"

	"book-newread/dirread"
	"book-newread/novel_chack"
)

type ListData struct {
	NowTime   time.Time  `json:Nowtime`
	Month     string     `json:Month`
	Year      string     `json:Year`
	Comic     []BookList `json:Comic`
	LiteNobel []BookList `json:LiteNobel`
}

type ReloadFlag struct {
	BookFlag     bool
	BookMarkFlag bool
}

type Status struct {
	BookNowTIme     time.Time `json:booknowtime`
	BookStatus      string    `json:bookstatus`
	BookMarkNowTime time.Time `json:bookmarknowtime`
	BookMarkStatus  string    `json:bookmarkstatus`
}

//環境変数で操作できる数値周
type SetupEnvData struct {
	MaxAccess  int //一つのサイトに対してアクセスできる限界数
	MaxBackDay int //今日を基準に表示する前日数
}

//デフォルト
const (
	DEF_ACCESS  = 2
	DEF_BACKDAY = 5
)

//グローバル数値化
var GrobalListData map[int]ListData

var Listdata []novel_chack.List

var GrobalStatus Status

var Reloadflag ReloadFlag

var EnvData SetupEnvData

func SetupEnv() SetupEnvData {
	var output SetupEnvData
	if true {
		output.MaxAccess = 3
	} else {
		output.MaxAccess = DEF_ACCESS
	}
	if true {
		output.MaxBackDay = 2
	} else {
		output.MaxBackDay = DEF_BACKDAY
	}
	return output
}

func main() {
	tmp := map[int]ListData{}
	GrobalListData = tmp
	var fpass dirread.Dirtype
	readflag := false
	ch1 := make(chan bool)
	EnvData = SetupEnv()

	fpass.Setup("bookmark")
	fpass.Read("/")
	if len(fpass.Data) == 0 {
		fmt.Println("err bookmark not file")
		return
	}
	GrobalStatus = Status{
		BookNowTIme:     time.Time{},
		BookStatus:      "Reload",
		BookMarkNowTime: time.Time{},
		BookMarkStatus:  "Reload",
	}
	go func() {
		log.Println("start novel data count")
		ch := novel_chack.Setup(EnvData.MaxAccess)
		data := []novel_chack.BookBark{}
		for {
			starttime := time.Now()
			GrobalStatus.BookMarkStatus = "Reload"

			listtmp := []novel_chack.List{}
			fpass.Read("/")
			if len(fpass.Data) != 0 {
				data = []novel_chack.BookBark{}
				for _, fd := range fpass.Data {
					pass := fd.RootPath + fd.Name
					ncr := novel_chack.ReadBookBark(pass)
					for _, tmpd := range ncr {
						data = append(data, tmpd)
					}
					data = novel_chack.BookBarkSout(data)
				}
			}

			//マルチ処理
			count := 0
			cpus := runtime.NumCPU()
			ch_th := make(chan int, cpus)
			tempdata := make([]novel_chack.List, len(data))
			for i, str := range data {
				// log.Println(str.Url)
				ch_th <- 1
				go func(num int, url string) {
					tempdata[num] = ch.ChackUrldata(url)
					count += <-ch_th
				}(i, str.Url)
				GrobalStatus.BookMarkStatus = "Reload:" + strconv.Itoa(count*100/len(data)) + "%"
				// time.Sleep(time.Millisecond * 100)
			}
			for {
				GrobalStatus.BookMarkStatus = "Reload:" + strconv.Itoa(count*100/len(data)) + "%"
				if len(ch_th) == 0 {
					for _, str := range tempdata {
						if str.Title != "" {
							// fmt.Println(str)
							listtmp = append(listtmp, str)
						}
					}
					break
				}
				time.Sleep(time.Millisecond * 100)
			}

			sort.Slice(listtmp, func(i, j int) bool { return listtmp[i].Lastdate.Unix() > listtmp[j].Lastdate.Unix() })
			endtime := time.Now()
			Listdata = listtmp
			if !readflag {
				ch1 <- true
				readflag = true
			}
			GrobalStatus.BookMarkStatus = "OK"
			GrobalStatus.BookMarkNowTime = time.Now()
			Reloadflag.BookMarkFlag = false
			log.Println("read novel data end", (endtime.Sub(starttime)).Seconds(), "s")
			for i := 0; i < 60*60; i++ {
				if Reloadflag.BookMarkFlag {
					break
				}
				time.Sleep(time.Second)
			}
			log.Println("reload novel data")
		}
	}()
	go func() {
		log.Println("start new book data count")
		for {
			starttime := time.Now()
			GrobalStatus.BookStatus = "Reload"
			t := time.Now()
			for i := 0; i < 3; i++ {
				var listdata ListData
				listdata.NowTime = time.Now()
				if (int(t.Month()) + i) <= 12 {
					listdata.Month = strconv.Itoa((int(t.Month()) + i))
					listdata.Year = strconv.Itoa(int(t.Year()))
				} else {
					listdata.Month = strconv.Itoa((int(t.Month()) + i) % 12)
					listdata.Year = strconv.Itoa(int(t.Year()) + 1)
				}

				listdata.LiteNobel = GetComicList(listdata.Year, listdata.Month, LITENOVEL)
				listdata.Comic = GetComicList(listdata.Year, listdata.Month, COMIC)
				listdata.LiteNobel = FilterComicList(listdata.LiteNobel, EnvData.MaxBackDay)
				listdata.Comic = FilterComicList(listdata.Comic, EnvData.MaxBackDay)
				GrobalListData[i] = listdata
			}
			endtime := time.Now()
			GrobalStatus.BookStatus = "OK"
			GrobalStatus.BookNowTIme = time.Now()
			Reloadflag.BookFlag = false
			log.Println("read new book data end", (endtime.Sub(starttime)).Seconds(), "s")
			for i := 0; i < 60*60*12; i++ {
				if Reloadflag.BookFlag {
					break
				}
				time.Sleep(time.Second)
			}
			log.Println("reload new book data")
		}
	}()
	<-ch1
	var web WebSetupData
	err := web.websetup()
	if err == nil {
		web.webstart()
	} else {
		log.Fatal(err)
	}
}
