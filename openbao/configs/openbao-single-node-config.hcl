
cluster_addr    = "http://127.0.0.1:8201"
api_addr        = "http://localhost:8200"

listener "tcp" {
    address     = "0.0.0.0:8200"
    tls_disable = true
}

storage "raft" {
    path = "/data"
}
