package getdata

import (
	"context"
	"ingresslist/getkube/client"
	"strconv"
	"sync"

	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
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

var kube kubernetes.Interface
var podtosertoing PodtoSertoIng
var tmpPodflag bool = false
var tmpPodMu sync.Mutex

func Setup(cfg *client.ClientGo) {
	podtosertoing = PodtoSertoIng{}
	kube = cfg.Client
}

func GetKubeData() error {
	//podデータの取得
	if err := getPodData(); err != nil {
		return err
	}
	//serviceデータの取得
	if err := getServiceData(); err != nil {
		return err
	}
	//ingressデータの取得
	if err := getIngressData(); err != nil {
		return err
	}
	// //nodeデータの取得
	if err := getNodeData(); err != nil {
		return err
	}
	tmpPodflag = true
	return nil
}

func getPodData() error {
	tmp := []PodData{}
	//podデータの取得
	pod, err := kube.CoreV1().Pods("").List(context.TODO(), meta_v1.ListOptions{})
	if err != nil {
		return err
	}
	for _, v := range pod.Items {
		port := []string{}
		if v.Spec.HostNetwork == true {
			for _, vv := range v.Spec.Containers {
				for _, vvv := range vv.Ports {
					if vvv.HostPort != 0 {
						port = append(port, strconv.Itoa(int(vvv.HostPort)))
					}
				}
			}
		}
		tmp = append(tmp, PodData{
			Name:        v.Name,
			Namespace:   v.Namespace,
			PodInfo:     v.Annotations["podinfo"],
			Selector:    v.Labels["app"],
			Ip:          v.Status.PodIP,
			HostNetwork: v.Spec.HostNetwork,
			Port:        port,
		})
	}
	tmpPodMu.Lock()
	podtosertoing.Pod = tmp
	tmpPodMu.Unlock()
	return nil
}

func getServiceData() error {
	var tmp []ServiceData
	//serviceデータの取得
	service, err := kube.CoreV1().Services("").List(context.TODO(), meta_v1.ListOptions{})
	if err != nil {
		return err
	}
	for _, v := range service.Items {
		tmp = append(tmp, ServiceData{
			Name:      v.Name,
			Namespace: v.Namespace,
			Aselector: v.Spec.Selector["app"],
			ClusterIP: v.Spec.ClusterIP,
		})
	}
	tmpPodMu.Lock()
	podtosertoing.Service = tmp
	tmpPodMu.Unlock()
	return nil
}

func getIngressData() error {
	var tmp []IngressData
	//ingressデータの取得
	ingress, err := kube.NetworkingV1().Ingresses("").List(context.TODO(), meta_v1.ListOptions{})
	if err != nil {
		return err
	}
	for _, v := range ingress.Items {
		for _, w := range v.Spec.Rules {
			for _, x := range w.HTTP.Paths {
				tmp = append(tmp, IngressData{
					Namespace: v.Namespace,
					Host:      w.Host,
					Path:      x.Path,
					Sselector: x.Backend.Service.Name,
				})
			}
		}
	}
	tmpPodMu.Lock()
	podtosertoing.Ingress = tmp
	tmpPodMu.Unlock()
	return nil
}

func getNodeData() error {
	var tmp []NodeData
	//nodeデータの取得
	node, err := kube.CoreV1().Nodes().List(context.TODO(), meta_v1.ListOptions{})
	if err != nil {
		return err
	}
	for _, v := range node.Items {
		tmp = append(tmp, NodeData{
			Name: v.Name,
			Ip:   v.Status.Addresses[0].Address,
		})
	}
	tmpPodMu.Lock()
	podtosertoing.Node = tmp
	tmpPodMu.Unlock()
	return nil
}
