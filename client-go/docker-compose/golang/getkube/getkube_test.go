package getkube

import (
	"ingresslist/config"
	"ingresslist/getkube/getdata"
	"os"
	"testing"
)

func TestKube(t *testing.T) {
	//環境変数KUBE_PATH_FLAGの設定
	os.Setenv("KUBE_PATH_FLAG", "true")
	os.Setenv("KUBE_PATH", "../kubeconfig")

	cfg, _ := config.Setup()
	kube := Setup(cfg)
	if kube.Flag == false {
		t.Errorf("Kubeconfig is not set")
	}
	getdata.GetKubeData()
	getdata.CreateData()
}
