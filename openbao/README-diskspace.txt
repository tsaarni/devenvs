
# Run container with limited disk space.
docker run --tmpfs /data:size=5m --volume ./:/input:ro --publish 8200:8200 --rm -it ghcr.io/openbao/openbao:2.2.0 ash

# Start OpenBAO server in the container.
bao server -config /input/configs/openbao-single-node-config.hcl -log-level=debug


### In other terminal

### Use case 1: Disk space exhaustion with service tokens.

# Initialize and unseal the server.
export BAO_ADDR=http://127.0.0.1:8200
http POST $BAO_ADDR/v1/sys/init secret_shares:=1 secret_threshold:=1 > init.json
export UNSEAL_KEY=$(jq -r .keys[0] init.json)
export ROOT_TOKEN=$(jq -r .root_token init.json)
http POST $BAO_ADDR/v1/sys/unseal key=$UNSEAL_KEY


# Enable userpass auth method (with service token).
http POST $BAO_ADDR/v1/sys/auth/userpass X-Vault-Token:$ROOT_TOKEN type=userpass

# Create user "joe" with password "joe",
http POST $BAO_ADDR/v1/auth/userpass/users/joe password=joe X-Vault-Token:$ROOT_TOKEN

# Login repeatedly to create masses of tokens until disk space is exhausted.
siege --concurrent=100 "$BAO_ADDR/v1/auth/userpass/login/joe POST password=joe"


### Workaround: Use batch token type to avoid disk space issues.

# Using batch token type avoids disk space issues.
http POST $BAO_ADDR/v1/sys/auth/userpass X-Vault-Token:$ROOT_TOKEN type=userpass config:='{"token_type": "batch"}'




### Use case 2: Disk space exhaustion with secrets.


# Initialize and unseal the server.
export BAO_ADDR=http://127.0.0.1:8200
http POST $BAO_ADDR/v1/sys/init secret_shares:=1 secret_threshold:=1 > init.json
export UNSEAL_KEY=$(jq -r .keys[0] init.json)
export ROOT_TOKEN=$(jq -r .root_token init.json)
http POST $BAO_ADDR/v1/sys/unseal key=$UNSEAL_KEY

# Enable k/v secrets engine.
http POST $BAO_ADDR/v1/sys/mounts/secret X-Vault-Token:$ROOT_TOKEN type=kv options:='{"version": "2"}'

# Enable userpass auth method (with service token).
http POST $BAO_ADDR/v1/sys/auth/userpass X-Vault-Token:$ROOT_TOKEN type=userpass

# Create policy "secret-writer" that allows creating secrets.
http POST "$BAO_ADDR/v1/sys/policy/secret-writer" X-Vault-Token:$ROOT_TOKEN policy="path \"secret/*\" { capabilities = [\"create\", \"read\", \"update\", \"delete\", \"list\"] }"

# Create user "joe" with password "joe",
http POST $BAO_ADDR/v1/auth/userpass/users/joe X-Vault-Token:$ROOT_TOKEN password=joe token_policies=secret-writer



# Login as "joe" to get a client token.
JOE_TOKEN=$(http POST $BAO_ADDR/v1/auth/userpass/login/joe password=joe | jq -r .auth.client_token)


# Create masses fo secrets until disk space is exhausted.
siege --concurrent=100 --content-type "application/json" --header "X-Vault-Token:$JOE_TOKEN" "$BAO_ADDR/v1/secret/data/mysecret POST {\"data\": {\"key\": \"value\"}}"



# Delete secrets to free up disk space.
# Delete all secrets in the "secret" path.
http DELETE "$BAO_ADDR/v1/secret/metadata/mysecret" X-Vault-Token:$JOE_TOKEN

# delete with root token
http DELETE "$BAO_ADDR/v1/secret/metadata/mysecret" X-Vault-Token:$ROOT_TOKEN

# Print raft status
http GET $BAO_ADDR/v1/sys/storage/raft/status X-Vault-Token:$ROOT_TOKEN

wireshark -i lo -k -f "port 8200" -Y http


# Check the seal status.
http GET $BAO_ADDR/v1/sys/seal-status

# Check the status
http GET $BAO_ADDR/v1/sys/health

# check mounts
http GET $BAO_ADDR/v1/sys/mounts


# Create a secret
http POST "$BAO_ADDR/v1/secret/data/mysecret" X-Vault-Token:$JOE_TOKEN data:='{"key": "value"}'








##############################################
#
# Clustered raft with limited disk space.
#

rm -rf certs
mkdir certs
certyaml -d certs configs/certs.yaml


docker compose -f docker-compose-limited-diskspace.yaml rm -f
docker compose -f docker-compose-limited-diskspace.yaml up




# Initialize.
http --verify=certs/ca.pem POST https://openbao.127-0-58-11.nip.io:8200/v1/sys/init secret_shares:=1 secret_threshold:=1 > init.json
export UNSEAL_KEY=$(jq -r .keys[0] init.json)
export ROOT_TOKEN=$(jq -r .root_token init.json)


http --verify=certs/ca.pem POST https://openbao.127-0-58-11.nip.io:8200/v1/sys/unseal key=$UNSEAL_KEY
http --verify=certs/ca.pem POST https://openbao.127-0-58-12.nip.io:8200/v1/sys/unseal key=$UNSEAL_KEY
http --verify=certs/ca.pem POST https://openbao.127-0-58-13.nip.io:8200/v1/sys/unseal key=$UNSEAL_KEY


# Enable userpass auth method (with service token).
http --verify=certs/ca.pem POST https://openbao.127-0-58-11.nip.io:8200/v1/sys/auth/userpass X-Vault-Token:$ROOT_TOKEN type=userpass

# Create user "joe" with password "joe",
http --verify=certs/ca.pem POST https://openbao.127-0-58-11.nip.io:8200/v1/auth/userpass/users/joe password=joe X-Vault-Token:$ROOT_TOKEN

# Login repeatedly to create masses of tokens until disk space is exhausted.
siege --concurrent=100 "https://openbao.127-0-58-11.nip.io:8200/v1/auth/userpass/login/joe POST password=joe"




# Fetch raft configuration.
http --verify=certs/ca.pem GET https://openbao.127-0-58-11.nip.io:8200/v1/sys/storage/raft/configuration X-Vault-Token:$ROOT_TOKEN
