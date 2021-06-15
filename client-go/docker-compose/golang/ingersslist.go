package main

import (
	"flag"
	"log"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type ClientGo struct {
	Client kubernetes.Interface
	flag   bool
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
