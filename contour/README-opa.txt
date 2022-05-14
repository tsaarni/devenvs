
export WORKDIR=/home/tsaarni/work/contour-devenv

# start new cluster
kind delete cluster --name contour
kind create cluster --config configs/kind-cluster-config.yaml --name contour

kubectl apply -f examples/contour


#
# generate certificates and store them into secrets
#
certyaml --destination certs configs/certs.yaml

kubectl -n projectcontour create secret generic envoycert --from-file=tls.crt=certs/envoy.pem --from-file=tls.key=certs/envoy-key.pem --from-file=ca.crt=certs/internal-root-ca.pem --dry-run=client -o yaml | kubectl apply -f -
kubectl -n projectcontour create secret generic contourcert --from-file=tls.crt=certs/contour.pem --from-file=tls.key=certs/contour-key.pem --from-file=ca.crt=certs/internal-root-ca.pem --dry-run=client -o yaml | cat - <(echo type: kubernetes.io/tls) | kubectl apply -f -
kubectl create secret tls ingress --cert=certs/ingress.pem --key=certs/ingress-key.pem --dry-run=client -o yaml | kubectl apply -f -
kubectl create secret generic opa-envoy-cert --from-file=tls.crt=certs/opa-envoy.pem --from-file=tls.key=certs/opa-envoy-key.pem --from-file=ca.crt=certs/internal-root-ca.pem --dry-run=client -o yaml | kubectl apply -f -


#
# generate up-to-date version of CRDs and deploy them
#
make generate
kubectl apply -f examples/contour/01-crds.yaml


#
# configure test services
#
export HOST_ADDRESS=$(docker network inspect kind | jq -r '.[0].IPAM.Config[0].Gateway')
export OAUTH_DUMMY_HMAC=secret
export OAUTH_DUMMY_CLIENT_SECRET=8dd885b0-8c6d-4314-8e6c-b1e21ea9637a


# create endpoints that directs traffic to host, to execute controllers directly from source code without deploying
envsubst < manifests/contour-endpoints-dev.yaml | kubectl apply -f -
envsubst < manifests/opa-envoy-endpoints-dev.yaml | kubectl apply -f -

# deploy test services
kubectl apply -f manifests/echoserver.yaml
kubectl apply -f manifests/echoserver-opa-basic-auth.yaml
envsubst < manifests/echoserver-opa-oauth.yaml | kubectl apply -f -


#
# run services locally
#

# contour

kubectl -n projectcontour get secret contourcert -o jsonpath='{..ca\.crt}' | base64 -d > ca.crt
kubectl -n projectcontour get secret contourcert -o jsonpath='{..tls\.crt}' | base64 -d > tls.crt
kubectl -n projectcontour get secret contourcert -o jsonpath='{..tls\.key}' | base64 -d > tls.key
kubectl -n projectcontour scale deployment --replicas=0 contour
kubectl -n projectcontour rollout restart daemonset envoy

go run github.com/projectcontour/contour/cmd/contour serve --xds-address=0.0.0.0 --xds-port=8001 --envoy-service-http-port=8080 --envoy-service-https-port=8443 --contour-cafile=ca.crt --contour-cert-file=tls.crt --contour-key-file=tls.key --config-path=$WORKDIR/configs/contour.yaml

# opa-envoy

go run ./cmd/opa-envoy-plugin/... run --server --addr=localhost:8181 --diagnostic-addr=0.0.0.0:8282 --set=plugins.envoy_ext_authz_grpc.addr=:9191 --set=plugins.envoy_ext_authz_grpc.query=data.envoy.authz.allow --set=decision_logs.console=true --ignore=.* --h2c --watch $WORKDIR/configs/opa-basic-auth.rego

go run ./cmd/opa-envoy-plugin/... run --server --addr=localhost:8181 --diagnostic-addr=0.0.0.0:8282 --set=plugins.envoy_ext_authz_grpc.addr=:9191 --set=plugins.envoy_ext_authz_grpc.query=data.envoy.authz.allow --set=decision_logs.console=true --ignore=.* --h2c --watch $WORKDIR/configs/opa-oauth.rego


# keycloak

cd ~/work/keycloak
export WORKDIR=/home/tsaarni/work/keycloak-devenv
mvn -f testsuite/utils/pom.xml exec:java -Pkeycloak-server -Dkeycloak.migration.action=import -Dkeycloak.migration.provider=dir -Dkeycloak.migration.dir=$WORKDIR/migrations/oauth -Djavax.net.ssl.trustStore=$WORKDIR/truststore.p12 -Djavax.net.ssl.trustStorePassword=password -Djavax.net.ssl.javax.net.ssl.trustStoreType="PKCS12" -Djavax.net.ssl.keyStore=$WORKDIR/admin-keystore.p12 -Djavax.net.ssl.keyStorePassword=password -Djavax.net.ssl.javax.net.ssl.keyStoreType="PKCS12" -Dresources -Dexec.args="-b 0.0.0.0"


        # - "--tls-ca-cert-file=/certs/ca.crt"
        # - "--tls-cert-file=/certs/tls.crt"
        # - "--tls-private-key-file=/certs/tls.key"




http -v http://unprotected.127-0-0-101.nip.io/
http -v --verify=certs/external-root-ca.pem https://protected-basic-auth.127-0-0-101.nip.io/
http -v --verify=certs/external-root-ca.pem https://protected-basic-auth.127-0-0-101.nip.io/ Authorization:"Basic charlie"

TOKEN=$(http --form POST http://localhost:8081/auth/realms/master/protocol/openid-connect/token username=admin password=admin grant_type=password client_id=admin-cli | jq -r .access_token)
echo $TOKEN | jwt -show -

TOKEN=$(http --form POST http://localhost:8081/auth/realms/master/protocol/openid-connect/token username=joe password=joe grant_type=password client_id=test-client client_secret=${OAUTH_DUMMY_CLIENT_SECRET} scope=myscope | jq -r .access_token)
echo $TOKEN | jwt -show -

http -v --verify=certs/external-root-ca.pem https://protected-oauth.127-0-0-101.nip.io/allowed Authorization:"Bearer $TOKEN"












# change envoy debug level (default: info)
kubectl -n projectcontour get daemonsets.apps envoy -o yaml | sed 's/--log-level \(.*\)/--log-level debug/' | kubectl apply -f -

kubectl -n projectcontour logs $(kubectl -n projectcontour get pod -l app=envoy -o jsonpath='{.items[0].metadata.name}') -c envoy -f




kubectl -n projectcontour port-forward $(kubectl -n projectcontour get pod -lapp=envoy -ojsonpath="{.items[0].metadata.name}") 9001

http http://localhost:9001/config_dump?include_eds | jq -C . | less
http http://localhost:9001/config_dump| jq '.configs[].dynamic_active_clusters'
http http://localhost:9001/config_dump?mask=EndpointsConfigDump
http http://localhost:9001/listeners?format=json


kubectl apply -f manifests/echoserver-opa.yaml


echo DECISION_FROM_OPA_LOGS | http "http://localhost:8181/v1/data/envoy/authz/allow?explain=full&pretty"






# install selenium-python

python3 -m venv selenium
. selenium/bin/activate
pip install selenium
pip install webdrivermanager
webdrivermanager chrome


# run login test
./keycloak-login-test.py
