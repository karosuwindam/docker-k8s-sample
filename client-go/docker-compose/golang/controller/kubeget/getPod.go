package kubeget

import (
	"context"
	"strconv"

	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type PodData struct {
	Name        string
	Namespace   string
	PodInfo     string
	Selector    string
	Ip          string
	HostNetwork bool
	Port        []string
}

func GetPodData(ctx context.Context) ([]PodData, error) {
	output := []PodData{}
	pods, err := kube.CoreV1().Pods("").List(ctx, meta_v1.ListOptions{})
	if err != nil {
		return output, err
	}
	for _, pod := range pods.Items {
		ports := []string{}
		if pod.Spec.HostNetwork == true {
			for _, containers := range pod.Spec.Containers {
				for _, port := range containers.Ports {
					if port.HostPort != 0 {
						ports = append(ports, strconv.Itoa(int(port.HostPort)))
					}
				}
			}
		}
		output = append(output, PodData{
			Name:        pod.Name,
			Namespace:   pod.Namespace,
			PodInfo:     pod.Annotations["podinfo"],
			Selector:    pod.Labels["app"],
			Ip:          pod.Status.PodIP,
			HostNetwork: pod.Spec.HostNetwork,
			Port:        ports,
		})
	}
	return output, nil
}
