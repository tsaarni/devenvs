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

    # Trigger snapshot after every 100 new log entries.
    snapshot_threshold = "100"

    # Also trigger snapshot every 5 seconds even if threshold is not reached.
    # Default is 120s. Set low to ensure time-based compaction will start too.
    snapshot_interval  = "5s"

    # Keep only 10 log entries after snapshot, for maximum log truncation.
    # Default is 10000.
    trailing_logs      = "10"
}

telemetry {
    # Enable Prometheus metrics at /v1/sys/metrics?format=prometheus.
    prometheus_retention_time = "24h"
    disable_hostname         = true
}
