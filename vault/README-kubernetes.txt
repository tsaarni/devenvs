

# Create Kind cluster
kind delete cluster --name vault
kind create cluster --config configs/kind-cluster-config.yaml --name vault

rm -rf certs
mkdir -p certs
certyaml -d certs/ configs/certs.yaml



##########
#
# Integrated storage - raft
#

kubectl create secret generic vault-certs --from-file=certs/ca.pem --from-file=certs/vault.pem --from-file=certs/vault-key.pem --dry-run=client -o yaml | kubectl apply -f -


kubectl delete -f manifests/vault-integrated-storage.yaml
kubectl apply -f manifests/vault-integrated-storage.yaml



# Initialize the vault on the first pod
kubectl exec -it vault-0 -- vault operator init -key-shares=1 -key-threshold=1 -format=json > vault-unseal-config.json


# Unseal the vaults
kubectl exec -it vault-0 -- vault operator unseal $(jq -r .unseal_keys_b64[0] vault-unseal-config.json)
kubectl exec -it vault-1 -- vault operator unseal $(jq -r .unseal_keys_b64[0] vault-unseal-config.json)
kubectl exec -it vault-2 -- vault operator unseal $(jq -r .unseal_keys_b64[0] vault-unseal-config.json)



# authenticate with the root token
kubectl exec -it vault-0 -- vault login $(jq -r .root_token vault-unseal-config.json)
kubectl exec -it vault-0 -- vault operator raft list-peers








##############
#
# Etcd
#

# Build version of etcd with shell inside the container
docker build -t localhost/etcd:latest docker/etcd/
kind load docker-image localhost/etcd:latest --name vault

kubectl apply -f manifests/etcd.yaml
