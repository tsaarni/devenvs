#!/bin/sh

exec /usr/local/bin/envoy -c /input/envoy.yaml --service-cluster mycluster --service-node envoy --log-level info --restart-epoch $RESTART_EPOCH
