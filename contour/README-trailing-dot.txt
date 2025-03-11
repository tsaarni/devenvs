

# https://github.com/projectcontour/contour/issues/6334
# https://github.com/projectcontour/contour/pull/6792

kubectl apply -f https://projectcontour.io/quickstart/contour.yaml

sudo nsenter --target $(pidof envoy) --net wireshark -i any -k -f "tcp port 8080" -Y http

http http://echoserver.127-0-0-101.nip.io./foo   # will work because httpie will remove the trailing dot
curl http://echoserver.127-0-0-101.nip.io./foo   # will not work because curl will not remove the trailing dot


make container
docker tag ghcr.io/projectcontour/contour:$(git rev-parse --short HEAD) localhost/contour:latest
kind load docker-image localhost/contour:latest --name contour
kubectl -n projectcontour set image deployment/contour contour=localhost/contour:latest


# create configuration file for contour

cat <<EOF | kubectl apply -f -
apiVersion: v1
kind: ConfigMap
metadata:
  name: contour
  namespace: projectcontour
data:
  contour.yaml: |
    network:
      strip-trailing-host-dot: true
EOF

# restart contour
kubectl -n projectcontour scale deployment contour --replicas=0
kubectl -n projectcontour scale deployment contour --replicas=1


#########

# create contourconfig crd
kubectl apply -f examples/contour/01-crds.yaml

cat <<EOF | kubectl apply -f -
apiVersion: projectcontour.io/v1alpha1
kind: ContourConfiguration
metadata:
  name: contour
  namespace: projectcontour
spec:
  envoy:
    network:
      stripTrailingHostDot: true
EOF

# set contour to use the CRD instead of the config file
kubectl -n projectcontour patch deployment contour --type='json' -p='[
  {
    "op": "replace",
    "path": "/spec/template/spec/containers/0/args",
    "value": [
      "serve",
      "--incluster",
      "--xds-address=0.0.0.0",
      "--xds-port=8001",
      "--contour-cafile=/certs/ca.crt",
      "--contour-cert-file=/certs/tls.crt",
      "--contour-key-file=/certs/tls.key",
      "--contour-config-name=contour"
    ]
  }
]'

# restart contour
kubectl -n projectcontour scale deployment contour --replicas=0
kubectl -n projectcontour scale deployment contour --replicas=1
