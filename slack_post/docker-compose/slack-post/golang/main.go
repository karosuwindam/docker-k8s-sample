package main

import (
	"time"
	_ "time/tzdata"
)

func main() {
	t := CSetup(1)
	slackdata := EnvConfRead()

	for {
		pods := t.GetPod("")
		datas := GetPodInfo(pods)
		output := output_data(GetSennserdata(AnariseData(datas)))
		loc, _ := time.LoadLocation(slackdata.TimeZone)
		nowtime := time.Now().In(loc)
		const layout2 = "2006-01-02 15:04:05"
		slackdata.PostSlack(nowtime.Format(layout2) + "\n" + output)
		for {
			t := time.Date(nowtime.Year(), nowtime.Month(), nowtime.Day(), nowtime.Hour(), nowtime.Minute(), 0, 0, nowtime.Location())
			if time.Now().Sub(t.Add(time.Minute*time.Duration(slackdata.CountTime))) > 0 {
				break
			}

			time.Sleep(500 * time.Microsecond)
		}

	}
}
