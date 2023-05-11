
# add unique local IPv6 address to access the cluster from host
sudo ip -6 addr add fd61:97d2:3f5a:16e0::1 dev lo
sudo ip -6 route add to local fd61:97d2:3f5a:16e0::/64 dev lo


kind delete cluster --name contour
kind create cluster --config ~/work/devenvs/contour/configs/kind-cluster-config-ipv6-only.yaml --name contour

kubectl apply -f https://projectcontour.io/quickstart/contour.yaml



# running contour locally on dev machine (IPAM.Config[1] is for IPv6 address)
sed "s/REPLACE_ADDRESS_HERE/$(docker network inspect kind | jq -r '.[0].IPAM.Config[1].Gateway')/" manifests/contour-endpoints-dev.yaml | kubectl apply -f -
kubectl -n projectcontour scale deployment --replicas=0 contour
kubectl -n projectcontour rollout restart daemonset envoy


kubectl -n projectcontour get secret contourcert -o jsonpath='{..ca\.crt}' | base64 -d > ca.crt
kubectl -n projectcontour get secret contourcert -o jsonpath='{..tls\.crt}' | base64 -d > tls.crt
kubectl -n projectcontour get secret contourcert -o jsonpath='{..tls\.key}' | base64 -d > tls.key

go run github.com/projectcontour/contour/cmd/contour serve --xds-address=:: --xds-port=8001 --envoy-service-http-address=::  --envoy-service-http-port=8080 --envoy-service-https-port=8443 --contour-cafile=ca.crt --contour-cert-file=tls.crt --contour-key-file=tls.key --config-path=$HOME/work/devenvs/contour/configs/contour-dscp.yaml



kubectl apply -f manifests/echoserver-ipv6.yaml

http http://echoserver.fd61-97d2-3f5a-16e0--1.sslip.io
