

# Compile two versions of Vault
git checkout v1.8.8
GO_CMD=go1.16.13 XC_OSARCH=linux/amd64 make
mv bin/vault vault-v1.8.8


git checkout 1.9.3-nordix
XC_OSARCH=linux/amd64 make
mv bin/vault vault-v1.9.3-nordix




# Create config for persistent storage
cat > config.hcl <<EOF
listener "tcp" {
  address = "127.0.0.1:8200"
  tls_disable = true
}

storage "file" {
  path = "/tmp/vault-data"
}

api_addr = "http://127.0.0.1:8200"

disable_mlock = true

EOF

rm -rf /tmp/vault-data



# Run old vault version
./vault-v1.8.8 server -config=config.hcl

vault operator init -key-shares=1 -key-threshold=1
export VAULT_TOKEN=<COPY ROOT TOKEN>
vault operator unseal <COPY UNSEAL KEY>
vault secrets enable pki


# Create role and double check that signature_bits is not there with old Vault version
vault write pki/roles/example-dot-com allowed_domains=my-website.com allow_subdomains=true max_ttl=72h
vault read pki/roles/example-dot-com -format=json | grep signature_bits




# Start new version 
./vault-v1.9.3-nordix server -config=config.hcl
vault operator unseal <COPY UNSEAL KEY>

# Read role and check that signature_bits is removed from response by the patch
vault read pki/roles/example-dot-com -format=json | grep signature_bits


# Create new role with new Vault and check that signature_bits (256) is there
# It will print
#    "signature_bits": 256,
vault read pki/roles/example-dot-com -format=json > output.json 
vault write pki/roles/example-dot-com-2 @output.json
vault read pki/roles/example-dot-com-2 -format=json | grep signature_bits








