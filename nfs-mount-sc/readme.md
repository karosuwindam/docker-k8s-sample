kubectl apply -f https://raw.githubusercontent.com/kubernetes-sigs/nfs-ganesha-server-and-external-provisioner/master/deploy/kubernetes/deployment.yaml
kubectl apply -f https://raw.githubusercontent.com/kubernetes-sigs/nfs-ganesha-server-and-external-provisioner/master/deploy/kubernetes/rbac.yaml
kubectl apply -f https://raw.githubusercontent.com/kubernetes-sigs/nfs-ganesha-server-and-external-provisioner/master/deploy/kubernetes/class.yaml
kubectl apply -f https://raw.githubusercontent.com/kubernetes-sigs/nfs-ganesha-server-and-external-provisioner/master/deploy/kubernetes/claim.yaml
