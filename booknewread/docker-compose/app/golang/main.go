package main

import (
	"log"
	"sort"
	"strconv"
	"time"

	"./dirread"
	"./novel_chack"
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

var GrobalListData map[int]ListData

var Listdata []novel_chack.List

var GrobalStatus Status

var Reloadflag ReloadFlag

func main() {
	tmp := map[int]ListData{}
	GrobalListData = tmp
	var fpass dirread.Dirtype
	fpass.Setup("bookmark")
	GrobalStatus = Status{
		BookNowTIme:     time.Time{},
		BookStatus:      "Reload",
		BookMarkNowTime: time.Time{},
		BookMarkStatus:  "Reload",
	}
	go func() {
		log.Println("start novel data count")
		for {
			data := []novel_chack.BookBark{}
			GrobalStatus.BookMarkStatus = "Reload"

			listtmp := []novel_chack.List{}
			// pass := "bookmarks_2021_05_18.html"
			fpass.Read("/")
			for _, fd := range fpass.Data {
				pass := fd.RootPath + fd.Name
				ncr := novel_chack.ReadBookBark(pass)
				for _, tmpd := range ncr {
					data = append(data, tmpd)
				}
				data = novel_chack.BookBarkSout(data)
			}

			for i, str := range data {
				tmp := novel_chack.ChackUrldata(str.Url)
				if tmp.Title != "" {
					listtmp = append(listtmp, tmp)
					// fmt.Println(tmp)
				}
				GrobalStatus.BookMarkStatus = "Reload:" + strconv.Itoa(i*100/len(data)) + "%"
			}
			sort.Slice(listtmp, func(i, j int) bool { return listtmp[i].Lastdate.Unix() > listtmp[j].Lastdate.Unix() })
			Listdata = listtmp
			GrobalStatus.BookMarkStatus = "OK"
			GrobalStatus.BookMarkNowTime = time.Now()
			Reloadflag.BookMarkFlag = false
			log.Println("read novel data end")
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
				listdata.LiteNobel = FilterComicList(listdata.LiteNobel)
				listdata.Comic = FilterComicList(listdata.Comic)
				GrobalListData[i] = listdata
			}
			GrobalStatus.BookStatus = "OK"
			GrobalStatus.BookNowTIme = time.Now()
			Reloadflag.BookFlag = false
			log.Println("read new book data end")
			for i := 0; i < 60*60*12; i++ {
				if Reloadflag.BookFlag {
					break
				}
				time.Sleep(time.Second)
			}
			log.Println("reload new book data")
		}
	}()

	var web WebSetupData
	err := web.websetup()
	if err == nil {
		web.webstart()
	} else {
		log.Fatal(err)
	}
}
