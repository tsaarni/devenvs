
#######################
#
# build & debug
#

make

# build debug version
make GCFLAGS="all=-N -l"

mkdir -p .vscode
cp ~/work/devenvs/vault/configs/launch.json .vscode/


# run dev mode (in-memory database)
bin/vault server -dev



# basic config with file storage
cat > config.hcl <<EOF
storage "file" {
  path = "/tmp/vault-test"
}
listener "tcp" {
    address = "127.0.0.1:8200"
    tls_disable = true
}
disable_mlock = true
EOF

bin/vault server -config=config.hcl

# remove file backend storage
rm -rf /tmp/vault-test/



# initialize vault

http POST http://127.0.0.1:8200/v1/sys/init secret_shares:=1 secret_threshold:=1
http POST http://127.0.0.1:8200/v1/sys/unseal key=PASTE_KEY_HERE
http http://127.0.0.1:8200/v1/sys/seal-status






# with tls
mkdir -p certs
certyaml -d certs/ configs/certs.yaml
bin/vault server -config=$HOME/work/devenvs/vault/configs/vault-config-with-tls.hcl



#######################
#
# vault CLI
#

export VAULT_ADDR=http://127.0.0.1:8200

vault login
vault secrets list   # lists mounts

# remove persisted vault token
rm ~/.vault-token


vault write cubbyhole/mycredentials username="joe" password="password"
vault read cubbyhole/mycredentials




###############################
#
# debug inside pod
#

kubectl cp ~/go/bin/dlv vault-0:/usr/bin/dlv
kubectl exec vault-0 -- chmod +x /usr/bin/dlv
kubectl exec vault-0 -- sh -c "mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2" # dlv compiled on ubuntu requires this on alpine

VAULT_API_ADDR=http://$POD_NAME.vault:8200 dlv --listen=:8181 --headless=true --api-version=2 exec /usr/bin/vault -- server -log-level=debug -config /config/config-ha.hcl
kubectl port-forward 8181
