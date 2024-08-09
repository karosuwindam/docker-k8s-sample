package controller

import (
	"ingresslist/controller/datastore"
	"ingresslist/controller/kubeget"
	"sort"
	"strings"
	"sync"
)

type OutputData struct {
	PName       string   `json:"Podname"`
	Namespace   string   `json:"Namspace"`
	Podinfo     string   `json:"Podinfo"`
	Port        []string `json:"Port"`
	Ip          string   `json:"Ip"`
	ClusterIP   string   `json:"ClusterIP"`
	HostNetwork bool     `json:"Hostnetwork"`
	SName       string   `json:"Servicename"`
	Selector    string   `json:"Selector"`
	URL         []string `json:"URL"`
}

type DomainOutputData struct {
	PName       string   `json:"Podname"`
	Namespace   string   `json:"Namspace"`
	Podinfo     string   `json:"Podinfo"`
	Port        []string `json:"Port"`
	Ip          string   `json:"Ip"`
	ClusterIP   string   `json:"ClusterIP"`
	HostNetwork bool     `json:"Hostnetwork"`
	SName       string   `json:"Servicename"`
	Selector    string   `json:"Selector"`
	URL         string   `json:"URL"`
}

type DomainData struct {
	Domain string             `json:"Domain"`
	Data   []DomainOutputData `json:"Data"`
}

type OuteJsonData struct {
	Data []DomainData `json:"Data"`
}

func getDomainData() (OuteJsonData, error) {
	var output OuteJsonData
	d, err := getDatabase()
	if err != nil {
		return output, err
	}
	//domainの一覧を取得
	var tmpdomainlist []string
	for _, v := range d {
		for _, vv := range v.URL {
			tmpdomainlist = append(tmpdomainlist, vv)
		}
	}
	var domainlist []string
	//からデータの追加
	domainlist = append(domainlist, "")
	for _, v := range tmpdomainlist {
		//ピリオド区切りのドメインを分割
		tmp := strings.Split(v, ".")
		if len(tmp) > 2 {
			ts := tmp[len(tmp)-2] + "." + tmp[len(tmp)-1]
			domainlist = append(domainlist, ts)
		} else {
			domainlist = append(domainlist, v)
		}
	}
	//重複を削除
	domainlist = removeDuplicate(domainlist)
	//ドメイン別にベースデータを作る
	for _, v := range domainlist {
		var tmpdata DomainData
		tmpdata.Domain = v
		output.Data = append(output.Data, tmpdata)
	}
	//ドメイン別にデータを振り分ける
	for _, v := range d {
		for i, vv := range output.Data {
			if (len(v.URL) == 0) && (vv.Domain == "") {
				tmp := DomainOutputData{
					PName:       v.PName,
					Namespace:   v.Namespace,
					Podinfo:     v.Podinfo,
					Port:        v.Port,
					Ip:          v.Ip,
					ClusterIP:   v.ClusterIP,
					HostNetwork: v.HostNetwork,
					SName:       v.SName,
					Selector:    v.Selector,
				}
				output.Data[i].Data = append(output.Data[i].Data, tmp)
				continue
			}
			for _, vvv := range v.URL {
				if (strings.Index(vvv, vv.Domain) > -1) && (vvv != "") && (vv.Domain != "") {
					tmp := DomainOutputData{
						PName:       v.PName,
						Namespace:   v.Namespace,
						Podinfo:     v.Podinfo,
						Port:        v.Port,
						Ip:          v.Ip,
						ClusterIP:   v.ClusterIP,
						HostNetwork: v.HostNetwork,
						SName:       v.SName,
						Selector:    v.Selector,
						URL:         vvv,
					}
					output.Data[i].Data = append(output.Data[i].Data, tmp)
				}
			}
		}
	}
	//ソート
	sort.Slice(output.Data, func(i, j int) bool {
		return output.Data[i].Domain > output.Data[j].Domain
	})
	for i := 0; i < len(output.Data); i++ {
		tmp := output.Data[i].Data
		sort.Slice(tmp, func(i, j int) bool {
			return tmp[i].PName < tmp[j].PName
		})
		output.Data[i].Data = tmp
	}
	for i, v := range output.Data {
		if v.Domain == "" || len(v.Domain) < 2 {
			output.Data[i].Domain = v.Domain
		} else if v.Domain[len(v.Domain)-1:len(v.Domain)] == "/" {
			output.Data[i].Domain = "*." + v.Domain[:len(v.Domain)-1]
		} else {
			output.Data[i].Domain = "*." + v.Domain
		}
	}

	return output, nil
}

// 重複を削除
func removeDuplicate(args []string) []string {
	results := make([]string, 0, len(args))
	encountered := map[string]bool{}
	for _, v := range args {
		if !encountered[v] {
			encountered[v] = true
			results = append(results, v)
		}
	}
	return results
}

func getDatabase() ([]OutputData, error) {
	output := []OutputData{}
	pods := datastore.GetPod()
	services := datastore.GetService()
	ingress := datastore.GetIngress()
	// nodes := datastore.GetNode()
	for _, pod := range pods {
		var tmpdata OutputData
		tmpdata.PName = pod.Name
		tmpdata.Namespace = pod.Namespace
		tmpdata.Podinfo = pod.PodInfo
		tmpdata.Port = pod.Port
		tmpdata.Ip = pod.Ip
		tmpdata.HostNetwork = pod.HostNetwork
		tmpdata.Selector = pod.Selector
		output = append(output, tmpdata)
	}

	var wq sync.WaitGroup
	for i := 0; i < len(output); i++ {
		wq.Add(1)
		go addServiceData(&output[i], services, &wq)
	}
	wq.Wait()
	for i := 0; i < len(output); i++ {
		wq.Add(1)
		go addIngressData(&output[i], ingress, &wq)
	}
	wq.Wait()

	return output, nil
}

// PodデータをもとにServiceデータを追加する
func addServiceData(pod *OutputData, service []kubeget.ServiceData, wq *sync.WaitGroup) {
	if pod != nil {
		for _, v := range service {
			if (pod.Selector == v.Aselector) && (pod.Namespace == v.Namespace) {
				pod.SName = v.Name
				pod.ClusterIP = v.ClusterIP
				break
			}
		}
	}
	wq.Done()
}

// podデータとServiceデータをもとにIngressデータを追加する
func addIngressData(pod *OutputData, ingress []kubeget.IngressData, wq *sync.WaitGroup) {
	if pod != nil {
		for _, v := range ingress {
			if (pod.SName == v.Sselector) && (pod.Namespace == v.Namespace) {
				pod.URL = append(pod.URL, v.Host+v.Path)
			}
		}
	}
	wq.Done()
}
