package kubeget

import (
	"context"

	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type IngressData struct {
	Namespace string
	Host      string
	Path      string
	Sselector string
}

func GetIngressData(ctx context.Context) ([]IngressData, error) {
	output := []IngressData{}
	//ingressデータの取得
	ingress, err := kube.NetworkingV1().Ingresses("").List(ctx, meta_v1.ListOptions{})
	if err != nil {
		return output, err
	}
	for _, item := range ingress.Items {
		for _, rule := range item.Spec.Rules {
			for _, path := range rule.HTTP.Paths {
				output = append(output, IngressData{
					Namespace: item.Namespace,
					Host:      rule.Host,
					Path:      path.Path,
					Sselector: path.Backend.Service.Name,
				})
			}
		}
	}
	return output, nil
}
