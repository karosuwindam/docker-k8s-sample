# Docker-composeやkubenetesを動かすためのデータ格納ファイル

* grafana-prometesu \
arm用のgrafana-prometesu監視がすべてできるようにしたファイル

* cadvisor-exporter \
  prometesu監視できるようなエージェント起動ファイル

* dockeri2c \
  prometesuが監視できるエージェントセット \
  CPU温度を測定する機能付き

* nextcloud
1. arm \
arm 用のnextcloud読み込みファイル

* other
1. raspi-cpu-temp \
  ラズベリーパイのCPU温度図るだけのprometesu拡張ファイル


inginx-ingressのインストール
 ```
  kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v0.46.0/deploy/static/provider/baremetal/deploy.yaml
```
動作確認
```
kubectl get all -n ingress-nginx
```

32bitのOSだとDockerの時刻取得に失敗するので以下のコマンドで本体に導入する
```
wget https://ftp.debian.org/debian/pool/main/libs/libseccomp/libseccomp2_2.5.4-1+b2_armhf.deb
sudo dpkg -i libseccomp2_2.5.4-1+b2_armhf.deb 
rm -rf libseccomp2_2.5.4-1+b2_armhf.deb
```


## 復旧手順について
```
kubectl label node/bookserver2 node-role.kubernetes.io/master=bookserver2
kubectl label node/bookserver2 type=bookserver2
kubectl label node/k8s-worker-1 node-role.kubernetes.io/worker=k8s-worker-1
kubectl label node/k8s-worker-1 type=k8s-worker-1
kubectl label node/k8s-worker-2 node-role.kubernetes.io/worker=k8s-worker-2
kubectl label node/k8s-worker-2 type=k8s-worker-2
kubectl label node/k8s-worker-3 type=k8s-worker-3
kubectl label node/k8s-worker-3 node-role.kubernetes.io/worker=k8s-worker-3
kubectl label node/raspberrypi5 node-role.kubernetes.io/worker=raspberrypi5
kubectl label node/raspberrypi5 type=raspberrypi5
```

```
kubectl apply -f pvd/kuberente-pv.yaml 
kubectl apply -f https://raw.githubusercontent.com/metallb/metallb/v0.13.9/config/manifests/metallb-native.yaml
kubectl apply -f metallb/metallb.yaml 
kubectl apply -f inggress/controller-v1.2.0-deploy.yaml 
```

```
kubectl apply -f docker-registry/k8s/namespace/
kubectl apply -f docker-registry/k8s/volume/
kubectl apply -f docker-registry/k8s/arm/
kubectl apply -f docker-registry/k8s/ingress/

kubectl apply -f grafana-prometesus/k8s/namespace/
kubectl apply -f grafana-prometesus/k8s/role/
kubectl apply -f grafana-prometesus/k8s/account/
kubectl apply -f grafana-prometesus/k8s/kube-state-metrics/account/
kubectl apply -f grafana-prometesus/k8s/kube-state-metrics/role/
kubectl apply -f victoriametrics/k8s/deployment/deployment.yml
kubectl apply -f grafana-prometesus/k8s/kube-state-metrics/pod/

kubectl label nodes raspberrypi5 i2c=true
kubectl label nodes k8s-worker-1 i2c=true
kubectl label nodes bookserver2 i2c=true

kubectl apply -f grafana-prometesus/k8s/volume/
kubectl apply -f grafana-prometesus/k8s/config/
kubectl apply -f grafana-prometesus/k8s/pod/
kubectl apply -f grafana-prometesus/k8s/service/
kubectl apply -f grafana-prometesus/k8s/ingress/ 
```

```
kubectl apply -f smb-csi/deployment/rbac-csi-smb.yaml
kubectl apply -f smb-csi/deployment/csi-smb-driver.yaml
kubectl apply -f smb-csi/deployment/csi-smb-controller.yaml
kubectl apply -f smb-csi/deployment/csi-smb-node.yaml

kubectl apply -f client-go/k8s/role/
kubectl apply -f client-go/k8s/acount/
kubectl apply -f client-go/k8s/

kubectl apply -f booknewread/k8s/namesapace/
kubectl apply -f booknewread/k8s/storage/
kubectl apply -f booknewread/k8s/pod/
kubectl apply -f booknewread/k8s/ingress/


kubectl apply -f isbm_server/k8s/namesapace/
kubectl apply -f isbm_server/k8s/volume/
kubectl apply -f isbm_server/k8s/pod/
kubectl apply -f isbm_server/k8s/ingress/

kubectl apply -f nextcloud/k8s/namespace/
kubectl apply -f nextcloud/k8s/volume/
kubectl apply -f nextcloud/k8s/pod/

kubectl apply -f loki/account
kubectl apply -f loki/configmap
kubectl apply -f loki/volume
kubectl apply -f loki/pod

kubectl apply -f buildkit/k8s

kubectl apply -f kube-web-view/k8s

kubectl apply -f pyroscorpe/deploy.yaml
```