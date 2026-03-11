#!/bin/bash

set -ex

# Initialize and save credentials to init.json
curl -s -X POST http://127.0.0.1:8200/v1/sys/init \
  -d '{"secret_shares":1,"secret_threshold":1}' | tee init.json

# Unseal
curl -s -X POST http://127.0.0.1:8200/v1/sys/unseal \
  -d "{\"key\":\"$(jq -r '.keys[0]' init.json)\"}"

# Enable KV engine
curl -s -X POST http://127.0.0.1:8200/v1/sys/mounts/secret \
  -H "X-Vault-Token: $(jq -r '.root_token' init.json)" \
  -d '{"type":"kv","options":{"version":"1"}}'
