storage "etcd" {
  etcd_api = "v3"
  address = "http://etcd0:2379,http://etcd1:2379,http://etcd2:2379"
  prefix = "vault/"
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
