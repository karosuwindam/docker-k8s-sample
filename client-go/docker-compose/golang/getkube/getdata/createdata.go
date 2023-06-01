package getdata

import (
	"errors"
	"sync"
)

type OutputData struct {
	PName       string   `json:"podname"`
	Namespace   string   `json:"namspace"`
	Podinfo     string   `json:"podinfo"`
	Port        []string `json:"port"`
	Ip          string   `json:"ip"`
	ClusterIP   string   `json:"ClusterIP"`
	HostNetwork bool     `json:"hostnetwork"`
	SName       string   `json:"servicename"`
	Selector    string   `json:"selector"`
	URL         []string `json:"url"`
}

func CreateData() ([]OutputData, error) {
	if tmpPodflag == false {
		return nil, errors.New("Pod Data is not set")
	}
	tmpPodMu.Lock()
	tmp := podtosertoing
	tmpPodMu.Unlock()

	output := []OutputData{}
	//podデータからベースの作成
	for _, pod := range tmp.Pod {
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
		go AddServiceData(&output[i], tmp.Service, &wq)
	}
	wq.Wait()
	for i := 0; i < len(output); i++ {
		wq.Add(1)
		go AddIngressData(&output[i], tmp.Ingress, &wq)
	}
	wq.Wait()

	return output, nil
}

// PodデータをもとにServiceデータを追加する
func AddServiceData(pod *OutputData, service []ServiceData, wq *sync.WaitGroup) {
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
func AddIngressData(pod *OutputData, ingress []IngressData, wq *sync.WaitGroup) {
	if pod != nil {
		for _, v := range ingress {
			if (pod.SName == v.Sselector) && (pod.Namespace == v.Namespace) {
				pod.URL = append(pod.URL, v.Host+v.Path)
			}
		}
	}
	wq.Done()
}
