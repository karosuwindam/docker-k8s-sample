## 感想
使ってみたけど、epubのビューアガweb版だと使いづらいのでこれはないと感じた

## 適用
```
kubectl apply -f caliber/k8s/namespace
kubectl apply -f caliber/k8s/env-config.yml
kubectl apply -f caliber/k8s/pvc.yml
kubectl apply -f caliber/k8s/pod.yml
kubectl apply -f caliber/k8s/service.yml
```
https://hub.docker.com/r/linuxserver/calibre