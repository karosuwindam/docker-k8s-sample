## 適用
```
kubectl apply -f jellyfin/k8s/namespace
kubectl -n jellyfin create secret generic smbcreds --from-literal username=USERNAME --from-literal password="PASSWORD"
kubectl -n jellyfin create secret generic smbpicreds --from-literal username=USERNAME --from-literal password="PASSWORD"
kubectl apply -f jellyfin/k8s/smb
kubectl apply -f jellyfin/k8s/pvc.yml
kubectl apply -f jellyfin/k8s/env-config.yml
kubectl apply -f jellyfin/k8s/pod-smb.yml
kubectl apply -f jellyfin/k8s/service.yml
```