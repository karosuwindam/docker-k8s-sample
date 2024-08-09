package kubeget

import (
	"flag"
	"ingresslist/config"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var kube kubernetes.Interface

func Init() error {
	if config.Kube.KubePathFlag {
		if cfg, err := k8sClientSetupcli(config.Kube.KubePath); err != nil {
			return err
		} else {
			kube = cfg
		}
	} else {
		if cfg, err := k8sClientSetupCon(); err != nil {
			return err
		} else {
			kube = cfg
		}
	}
	return nil
}

func k8sClientSetupCon() (kubernetes.Interface, error) {
	cfg, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}
	return kubernetes.NewForConfig(cfg)
}

func k8sClientSetupcli(k8sConfPass string) (kubernetes.Interface, error) {
	var kubeConfig *string
	kubeConfig = flag.String("kubeconfig", k8sConfPass, "(optional) absolute path to the kubeconfig file")

	cfg, err := clientcmd.BuildConfigFromFlags("", *kubeConfig)
	if err != nil {
		return nil, err
	}
	return kubernetes.NewForConfig(cfg)
}
