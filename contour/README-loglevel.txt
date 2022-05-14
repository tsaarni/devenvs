


go run github.com/projectcontour/contour/cmd/contour serve --xds-address=0.0.0.0 --xds-port=8001 --envoy-service-http-port=8080 --envoy-service-https-port=8443 --contour-cafile=ca.crt --contour-cert-file=tls.crt --contour-key-file=tls.key --config-path=$WORKDIR/configs/contour-accesslog-level.yaml



kubectl apply -f manifests/echoserver.yaml

kubectl apply -f manifests/shell.yaml
kubectl exec -it shell -- ash


python3

import socket
s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
s.bind(("",8000))
s.listen(0)

# connect one client to fill in the listen queue
c = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
c.connect(("127.0.0.1",8000))





http http://echoserver.127-0-0-101.nip.io
http http://echoserver.127-0-0-101.nip.io/status/418
http http://shell.127-0-0-101.nip.io   # cancel request with ctrl+c



# Test ContourConfiguration CRD
kubectl apply -f examples/contour/01-crds.yaml


kubectl apply -f manifests/contourconfig-accesslog-level.yaml

go run github.com/projectcontour/contour/cmd/contour serve --xds-address=0.0.0.0 --xds-port=8001 --envoy-service-http-port=8080 --envoy-service-https-port=8443 --contour-cafile=ca.crt --contour-cert-file=tls.crt --contour-key-file=tls.key --contour-config-name=contour
