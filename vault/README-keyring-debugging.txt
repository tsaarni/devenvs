
https://github.com/hashicorp/vault/issues/21521


vault login
vault operator key-status




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



# initialize vault

rm -rf /tmp/vault-test


export VAULT_ADDR=http://127.0.0.1:8200

http POST $VAULT_ADDR/v1/sys/init secret_shares:=1 secret_threshold:=1 > init.json
export UNSEAL_KEY=$(jq -r .keys[0] init.json)
export ROOT_TOKEN=$(jq -r .root_token init.json)


http POST http://127.0.0.1:8200/v1/sys/unseal key=$UNSEAL_KEY



echo '{"Value":"AAAAAQJCRn8M5QevIZCvxeC3g8LUxEXxMdgLx4F2f7ZN0BC176YsuFWbtc4wEawhTnPxPX9XqkYK0qoGUsyecr4="}' > /tmp/vault-test/core/_shamir-kek

echo '{"Value":"CkxKpdrrWn20B0DLUqFGPoaNPlYMQbEf1jwn4vm3lIqo5/Ye5xcU1b4swjlEIlra7Z1DFIb/Bd7ST8CTqmEi0dNGgshB96T/p0ZDdaDMKgA="}' > /tmp/vault-test/core/hsm/_barrier-unseal-keys


# broken keyring
echo '{"Value":"AAAAAQL30OLWEhRr7Obfl2BJRzMK0GdKNZqQ0Dy9l9G/4jmXkd/Ob6ogt9nahUfyAFrsk3un0wYZu/eO0agTmj8nQMHBTYJEoMqby3OSegnNf/bXFdOQaW+sJp4wQKbrvQ+AZ5d3C2fS4ssn5MJi4KYOvS/Es0nBhcTo7hrAy21xajkhYI/zNj5xkRGGAq2HkbqYccuj+MDKvnKEHW9AtWZKj7S5ZappgME8wk1y6eyU1JoVb1fjAZhyDSFeB0Fr7zMPNBKdik4wmUxk8klGp08P3ws0ovmDt9uqgVfhLU6uS1795igD34fklsFUsc9NwCfdHRjNTqWpWvmlVbqmvBkD5zuz9tsMfySdPsZBcpcQmiNHX5H2A7ZpjKQx+vVGiYw+RozaMQPNDNPa/YRwXal00gODriQA0wnXKaRy6KPQHsEC6pk="}' > /tmp/vault-test/core/_keyring


# bad unseal response with broken keyring
export UNSEAL_KEY=a7d165282c88f6e56713a43dcced935fd12ba3b34078900a304da6858c71a997
http POST http://127.0.0.1:8200/v1/sys/unseal key=$UNSEAL_KEY

HTTP/1.1 400 Bad Request
Cache-Control: no-store
Content-Length: 42
Content-Type: application/json
Date: Fri, 30 Jun 2023 06:54:11 GMT
Strict-Transport-Security: max-age=31536000; includeSubDomains

{
    "errors": [
        "Unseal failed, invalid key"
    ]
}



# valid keyring
echo '{"Value":"AAAAAQKQP5FoB6CDgUN6SDe1LE2mFhCfSvT7Z7FEdv9iRQ2V2PETY4RCqfp+JU7UmGiAJGXX6DYHHtcDG4KlAZ+39IqxOsU+aM04BfG7VEYoorSsdWZpSXfwmAq79F2sS24ByQegcapZdxJBfjjZdbzLuJKtGAnfGfwq3pAyC8IWqKIcAgEbs9Hb8nQLsP53zd7XMs32zXO3Ymw+Vc4WArOrQCzTTcm//XeNDj6Nm4K+2HWEvJ0VbdfW6UUPILVGHcPkG/u5DxP2RGkeAXw9HJI4V1RDlhcbJ7gOAhMZDjMoNMzuu1n09uYSZ6XULRiDif3vopSyWf8+x/uFuLi0JQJoBWI1NeBO+Q4FiKxdt1FpneW+yxke82q1OKak/LnLVIveZYbwhs00mdlYw16GOWerFlPlaScV/YjkXb0YqOniTuJ8EXU="}' > /tmp/vault-test/core/_keyring


### corrupt keyring manually

jq -r .Value < /tmp/vault-test/core/_keyring | base64 -d | xxd -g1 > dump.hex
# insert fault by modifying dump.hex in editor and write it back to file store
echo "{\"Value\":\"$(xxd -r dump.hex | base64 -w0)\"}" > /tmp/vault-test/core/_keyring







# issue


On a couple of occasions, we have seen Vault fail to unseal with the error `Unseal failed, invalid key`.
The issue is permanent and the only way to recover is to restore a working backup.

The unseal keys seem to be correct because the error originates from decrypting the keyring, which becomes after the unseal key was already used to decrypt the storage root key successfully (`core/hsm/barrier-unseal-keys`)

The callstack is below (the topmost is the most reason)

```
go/src/crypto/aes/aes_gcm.go:gcmAsm.Open()    # https://cs.opensource.google/go/go/+/refs/tags/go1.20.5:src/crypto/aes/aes_gcm.go;l=182
vault/barrier_aes_gcm.go:AESGCMBarrier:decrypt()    # https://github.com/hashicorp/vault/blob/325233ea7dba833e987909b21af547d0933751e3/vault/barrier_aes_gcm.go#L1037
vault/barrier_aes_gcm.go:AESGCMBarrier.Unseal()    # https://github.com/hashicorp/vault/blob/325233ea7dba833e987909b21af547d0933751e3/vault/barrier_aes_gcm.go#L453)
```

The error message from `gcmAsm.Open()` is `cipher: message authentication failed`.

The data in the keyring `core/keyring` does not seem to be corrupted, since its length is the same as in the working backup and it begins with a valid header that Vault appends to encrypted values

```yaml
{"Value":"AAAAAQ....
```

We are using etcd as the storage backend.

We have not been able to reproduce the issue on demand, but we have speculated with the following scenarios
- Vault is initialized & unsealed, then the storage is restored from backup while Vault is still running. Vault overwrites the restored keyring with the current one, which is encrypted with a different key.
- Storage root key is somehow corrupted in runtime memory (e.g. due to a bug), and when Vault periodically writes the keyring, it gets encrypted with the corrupted key.


Questions:

Have you seen this kind of issue before?
Can you think of any theories on how this issue could happen?
Do you have any suggestions on how to go about troubleshooting this kind of issue and avoid it in the future?
