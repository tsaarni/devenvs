

# Following instructions are for testing the token review API manually.



# create a pod that we use to generate a token

kubectl apply -f - <<EOF
apiVersion: v1
kind: Pod
metadata:
  name: shell
spec:
  serviceAccountName: shell
  containers:
  - name: shell
    image: alpine:3.20
    command:
    - /bin/sh
    - -c
    - sleep infinity;
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: shell
EOF


# get the token

TOKEN=$(kubectl exec shell -- cat /var/run/secrets/kubernetes.io/serviceaccount/token)


# create a POST request towards token review API (should succeed)

kubectl create --raw /apis/authentication.k8s.io/v1/tokenreviews -f -  <<EOF | jq .
{
  "kind": "TokenReview",
  "apiVersion": "authentication.k8s.io/v1",
  "spec": {
    "token": "$TOKEN"
  }
}
EOF


# delete the service account

kubectl delete serviceaccount shell



# create a POST request towards token review API (should fail)

kubectl create --raw /apis/authentication.k8s.io/v1/tokenreviews -f -  <<EOF | jq .
{
  "kind": "TokenReview",
  "apiVersion": "authentication.k8s.io/v1",
  "spec": {
    "token": "$TOKEN"
  }
}
EOF


# delete the pod
kubectl delete pod shell --force
