
# Create Kind cluster
kind delete cluster --name vault
kind create cluster --config configs/kind-cluster-config.yaml --name vault

# generate certs
rm -rf certs
mkdir -p certs
certyaml -d certs configs/certs.yaml


kubectl create secret generic vault-certs --dry-run=client -o yaml --from-file=certs/ca.pem --from-file=certs/client.pem --from-file=certs/client-key.pem | kubectl apply -f -
kubectl create secret generic etcd-certs --dry-run=client -o yaml --from-file=certs/ca.pem --from-file=certs/etcd.pem --from-file=certs/etcd-key.pem | kubectl apply -f -



# Deploy Etcd cluster and empty placeholder pods for Vault
kubectl apply -f manifests/etcd.yaml
kubectl apply -f manifests/vault.yaml

# Copy Vault binary to the pods (run in vault source dir)
tar cf - bin/vault | kubectl exec -i vault-0 -- tar xf - -C /usr/
tar cf - bin/vault | kubectl exec -i vault-1 -- tar xf - -C /usr/

# Exec into Vault containers
kubectl exec -it vault-0 -- ash
kubectl exec -it vault-1 -- ash



# Run Vault in both containers
VAULT_API_ADDR=http://$POD_NAME.vault:8200 vault server -log-level=debug -config /config/config-etcd-ha.hcl
VAULT_API_ADDR=http://$POD_NAME.vault:8200 vault server -log-level=debug -config /config/config-etcd.hcl



kubectl apply -f manifests/shell.yaml
kubectl exec -it shell -- ash

# Initialize one of the Vaults if starting from scratch without persistent state
http POST http://vault-0.vault:8200/v1/sys/init secret_shares:=1 secret_threshold:=1 > init.json

export UNSEAL_KEY=$(jq -r .keys[0] init.json)
export ROOT_TOKEN=$(jq -r .root_token init.json)

# Unseal Vault instances, replace key with the deployment specific key
http -v POST http://vault-0.vault:8200/v1/sys/unseal key=$UNSEAL_KEY
http -v POST http://vault-1.vault:8200/v1/sys/unseal key=$UNSEAL_KEY

# Enable k/v secrets engine if starting from scratch without persistent state
http -v POST http://vault:8200/v1/sys/mounts/secret X-Vault-Token:$ROOT_TOKEN type=kv-v2
http -v POST http://vault:8200/v1/secret/config X-Vault-Token:$ROOT_TOKEN



# Monitor the status of Vault instances
#
# responses:
#   - master: "200 OK"
#   - slave: "429 Too Many Requests"
while true; do echo; date +%FT%T; for i in vault-0 vault-1; do printf "   $i: $(http -h http://$i.vault:8200/v1/sys/health | grep HTTP)\n"; done; sleep 1; done

# Monitor the status of all pods
watch -n .5 kubectl get pod

# Print status of etcd instances
echo -n etcd-0 etcd-1 etcd-2 | xargs -d' '  -n1 -I% -P0 sh -c 'printf "%: $(kubectl exec % -- etcdctl endpoint status)\n"' | sort

# Kill leader in etcd cluster
kubectl exec etcd-0 -- etcdctl endpoint status --cluster=true | grep true | printf "kubectl delete --force pod $(cut -c8-13)" | sh

# Periodically write data to k/v
count=0 ; while true; do printf "$(date +%FT%T) writing: $count: $(http -h POST http://vault:8200/v1/secret/data/mysecret X-Vault-Token:$ROOT_TOKEN data:={\"key\":\"$count\"} | grep HTTP)\n"; count=$((count+1)); done

# write & read data from k/v
http -h POST http://vault:8200/v1/secret/data/mysecret X-Vault-Token:$ROOT_TOKEN data:={\"key\":\"123\"}
http -v GET http://vault:8200/v1/secret/data/mysecret X-Vault-Token:$ROOT_TOKEN



# Restart all etcd instances
kubectl rollout restart statefulset/etcd

# Abruptly delete
kubectl delete pod -l app=etcd --force

# To abruptly shutdown etcd cluster
kubectl delete statefulset etcd --force
kubectl delete pod -l app=etcd --force




# To remove etcd deployment and PVCs (WARNING: this will delete all data)
kubectl delete -f manifests/etcd.yaml --force
echo etcd-etcd-0 etcd-etcd-1 etcd-etcd-2 | xargs -n1 kubectl delete pvc --force


######################################################
#
# Random tips and tricks
#

# List keys from etcd
ETCDCTL_API=3 etcdctl get --keys-only --prefix ""

# Capture traffic
#   in wireshark, select packet between Vault and Etcd, decode as HTTP2
#   filter: protobuf
sudo nsenter -t $(pidof -s vault) -n wireshark
