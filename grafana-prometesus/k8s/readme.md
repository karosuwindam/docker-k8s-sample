kubectl label nodes raspberrypi5 i2c-
kubectl label nodes raspberrypi5 i2c=true
kubectl label nodes k8s-worker-1 gpio=true
kubectl label nodes k8s-worker-4 uart=true

kubectl get node -o jsonpath="Name{'\t'}{'\t'}i2c{'\t'}gpio{'\t'}uart{'\n'}{range .items[*]}{.metadata.name}{'\t'}{.metadata.labels.i2c}{'\t'}{.metadata.labels.gpio}{'\t'}{.metadata.labels.uart}{'\n'}{end}"

参考URL
https://hub.docker.com/r/prom/node-exporter

https://hub.docker.com/r/grafana/grafana

https://hub.docker.com/r/prom/prometheus