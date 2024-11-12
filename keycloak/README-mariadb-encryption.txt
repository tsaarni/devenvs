

# start new cluster
kind delete cluster --name keycloak
kind create cluster --config configs/kind-cluster-config.yaml --name keycloak

# Build mariadb image with vault plugin
docker build docker/mariadb/ -t localhost/mariadb:latest
kind load docker-image --name keycloak localhost/mariadb:latest


kubectl apply -f https://projectcontour.io/quickstart/contour.yaml

kubectl create secret tls keycloak-external --cert=certs/keycloak-server.pem --key=certs/keycloak-server-key.pem --dry-run=client -o yaml | kubectl apply -f -

kubectl apply -f manifests/mariadb.yaml
kubectl apply -f manifests/keycloak-26.yaml
