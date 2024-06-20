


# https://github.com/projectcontour/contour/pull/6269



# https://prometheus-operator.dev/docs/prologue/quick-start/

git clone https://github.com/prometheus-operator/kube-prometheus.git

kubectl create -f manifests/setup
# Wait until the "servicemonitors" CRD is created. The message "No resources found" means success in this context.
until kubectl get servicemonitors --all-namespaces ; do date; sleep 1; echo ""; done

kubectl create -f manifests/



kubectl apply -f examples/prometheus/httpproxy.yaml
kubectl apply -f examples/grafana/httpproxy.yaml


