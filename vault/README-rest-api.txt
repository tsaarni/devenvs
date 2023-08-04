

export VAULT_ADDR=http://127.0.0.1:8200

# initialize vault

http POST $VAULT_ADDR/v1/sys/init secret_shares:=1 secret_threshold:=1

export UNSEAL_KEY=NNNNNNNNNNNNN
export ROOT_TOKEN=NNNNNNNNNNNNN

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

http POST $VAULT_ADDR/v1/sys/mounts/transit type=transit X-Vault-Token:$ROOT_TOKEN # enable
http POST $VAULT_ADDR/v1/transit/keys/foo X-Vault-Token:$ROOT_TOKEN                # create key
http POST $VAULT_ADDR/v1/transit/encrypt/foo X-Vault-Token:$ROOT_TOKEN plaintext=$(base64 <<< "mysecret") # encrypt
http POST $VAULT_ADDR/v1/transit/decrypt/foo X-Vault-Token:$ROOT_TOKEN ciphertext="vault:v1:<BASE64_DATA_FROM_ENCRYPT_RESPONSE>" # decrypt
