

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











Vault kubernetes auth method


make test
make testacc


# test compile
XC_OSARCH=linux/amd64 make
