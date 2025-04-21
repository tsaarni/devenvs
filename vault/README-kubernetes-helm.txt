



# Create Kind cluster
kind delete cluster --name vault
kind create cluster --config configs/kind-cluster-config.yaml --name vault

rm -rf certs
mkdir -p certs
certyaml -d certs/ configs/certs.yaml


##########

helm repo list
helm repo add hashicorp https://helm.releases.hashicorp.com
helm repo update hashicorp

helm search repo hashicorp/vault

# Show all configuration options
helm show values hashicorp/vault



####
# Create namespace
kubectl delete namespace vault
kubectl create namespace vault

# Create secret for certs
kubectl create secret generic vault-certs --from-file=certs/ca.pem --from-file=certs/vault.pem --from-file=certs/vault-key.pem --namespace vault --dry-run=client -o yaml | kubectl apply -f -

helm install vault hashicorp/vault --namespace vault --values configs/helm-override-values.yml

kubectl -n vault logs -f vault-0



# Initialize the vault on the first pod
kubectl -n vault exec -it vault-0 -- vault operator init -key-shares=1 -key-threshold=1 -format=json > vault-unseal-config.json




# Unseal the vaults
kubectl -n vault exec -it vault-0 -- vault operator unseal $(jq -r .unseal_keys_b64[0] vault-unseal-config.json)
kubectl -n vault exec -it vault-1 -- vault operator unseal $(jq -r .unseal_keys_b64[0] vault-unseal-config.json)
kubectl -n vault exec -it vault-2 -- vault operator unseal $(jq -r .unseal_keys_b64[0] vault-unseal-config.json)

# authenticate with the root token
kubectl -n vault exec -it vault-0 -- vault login $(jq -r .root_token vault-unseal-config.json)
kubectl -n vault exec -it vault-0 -- vault operator raft list-peers
