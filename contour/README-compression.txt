

https://github.com/projectcontour/contour/pull/6546



kind delete cluster --name contour
kind create cluster --config ~/work/devenvs/contour/configs/kind-cluster-config.yaml --name contour
kubectl apply -f https://projectcontour.io/quickstart/contour.yaml

kubectl apply -f manifests/echoserver.yaml



make container
docker tag ghcr.io/projectcontour/contour:$(git rev-parse --short HEAD) localhost/contour:latest
kind load docker-image localhost/contour:latest --name contour


kubectl apply -f examples/contour/01-crds.yaml


cat <<EOF | kubectl -n projectcontour patch deployment contour --patch-file=/dev/stdin
spec:
  template:
    spec:
      containers:
      - name: contour
        image: localhost/contour:latest
        imagePullPolicy: Never
EOF

cat <<EOF | kubectl -n projectcontour patch daemonset envoy --patch-file=/dev/stdin
spec:
  template:
    spec:
      containers:
      - name: shutdown-manager
        image: localhost/contour:latest
        imagePullPolicy: Never
      initContainers:
      - name: envoy-initconfig
        image: localhost/contour:latest
        imagePullPolicy: Never
EOF



#########
#
# Config file
#

http -v http://echoserver.127-0-0-101.nip.io
# check that it returns:
#   content-encoding: gzip

cat <<EOF | kubectl --namespace projectcontour create configmap contour --from-file=contour.yaml=/dev/stdin --dry-run=client -o yaml | kubectl apply -f -
compression:
  algorithm: brotli
EOF


kubectl scale deployment contour --replicas=0 -n projectcontour
kubectl scale deployment contour --replicas=2 -n projectcontour


http -v http://echoserver.127-0-0-101.nip.io
# check that it returns:
#    content-encoding: br



#########
#
# ContourConfiguration CRD
#


kubectl apply -f - <<EOF
apiVersion: projectcontour.io/v1alpha1
kind: ContourConfiguration
metadata:
  name: contour
  namespace: projectcontour
spec:
  envoy:
    listener:
      compression:
        algorithm: disabled
EOF


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


http -v http://echoserver.127-0-0-101.nip.io
# check that it returns:
#  NO content-encoding header




kubectl apply -f - <<EOF
apiVersion: projectcontour.io/v1alpha1
kind: ContourConfiguration
metadata:
  name: contour
  namespace: projectcontour
spec:
  envoy:
    listener: {}
EOF


kubectl scale deployment contour --replicas=0 -n projectcontour
kubectl scale deployment contour --replicas=2 -n projectcontour


http -v http://echoserver.127-0-0-101.nip.io
# check that it returns:
#   content-encoding: gzip



kubectl explain ContourConfiguration.spec.envoy.listener
kubectl explain ContourConfiguration.spec.envoy.listener.compression
