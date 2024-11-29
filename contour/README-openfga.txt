



cd ~/work/openfga-envoy
make build
docker build -t localhost/openfga-envoy:latest -f extauthz/Dockerfile extauthz
kind load docker-image --name contour localhost/openfga-envoy:latest


kubectl create secret tls ingress --cert=certs/ingress.pem --key=certs/ingress-key.pem --dry-run=client -o yaml | kubectl apply -f -

kubectl apply -f https://projectcontour.io/quickstart/contour.yaml
kubectl apply -f manifests/openfga-demo.yaml
kubectl apply -f manifests/keycloak.yaml



kubectl logs deployment/openfga -f

# Create store.
# https://openfga.dev/api/service#/Stores/CreateStore

kubectl port-forward $(kubectl get pod -l app=openfga -o jsonpath='{.items[0].metadata.name}') 8080:8080
http POST localhost:8080/stores name=openfga

#output:
# {
#     "created_at": "2024-11-18T13:36:11.969044347Z",
#     "id": "01JCZQSV613FYED31C7M4242JZ",
#     "name": "openfga",
#     "updated_at": "2024-11-18T13:36:11.969044347Z"
# }


# Edit the ID in the store_id field
kubectl edit configmap openfga-envoy-config




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


get_token() {
  http --form POST http://keycloak.127-0-0-101.nip.io/realms/contour/protocol/openid-connect/token username=joe password=password grant_type=password client_id=admin-cli | jq -r .access_token
}

http --verify=certs/external-root-ca.pem https://protected.127-0-0-101.nip.io Authorization:"Bearer $(get_token)"


##########################
#
# Debugging
#

### Start with configuration file

"args": ["serve", "--xds-address=0.0.0.0", "--xds-port=8001", "--envoy-service-http-port=8080", "--envoy-service-https-port=8443", "--contour-cafile=ca.crt", "--contour-cert-file=tls.crt", "--contour-key-file=tls.key", "--debug", "--config-path=/home/tsaarni/work/devenvs/contour/configs/contour-config-global-ext-authz.yaml"]




### Start with ContourConfiguration

kubectl apply -f manifests/contourconfig-global-ext-authz.yaml

"args": ["serve", "--xds-address=0.0.0.0", "--xds-port=8001", "--envoy-service-http-port=8080", "--envoy-service-https-port=8443", "--contour-cafile=ca.crt", "--contour-cert-file=tls.crt", "--contour-key-file=tls.key", "--debug", "--contour-config-name=contour"]








### Disable the external authz filter

# https://github.com/projectcontour/contour/pull/6661


# Build Contour

make container
docker tag ghcr.io/projectcontour/contour:$(git rev-parse --short HEAD) localhost/contour:latest
kind load docker-image localhost/contour:latest --name contour

cat <<EOF | kubectl -n projectcontour patch deployment contour --patch-file=/dev/stdin
spec:
  template:
    spec:
      containers:
      - name: contour
        image: localhost/contour:latest
        imagePullPolicy: Never
EOF


# Create global policy and restart Contour

kubectl edit -n projectcontour configmaps contour

    globalExtAuth:
      extensionService: "default/openfga-envoy"
      authPolicy:
        #disabled: false
        disabled: true

kubectl -n projectcontour delete pod -l app=contour


# Global policy takes effect: there will be no call to the external authz service

kubectl apply -f - <<EOF
apiVersion: projectcontour.io/v1
kind: HTTPProxy
metadata:
  name: echoserver
spec:
  virtualhost:
    fqdn: protected.127-0-0-101.nip.io
    tls:
      secretName: ingress
  routes:
    - services:
        - name: echoserver
          port: 80
EOF

http --verify=certs/external-root-ca.pem https://protected.127-0-0-101.nip.io


# Virtual host level override: external authz is called

kubectl apply -f - <<EOF
apiVersion: projectcontour.io/v1
kind: HTTPProxy
metadata:
  name: echoserver
spec:
  virtualhost:
    fqdn: protected.127-0-0-101.nip.io
    tls:
      secretName: ingress
    authorization:
      authPolicy:
        disabled: false
  routes:
    - services:
        - name: echoserver
          port: 80

EOF

http --verify=certs/external-root-ca.pem https://protected.127-0-0-101.nip.io


# Route level override: external authz is called

kubectl apply -f - <<EOF
apiVersion: projectcontour.io/v1
kind: HTTPProxy
metadata:
  name: echoserver
spec:
  virtualhost:
    fqdn: protected.127-0-0-101.nip.io
    tls:
      secretName: ingress
  routes:
    - conditions:
        - prefix: /disabled
      services:
        - name: echoserver
          port: 80
    #  authPolicy:
    #    disabled: true
    - conditions:
        - prefix: /enabled
      services:
        - name: echoserver
          port: 80
      authPolicy:
        disabled: false
EOF

http --verify=certs/external-root-ca.pem https://protected.127-0-0-101.nip.io/disabled
http --verify=certs/external-root-ca.pem https://protected.127-0-0-101.nip.io/enabled



#### CHANGE CONFIG FILE TO ENABLE AUTHZ BY DEFAULT


kubectl apply -f - <<EOF
apiVersion: projectcontour.io/v1
kind: HTTPProxy
metadata:
  name: echoserver
spec:
  virtualhost:
    fqdn: protected.127-0-0-101.nip.io
    tls:
      secretName: ingress
    authorization:
      authPolicy:
        disabled: true
  routes:
    - services:
        - name: echoserver
          port: 80

EOF

http --verify=certs/external-root-ca.pem https://protected.127-0-0-101.nip.io



kubectl apply -f - <<EOF
apiVersion: projectcontour.io/v1
kind: HTTPProxy
metadata:
  name: echoserver
spec:
  virtualhost:
    fqdn: protected.127-0-0-101.nip.io
    tls:
      secretName: ingress
  routes:
    - conditions:
        - prefix: /disabled
      services:
        - name: echoserver
          port: 80
      authPolicy:
        disabled: true
    - conditions:
        - prefix: /enabled
      services:
        - name: echoserver
          port: 80
      authPolicy:
        disabled: false
EOF

http --verify=certs/external-root-ca.pem https://protected.127-0-0-101.nip.io/disabled
http --verify=certs/external-root-ca.pem https://protected.127-0-0-101.nip.io/enabled


kubectl apply -f manifests/contourconfig-global-ext-authz.yaml





kubectl apply -f - <<EOF
apiVersion: projectcontour.io/v1
kind: HTTPProxy
metadata:
  name: echoserver
spec:
  virtualhost:
    fqdn: protected.127-0-0-101.nip.io
    tls:
      secretName: ingress
    jwtProviders:
      - name: keycloak
        issuer: "http://keycloak.127-0-0-101.nip.io/realms/contour"
        remoteJWKS:
          # Must be FQDN if not running in the same namespace as Envoy.
          uri: "http://keycloak.default.svc.cluster.local:8080/realms/contour/protocol/openid-connect/certs"
    #authorization:
    #  authPolicy:
    #    disabled: false
    #  extensionRef:
    #    name: openfga-envoy
    #    namespace: default

  routes:
    - services:
        - name: echoserver
          port: 80
      jwtVerificationPolicy:
        require: keycloak
      #authPolicy:
      #  disabled: false
EOF

get_token() {
  http --form POST http://keycloak.127-0-0-101.nip.io/realms/contour/protocol/openid-connect/token username=joe password=password grant_type=password client_id=admin-cli | jq -r .access_token
}

http --verify=certs/external-root-ca.pem https://protected.127-0-0-101.nip.io Authorization:"Bearer $(get_token)"
