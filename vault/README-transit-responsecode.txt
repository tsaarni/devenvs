
Transit UX improvements: show key policy, configs on write 
https://github.com/hashicorp/vault/pull/20652



make

bin/vault server -dev

export VAULT_ADDR=http://127.0.0.1:8200
export ROOT_TOKEN=              # copy from vault output



http POST $VAULT_ADDR/v1/sys/mounts/transit type=transit X-Vault-Token:$ROOT_TOKEN  # enable
http POST $VAULT_ADDR/v1/transit/keys/foo X-Vault-Token:$ROOT_TOKEN                 # create encryption key


# old version responds
#    HTTP/1.1 204 No Content

# after 1.14.0 it responds
#    HTTP/1.1 200 OK
# with document in response body

