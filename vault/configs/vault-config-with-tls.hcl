storage "file" {
  path = "/tmp/vault-test"
}
listener "tcp" {
    address = "127.0.0.1:8200"
    tls_cert_file = "/home/tsaarni/work/devenvs/vault/certs/vault.pem"
    tls_key_file  = "/home/tsaarni/work/devenvs/vault/certs/vault-key.pem"

}
disable_mlock = true
