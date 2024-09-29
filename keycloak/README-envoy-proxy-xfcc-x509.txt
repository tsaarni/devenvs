


# Start envoy proxy
docker run --rm --volume $PWD:/input:ro --network host envoyproxy/envoy:v1.31.0 envoy --config-path /input/configs/envoy-xfcc.yaml

# run echoserver backend for testing
docker run --rm --network host --env HTTP_PORT=8080  gcr.io/k8s-staging-ingressconformance/echoserver:v20210922-cec7cf2

# test making requests with and without client certificate
http --cert certs/x509client.pem --cert-key certs/x509client-key.pem --verify certs/ca.pem https://localhost:8443/
http --verify certs/ca.pem https://localhost:8443/





# start new cluster
kind delete cluster --name keycloak
kind create cluster --config configs/kind-cluster-config.yaml --name keycloak


# create certs
rm -rf certs
mkdir -p certs
certyaml -d certs configs/certs.yaml
keytool -importcert -storetype PKCS12 -keystore certs/truststore.p12 -storepass secret -noprompt -alias ca -file certs/ca.pem

kubectl create secret tls keycloak-external --cert=certs/keycloak-server.pem --key=certs/keycloak-server-key.pem --dry-run=client -o yaml | kubectl apply -f -
kubectl create secret tls keycloak-internal --cert=certs/keycloak-internal.pem --key=certs/keycloak-internal-key.pem --dry-run=client -o yaml | kubectl apply -f -
kubectl create secret generic internal-ca --from-file=ca.crt=certs/internal-ca.pem --dry-run=client -o yaml | kubectl apply -f -
kubectl create secret generic external-ca --from-file=certs/truststore.p12 --dry-run=client -o yaml | kubectl apply -f -
kubectl create secret generic client-ca --from-file=ca.crt=certs/client-ca.pem --dry-run=client -o yaml | kubectl apply -f -


chmod +r certs/*







# Add following to vscode debbugger launch.json

    "args": "start-dev --verbose --spi-x509cert-lookup-provider=envoy"




X.509 Authenticator Type and OAuth2 Client Credentials Grant
https://datatracker.ietf.org/doc/html/rfc8705


# Create client in admin gui
http://keycloak.127-0-0-1.nip.io:8080


In Keycloak admin console:

1. In step "General settings" fill in:
    - Client ID: x509test
2. In step "Capabilicy config" fill in:
    - Client Authentication: On
    - Select "Service accounts roles"

Enable X509 authenticator:

3. Go to "Credentials" tab
    - In "Client Authenticator" select "X509 Certificate"
    - Fill in Subject DN "cn=x509client"


https://www.keycloak.org/server/keycloak-truststore




# Make request throught proxy
http --cert certs/x509client.pem --cert-key certs/x509client-key.pem --verify certs/ca.pem --form POST https://keycloak.127-0-0-1.nip.io:8443/realms/master/protocol/openid-connect/token grant_type=client_credentials client_id=x509test







wireshark -k -i lo -Y "http" -f "tcp port 8080"



GET / HTTP/1.1
host: localhost:8443
accept-encoding: gzip, deflate, br
accept: */*
user-agent: HTTPie/3.2.2
x-forwarded-proto: https
x-request-id: 09fb2264-7d85-4c59-bdb8-654ec9a91438
x-forwarded-client-cert: Hash=ddb484787c2db45582770febb6d899606fbdd32099386d85a73a17f9ffce9627;Cert="-----BEGIN%20CERTIFICATE-----%0AMIICOTCCAd%2BgAwIBAgIIF7yk8McAWq8wCgYIKoZIzj0EAwIwFDESMBAGA1UEAxMJ%0AY2xpZW50LWNhMB4XDTI0MDMxNDEzMzUxMloXDTI1MDMxNDEzMzUxMlowMTEQMA4G%0AA1UEChMHZXhhbXBsZTEOMAwGA1UECxMFdXNlcnMxDTALBgNVBAMTBHVzZXIwggEi%0AMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQC%2FKaRbk44oilMLPjzOz6%2Fo7s7I%0AmD3HMtTX4mX9%2F3dpBdFGeu9vIKuBce0wN2VLptJyk31j1rEy1Gb%2Ba4HbMBhXjk8D%0Amnpo2rbFMJgP4joh3i0Ga6jeWDzHFuqIFXeeXHxqDwFE6DSMnc8rzZ81CEl1Xs4N%0AFDmRSiOy7CVhGlQWcjsgCxLuZfcjI%2FSrzrGpO3WgEy43y%2BhDG7zqp%2FEDJhKK8Xjs%0AYmQFS8mbWsaDUbRCiHnPHgtoy6tW0wpAeFl5UZDe2iKee3JtBIfP14hYpFn%2FlYCi%0AZH1qh8dzBVhY6SG5v8KbB69GJo1tzrCggGfeizl6eUOPFsIiI3LJ8Te9Ez1lAgMB%0AAAGjMzAxMA4GA1UdDwEB%2FwQEAwIFoDAfBgNVHSMEGDAWgBQSWkBGntMw9K%2FbixHZ%0AvvTjpqLRjDAKBggqhkjOPQQDAgNIADBFAiAEINgP%2BzOdXFSw05%2FNb6hb4OI2GUv6%0AyAmWGLIDmkFfpAIhAOvvdCrGDUpedZNAw%2FBZ2%2FGUaScWimjmu3TllKdZpSv9%0A-----END%20CERTIFICATE-----%0A";Chain="-----BEGIN%20CERTIFICATE-----%0AMIICOTCCAd%2BgAwIBAgIIF7yk8McAWq8wCgYIKoZIzj0EAwIwFDESMBAGA1UEAxMJ%0AY2xpZW50LWNhMB4XDTI0MDMxNDEzMzUxMloXDTI1MDMxNDEzMzUxMlowMTEQMA4G%0AA1UEChMHZXhhbXBsZTEOMAwGA1UECxMFdXNlcnMxDTALBgNVBAMTBHVzZXIwggEi%0AMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQC%2FKaRbk44oilMLPjzOz6%2Fo7s7I%0AmD3HMtTX4mX9%2F3dpBdFGeu9vIKuBce0wN2VLptJyk31j1rEy1Gb%2Ba4HbMBhXjk8D%0Amnpo2rbFMJgP4joh3i0Ga6jeWDzHFuqIFXeeXHxqDwFE6DSMnc8rzZ81CEl1Xs4N%0AFDmRSiOy7CVhGlQWcjsgCxLuZfcjI%2FSrzrGpO3WgEy43y%2BhDG7zqp%2FEDJhKK8Xjs%0AYmQFS8mbWsaDUbRCiHnPHgtoy6tW0wpAeFl5UZDe2iKee3JtBIfP14hYpFn%2FlYCi%0AZH1qh8dzBVhY6SG5v8KbB69GJo1tzrCggGfeizl6eUOPFsIiI3LJ8Te9Ez1lAgMB%0AAAGjMzAxMA4GA1UdDwEB%2FwQEAwIFoDAfBgNVHSMEGDAWgBQSWkBGntMw9K%2FbixHZ%0AvvTjpqLRjDAKBggqhkjOPQQDAgNIADBFAiAEINgP%2BzOdXFSw05%2FNb6hb4OI2GUv6%0AyAmWGLIDmkFfpAIhAOvvdCrGDUpedZNAw%2FBZ2%2FGUaScWimjmu3TllKdZpSv9%0A-----END%20CERTIFICATE-----%0A";Subject="CN=user,OU=users,O=example"
x-envoy-expected-rq-timeout-ms: 15000




## Build container


mvn clean install -DskipTests
cp ./quarkus/dist/target/keycloak-999.0.0-SNAPSHOT.tar.gz quarkus/container/
(cd quarkus/container/; docker build --build-arg KEYCLOAK_DIST=keycloak-999.0.0-SNAPSHOT.tar.gz -t keycloak:latest .)
kind load docker-image --name keycloak keycloak:latest



## Contour setup

kubectl apply -f https://projectcontour.io/quickstart/contour.yaml

kubectl create secret tls keycloak-external --cert=certs/keycloak-server.pem --key=certs/keycloak-server-key.pem --dry-run=client -o yaml | kubectl apply -f -
kubectl create secret generic client-ca --from-file=ca.crt=certs/client-ca.pem --dry-run=client -o yaml | kubectl apply -f -


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
      clientValidation:
        caSecret: client-ca
        optionalClientCertificate: true
        forwardClientCertificate:
          cert: true
  routes:
    - services:
        - name: keycloak
          port: 8080
EOF



http --cert certs/x509client.pem --cert-key certs/x509client-key.pem --verify certs/ca.pem --form POST https://keycloak.127-0-0-121.nip.io/realms/master/protocol/openid-connect/token grant_type=client_credentials client_id=x509test



## Istio setup


wget https://github.com/istio/istio/releases/download/1.23.1/istio-1.23.1-linux-amd64.tar.gz
tar zxvf istio-*.tar.gz && rm istio-*.tar.gz
export PATH=$PATH:$(cd istio-*/bin; pwd)
istioctl install --set profile=demo


# edit the nodePorts for http2 and https to 80 and 443
cat <<EOF | kubectl -n istio-system patch service istio-ingressgateway --patch-file /dev/stdin
spec:
  ports:
  - name: http2
    nodePort: 80
    port: 80
    protocol: TCP
    targetPort: 8080
  - name: https
    nodePort: 443
    port: 443
    protocol: TCP
    targetPort: 8443
EOF



kubectl -n istio-system create secret generic keycloak-external --from-file=tls.key=certs/keycloak-server-key.pem --from-file=tls.crt=certs/keycloak-server.pem --from-file=ca.crt=certs/client-ca.pem --dry-run=client -o yaml | kubectl apply -f -
st



cat <<EOF | kubectl apply -f -
apiVersion: networking.istio.io/v1
kind: Gateway
metadata:
  name: mygateway
spec:
  selector:
    istio: ingressgateway
  servers:
  - port:
      number: 80
      name: http
      protocol: HTTP
    hosts:
    - "*"
  - port:
      number: 443
      name: https
      protocol: HTTPS
    tls:
      #mode: SIMPLE
      mode: OPTIONAL_MUTUAL
      credentialName: keycloak-external
    hosts:
    - "*"
---
apiVersion: networking.istio.io/v1
kind: VirtualService
metadata:
  name: keycloak
spec:
  hosts:
  - keycloak.127-0-0-121.nip.io
  gateways:
  - mygateway
  http:
  - route:
    - destination:
        host: keycloak
EOF




http --cert certs/x509client.pem --cert-key certs/x509client-key.pem --verify certs/ca.pem --form POST https://keycloak.127-0-0-121.nip.io/realms/master/protocol/openid-connect/token grant_type=client_credentials client_id=x509test





## Deploy Keycloak


kubectl apply -f manifests/postgresql.yaml
kubectl apply -f manifests/keycloak-local.yaml


# Add following to command line
#   --spi-x509cert-lookup-provider=envoy
# Change image name to
#   keycloak:latest



## Keycloak config

# First configure either contour or istio

# Create client in admin gui
https://keycloak.127-0-0-121.nip.io


In Keycloak admin console:

1. In step "General settings" fill in:
    - Client ID: x509test
2. In step "Capabilicy config" fill in:
    - Client Authentication: On
    - Select "Service accounts roles"

Enable X509 authenticator:

3. Go to "Credentials" tab
    - In "Client Authenticator" select "X509 Certificate"
    - Fill in Subject DN "cn=x509client"





# https://istio.io/latest/docs/tasks/traffic-management/ingress/secure-ingress/


cat <<EOF | kubectl apply -f -
apiVersion: apps/v1
kind: Deployment
metadata:
  name: echoserver
spec:
  selector:
    matchLabels:
      app.kubernetes.io/name: echoserver
  template:
    metadata:
      labels:
        app.kubernetes.io/name: echoserver
    spec:
      containers:
      - name: echoserver
        image: gcr.io/k8s-staging-ingressconformance/echoserver:v20210922-cec7cf2
        env:
        - name: POD_NAME
          valueFrom:
            fieldRef:
              fieldPath: metadata.name
        - name: NAMESPACE
          valueFrom:
            fieldRef:
              fieldPath: metadata.namespace
        ports:
        - name: http-api
          containerPort: 3000
---
apiVersion: v1
kind: Service
metadata:
  name: echoserver
spec:
  ports:
  - name: http
    port: 80
    targetPort: http-api
  selector:
    app.kubernetes.io/name: echoserver
  type: ClusterIP
---
apiVersion: networking.istio.io/v1
kind: VirtualService
metadata:
  name: echoserver
spec:
  hosts:
  - keycloak.127-0-0-121.nip.io
  gateways:
  - mygateway
  http:
  - route:
    - destination:
        host: echoserver
        port:
          number: 80
---
apiVersion: networking.istio.io/v1
kind: Gateway
metadata:
  name: mygateway
spec:
  selector:
    istio: ingressgateway
  servers:
  - port:
      number: 80
      name: http
      protocol: HTTP
    hosts:
    - "*"
  - port:
      number: 443
      name: https
      protocol: HTTPS
    tls:
      # mode: SIMPLE
      mode: MUTUAL
      credentialName: keycloak-external
    hosts:
    - "*"
EOF
