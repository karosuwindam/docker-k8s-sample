package main

import (
	"strconv"

	v1 "k8s.io/api/core/v1"
)

type PodPort struct {
	Protocol      string `json:protocol`
	ContainerPort int    `json:containerport`
	HostPort      int    `json:hostport`
}

type PodInfo struct {
	Name        string            `json:name`
	Namespace   string            `json:namespace`
	Node        string            `json:node`
	Ip          string            `json:ip`
	Ports       []PodPort         `json:ports`
	Labels      map[string]string `json:labels`
	Annotations map[string]string `json:annotations`
}

type UrlData struct {
	Url  string `json:url`
	Node string `json:node`
}

// Podデータからslackpostでtrueになっているpod情報を取得
func AnariseData(datas []PodInfo) []UrlData {
	output := []UrlData{}
	for _, data := range datas {
		if data.Annotations["slackpost"] == "true" {
			for _, port := range data.Ports {
				if port.Protocol == "TCP" {
					var tmp UrlData
					tmp.Node = data.Node
					tmp.Url = "http://" + data.Ip + ":" + strconv.Itoa(int(port.ContainerPort)) + "/json"
					output = append(output, tmp)
				}
			}
		}
	}
	return output
}

//取得したPod情報からセンサーデータを出力
func GetPodInfo(pods []v1.Pod) []PodInfo {
	output := []PodInfo{}
	for _, pod := range pods {
		var tmp PodInfo
		tmp.Name = pod.Name
		tmp.Namespace = pod.Namespace
		tmp.Ip = pod.Status.PodIP
		tmpport := []PodPort{}
		for _, container := range pod.Spec.Containers {
			for _, port := range container.Ports {
				var tport PodPort
				tport.Protocol = string(port.Protocol)
				tport.HostPort = int(port.HostPort)
				tport.ContainerPort = int(port.ContainerPort)
				tmpport = append(tmpport, tport)
			}
		}
		tmp.Node = pod.Spec.NodeName
		tmp.Ports = tmpport
		tmp.Labels = pod.Labels
		tmp.Annotations = pod.Annotations
		output = append(output, tmp)
	}
	return output
}
