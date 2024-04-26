storage "etcd" {
  etcd_api = "v3"
  address = "https://etcd0:2379,https://etcd1:2379,https://etcd2:2379"
  prefix = "vault/"
  tls_ca_file = "/home/vscode/work/devenvs/vault/certs/ca.pem"
  tls_cert_file = "/home/vscode/work/devenvs/vault/certs/vault.pem"
  tls_key_file = "/home/vscode/work/devenvs/vault/certs/vault-key.pem"
}

listener "tcp" {
  address = "0.0.0.0:8200"
  tls_disable = true
}

api_addr = "http://localhost:8200"

default_lease_ttl = "1m"
disable_mlock = true

telemetry {
  prometheus_retention_time = "12h"
}
