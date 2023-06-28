# SMB CSI Driver for Kubernetes

## ソースについて
v1.9.0をダウンロードする
なお、v1.10.0だとarmは動かない

## install
```
kubectl apply -f ./rbac-csi-smb.yaml
kubectl apply -f ./csi-smb-driver.yaml
kubectl apply -f ./csi-smb-controller.yaml
kubectl apply -f ./csi-smb-node.yaml
```

以下のコマンドで、StorageCrassを作成
```
kubectl apply -f ./sc/sc.yml
```

## 使い方
以下のコマンドで、パスワードを登録する
```
kubectl -n default create secret generic smbcreds --from-literal username=USERNAME --from-literal password="PASSWORD"
kubectl -n default create secret generic k8s-1-smb --from-literal username=USERNAME --from-literal password="PASSWORD"
```

pv-pvcで設定することで、deploymentで定義したPodにマウントできる


## 参考
[github](https://github.com/kubernetes-csi/csi-driver-smb)