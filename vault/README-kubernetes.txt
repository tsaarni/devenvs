

# Create Kind cluster
kind delete cluster --name vault
kind create cluster --config configs/kind-cluster-config.yaml --name vault


# Build version of etcd with shell inside the container
docker build -t localhost/etcd:latest docker/etcd/
kind load docker-image localhost/etcd:latest --name vault




kubectl apply -f manifests/etcd.yaml
