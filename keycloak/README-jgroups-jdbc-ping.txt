
# start new cluster
kind delete cluster --name keycloak
kind create cluster --config configs/kind-cluster-config.yaml --name keycloak

kubectl apply -f https://projectcontour.io/quickstart/contour.yaml


rm -rf certs
mkdir -p certs
certyaml -d certs configs/certs.yaml
# convert pem key to der
openssl pkcs8 -topk8 -inform PEM -outform DER -in certs/keycloak-internal-key.pem -out certs/keycloak-internal-key.der -nocrypt



kubectl create secret generic keycloak-internal --from-file=certs/internal-ca.pem --from-file=certs/keycloak-internal.pem --from-file=certs/keycloak-internal-key.der --dry-run=client -o yaml | kubectl apply -f -
kubectl create secret generic postgres-certs --from-file=certs/internal-ca.pem --from-file=certs/postgres-internal.pem --from-file=certs/postgres-internal-key.pem --dry-run=client -o yaml | kubectl apply -f -


kubectl apply -f manifests/postgresql.yaml
kubectl apply -f manifests/keycloak-26-jgroups-jdbc-ping.yaml
