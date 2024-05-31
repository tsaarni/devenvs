
# Namespaced contour, similar to
# https://docs.nginx.com/nginx-ingress-controller/installation/running-multiple-ingress-controllers/
# -watch-namespace
#  "This can be useful if you want to use different NGINX Ingress Controllers for different applications, both in terms of isolation and/or operation."

kind delete cluster --name contour
kind create cluster --config configs/kind-cluster-config.yaml --name contour


# deploy with customized Roles and RoleBindings
kubectl kustomize examples/namespaced/ | kubectl apply -f -



#### Run locally built container

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


make container
docker tag ghcr.io/projectcontour/contour:$(git rev-parse --short HEAD) localhost/contour:latest
kind load docker-image localhost/contour:latest --name contour

kubectl -n projectcontour scale deployment --replicas=0 contour  # NOTE:  WAIT FOR CONTOUR PODS TO TERMINATE
kubectl -n projectcontour scale deployment --replicas=2 contour


kubectl -n projectcontour get pod
kubectl -n projectcontour logs deployment/contour -f





#### Run contour outside cluster with the debugger

# point contour service towards the host
sed "s/REPLACE_ADDRESS_HERE/$(docker network inspect kind | jq -r '.[0].IPAM.Config[0].Gateway')/" ~/work/devenvs/contour/manifests/contour-endpoints-dev.yaml | kubectl apply -f -

# shutdown contour inside the cluster
kubectl -n projectcontour scale deployment --replicas=0 contour
kubectl -n projectcontour rollout restart daemonset envoy


# Download contour certs
kubectl -n projectcontour get secret contourcert -o jsonpath='{..ca\.crt}' | base64 -d > ca.crt
kubectl -n projectcontour get secret contourcert -o jsonpath='{..tls\.crt}' | base64 -d > tls.crt
kubectl -n projectcontour get secret contourcert -o jsonpath='{..tls\.key}' | base64 -d > tls.key


# Create secret with contour service account token
cat <<EOF | kubectl  apply -f -
apiVersion: v1
kind: Secret
metadata:
  name: contour-serviceaccount-secret
  namespace: projectcontour
  annotations:
    kubernetes.io/service-account.name: contour
type: kubernetes.io/service-account-token
EOF

# Create kubeconfig file for service account
kubectl config view --raw=true > contour.kubeconfig
sed -i -e '/client-certificate-data/d' -e '/client-key-data/d' contour.kubeconfig
echo "    token: $(kubectl -n projectcontour get secret contour-serviceaccount-secret -o jsonpath='{..token}' | base64 -d )" >> contour.kubeconfig

# test
KUBECONFIG=contour.kubeconfig kubectl -n projectcontour get secret  # succeeds
KUBECONFIG=contour.kubeconfig kubectl -n kube-system get secret     # fails


# add following to to contour serve
#    --kubeconfig=contour.kubeconfig
#

kubectl create namespace empty

go run github.com/projectcontour/contour/cmd/contour serve --xds-address=0.0.0.0 --xds-port=8001 --envoy-service-http-port=8080 --envoy-service-https-port=8443 --contour-cafile=ca.crt --contour-cert-file=tls.crt --contour-key-file=tls.key --watch-namespaces=empty



#### Test with workload

mkdir -p certs
certyaml --destination certs configs/certs.yaml

# should not work
kubectl create secret tls ingress --cert=certs/ingress.pem --key=certs/ingress-key.pem --dry-run=client -o yaml | kubectl apply -f -
kubectl apply -f manifests/echoserver-tls.yaml
kubectl get httpproxy

# should work
kubectl -n projectcontour create secret tls ingress --cert=certs/ingress.pem --key=certs/ingress-key.pem --dry-run=client -o yaml | kubectl apply -f -
kubectl -n projectcontour apply -f manifests/echoserver-tls.yaml
kubectl -n projectcontour get httpproxy

http --verify=certs/external-root-ca.pem https://protected.127-0-0-101.nip.io








