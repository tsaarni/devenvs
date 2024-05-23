# https://github.com/projectcontour/contour/issues/6291
# https://github.com/projectcontour/contour/pull/6295


kind delete cluster --name contour
kind create cluster --config configs/kind-cluster-config.yaml --name contour


kubectl apply -f https://projectcontour.io/quickstart/contour.yaml

# create empty namespace
kubectl create ns empty

# confirm no secrets in empty namespace
kubectl get secrets empty

# confirm that both contour replicas are ready: should be READY 2/2
kubectl -n projectcontour get deployment contour

# edit deployment to add watch namespace and debug flags
kubectl -n projectcontour edit deployment contour


# add
        - --watch-namespaces=empty
        - --debug


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
