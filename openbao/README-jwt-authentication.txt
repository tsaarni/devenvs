
# Create Kind cluster
kind delete cluster --name openbao
kind create cluster --config configs/kind-cluster-config.yaml --name openbao


# Generate certs
rm -rf certs
mkdir -p certs
certyaml -d certs configs/certs.yaml


# Create the certs secret
kubectl create secret generic openbao-certs --dry-run=client -o yaml --from-file=certs/ca.pem --from-file=certs/bao-server.pem --from-file=certs/bao-server-key.pem | kubectl apply -f -

# Deploy openbao
kubectl apply -f manifests/openbao.yaml


# Initialize the openbao on the first pod
kubectl exec -it openbao-0 -- bao operator init -key-shares=1 -key-threshold=1 -format=json > openbao-unseal-config.json


# Unseal
kubectl exec -it openbao-0 -- bao operator unseal $(jq -r .unseal_keys_b64[0] openbao-unseal-config.json)
kubectl exec -it openbao-1 -- bao operator unseal $(jq -r .unseal_keys_b64[0] openbao-unseal-config.json)
kubectl exec -it openbao-2 -- bao operator unseal $(jq -r .unseal_keys_b64[0] openbao-unseal-config.json)



# Authenticate with the root token
kubectl exec -it openbao-0 -- bao login $(jq -r .root_token openbao-unseal-config.json)
kubectl exec -it openbao-0 -- bao operator raft list-peers





# Enable JWT authentication
kubectl exec -it openbao-0 -- bao auth enable jwt


###########################
#
# Configure Kubernetes sa.pem manually
#


# Deploy shell pod
kubectl apply -f manifests/shell.yaml


kubectl exec -it shell -- ash

http --verify=/var/run/secrets/kubernetes.io/serviceaccount/ca.crt https://kubernetes.default.svc/.well-known/openid-configuration Authorization:"Bearer $(cat /var/run/secrets/kubernetes.io/serviceaccount/token)"
http --verify=/var/run/secrets/kubernetes.io/serviceaccount/ca.crt https://kubernetes.default.svc/openid/v1/jwks Authorization:"Bearer $(cat /var/run/secrets/kubernetes.io/serviceaccount/token)"

kubectl exec -it shell -- cat /var/run/secrets/kubernetes.io/serviceaccount/token | jq -R 'split(".") | .[1] | @base64d | fromjson | .exp = (.exp | todate) | .iat = (.iat | todate)'


kubectl exec -it shell -- /host/apps/jwks-to-pem.py > sa.pub
kubectl cp sa.pub openbao-0:.


kubectl exec -it openbao-0 -- bao write auth/jwt/config jwt_validation_pubkeys=@sa.pub bound_issuer="https://kubernetes.default.svc.cluster.local" bound_audience="https://kubernetes.default.svc.cluster.local"

# Create roles and policies
kubectl exec -it openbao-0 -- bao write auth/jwt/role/secret-writer - <<EOF
{
    "role_type": "jwt",
    "user_claim": "sub",
    "bound_subject": "system:serviceaccount:default:openbao-client",
    "bound_audiences": "https://kubernetes.default.svc.cluster.local",
    "token_policies": "secret-writer"
}
EOF

kubectl exec -it openbao-0 -- bao write auth/jwt/role/secret-reader - <<EOF
{
    "role_type": "jwt",
    "user_claim": "sub",
    "bound_audiences": "https://kubernetes.default.svc.cluster.local",
    "bound_claims": {
        "/kubernetes.io/namespace": "default"
    },
    "token_policies": "secret-reader"
}
EOF

kubectl exec -it openbao-0 -- bao policy write secret-reader -<<EOF
path "secret/data/my-secret" {
  capabilities = ["read"]
}
EOF

kubectl exec -it openbao-0 -- bao policy write secret-writer -<<EOF
path "secret/data/my-secret" {
  capabilities = ["read", "create", "update"]
}
EOF



# Create a test secret
kubectl exec -it openbao-0 -- bao secrets enable -path=secret kv
kubectl exec -it openbao-0 -- bao write secret/data/my-secret data='{"foo":"bar"}'




# Test login and read secret
kubectl exec -it shell -- ash

# Read secret

http --verify=no POST https://openbao:8200/v1/auth/jwt/login jwt="$(cat /var/run/secrets/kubernetes.io/serviceaccount/token)" role=secret-reader | jq -r .auth.client_token > token.txt
http --verify=no GET https://openbao:8200/v1/secret/data/my-secret X-Vault-Token:"$(cat token.txt)"

# Write secret

http --verify=no POST https://openbao:8200/v1/auth/jwt/login jwt="$(cat /var/run/secrets/kubernetes.io/serviceaccount/token)" role=secret-writer | jq -r .auth.client_token > token.txt
http --verify=no POST https://openbao:8200/v1/secret/data/my-secret X-Vault-Token:"$(cat token.txt)" data["foo"]=baz




###########################
#
# Run under debugger
#


# Create Kind cluster
kind delete cluster --name openbao
kind create cluster --config configs/kind-cluster-config.yaml --name openbao


kubectl apply -f manifests/shell.yaml


kubectl port-forward -n kube-system pod/kube-apiserver-openbao-control-plane 16443:6443

kubectl exec shell -- cat /var/run/secrets/kubernetes.io/serviceaccount/token > token
kubectl exec shell -- cat /var/run/secrets/kubernetes.io/serviceaccount/ca.crt > ca.crt

jq -R 'split(".") | .[1] | @base64d | fromjson | .exp = (.exp | todate) | .iat = (.iat | todate)' < token




wireshark -i lo -k -Y http2 -o tls.keylog_file:$HOME/work/openbao/wireshark-keys.log



### With OIDC discovery

export BAO_ADDR=http://localhost:8200

bin/bao login root
bin/bao auth enable jwt
bin/bao write auth/jwt/config - <<EOF
{
  "provider_config": {
    "provider": "kubernetes"
  }
}
EOF

# create role and policy
bin/bao write auth/jwt/role/secret-writer - <<EOF
{
    "role_type": "jwt",
    "user_claim": "sub",
    "bound_subject": "system:serviceaccount:default:openbao-client",
    "bound_audiences": "https://kubernetes.default.svc.cluster.local",
    "token_policies": "secret-writer"
}
EOF

bin/bao policy write secret-writer -<<EOF
path "secret/data/my-secret" {
  capabilities = ["read", "create", "update"]
}
EOF

http --verify=no POST http://localhost:8200/v1/auth/jwt/login jwt="$(cat token)" role=secret-writer




http POST http://openbao:8200/v1/auth/jwt/login jwt="$(cat /var/run/secrets/kubernetes.io/serviceaccount/token)" role=secret-writer


sudo nsenter --target $(kindps openbao openbao-0 -o json | jq -r .[0].pids[0].pid) --net wireshark -i any -k -Y http2 -o tls.keylog_file:/proc/$(kindps openbao openbao-0 -o json | jq -r .[0].pids[0].pid)/root/tmp/wireshark-keys.log







// Global variables instead of const to allow test cases to overwrite them.
var (
	// localJWTPath is the path to the Kubernetes Service Account token file.
	localJWTPath = "/var/run/secrets/kubernetes.io/serviceaccount/token"
	// localJWTPath = "token"

	// localCACertPath is the path to the Kubernetes CA certificate.
	localCACertPath = "/var/run/secrets/kubernetes.io/serviceaccount/ca.crt"
	// localCACertPath = "ca.crt"
)



	keylogger := os.Getenv("SSLKEYLOGFILE")
	var keylogWriter io.Writer
	if keylogger != "" {
		f, err := os.OpenFile(keylogger, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0o600)
		if err == nil {
			keylogWriter = f
		}
	}

	tlsConfig := &tls.Config{
		RootCAs: certPool,
		KeyLogWriter: keylogWriter,
	}




############
#
# Run under kind
#

make
cp -a bin/bao ~/work/devenvs/openbao/


kubectl exec -it shell -- /host/bao server -dev -dev-root-token-id=root -log-level=debug
kubectl port-forward pod/shell 8200:8200


## follow the same sequence as "With OIDC discovery" above
