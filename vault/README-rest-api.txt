

export VAULT_ADDR=http://127.0.0.1:8200

# initialize vault

http POST $VAULT_ADDR/v1/sys/init secret_shares:=1 secret_threshold:=1 > init.json
export UNSEAL_KEY=$(jq -r .keys[0] init.json)
export ROOT_TOKEN=$(jq -r .root_token init.json)


http POST $VAULT_ADDR/v1/sys/unseal key=$UNSEAL_KEY


### System backend

http $VAULT_ADDR/v1/sys/seal-status                          # seal status
http $VAULT_ADDR/v1/sys/mounts X-Vault-Token:$ROOT_TOKEN     # read mounts
http $VAULT_ADDR/v1/sys/key-status X-Vault-Token:$ROOT_TOKEN # encryption key status


### kv engine

http POST $VAULT_ADDR/v1/sys/mounts/secret type=kv X-Vault-Token:$ROOT_TOKEN # enable
http POST $VAULT_ADDR/v1/secret/foo X-Vault-Token:$ROOT_TOKEN mysecret=foo   # write
http $VAULT_ADDR/v1/secret/foo X-Vault-Token:$ROOT_TOKEN                     # read


### transit engine

http POST $VAULT_ADDR/v1/sys/mounts/transit type=transit X-Vault-Token:$ROOT_TOKEN  # enable
http POST $VAULT_ADDR/v1/transit/keys/foo X-Vault-Token:$ROOT_TOKEN                 # create encryption key
http POST $VAULT_ADDR/v1/transit/keys/bar X-Vault-Token:$ROOT_TOKEN type=ed25519    # create signing key

# encrypt & decrypt
ciphertext=$(http POST $VAULT_ADDR/v1/transit/encrypt/foo X-Vault-Token:$ROOT_TOKEN plaintext=$(base64 <<< "mysecret") | jq -r .data.ciphertext)
http POST $VAULT_ADDR/v1/transit/decrypt/foo X-Vault-Token:$ROOT_TOKEN ciphertext=$ciphertext | jq -r .data.plaintext | base64 -d

# sign & verify
signature=$(http POST $VAULT_ADDR/v1/transit/sign/bar X-Vault-Token:$ROOT_TOKEN input=$(base64 <<< "my data") | jq -r .data.signature)
http POST $VAULT_ADDR/v1/transit/verify/bar X-Vault-Token:$ROOT_TOKEN input=$(base64 <<< "my data") signature=$signature



### cubbyhole

http POST $VAULT_ADDR/v1/cubbyhole/mysecret X-Vault-Token:$ROOT_TOKEN mysecret=foo # write
http $VAULT_ADDR/v1/cubbyhole/mysecret X-Vault-Token:$ROOT_TOKEN                   # read
