



cd ~/work/openfga-envoy
make build
docker build -t localhost/openfga-envoy:latest -f extauthz/Dockerfile extauthz
kind load docker-image --name contour localhost/openfga-envoy:latest


kubectl create secret tls ingress --cert=certs/ingress.pem --key=certs/ingress-key.pem --dry-run=client -o yaml | kubectl apply -f -

kubectl apply -f https://projectcontour.io/quickstart/contour.yaml
kubectl apply -f manifests/openfga-demo.yaml
kubectl apply -f manifests/keycloak.yaml



kubectl logs deployment/openfga -f
kubectl logs deployment/openfga-envoy -f
kubectl -n projectcontour logs daemonsets/envoy -c envoy -f


http --verify=certs/external-root-ca.pem https://protected.127-0-0-101.nip.io

# change debug level from "info" to "debug"
kubectl -n projectcontour edit daemonsets.apps envoy





http keycloak.127-0-0-101.nip.io

http http://keycloak.127-0-0-101.nip.io/realms/master/.well-known/openid-configuration
http http://keycloak.127-0-0-101.nip.io/realms/master/protocol/openid-connect/certs

http http://keycloak.127-0-0-101.nip.io/realms/contour/.well-known/openid-configuration
http http://keycloak.127-0-0-101.nip.io/realms/contour/protocol/openid-connect/certs


http --form POST http://keycloak.127-0-0-101.nip.io/realms/master/protocol/openid-connect/token username=admin password=admin grant_type=password client_id=admin-cli
http --form POST http://keycloak.127-0-0-101.nip.io/realms/contour/protocol/openid-connect/token username=joe password=password grant_type=password client_id=admin-cli


# client secret flow
http --form POST http://keycloak.127-0-0-101.nip.io/realms/contour/protocol/openid-connect/token client_id=contour client_secret= grant_type=client_credentials


# decode jwt by running the following command and pasting the access token on the console
jq -R 'split(".") | .[1] | @base64d | fromjson'


TOKEN=$(http --form POST http://keycloak.127-0-0-101.nip.io/realms/contour/protocol/openid-connect/token username=joe password=password grant_type=password client_id=admin-cli | jq -r .access_token)
http --verify=certs/external-root-ca.pem https://protected.127-0-0-101.nip.io Authorization:"Bearer $TOKEN"


##########################
#
# Debugging
#

### Start with configuration file

"args": ["serve", "--xds-address=0.0.0.0", "--xds-port=8001", "--envoy-service-http-port=8080", "--envoy-service-https-port=8443", "--contour-cafile=ca.crt", "--contour-cert-file=tls.crt", "--contour-key-file=tls.key", "--debug", "--config-path=/home/tsaarni/work/devenvs/contour/configs/contour-config-global-ext-authz.yaml"]




### Start with ContourConfiguration

kubectl apply -f manifests/contourconfig-global-ext-authz.yaml

"args": ["serve", "--xds-address=0.0.0.0", "--xds-port=8001", "--envoy-service-http-port=8080", "--envoy-service-https-port=8443", "--contour-cafile=ca.crt", "--contour-cert-file=tls.crt", "--contour-key-file=tls.key", "--debug", "--contour-config-name=contour"]
