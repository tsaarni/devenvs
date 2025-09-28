docker run --rm -it -p 8200:8200 hashicorp/vault:1.20.3

export ROOT_TOKEN=<token>
export VAULT_ADDR=http://localhost:8200

# Create policy
http POST ${VAULT_ADDR}/v1/sys/policy/my-policy X-Vault-Token:${ROOT_TOKEN} policy="path \"secret/*\" { capabilities = [\"create\", \"read\", \"update\", \"delete\", \"list\"] }"









kind delete cluster --name vault
kind create cluster --config configs/kind-cluster-config.yaml --name vault

kubectl apply -f manifests/vault-dev.yaml

# Copy root token.
kubectl exec $(kubectl get pod -l app=vault -o jsonpath='{.items[0].metadata.name}') -c vault-configurator -- cat /unseal/init.json

export ROOT_TOKEN=<token>
export VAULT_ADDR=http://127.0.0.195

# Create policy
http POST ${VAULT_ADDR}/v1/sys/policy/my-policy X-Vault-Token:${ROOT_TOKEN} policy="path \"secret/*\" { capabilities = [\"create\", \"read\", \"update\", \"delete\", \"list\"] }"

# Create role under kubernetes auth
http POST ${VAULT_ADDR}/v1/auth/kubernetes/role/my-role X-Vault-Token:${ROOT_TOKEN} \
    bound_service_account_names=tokenreview \
    bound_service_account_namespaces=default \
    policies=my-policy \
    ttl=1h







HTTP/1.1 200 OK
Cache-Control: no-store
Content-Length: 294
Content-Type: application/json
Date: Tue, 02 Sep 2025 12:11:42 GMT
Strict-Transport-Security: max-age=31536000; includeSubDomains

{
    "auth": null,
    "data": null,
    "lease_duration": 0,
    "lease_id": "",
    "mount_type": "kubernetes",
    "renewable": false,
    "request_id": "1bd1063e-d59a-f8ba-cefe-4a0075dd0893",
    "warnings": [
        "Role my-role does not have an audience. In Vault v1.21+, specifying an audience on roles will be required."
    ],
    "wrap_info": null
}

-VS-

HTTP/1.1 204 No Content
Cache-Control: no-store
Content-Type: application/json
Date: Tue, 02 Sep 2025 12:14:41 GMT
Strict-Transport-Security: max-age=31536000; includeSubDomains
