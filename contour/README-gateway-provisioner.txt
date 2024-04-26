
# start a new cluster

make install-provisioner-working

# or

kind delete cluster --name contour
kind create cluster --config ~/work/devenvs/contour/configs/kind-cluster-config.yaml --name contour
kubectl apply -f https://projectcontour.io/quickstart/contour-gateway-provisioner.yaml





kubectl apply -f manifests/gateway-provisioner.yaml
kubectl apply -f manifests/echoserver-gatewayapi.yaml

# address is from metallb ip pool
http http://echoserver-172.20.255.200.nip.io


$ kubectl -n metallb-system get ipaddresspools.metallb.io pool
NAME   AUTO ASSIGN   AVOID BUGGY IPS   ADDRESSES
pool   true          false             ["172.20.255.200-172.20.255.250"]

$ kubectl -n projectcontour get service envoy-mygateway
NAME              TYPE           CLUSTER-IP     EXTERNAL-IP      PORT(S)        AGE
envoy-mygateway   LoadBalancer   10.96.135.70   172.20.255.200   80:30093/TCP   34m






