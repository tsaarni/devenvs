

Manual test procedure for Vault Kubernetes Auth Backend
https://gist.github.com/tsaarni/66b085e0244fe1b0635cc1f806b5c334



# Run unittests
cd ~/work/vault-plugin-auth-kubernete

make test
make testacc


###############################################################################
#
# Test in Kubernetes cluster
#


# Create Kind cluster
kind delete cluster --name vault
kind create cluster --config configs/kind-cluster-config-1h-token-expiry.yaml --name vault

# Deploy empty placeholder pods for Vault
kubectl apply -f manifests/vault.yaml
kubectl apply -f manifests/shell.yaml



# Patch Vault with local version of Kubernetes Auth Method
cd ~/work/vault
echo "replace github.com/hashicorp/vault-plugin-auth-kubernetes => /home/tsaarni/work/vault-plugin-auth-kubernetes" >> go.mod

# Compile Vault binary
XC_OSARCH=linux/amd64 make

# Copy Vault binary to the pod (run in vault source dir)
tar cf - bin/vault | kubectl exec -i vault-0 -- tar xf - -C /usr/

# Exec into Vault container
kubectl exec -it vault-0 -- ash



# Run Vault in dev mode
vault server -dev -dev-listen-address=0.0.0.0:8200 -log-level=debug



# Run with persistent storage

cat > config.hcl <<EOF
listener "tcp" {
  address = "0.0.0.0:8200"
  tls_disable = true
}

storage "file" {
  path = "/tmp/vault-data"
}

disable_mlock = true

api_addr = "http://vault-0.vault:8200"
EOF

vault server -config=config.hcl


rm -rf /tmp/vault-data


kubectl exec -it shell -- ash


http -v POST http://vault-0.vault:8200/v1/sys/init secret_shares:=1 secret_threshold:=1
http -v POST http://vault-0.vault:8200/v1/sys/unseal key=    # from "keys" field in init response


# Set secret in env var (update the values to the specific deployment)
export ROOT_TOKEN=                   # from "root_token" field in init response


# Configure Kubernets Auth method
http POST http://vault-0.vault:8200/v1/sys/auth/kubernetes X-Vault-Token:$ROOT_TOKEN type=kubernetes
http POST http://vault-0.vault:8200/v1/auth/kubernetes/config X-Vault-Token:$ROOT_TOKEN kubernetes_host=https://kubernetes disable_iss_validation=true


# Configure policy and role
http POST http://vault-0.vault:8200/v1/sys/policy/my-policy X-Vault-Token:$ROOT_TOKEN policy="path \"secret/foo\" { capabilities = [\"create\", \"read\", \"delete\"] }"
http POST http://vault-0.vault:8200/v1/auth/kubernetes/role/my-role X-Vault-Token:$ROOT_TOKEN bound_service_account_names="*" bound_service_account_namespaces=default policies=my-policy
http POST http://vault-0.vault:8200/v1/auth/kubernetes/role/my-role-with-audience X-Vault-Token:$ROOT_TOKEN bound_service_account_names="*" bound_service_account_namespaces=default policies=my-policy audience=https://my-own-audience

# Authenticate successfully
http POST http://vault-0.vault:8200/v1/auth/kubernetes/login role=my-role jwt=@/var/run/secrets/kubernetes.io/serviceaccount/token
http POST http://vault-0.vault:8200/v1/auth/kubernetes/login role=my-role jwt=@/projected/token

# Authenticate unsuccessfully (wrong audience)
http POST http://vault-0.vault:8200/v1/auth/kubernetes/login role=my-role-with-audience jwt=@/var/run/secrets/kubernetes.io/serviceaccount/token




# Configure Kubernetes Auth method: authenticate by client's JWT (not local SA token!!!) to token-review API
http POST http://vault-0.vault:8200/v1/auth/kubernetes/config X-Vault-Token:$ROOT_TOKEN kubernetes_host=https://kubernetes disable_iss_validation=true disable_local_ca_jwt=true kubernetes_ca_cert=@/run/secrets/kubernetes.io/serviceaccount/ca.crt

# Authenticate successfully using client's JWT
kubectl exec -it vault-0 -- http POST http://vault-0.vault:8200/v1/auth/kubernetes/login role=my-role jwt=@/var/run/secrets/kubernetes.io/serviceaccount/token

# Authenticate unsuccesfully using client's JWT
kubectl exec -it shell -- http POST http://vault-0.vault:8200/v1/auth/kubernetes/login role=my-role jwt=@/var/run/secrets/kubernetes.io/serviceaccount/token












###############################################################################
#
# Test outside Kubernetes cluster
#




cat > config.hcl <<EOF
listener "tcp" {
  tls_disable = true
}

storage "file" {
  path = "/tmp/vault-data"
}

disable_mlock = true
EOF

vault server -config=config.hcl


http -v POST http://localhost:8200/v1/sys/init secret_shares:=1 secret_threshold:=1
http -v POST http://localhost:8200/v1/sys/unseal key=    # from "keys" field in init response

# Set secret in env var (update the values to the specific deployment)
export ROOT_TOKEN=                   # from "root_token" field in init response

http POST http://localhost:8200/v1/sys/auth/kubernetes X-Vault-Token:$ROOT_TOKEN type=kubernetes

APISERVER_ADDRESS=$(kubectl config view -o json | jq -r '.clusters[] | select(.name|contains("kind-vault")) | .cluster.server')
kubectl cp vault-0:/var/run/secrets/kubernetes.io/serviceaccount/..data/token vaulttoken
kubectl cp shell:/var/run/secrets/kubernetes.io/serviceaccount/..data/ca.crt ca.crt
kubectl cp shell:/var/run/secrets/kubernetes.io/serviceaccount/..data/token clienttoken
kubectl cp shell:/projected/..data/token projectedtoken


# Configure Kubernets Auth method
http POST http://localhost:8200/v1/auth/kubernetes/config X-Vault-Token:$ROOT_TOKEN kubernetes_host=$APISERVER_ADDRESS disable_iss_validation=true kubernetes_ca_cert=@ca.crt

# Or alternatively
http POST http://localhost:8200/v1/auth/kubernetes/config X-Vault-Token:$ROOT_TOKEN kubernetes_host=$APISERVER_ADDRESS disable_iss_validation=true kubernetes_ca_cert=@ca.crt disable_local_ca_jwt=true

# Configure policy and role
http POST http://localhost:8200/v1/sys/policy/my-policy X-Vault-Token:$ROOT_TOKEN policy="path \"secret/foo\" { capabilities = [\"create\", \"read\", \"delete\"] }"
http POST http://localhost:8200/v1/auth/kubernetes/role/my-role X-Vault-Token:$ROOT_TOKEN bound_service_account_names="*" bound_service_account_namespaces=default policies=my-policy
http POST http://localhost:8200/v1/auth/kubernetes/role/my-role-with-audience X-Vault-Token:$ROOT_TOKEN bound_service_account_names="*" bound_service_account_namespaces=default policies=my-policy audience=https://my-own-audience




# Unsuccessful: "open /var/run/secrets/kubernetes.io/serviceaccount/ca.crt: no such file or directory"
http POST http://localhost:8200/v1/auth/kubernetes/config X-Vault-Token:$ROOT_TOKEN kubernetes_host=$APISERVER_ADDRESS disable_iss_validation=true

# Unsuccessful: kubernetes_ca_cert must be given when disable_local_ca_jwt is true
http POST http://localhost:8200/v1/auth/kubernetes/config X-Vault-Token:$ROOT_TOKEN kubernetes_host=$APISERVER_ADDRESS disable_iss_validation=true disable_local_ca_jwt=true


# Unsuccessful auth by relaying client token that has not token-review API permissions
http POST http://localhost:8200/v1/auth/kubernetes/login role=my-role jwt=@clienttoken
http POST http://localhost:8200/v1/auth/kubernetes/login role=my-role jwt=@projectedtoken   # NOTE: token expires soon!

# Successful auth by relaying client token
http POST http://localhost:8200/v1/auth/kubernetes/login role=my-role jwt=@vaulttoken


# Reconfigure Auth Method with a copy of vault token
http POST http://localhost:8200/v1/auth/kubernetes/config X-Vault-Token:$ROOT_TOKEN kubernetes_host=$APISERVER_ADDRESS disable_iss_validation=true kubernetes_ca_cert=@ca.crt token_reviewer_jwt=@vaulttoken


# Successful auth
http POST http://localhost:8200/v1/auth/kubernetes/login role=my-role jwt=@clienttoken
http POST http://localhost:8200/v1/auth/kubernetes/login role=my-role jwt=@projectedtoken   # NOTE: token expires soon!





##################################################################################
#
# Testing service account token expiration
#

Start a new cluster with a configuration file that enables service account expiration.

```console
kind delete cluster --name exptest
kind create cluster --config configs/configs/kind-cluster-config-1h-token-expiry.yaml --name exptest
```

Observe the default service account expiration time.
The minimum that Kubernetes allows is 1 hour, so it is bit tedious from the perspective of running tests.

```console
kubectl apply -f manifests/pod-with-default-sa.yaml
kubectl exec -it sa-defaults -- python3 -c "import jwt, time, sys; t = jwt.decode(open(sys.argv[1]).read(), verify=False); print('Expires in: {}\nTime now: {}'.format(time.ctime(t['exp']), time.ctime()))" /var/run/secrets/kubernetes.io/serviceaccount/token
kubectl exec -it sa-defaults -- python3 -c "import jwt, time, sys, pprint; t = jwt.decode(open('/var/run/secrets/kubernetes.io/serviceaccount/token').read(), verify=False); pprint.pprint(t)"
```

Observe the expiration of "bound service account" tokens.
The minimum that Kubernetes allows is 10 minutes.

```console
kubectl apply -f manifests/manifests/pod-with-sa-audience-and-expiration.yaml
kubectl exec -it audience -- python3 -c "import jwt, time, sys; t = jwt.decode(open(sys.argv[1]).read(), verify=False); print('Expires in: {}\nTime now: {}'.format(time.ctime(t['exp']), time.ctime()))" /projected/token
```

To make manual API server requests run

```console
kubectl exec -it audience -- ash
http --verify=/var/run/secrets/kubernetes.io/serviceaccount/ca.crt https://$KUBERNETES_SERVICE_HOST/api/v1/namespaces/default/pods Authorization:"Bearer $(cat /var/run/secrets/kubernetes.io/serviceaccount/token)"
http --verify=/var/run/secrets/kubernetes.io/serviceaccount/ca.crt https://$KUBERNETES_SERVICE_HOST/api/v1/namespaces/default/pods Authorization:"Bearer $(cat /projected/token)"
```



## Using the new expiring id token with vault

Pre-condition https://github.com/hashicorp/vault-plugin-auth-kubernetes/issues/121

Deploy

```console
kubectl apply -f manifests/token-review-example.yaml
```

Check vault server tokens

```console
kubectl exec -it vault -- cat /var/run/secrets/kubernetes.io/serviceaccount/token | python3 -c "import jwt, time, sys, pprint; t = jwt.decode(sys.stdin.read(), verify=False); pprint.pprint(t)"

# monitor exipration of vault token
while true; do kubectl exec -it vault -- cat /var/run/secrets/kubernetes.io/serviceaccount/token | python3 -c "import jwt, time, sys, pprint; t = jwt.decode(sys.stdin.read(), verify=False); print('{}\n  Issued at: {}\n  Expires: {}\n'.format(time.ctime(), time.ctime(t['iat']), time.ctime(t['exp'])))"; sleep 60; done
```

Check vault-client tokens

```console
kubectl exec -it vault-client -- python3 -c "import jwt, time, sys, pprint; t = jwt.decode(open('/var/run/secrets/kubernetes.io/serviceaccount/token').read(), verify=False); pprint.pprint(t)"
kubectl exec -it vault-client -- python3 -c "import jwt, time, sys, pprint; t = jwt.decode(open('/projected/token').read(), verify=False); pprint.pprint(t)"
```

Copy vault binary to pod and port forward

```console
tar cf - bin/vault | kubectl exec -i vault -- tar xf - -C /usr/  # run in vault source code dir with the patch
```


Run dev server

```console
kubectl exec -it vault -- vault server -dev -dev-listen-address=0.0.0.0:8200 -log-level=debug
```

Or with file backend

```console
cat >config.hcl <<EOF

listener "tcp" {
  address     = "0.0.0.0:8200"
  tls_disable = 1
}

storage "file" {
  path = "/data"
}

disable_mlock = true
EOF

vault server -log-level=debug -config config.hcl

http -v POST http://vault:8200/v1/sys/init secret_shares:=1 secret_threshold:=1
http -v POST http://vault:8200/v1/sys/unseal key=    # from "keys" field in init response
```

Configure Kubernetes Auth Method

```console
kubectl exec -it vault-client -- ash
export ROOT_TOKEN=   # from "root_token" field in init response
http -v POST http://vault:8200/v1/sys/auth/kubernetes X-Vault-Token:$ROOT_TOKEN type=kubernetes
http -v POST http://vault:8200/v1/auth/kubernetes/config X-Vault-Token:$ROOT_TOKEN kubernetes_host=https://kubernetes disable_iss_validation=True
http -v POST http://vault:8200/v1/sys/policy/demo X-Vault-Token:$ROOT_TOKEN policy="path \"secret/demo/*\" { capabilities = [\"create\", \"read\", \"update\", \"delete\", \"list\"] }"
http -v POST http://vault:8200/v1/auth/kubernetes/role/demo-role X-Vault-Token:$ROOT_TOKEN audience=http://vault bound_service_account_names=vault-client bound_service_account_namespaces=default policies=demo
```

Use the projected service account with vault as audience

```console
http -v POST http://vault:8200/v1/auth/kubernetes/login role=demo-role jwt=@/var/run/secrets/kubernetes.io/serviceaccount/token
http -v POST http://vault:8200/v1/auth/kubernetes/login role=demo-role jwt=@/projected/token

```

## Accessing the legacy token

```console
# Create service account
$ cat <<EOF | kubectl apply -f -
apiVersion: v1
kind: ServiceAccount
metadata:
  name: my-service-account
EOF

# Find out the generated secret name
$ kubectl get secrets | grep my-service-account

# Access the token via the secret
$ kubectl get secret my-service-account-token-chdd8 -o jsonpath={.data.token} | base64 -d | python3 -c "import jwt, time, sys, pprint; t = jwt.decode(sys.stdin.read(), verify=False); pprint.pprint(t)"
{'iss': 'kubernetes/serviceaccount',
 'kubernetes.io/serviceaccount/namespace': 'default',
 'kubernetes.io/serviceaccount/secret.name': 'my-service-account-token-chdd8',
 'kubernetes.io/serviceaccount/service-account.name': 'my-service-account',
 'kubernetes.io/serviceaccount/service-account.uid': 'd5a74d52-9f1f-45d7-8240-3db9946fcd1b',
 'sub': 'system:serviceaccount:default:my-service-account'}
```

