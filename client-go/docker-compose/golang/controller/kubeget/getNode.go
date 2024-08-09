package kubeget

import (
	"context"

	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type NodeData struct {
	Name string
	Ip   string
}

func GetNodeData(ctx context.Context) ([]NodeData, error) {
	output := []NodeData{}
	//nodeデータの取得
	node, err := kube.CoreV1().Nodes().List(ctx, meta_v1.ListOptions{})
	if err != nil {
		return output, err
	}
	for _, item := range node.Items {
		output = append(output, NodeData{
			Name: item.Name,
			Ip:   item.Status.Addresses[0].Address,
		})
	}
	return output, nil
}
