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

* CiliumをHelmでインストールするために以下コマンドでリポジトリを登録する

```Bash
helm repo add cilium https://helm.cilium.io/
helm repo update
```

以下コマンドで、cilium-systemのネームスペースを作成してから、helmコマンドでインストールを実行します。
以下のコマンドは、kube-poroxyの機能置き換えやIngress Controllerの配置をciliumが行うことや、共有ロードバランサモードを設置するようにします。また、Prometheusによる監視を有効にしています。

```bash
helm install cilium cilium/cilium --version 1.16.5 \
--namespace cilium-system \
--create-namespace \
--set kubeProxyReplacement=true \
--set ingressController.enabled=true \
--set ingressController.loadbalancerMode=shared
# オプション: Ingress ControllerとLoad Balancer機能を使う場合: 
```

```bash
helm install cilium cilium/cilium --version 1.16.5 \
--namespace cilium-system \
--create-namespace \
--set kubeProxyReplacement=true \
--set ingressController.enabled=true \
--set ingressController.loadbalancerMode=shared \
--set monitorAggregation=none \
--set metrics.enabled=true \
--set prometheus.enabled=true \
--set hubble.enabled=true \
--set hubble.relay.enabled=true \
--set hubble.ui.enabled=true \
--set hubble.ui.standalone.enabled=true
```

```
kubectl -n cilium-system get pod -w
```

以下の通りの設定で、Prometheusによる監視を有効にする
```bash
helm upgrade cilium cilium/cilium --version 1.16.5 \
--namespace cilium-system \
--reuse-values \
--set monitorAggregation=none \
--set metrics.enabled=true \
--set prometheus.enabled=true
```
```bash
helm upgrade cilium cilium/cilium --version 1.16.5 \
  --namespace cilium-system \
  --reuse-values \
  --set hubble.enabled=true \
  --set hubble.relay.enabled=true \
  --set hubble.ui.enabled=true \
  --set hubble.ui.standalone.enabled=true
```

* inginx-ingressのインストール

flunnel用のデータ設定

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

動いているpodのイメージリストを確認するコマンド
```
kubectl get deployment -A -o  jsonpath="Name{'\t'}Spec{'\t'}Pods{'\t'}image{'\t'}terminationGracePeriodSeconds{'\t'}nodeSelector{'\n'}{range .items[*]}{.metadata.name}{'\t'}{.spec.template.spec.containers[].resources.requests}{'\t'}{.status.replicas}{'\t'}{.spec.template.spec.containers[].image}{'\t'}{.spec.template.spec.terminationGracePeriodSeconds}{'\t'}{.spec.template.spec.nodeSelector.*}{'\n'}{end}"
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
kubectl label node/k8s-worker-4 node-role.kubernetes.io/worker=k8s-worker-4
kubectl label node/k8s-worker-4 type=k8s-worker-4
```


* flunnel用
```
kubectl apply -f pvd/kuberente-pv.yaml 
kubectl apply -f https://raw.githubusercontent.com/metallb/metallb/v0.15.2/config/manifests/metallb-native.yaml
```
```
kubectl apply -f metallb/metallb.yaml 
kubectl apply -f inggress/controller-v1.10.1-deploy.yaml 
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

kubectl label nodes k8s-worker-4 i2c=true
kubectl label nodes k8s-worker-1 i2c=true
kubectl label nodes k8s-worker-2 i2c=true
kubectl label nodes bookserver2 i2c=true
kubectl label nodes bookserver2 gpio=false
kubectl label nodes bookserver2 uart=false
kubectl label nodes k8s-worker-1 gpio=true
kubectl label nodes k8s-worker-1 uart=false
kubectl label nodes k8s-worker-2 gpio=false
kubectl label nodes k8s-worker-2 uart=false
kubectl label nodes k8s-worker-4 gpio=false
kubectl label nodes k8s-worker-4 uart=true


kubectl apply -f grafana-prometesus/k8s/volume/
kubectl apply -f grafana-prometesus/k8s/config/
kubectl apply -f grafana-prometesus/k8s/pod/
kubectl apply -f grafana-prometesus/k8s/service/
kubectl apply -f grafana-prometesus/k8s/ingress/ 
```

```
kubectl apply -f nfs-sc/k8s/deploy.yml
```

```
kubectl apply -f grafana-jaeger/otel-config.yaml
kubectl apply -f grafana-jaeger/otel-ingreass.yml
kubectl apply -f grafana-jaeger/volume
kubectl apply -f grafana-jaeger/pod
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

kubectl apply -f nextcloud/k8s/namespace/
kubectl apply -f nextcloud/k8s/volume/
kubectl apply -f nextcloud/k8s/config/
kubectl apply -f nextcloud/k8s/pod/

kubectl apply -f loki/account
kubectl apply -f loki/configmap
kubectl apply -f loki/volume
kubectl apply -f loki/pod

kubectl apply -f buildkit/k8s

kubectl apply -f kube-web-view/k8s

kubectl apply -f haproxy/ns.yaml
kubectl apply -f haproxy/configmap.yaml
kubectl apply -f haproxy/deploy.yaml
```

longhorn
```
helm repo add longhorn https://charts.longhorn.io
helm repo update

helm install longhorn longhorn/longhorn --namespace longhorn-system --create-namespace \
  --set persistence.defaultDataPath=/mnt/longhorn-data \
  --set persistence.defaultSettings.defaultReplicaCount=3 \
  --set monitoring.enabled=true \
  --set monitoring.prometheus.serviceMonitor.enabled=true

kubectl apply -f longhorn/ingress.yaml
```

nfs://
バックアップ設定する場合は以下の通りのものをURLへ
```
nfs://k8s-worker-2:/home/pi/usb/usb1/nfs/kubernetes/longhorn
```


other app
```
kubectl apply -f caliber/k8s/namespace/
kubectl apply -f caliber/k8s/

kubectl apply -f kavita/k8s/namesapce/
kubectl apply -f kavita/k8s/volume/
kubectl apply -f kavita/k8s/pod
kubectl apply -f kavita/k8s/service
```

jellyfin
```
kubectl apply -f jellyfin/k8s/namespace/

# NAS用
kubectl -n jellyfin create secret generic smbcreds --from-literal username=USERNAME --from-literal password="PASSWORD"
# Next用
kubectl -n jellyfin create secret generic smbpicreds --from-literal username=USERNAME --from-literal password="PASSWORD"

kubectl apply -f jellyfin/k8s/
```


## イメージリスト確認

```
kubectl get pods -A -o=jsonpath='{range .items[*]}{.metadata.namespace}{"\t"}{.metadata.name}{"\t"}{.spec.nodeName}{"\t"}{"\t"}{.spec.containers[].image}{"\n"}{end}'
```