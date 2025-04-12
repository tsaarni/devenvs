
kind delete cluster --name kind

cd ~/work/cert-manager


# Add --namespace=cert-manager to make cert-manager list only resources in the cert-manager namespace.
# Add leaderElection.namespace to make cert-manager use the cert-manager namespace for leader election.
patch -p1 < ~/work/devenvs/cert-manager/patches/namespaced-e2e-setup.patch



# Change ClusterRoles and ClusterRoleBindings to Roles and RoleBindings
#   Note: this does not make perfect modifications but it is good enough for testing
sed -i '
  s/kind: ClusterRole$/kind: Role/;
  s/kind: ClusterRoleBinding$/kind: RoleBinding/;
  /metadata:/!b;n;/namespace:/!a\  namespace: {{ include "cert-manager.namespace" . }}
  /roleRef:/!b;n;s/kind: ClusterRole/kind: Role/
' deploy/charts/cert-manager/templates/rbac.yaml


# Remove the webhook
rm deploy/charts/cert-manager/templates/webhook-*


### Building


# Force rebuild for everything
make clean

# Start kind cluster, compile and deploy (or redeloy after code change)
make e2e-setup-certmanager


# Check that issuing a certificate works
kubectl apply -f manifests/certificates.yaml
kubectl -n cert-manager get secrets my-end-entity-cert-secret -o yaml



kubectl -n cert-manager delete pod -l app=cert-manager

kubectl -n cert-manager logs deployments/cert-manager


helm uninstall cert-manager -n cert-manager
kubectl delete namespaces cert-manager  --force



# Check that cert-manager is executed with argument
#    --namespace=cert-manager
kindps kind cert-manager

# Check that logs have this
I0409 13:36:53.211597       1 controller.go:231] "skipping as cert-manager is scoped to a single namespace" logger="cert-manager.controller" controller="clusterissuers"


######

# Run within vscode debugger



# Shut down the cert-manager instance running in the kind cluster
kubectl -n cert-manager scale deployment cert-manager --replicas 0


# Create long-lived service account token for cert-manager service account

kubectl apply -f - <<EOF
apiVersion: v1
kind: Secret
metadata:
  name: cert-manager-sa-token
  namespace: cert-manager
  annotations:
    kubernetes.io/service-account.name: cert-manager
type: kubernetes.io/service-account-token
EOF

# Create copy of the kube config file and modify it
cp ~/.kube/config kubeconfig-cert-manager-sa.yaml

#remove client certificate and key
kubectl config --kubeconfig=kubeconfig-cert-manager-sa.yaml unset users.kind-kind.client-certificate-data
kubectl config --kubeconfig=kubeconfig-cert-manager-sa.yaml unset users.kind-kind.client-key-data

TOKEN=$(kubectl get secret cert-manager-sa-token -n cert-manager -o jsonpath='{.data.token}' | base64 --decode)

kubectl config --kubeconfig=kubeconfig-cert-manager-sa.yaml set-credentials kind-kind --token="${TOKEN}"
