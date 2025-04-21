
listener "tcp" {
    address       = "0.0.0.0:8200"
    tls_cert_file = "/input/certs/bao-server.pem"
    tls_key_file  = "/input/certs/bao-server-key.pem"
}

storage "raft" {
    path = "/data"

    retry_join {
        leader_api_addr         = "https://openbao-node-1:8200"
        leader_ca_cert_file     = "/input/certs/ca.pem"
        leader_client_cert_file = "/input/certs/bao-client.pem"
        leader_client_key_file  = "/input/certs/bao-client-key.pem"
    }

    retry_join {
        leader_api_addr         = "https://openbao-node-2:8200"
        leader_ca_cert_file     = "/input/certs/ca.pem"
        leader_client_cert_file = "/input/certs/bao-client.pem"
        leader_client_key_file  = "/input/certs/bao-client-key.pem"
    }

    retry_join {
        leader_api_addr         = "https://openbao-node-3:8200"
        leader_ca_cert_file     = "/input/certs/ca.pem"
        leader_client_cert_file = "/input/certs/bao-client.pem"
        leader_client_key_file  = "/input/certs/bao-client-key.pem"
    }

}
