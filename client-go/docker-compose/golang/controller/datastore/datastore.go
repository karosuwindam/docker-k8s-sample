package datastore

import (
	"errors"
	"ingresslist/controller/kubeget"
	"sync"
)

type dataStore struct {
	pod     []kubeget.PodData
	ingress []kubeget.IngressData
	node    []kubeget.NodeData
	service []kubeget.ServiceData
	mux     sync.Mutex
}

var database dataStore

func Init() error {
	database = dataStore{}
	return nil
}

func Write(data interface{}) error {
	database.mux.Lock()
	defer database.mux.Unlock()
	switch data.(type) {
	case []kubeget.IngressData:
		database.ingress = data.([]kubeget.IngressData)
	case []kubeget.NodeData:
		database.node = data.([]kubeget.NodeData)
	case []kubeget.PodData:
		database.pod = data.([]kubeget.PodData)
	case []kubeget.ServiceData:
		database.service = data.([]kubeget.ServiceData)
	default:
		return errors.New("data input type error")
	}
	return nil
}

func GetIngress() []kubeget.IngressData {
	database.mux.Lock()
	defer database.mux.Unlock()
	return database.ingress
}

func GetPod() []kubeget.PodData {
	database.mux.Lock()
	defer database.mux.Unlock()
	return database.pod
}
func GetService() []kubeget.ServiceData {
	database.mux.Lock()
	defer database.mux.Unlock()
	return database.service
}
func GetNode() []kubeget.NodeData {
	database.mux.Lock()
	defer database.mux.Unlock()
	return database.node
}
