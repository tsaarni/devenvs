

https://github.com/hashicorp/vault-plugin-auth-kubernetes/pull/300




# Create service account and pod using the service account
$ kubectl apply -f - <<EOF
apiVersion: v1
kind: ServiceAccount
metadata:
  name: my-sa
---
apiVersion: v1
kind: Pod
metadata:
  name: my-pod
spec:
  serviceAccountName: my-sa
  containers:
  - name: my-container
    image: alpine:3
    command: ["sleep", "9999999"]
EOF

# Decode the service account token
$ kubectl exec my-pod -- cat /var/run/secrets/kubernetes.io/serviceaccount/token | jq -R 'split(".") | .[1] | @base64d | fromjson'
{
  "aud": [
    "https://kubernetes.default.svc.cluster.local"
  ],
  ...
}



#
# Create a secret of type kubernetes.io/service-account-token referring to my-sa
#  https://kubernetes.io/docs/concepts/configuration/secret/#serviceaccount-token-secrets
#
$ kubectl apply -f - <<EOF
apiVersion: v1
kind: Secret
metadata:
  name: my-sa-token
  annotations:
    kubernetes.io/service-account.name: my-sa
type: kubernetes.io/service-account-token
EOF

# Decode the service account token from secret
$ kubectl get secret my-sa-token -o jsonpath='{.data.token}' | base64 --decode | jq -R 'split(".") | .[1] | @base64d | fromjson'
{
  ... there is no "aud" claim at all
  ...
}



Create service account using TokenRequest API

$ kubectl create token my-sa --audience=whatever-i-want | jq -R 'split(".") | .[1] | @base64d | fromjson'
{
  "aud": [
    "whatever-i-want"
  ],
  ...
}


Create token using projected volume

$ kubectl apply -f - <<EOF
apiVersion: v1
kind: Pod
metadata:
  name: my-pod-projected
spec:
  serviceAccountName: my-sa
  containers:
  - name: my-container
    image: alpine:3
    command: ["sleep", "9999999"]
    volumeMounts:
    - name: token-volume
      mountPath: /var/run/secrets/kubernetes.io/my-projected-serviceaccount
  volumes:
  - name: token-volume
    projected:
      sources:
      - serviceAccountToken:
          path: token
          audience: bar
EOF

$ kubectl exec my-pod-projected -- cat /var/run/secrets/kubernetes.io/my-projected-serviceaccount/token | jq -R 'split(".") | .[1] | @base64d | fromjson'
{
  "aud": [
    "bar"
  ],
  ...
}
