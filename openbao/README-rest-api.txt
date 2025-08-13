
export BAO_ADDR=http://127.0.0.1:8200

# initialize vault

http POST $BAO_ADDR/v1/sys/init secret_shares:=1 secret_threshold:=1 | tee init.json
export UNSEAL_KEY=$(jq -r .keys[0] init.json)
export BAO_TOKEN=$(jq -r .root_token init.json)


http POST $BAO_ADDR/v1/sys/unseal key=$UNSEAL_KEY


### System backend

http $BAO_ADDR/v1/sys/seal-status                          # seal status
http $BAO_ADDR/v1/sys/mounts X-Vault-Token:$BAO_TOKEN      # read mounts
http $BAO_ADDR/v1/sys/key-status X-Vault-Token:$BAO_TOKEN  # encryption key status



### kv v1 engine

http POST $BAO_ADDR/v1/sys/mounts/secret type=kv X-Vault-Token:$BAO_TOKEN # enable
http POST $BAO_ADDR/v1/secret/foo X-Vault-Token:$BAO_TOKEN mysecret=foo   # write
http $BAO_ADDR/v1/secret/foo X-Vault-Token:$BAO_TOKEN                     # read


### kv v2 engine

http POST $BAO_ADDR/v1/sys/mounts/kv2 type=kv X-Vault-Token:$BAO_TOKEN options:='{"version":"2"}' # enable
http POST $BAO_ADDR/v1/kv2/data/foo X-Vault-Token:$BAO_TOKEN data:='{"mysecret":"foo"}'           # write
http $BAO_ADDR/v1/kv2/data/foo X-Vault-Token:$BAO_TOKEN                                           # read current version
http POST $BAO_ADDR/v1/kv2/data/foo X-Vault-Token:$BAO_TOKEN data:='{"mysecret":"new-value"}'     # update
http $BAO_ADDR/v1/kv2/data/foo version==1 X-Vault-Token:$BAO_TOKEN                                # read specific version


### transit engine

http POST $BAO_ADDR/v1/sys/mounts/transit type=transit X-Vault-Token:$BAO_TOKEN  # enable
http POST $BAO_ADDR/v1/transit/keys/foo X-Vault-Token:$BAO_TOKEN                 # create encryption key
http POST $BAO_ADDR/v1/transit/keys/bar X-Vault-Token:$BAO_TOKEN type=ed25519    # create signing key

# encrypt & decrypt
ciphertext=$(http POST $BAO_ADDR/v1/transit/encrypt/foo X-Vault-Token:$BAO_TOKEN plaintext=$(base64 <<< "mysecret") | jq -r .data.ciphertext)
http POST $BAO_ADDR/v1/transit/decrypt/foo X-Vault-Token:$BAO_TOKEN ciphertext=$ciphertext | jq -r .data.plaintext | base64 -d

# sign & verify
signature=$(http POST $BAO_ADDR/v1/transit/sign/bar X-Vault-Token:$BAO_TOKEN input=$(base64 <<< "my data") | jq -r .data.signature)
http POST $BAO_ADDR/v1/transit/verify/bar X-Vault-Token:$BAO_TOKEN input=$(base64 <<< "my data") signature=$signature



### cubbyhole

http POST $BAO_ADDR/v1/cubbyhole/mysecret X-Vault-Token:$BAO_TOKEN mysecret=foo # write
http $BAO_ADDR/v1/cubbyhole/mysecret X-Vault-Token:$BAO_TOKEN                   # read


### Re-key (rotate unseal keys)
http POST $BAO_ADDR/v1/sys/rekey/init X-Vault-Token:$BAO_TOKEN secret_shares:=1 secret_threshold:=1 | tee rekey1.json                # initialize rekey
http POST $BAO_ADDR/v1/sys/rekey/update X-Vault-Token:$BAO_TOKEN key=$UNSEAL_KEY nonce=$(jq -r .nonce rekey1.json) | tee rekey2.json # update rekey with unseal key
export UNSEAL_KEY=$(jq -r .keys[0] rekey2.json)
http $BAO_ADDR/v1/sys/rekey/verify X-Vault-Token:$BAO_TOKEN   # verify rekey status


### Rotate encryption key (the data encryption key)
http POST $BAO_ADDR/v1/sys/rotate X-Vault-Token:$BAO_TOKEN    # rotate encryption key
http $BAO_ADDR/v1/sys/key-status X-Vault-Token:$BAO_TOKEN     # encryption key status



### metrics

http $BAO_ADDR/v1/sys/metrics X-Vault-Token:$BAO_TOKEN
