cluster_addr = "http://127.0.0.1:8201"
api_addr     = "http://localhost:8200"

listener "tcp" {
    address     = "0.0.0.0:8200"
    tls_disable = true
}

storage "raft" {
    path = "/data"

    # Documentation for the parameters
    # https://openbao.org/docs/configuration/storage/raft/

    # Effectively disable count-based snapshots: require 999999 new log entries before creating a snapshot.
    snapshot_threshold = "999999"

    # Effectively disable time-based snapshots: wait roughly 11.5 days (999999s) between snapshots.
    snapshot_interval  = "999999s"
}

telemetry {
    # Enable Prometheus metrics at /v1/sys/metrics?format=prometheus.
    prometheus_retention_time = "24h"
    disable_hostname         = true
}
