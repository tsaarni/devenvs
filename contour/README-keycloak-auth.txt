# start new cluster
kind delete cluster --name contour
kind create cluster --config ~/work/devenvs/contour/configs/kind-cluster-config.yaml --name contour


# generate certificates and create secrets
rm -f certs/*
mkdir -p certs
certyaml --destination certs configs/certs.yaml

kubectl create secret tls ingress --cert=certs/ingress.pem --key=certs/ingress-key.pem --dry-run=client -o yaml | kubectl apply -f -
kubectl --namespace projectcontour create secret generic contour-authserver-cert --from-file=tls.crt=certs/contour-authserver.pem --from-file=tls.key=certs/contour-authserver-key.pem --from-file=ca.crt=certs/internal-root-ca.pem --dry-run=client -o yaml | kubectl apply -f -
kubectl --namespace projectcontour create secret generic internal-root-ca --from-file=ca.crt=certs/internal-root-ca.pem --dry-run=client -o yaml | kubectl apply -f -

# deploy contour
kubectl apply -f https://projectcontour.io/quickstart/contour.yaml

# deploy keycloak and postgres
kubectl apply -f manifests/keycloak.yaml
kubectl apply -f manifests/postgres.yaml



# deploy echoserver
kubectl apply -f manifests/echoserver-extauth.yaml



# configure keycloak

1. Login to keyclaok as admin:admin https://keycloak.127-0-0-101.nip.io/
2. Clients -> Create client
     set Client ID: contour-authserver



kubectl apply -f manifests/contour-authserver.yaml





http --verify certs/external-root-ca.pem https://protected-oauth.127-0-0-101.nip.io


    kubectl -n projectcontour logs daemonset/envoy -c envoy -f
kubectl -n projectcontour logs deployment/contour-authserver -f
kubectl logs statefulset/keycloak -f





*** Run locally for development

kubectl -n projectcontour scale deployment --replicas=0 contour-authserver

# create endpoints that directs traffic to host, to execute controllers directly from source code without deploying
sed "s/REPLACE_ADDRESS_HERE/$(docker network inspect kind | jq -r '.[0].IPAM.Config[0].Gateway')/" manifests/contour-authserver-endpoints-dev.yaml | kubectl apply -f -

mkdir -p .vscode
cp ~/work/devenvs/contour/configs/contour-authserver-vscode-launch.json .vscode/launch.json


cat <<EOF >contour-authserver-config.yaml

address: ":19443"
issuerURL: "http://keycloak.127-0-0-101.nip.io/realms/master"
redirectURL: "https://protected-oauth.127-0-0-101.nip.io"
redirectPath: "/callback"
allowEmptyClientSecret: false
scopes:
- openid
- profile
- email
- offline_access
clientID: "contour-authserver"
clientSecret: "llh6qjqUCY9zcTLbRP5eZcRa08T3ZKyE"

EOF
