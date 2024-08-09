package kubeget

import (
	"context"

	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ServiceData struct {
	Name      string
	Namespace string
	Aselector string
	ClusterIP string
}

func GetServiceData(ctx context.Context) ([]ServiceData, error) {
	output := []ServiceData{}
	//serviceデータの取得
	service, err := kube.CoreV1().Services("").List(ctx, meta_v1.ListOptions{})
	if err != nil {
		return output, err
	}
	for _, item := range service.Items {
		output = append(output, ServiceData{
			Name:      item.Name,
			Namespace: item.Namespace,
			Aselector: item.Spec.Selector["app"],
			ClusterIP: item.Spec.ClusterIP,
		})
	}
	return output, nil
}
