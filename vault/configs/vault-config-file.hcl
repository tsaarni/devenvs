
cluster_name = "vault"
cluster_addr = "http://127.0.0.1:8201"
disable_mlock = true

listener "tcp" {
  address = "127.0.0.1:8200"
  tls_disable = true
}

ui = true

api_addr = "http://localhost:8200"

storage "file" {
  path = "/tmp/vault/data"
}
