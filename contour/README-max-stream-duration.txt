https://github.com/projectcontour/contour/issues/6873
https://github.com/projectcontour/contour/pull/6895





cat <<EOF | kubectl apply -f -
apiVersion: v1
kind: ConfigMap
metadata:
  name: contour
  namespace: projectcontour
data:
  contour.yaml: |
    timeouts:
      max-stream-duration: infinite
EOF



kubectl apply -f examples/contour/01-crds.yaml

cat <<EOF | kubectl apply -f -
apiVersion: projectcontour.io/v1alpha1
kind: ContourConfiguration
metadata:
  name: contour
  namespace: projectcontour
spec:
  envoy:
    timeouts:
      maxStreamDuration: 3s
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


make container
docker tag ghcr.io/projectcontour/contour:$(git rev-parse --short HEAD) localhost/contour:latest
kind load docker-image localhost/contour:latest --name contour

kubectl -n projectcontour set image deployment/contour contour=localhost/contour:latest



# restart contour
kubectl -n projectcontour scale deployment contour --replicas=0
kubectl -n projectcontour scale deployment contour --replicas=1


http http://echoserver.127-0-0-101.nip.io/sse




cat <<EOF | kubectl apply -f -
apiVersion: projectcontour.io/v1
kind: HTTPProxy
metadata:
  name: echoserver-sse
spec:
  virtualhost:
    fqdn: echoserver-sse.127-0-0-101.nip.io
  routes:
    - services:
        - name: echoserver
          port: 80
      timeoutPolicy:
        response: infinite
        maxStreamDuration: 6s
EOF


http http://echoserver-sse.127-0-0-101.nip.io/sse
