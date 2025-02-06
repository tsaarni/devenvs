
cluster_name = "vault"
cluster_addr = "http://127.0.0.1:8201"
disable_mlock = true

listener "tcp" {
  address = "127.0.0.1:8200"
  tls_disable = true
}

api_addr = "http://localhost:8200"

storage "raft" {
  path = "/home/tsaarni/work/devenvs/vault/blocking/"
  node_id = "node1"
}
