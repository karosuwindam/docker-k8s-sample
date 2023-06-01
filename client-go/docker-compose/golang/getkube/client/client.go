package client

import "k8s.io/client-go/kubernetes"

type ClientGo struct {
	Client kubernetes.Interface
	Flag   bool
}
