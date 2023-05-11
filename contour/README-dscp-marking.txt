
https://github.com/projectcontour/contour/issues/4605

kind delete cluster --name contour
kind create cluster --config ~/work/devenvs/contour/configs/kind-cluster-config.yaml --name contour

kubectl apply -f https://projectcontour.io/quickstart/contour.yaml

sed "s/REPLACE_ADDRESS_HERE/$(docker network inspect kind | jq -r '.[0].IPAM.Config[0].Gateway')/" manifests/contour-endpoints-dev.yaml | kubectl apply -f -
kubectl -n projectcontour scale deployment --replicas=0 contour
kubectl -n projectcontour rollout restart daemonset envoy

kubectl apply -f manifests/echoserver.yaml

kubectl -n projectcontour get secret contourcert -o jsonpath='{..ca\.crt}' | base64 -d > ca.crt
kubectl -n projectcontour get secret contourcert -o jsonpath='{..tls\.crt}' | base64 -d > tls.crt
kubectl -n projectcontour get secret contourcert -o jsonpath='{..tls\.key}' | base64 -d > tls.key


go run github.com/projectcontour/contour/cmd/contour serve --xds-address=0.0.0.0 --xds-port=8001 --envoy-service-http-port=8080 --envoy-service-https-port=8443 --contour-cafile=ca.crt --contour-cert-file=tls.crt --contour-key-file=tls.key --config-path=$HOME/work/devenvs/contour/configs/contour-dscp.yaml

kubectl -n projectcontour logs daemonsets/envoy envoy -f

http http://echoserver.127-0-0-101.nip.io


# DSCP value will be visible if capturing from the envoy network namespace
sudo nsenter --target $(pidof envoy) --net wireshark  -f "port 8080" -k -Y http

# DSCP value will not be visible if capturing from the host network namespace
wireshark -i lo -f "port 80" -k -Y http



kubectl -n projectcontour port-forward daemonset/envoy 9001:9001
http http://localhost:9001/config_dump | jq -C . | less




kubectl -n projectcontour logs daemonsets/envoy -c envoy



kind delete cluster --name contour
kind create cluster --config ~/work/devenvs/contour/configs/kind-cluster-config-ipv6-only.yaml --name contour


sed "s/REPLACE_ADDRESS_HERE/$(docker network inspect kind | jq -r '.[0].IPAM.Config[1].Gateway')/" manifests/contour-endpoints-dev.yaml | kubectl apply -f -
kubectl -n projectcontour scale deployment --replicas=0 contour
kubectl -n projectcontour rollout restart daemonset envoy


go run github.com/projectcontour/contour/cmd/contour serve  --xds-address=:: --xds-port=8001 --envoy-service-http-address=::  --envoy-service-http-port=8080 --envoy-service-https-port=8443 --contour-cafile=ca.crt --contour-cert-file=tls.crt --contour-key-file=tls.key --config-path=$HOME/work/devenvs/contour/configs/contour-dscp.yaml

http http://echoserver.fd61-97d2-3f5a-16e0--1.sslip.io
