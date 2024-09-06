


#### flannel

kind delete cluster --name echo
kind create cluster --config configs/kind-cluster-config-for-flannel.yaml --name echo

kubectl apply -f https://github.com/flannel-io/flannel/releases/latest/download/kube-flannel.yml

# ERROR
#   pod will not start and report error:   failed to find plugin "bridge" in path [/opt/cni/bin]
#
# WORKAROUND
#   install bridge plugin manually
#   -   https://routemyip.com/posts/k8s/setup/flannel/
#   -   https://github.com/containernetworking/plugins/releases

CNI_PLUGIN_URL=https://github.com/containernetworking/plugins/releases/download/v1.3.0/cni-plugins-linux-amd64-v1.3.0.tgz

docker exec echo-control-plane sh -c "curl -L -o /tmp/cni-plugins.tgz $CNI_PLUGIN_URL && tar -C /opt/cni/bin -xzvf /tmp/cni-plugins.tgz && rm /tmp/cni-plugins.tgz"
docker exec echo-worker sh -c "curl -L -o /tmp/cni-plugins.tgz $CNI_PLUGIN_URL && tar -C /opt/cni/bin -xzvf /tmp/cni-plugins.tgz && rm /tmp/cni-plugins.tgz"
docker exec echo-worker2 sh -c "curl -L -o /tmp/cni-plugins.tgz $CNI_PLUGIN_URL && tar -C /opt/cni/bin -xzvf /tmp/cni-plugins.tgz && rm /tmp/cni-plugins.tgz"

# check status after install
kubectl get pod -A


kubectl apply -f manifests/echo.yaml
kubectl get pod -o wide

kubectl delete pod -l app=client --force
kubectl delete pod -l app=server --force

kubectl logs -l app=client -f


# Backends
#    https://github.com/flannel-io/flannel/blob/master/Documentation/backends.md
#
# vxlan
#    - UDP￼port 8472
# ipip

kubectl -n kube-flannel get configmap kube-flannel-cfg -o yaml
kubectl -n kube-flannel edit configmap kube-flannel-cfg

# restart after config change
kubectl -n kube-flannel delete pod -l app=flannel --force




#### cilium

kind delete cluster --name echo
kind create cluster --config configs/kind-cluster-config-for-cilium.yaml --name echo

# install
#    https://docs.cilium.io/en/stable/installation/kind/
helm repo add cilium https://helm.cilium.io/
helm repo update

helm install cilium cilium/cilium --version 1.14.3 --namespace kube-system --set image.pullPolicy=IfNotPresent --set ipam.mode=kubernetes

# check status after install
kubectl -n kube-system get pods

kubectl apply -f manifests/echo.yaml
kubectl get pod -o wide


kubectl delete pod -l app=client --force
kubectl delete pod -l app=server --force

kubectl logs -l app=client -f




### antrea

# https://antrea.io/
# https://github.com/antrea-io/antrea
# TODO (already tested by L. that it also works expectedly





#######



# capture traffic from workers
sudo nsenter -t $(docker inspect --format '{{.State.Pid}}' echo-worker) --net wireshark -i any -k -o gui.window_title:worker -Y "tcp.port==8000"
sudo nsenter -t $(docker inspect --format '{{.State.Pid}}' echo-worker2) --net wireshark -i any -k -o gui.window_title:worker2 -Y "tcp.port==8000"
