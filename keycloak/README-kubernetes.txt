


# Create kind cluster and deploy services

kind delete cluster --name keycloak
kind create cluster --config configs/kind-cluster-config.yaml --name keycloak



# Create certs
rm -rf certs
mkdir -p certs
certyaml -d certs configs/certs.yaml


kubectl create secret tls keycloak-external --cert=certs/keycloak-server.pem --key=certs/keycloak-server-key.pem --dry-run=client -o yaml | kubectl apply -f -




kubectl apply -f https://projectcontour.io/quickstart/contour.yaml
kubectl apply -f manifests/postgresql.yaml



docker buildx build docker/openldap/ -t localhost/openldap:latest
kind load docker-image --name keycloak localhost/openldap:latest
kubectl apply -f manifests/openldap.yaml

docker buildx build docker/ldap-client/ -t localhost/ldap-client:latest
kind load docker-image --name keycloak localhost/ldap-client:latest
kubectl apply -f manifests/sshd.yaml



# URLs

http://keycloak.127-0-0-121.nip.io/
https://keycloak.127-0-0-121.nip.io/




##################################
#
# Running keycloak locally with mirrord to debug against cluster services
#
#  - Requires vscode extension: MetalBear.mirrord
#    Downloads binary here on demaind
#    ~/.config/Code/User/globalStorage/metalbear.mirrord/mirrord
#
#  - verify config
#    /home/tsaarni/.config/Code/User/globalStorage/metalbear.mirrord/mirrord verify-config /home/tsaarni/work/devenvs/keycloak/configs/mirrord.yaml

kubectl apply -f manifests/keycloak-dev.yaml


# Copy VS Code configs to keycloak source tree
mkdir -p .vscode
cp -a ~/work/devenvs/keycloak/configs/launch.json ~/work/devenvs/keycloak/configs/settings.json .vscode


# Launch VS Code and start debugging:
#   - Select "Debug Quarkus with mirrord (keycloak-junit5)" launch config
#
# mirrord config: configs/mirrord.yaml


# Wireshark

sudo nsenter --target $(pidof slapd) --net wireshark -f "port 389" -k
sudo nsenter --target $(pgrep -f quarkus) --net wireshark -f "port 8080" -Y http -k


# Delete postgres volumes to start fresh

kubectl scale statefulset keycloak --replicas 0
kubectl scale statefulset postgres --replicas 0
kubectl delete persistentvolumeclaims data-postgres-0
kubectl scale statefulset postgres --replicas 1
kubectl scale statefulset keycloak --replicas 1


# Operator
# https://www.keycloak.org/operator/installation
# https://www.keycloak.org/operator/basic-deployment

kubectl create secret generic keycloak-db-secret --from-literal=username=keycloak --from-literal=password=keycloak

cat <<EOF | kubectl apply -f -
apiVersion: k8s.keycloak.org/v2alpha1
kind: Keycloak
metadata:
  name: keycloak
spec:
  instances: 1
  db:
    vendor: postgres
    host: postgres
    usernameSecret:
      name: keycloak-db-secret
      key: username
    passwordSecret:
      name: keycloak-db-secret
      key: password
  http:
    httpEnabled: true
    tlsSecret: keycloak-internal
  hostname:
    hostname: keycloak.127-0-0-121.nip.io
  proxy:
    headers: xforwarded
EOF

cat <<EOF | kubectl apply -f -
apiVersion: projectcontour.io/v1
kind: HTTPProxy
metadata:
  name: keycloak
spec:
  virtualhost:
    fqdn: keycloak.127-0-0-121.nip.io
    tls:
      secretName: keycloak-external
  routes:
    - services:
        - name: keycloak-service
          port: 8080
      permitInsecure: true
EOF
