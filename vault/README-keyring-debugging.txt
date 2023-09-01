# issue

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
