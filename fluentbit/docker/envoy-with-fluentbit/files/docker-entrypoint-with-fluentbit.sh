#!/usr/bin/env bash
set -e

# Start fluentbit in a subshell and forward all stdout and stderr to fluentbit process.
# If fluentbit exits e.g. due to crash, execute kill to exit the parent process (which will be catatonit).
# The reason for this is to restart the container since otherwise we would lose the logs from envoy process.
exec > >(/opt/fluent-bit/bin/fluent-bit --quiet --config=/configs/fluentbit-envoy.conf ; kill -SIGTERM $$ ) 2>&1

# Start envoy with catatonit as the init process.
# Exec will replace the current (bash) process with the catatonit process.
exec /usr/local/bin/envoy --config-path /etc/envoy/envoy-httpbingo-config.yaml
