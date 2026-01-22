
# KVv2 secrets disappeared after updating from 1.20.2 to 1.21.0
# https://github.com/hashicorp/vault/issues/31614




git checkout v1.20.4
make bootstrap
make static-dist dev-ui
cp -a bin/vault vault-v1.20.4

git checkout v1.21.0
make bootstrap
make static-dist dev-ui
cp -a bin/vault vault-v1.21.0



rm -rf /tmp/vault
mkdir -p /tmp/vault/data
./vault-v1.20.4 server -config=$HOME/work/devenvs/vault/configs/vault-config-file.hcl
./vault-v1.20.4 server -config=$HOME/work/devenvs/vault/configs/vault-config-raft.hcl


# In another terminal

# init vault
export VAULT_ADDR='http://localhost:8200'

~/work/vault/bin/vault operator init -key-shares=1 -key-threshold=1 -format=json > /tmp/vault/init.json

export VAULT_UNSEAL_KEY=$(cat /tmp/vault/init.json | jq -r '.unseal_keys_b64[0]')
export VAULT_ROOT_TOKEN=$(cat /tmp/vault/init.json | jq -r '.root_token')

~/work/vault/bin/vault operator unseal $VAULT_UNSEAL_KEY
~/work/vault/bin/vault login $VAULT_ROOT_TOKEN

# enable kv v2 at secret/
~/work/vault/bin/vault secrets enable -path=my:path -version=2 kv

# store secret in kv v2
~/work/vault/bin/vault kv put my:path/foo foo=bar

#  get secret
~/work/vault/bin/vault kv get my:path/foo


# check ui
echo $VAULT_ROOT_TOKEN
http://localhost:8200/ui/

./vault-v1.21.0 server -config=$HOME/work/devenvs/vault/configs/vault-config-file.hcl
./vault-v1.21.0 server -config=$HOME/work/devenvs/vault/configs/vault-config-raft.hcl


~/work/vault/bin/vault operator unseal $VAULT_UNSEAL_KEY
~/work/vault/bin/vault login $VAULT_ROOT_TOKEN

# list all mounts
~/work/vault/bin/vault secrets list
