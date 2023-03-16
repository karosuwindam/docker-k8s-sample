package main

import (
	"context"
	"encoding/json"
	"log"
	"sort"
	"strconv"
	"strings"
	"time"

	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type NodeData struct {
	Name string
	Ip   string
}

type PodData struct {
	Name        string
	Namespace   string
	PodInfo     string
	Selector    string
	Ip          string
	HostNetwork bool
	Port        []string
}

type ServiceData struct {
	Name      string
	Namespace string
	Aselector string
	ClusterIP string
}

type IngressData struct {
	Namespace string
	Host      string
	Path      string
	Sselector string
}

type PodtoSertoIng struct {
	Pod     []PodData
	Service []ServiceData
	Ingress []IngressData
	Node    []NodeData
}

type OutputData struct {
	PName       string   `json:podname`
	Namespace   string   `json:namspace`
	Podinfo     string   `json:podinfo`
	Port        []string `json:port`
	Ip          string   `json:ip`
	ClusterIP   string   `json:ClusterIP`
	HostNetwork bool     `json:hostnetwork`
	SName       string   `json:servicename`
	Selector    string   `json:selector`
	URL         string   `json:url`
}

type ListOutputData struct {
	Domain string       `json:domain`
	Data   []OutputData `json:data`
}

func AnaListUrl(data []OutputData) []string {
	output := []string{}
	for _, tmp := range data {
		if tmp.URL != "" {
			i := strings.Index(tmp.URL, ".")
			if i >= 0 {
				j := strings.Index(tmp.URL, "/")
				var temp string
				if j >= 0 {
					temp = "*" + tmp.URL[i:j]
				} else {
					temp = "*" + tmp.URL[i:]
				}
				flag := true
				if len(output) != 0 {
					for _, str := range output {
						if temp == str {
							flag = false
							break
						}
					}
				}
				if flag {
					output = append(output, temp)
				}
			}

		}
	}
	return output
}

func SetPoddata(inputData PodtoSertoIng) []OutputData {
	pods := inputData.Pod
	services := inputData.Service
	ingresses := inputData.Ingress
	output := []OutputData{}

	for _, pod := range pods {
		var tmp OutputData
		tmp.PName = pod.Name
		tmp.Namespace = pod.Namespace
		tmp.Podinfo = pod.PodInfo
		tmp.Selector = pod.Selector
		tmp.Ip = pod.Ip
		tmp.HostNetwork = pod.HostNetwork
		tmp.Port = pod.Port
		output = append(output, tmp)
	}
	ch := make(chan bool, 4)
	suboutput := []OutputData{}
	for i, data := range output {
		ch <- true
		// num := i

		go func(num int, datatmp OutputData) {
			for _, service := range services {
				if (service.Aselector == datatmp.Selector) && (service.Namespace == datatmp.Namespace) {
					output[num].SName = service.Name
					output[num].ClusterIP = service.ClusterIP
					backup_output := output[num]
					count := 0
					for _, ingress := range ingresses {
						if (ingress.Sselector == service.Name) && (ingress.Namespace == datatmp.Namespace) {
							if count != 0 {
								backup_output.URL = ingress.Host + ingress.Path
								suboutput = append(suboutput, backup_output)
							} else {
								output[num].URL = ingress.Host + ingress.Path
							}
							count++
						}
					}
					break
				}
			}
			<-ch
		}(i, data)
	}
	for {
		if len(ch) == 0 {
			break
		}
		time.Sleep(time.Nanosecond)
	}
	output = append(output, suboutput...)
	return output
}

func GetPodIngres(t ClientGo) []OutputData {

	pods, err := t.Client.CoreV1().Pods("").List(context.TODO(), meta_v1.ListOptions{})
	if err != nil {
		log.Fatal(err)
	}
	// _ = SetPoddata(pods)
	ingresses, err := t.Client.NetworkingV1().Ingresses("").List(context.TODO(), meta_v1.ListOptions{})
	// ingresses, err := t.Client.ExtensionsV1beta1().Ingresses("").List(context.TODO(), meta_v1.ListOptions{})

	if err != nil {
		log.Fatal(err)
	}
	services, err := t.Client.CoreV1().Services("").List(context.TODO(), meta_v1.ListOptions{})
	if err != nil {
		log.Fatal(err)
	}

	nodes, err1 := t.Client.CoreV1().Nodes().List(context.TODO(), meta_v1.ListOptions{})
	if err1 != nil {
		log.Fatal(err1)
	}
	ingressTmp := []IngressData{}
	servicesTmp := []ServiceData{}
	podTmp := []PodData{}
	nodedata := []NodeData{}

	ch := make(chan bool, 4)
	go func() {
		for _, ingress := range ingresses.Items {
			for _, rule := range ingress.Spec.Rules {
				for _, path := range rule.IngressRuleValue.HTTP.Paths {
					var tmp IngressData
					tmp.Namespace = ingress.Namespace
					tmp.Host = rule.Host
					tmp.Path = path.Path
					tmp.Sselector = path.Backend.Service.Name
					ingressTmp = append(ingressTmp, tmp)
					// fmt.Println(rule.Host, path.Path, path.Backend.ServiceName)
				}
			}
		}
		ch <- true
	}()
	go func() {
		for _, service := range services.Items {
			var tmp ServiceData
			tmp.Name = service.ObjectMeta.Name
			tmp.Namespace = service.Namespace
			if service.Spec.ClusterIP != "None" {
				tmp.ClusterIP = service.Spec.ClusterIP
			}
			tmp.Aselector = service.Spec.Selector["app"]
			servicesTmp = append(servicesTmp, tmp)
			// fmt.Println(service.ObjectMeta.Name, service.Spec.Selector["app"])
		}
		ch <- true
	}()
	go func() {
		for _, pod := range pods.Items {
			var tmp PodData
			tmp.Name = pod.Name
			tmp.Namespace = pod.Namespace
			tmp.Selector = pod.Labels["app"]
			tmp.PodInfo = pod.Annotations["podinfo"]
			tmp.Ip = pod.Status.PodIP
			tmp_port := []string{}
			for _, container := range pod.Spec.Containers {
				for _, port := range container.Ports {
					if port.ContainerPort != 0 {
						tmp_port = append(tmp_port, strconv.Itoa(int(port.ContainerPort)))
					}
				}
			}
			tmp.Port = tmp_port
			tmp.HostNetwork = pod.Spec.HostNetwork
			podTmp = append(podTmp, tmp)
			// fmt.Println(pod.Name, pod.Labels["app"], pod.Annotations["podinfo"])
		}
		ch <- true
	}()
	go func() {
		for _, node := range nodes.Items {
			var tmp NodeData
			tmp.Name = node.Name
			for _, address := range node.Status.Addresses {
				if address.Type == "InternalIP" {
					tmp.Ip = address.Address
					break
				}
			}
			// fmt.Println(node)
			nodedata = append(nodedata, tmp)
		}
		ch <- true
		// fmt.Println(nodedata)
	}()
	for {
		if len(ch) == 4 {
			break
		}
		time.Sleep(time.Microsecond * 10)
	}
	var inputData PodtoSertoIng
	inputData.Pod = podTmp
	inputData.Node = nodedata
	inputData.Service = servicesTmp
	inputData.Ingress = ingressTmp
	return SetPoddata(inputData)
}

func MargeData(data []OutputData, list []string) []ListOutputData {
	output := []ListOutputData{}
	for _, str := range list {
		tmp := ListOutputData{}
		tmp.Data = []OutputData{}
		tmp.Domain = str
		for _, datat := range data {
			if strings.Index(datat.URL, str[1:]) >= 0 {
				tmp.Data = append(tmp.Data, datat)
			}
		}
		output = append(output, tmp)
	}
	tmp2 := ListOutputData{Domain: ""}
	for _, datat := range data {
		if datat.URL == "" {
			tmp2.Data = append(tmp2.Data, datat)
		}
	}
	output = append(output, tmp2)
	return output
}

func ckdata(a, b []OutputData) bool {
	sort.Slice(a, func(i, j int) bool { return a[i].PName < a[j].PName })
	sort.Slice(b, func(i, j int) bool { return b[i].PName < b[j].PName })
	if len(a) != len(b) {
		return false
	} else {
		for num, tmp := range a {
			if tmp.PName != b[num].PName {
				return false
			}
		}
	}
	return true
}

var server WebSetupData

const (
	KUBECONFIG_ON  = 0
	KUBECONFIG_OFF = 1
)

func main() {
	ch := make(chan bool, 1)
	t := CSetup(KUBECONFIG_OFF)
	output_bak := []OutputData{}

	go func() {
		for {
			output := GetPodIngres(t)
			list := AnaListUrl(output)
			// for _, data := range output {
			// 	fmt.Println(data)
			// }
			jsond := MargeData(output, list)
			server.Output = jsond
			jsondata, _ := json.Marshal(jsond)
			if !ckdata(output, output_bak) {
				log.Println(string(jsondata))
			}
			// for _, str := range list {
			// 	log.Println(str)
			// 	for _, data := range output {
			// 		if strings.Index(data.URL, str[1:]) >= 0 {
			// 			log.Println(data)
			// 		}
			// 	}
			// }
			if len(ch) == 0 {
				ch <- true
			}
			time.Sleep(time.Second * 10)
			output_bak = output
		}
	}()
	err := server.websetup()
	if err != nil {
		log.Println(err.Error())
		return
	}
	<-ch
	server.webstart()

}
