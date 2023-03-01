

make

# dev
bin/vault server -dev



# basic config

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

http POST http://127.0.0.1:8200/v1/sys/init secret_shares:=1 secret_threshold:=1
http POST http://127.0.0.1:8200/v1/sys/unseal key=PASTE_KEY_HERE
http http://127.0.0.1:8200/v1/sys/seal-status

http -v POST http://127.0.0.1:8200/v1/sys/seal X-Vault-Token:PASTE_ROOT_TOKEN_HERE

rm -rf /tmp/vault-test/



# with tls
mkdir -p certs
certyaml -d certs/ configs/certs.yaml
bin/vault server -config=$HOME/work/devenvs/vault/configs/vault-config-with-tls.hcl




# build debug version
make GCFLAGS="all=-N -l"

# debug inside pod
kubectl cp ~/go/bin/dlv vault-0:/usr/bin/dlv
kubectl exec vault-0 -- chmod +x /usr/bin/dlv
kubectl exec vault-0 -- sh -c "mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2" # dlv compiled on ubuntu requires this on alpine

VAULT_API_ADDR=http://$POD_NAME.vault:8200 dlv --listen=:8181 --headless=true --api-version=2 exec /usr/bin/vault -- server -log-level=debug -config /config/config-ha.hcl
kubectl port-forward 8181

cat >.vscode/launch.json <<EOF
{
    "version": "0.2.0",
    "configurations": [
        {
            "name": "Attach remote",
            "type": "go",
            "request": "attach",
            "mode": "remote",
            "port": 8181,
            "host": "127.0.0.1",
        },
    ]
}
EOF
