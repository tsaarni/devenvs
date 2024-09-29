# https://github.com/projectcontour/contour/issues/6291
# https://github.com/projectcontour/contour/pull/6295


kind delete cluster --name contour
kind create cluster --config configs/kind-cluster-config.yaml --name contour


kubectl apply -f https://projectcontour.io/quickstart/contour.yaml

# create empty namespace
kubectl create ns empty

# confirm that both contour replicas are ready: should be READY 2/2
kubectl -n projectcontour get deployment contour

# edit deployment to add watch namespace and debug flags
kubectl patch deployment contour -n projectcontour \
  --type='json' \
  -p='[{"op": "add", "path": "/spec/template/spec/containers/0/args/-", "value": "--watch-namespaces=empty"},
       {"op": "add", "path": "/spec/template/spec/containers/0/args/-", "value": "--debug"}]'


# confirm that both contour replicas are ready: should be READY 2/2
kubectl -n projectcontour get deployment contour

# confirm that contour is running
kubectl -n projectcontour get pods



# create dummy configmaps in empty namespace
for i in {1..1000}; do cat <<EOF | kubectl apply -f -
apiVersion: v1
kind: ConfigMap
metadata:
  name: dummy-$i
  namespace: empty
data:
  dummy: "dummy"
EOF
done



kubectl -n projectcontour scale deployment --replicas=0 contour
kubectl -n projectcontour scale deployment --replicas=2 contour


# confirm that both configmaps are created in empty namespace
kubectl get configmaps -n empty



#### Problem: only the leader contour becomes ready



# create a dummy secret in empty namespace
kubectl -n empty create secret generic dummy-secret --from-literal=dummy=dummy

# force-restart contour
kubectl -n projectcontour delete pod -l app=contour

# check that only the leader contour is ready
kubectl -n projectcontour get deployment contour








kubectl -n projectcontour scale deployment --replicas=1 contour

kubectl -n projectcontour get secret contourcert -o jsonpath='{..ca\.crt}' | base64 -d > ca.crt
kubectl -n projectcontour get secret contourcert -o jsonpath='{..tls\.crt}' | base64 -d > tls.crt
kubectl -n projectcontour get secret contourcert -o jsonpath='{..tls\.key}' | base64 -d > tls.key


go run github.com/projectcontour/contour/cmd/contour serve --xds-address=0.0.0.0 --xds-port=8001 --contour-cafile=ca.crt --contour-cert-file=tls.crt --contour-key-file=tls.key --watch-namespaces=empty --debug
