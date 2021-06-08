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