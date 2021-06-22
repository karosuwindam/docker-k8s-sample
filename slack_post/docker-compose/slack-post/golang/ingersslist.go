package main

import (
	"context"
	"flag"
	"log"
	"path/filepath"

	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ClientGo struct {
	Client kubernetes.Interface
	flag   bool
}

func (t *ClientGo) GetPod(namespace string) []v1.Pod {
	pods, err := t.Client.CoreV1().Pods(namespace).List(context.TODO(), meta_v1.ListOptions{})
	if err != nil {
		log.Fatal(err)
	}
	return pods.Items
}

func CSetup(num int) ClientGo {
	output := ClientGo{}
	if num == 0 {
		tmp, err := newClient_cli()
		if err == nil {
			output.Client = tmp
			output.flag = true
		} else {
			log.Println(err.Error())
		}

	} else {
		tmp, err := newClient_con()
		if err == nil {
			output.Client = tmp
			output.flag = true
		} else {
			log.Println(err.Error())
		}
	}
	return output
}

func newClient_con() (kubernetes.Interface, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}

	return kubernetes.NewForConfig(config)
}

func newClient_cli() (kubernetes.Interface, error) {
	// kubeConfig, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	var kubeconfig *string
	kubeconfig = flag.String("kubeconfig", filepath.Join("../", ".kube", "config"), "(optional) absolute path to the kubeconfig file")

	kubeConfig, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		return nil, err
	}

	return kubernetes.NewForConfig(kubeConfig)
}
