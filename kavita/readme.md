# kavita

## 概要
図書管理ソフト

## デプロイ
```
kubectl apply -f kavita/k8s/namesapce
kubectl -n jellyfin create secret generic smbcreds --from-literal username=USERNAME --from-literal password="PASSWORD"
kubectl apply -f kavita/k8s/volume
kubectl apply -f kavita/k8s/pod
kubectl apply -f kavita/k8s/service
```

参考
[ホームページ](https://www.kavitareader.com/)
[github](https://github.com/Kareadita/Kavita)
[dockerhub](https://hub.docker.com/r/kizaing/kavita)
[マニュアル](https://wiki.kavitareader.com/en)