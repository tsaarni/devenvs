
https://github.com/keycloak/keycloak/issues/23729



# (see presentation for full text)
#
# Current status​
#
# - Keycloak has revocation check support only for X.509 client certificate user authentication https://www.keycloak.org/docs/latest/server_admin/#_x509
#   - Can use CRL (file path), CRL Distribution Points (download) or using OCSP (Online Certificate Status Protocol).
# - Keycloak lacks support for other use cases such as (There are other use cases besides those listed, such as SQL - check PostgreSQL server certificate)
#   - IDP brokering – check revocation status of HTTP server certificate
#   - LDAP federation – check revocation status of LDAP server certificate
#
# IDP brokering – HTTP Client​
#
# - Keycloak has internal HTTP client, which it accesses via HTTP Client SPI interface​
#   - Used in various places, not limited only to identity brokering​
# - There is built-in implementation called "default"​
#   - Depends on Apache HTTP Client​
#   - Has various global configuration parameters https://www.keycloak.org/server/outgoinghttp, potentially we could suggest CRL parameter there as well​
#   - It is possible to replace the "default" provider with a custom one
# - Caveat​s
#   - The code using HTTP Client does not pass any configuration / context to the HTTP Client implementation​
#     - All configuration needs to be provided globally as provider config – client must behave the same way for all use cases​
#     - It can be difficult to ensure nothing breaks when all HTTPS connections suddenly require CRL check​
#     - Asking for SPI clients to pass the TLS configuration would require a bigger architectural and API change in Keycloak​
#   - Default implementation currently constructs (singleton) client instance at program start – no hot-reload​
#   - The HTTP Client SPI is marked private in Keycloak – can change without notice
#
# LDAP federation – LDAP client​
#
# - The LDAP federation implementation uses LDAP client from JDK​
#   - Keycloak defines custom SSL Socket Factory and gives class reference to JDK LDAP client​
# - Configuration of LDAP federation is split into two:​
#   1. Most parameters are configured via REST API (=model, stored in SQL), each realm has separate config​
#   2. Credentials are configured globally via SPI provider parameters – all realms use same credentials​
# - CRL could be proposed to (1) LDAP federation config model or (2) provider config of SPI implementation (but parameter of which SPI?)​
#
# PoC: Custom HTTP Client SPI provider​ with CRL check
#   - See providers/httpclient-with-crl/
#
# PoC: Adding CRL to LDAP storage provider​
#   - https://github.com/Nordix/keycloak/pull/new/keystore-spi-ldap-crl
#


# create new cluster
kind delete cluster --name keycloak
kind create cluster --config configs/kind-cluster-config.yaml --name keycloak

# generate certs
mkdir -p certs
rm -f certs/*
(cd providers/httpclient-with-crl; rm *pem *p12; mvn exec:java; cp *pem *p12 ../../certs/)


kubectl create secret tls keycloak-external --cert=certs/keycloak-server.pem --key=certs/keycloak-server-key.pem --dry-run=client -o yaml | kubectl apply -f -
kubectl create secret tls keycloak-server-certs --cert=certs/idp-server.pem --key=certs/idp-server-key.pem --dry-run=client -o yaml | kubectl apply -f -

# CA and CRL for NOT revoked IDP cert
kubectl create secret generic keycloak-client-certs --from-file=certs/truststore.p12 --from-file=crl.pem=certs/crl-server-not-revoked.pem --dry-run=client -o yaml | kubectl apply -f -

# CA for CRL for revoked IDP cert
kubectl create secret generic keycloak-client-certs --from-file=certs/truststore.p12 --from-file=crl.pem=certs/crl-server-revoked.pem --dry-run=client -o yaml | kubectl apply -f -

# CA and CRL for unrelated CRL (triggers automatic download of CRL from http://shell/crl.pem)
kubectl create secret generic keycloak-client-certs --from-file=certs/truststore.p12 --from-file=crl.pem=certs/crl-unrelated.pem --dry-run=client -o yaml | kubectl apply -f -

# No CRL file at all
kubectl patch secret keycloak-client-certs -p '{"data": {"crl.pem": null}}'



sudo openssl crl -in /proc/$(pgrep -f QuarkusEntryPoint)/root/run/secrets/client-certs/crl.pem -text


# Compile HttpClient provider
(cd providers/httpclient-with-crl; mvn package)
rm providers/*.jar
cp providers/httpclient-with-crl/target/httpclient-with-crl-1.0-SNAPSHOT.jar providers


# Deploy
kubectl apply -f https://projectcontour.io/quickstart/contour.yaml
kubectl apply -f manifests/postgresql.yaml
kubectl apply -f manifests/keycloak-with-crl.yaml



# Restart
kubectl rollout restart statefulset keycloak

# Logs
kubectl logs -f statefulset/keycloak


https://keycloak.127-0-0-121.nip.io/

https://idp/realms/master/.well-known/openid-configuration




docker-compose up openldap



## Run in vscode


mvnd clean install -DskipTestsuite -DskipExamples -DskipTests

mkdir -p .vscode
cat >.vscode/launch.json <<EOF
{
    "version": "0.2.0",
    "configurations": [
        {
            "type": "java",
            "name": "Debug Quarkus",
            "request": "launch",
            "mainClass": "org.keycloak.quarkus._private.IDELauncher",
            "projectName": "keycloak-quarkus-server-app",
            "args": "start-dev --hostname=keycloak.127-0-0-1.nip.io --https-certificate-file=/home/tsaarni/work/devenvs/keycloak/certs/keycloak-server.pem --https-certificate-key-file=/home/tsaarni/work/devenvs/keycloak/certs/keycloak-server-key.pem --spi-truststore-file-file=/home/tsaarni/work/devenvs/keycloak/certs/truststore.p12 --spi-truststore-file-password=secret --spi-storage-ldap-crl-file=/home/tsaarni/work/devenvs/keycloak/certs/crl.pem",
            //"vmArgs": "-Djavax.net.debug=ssl:handshake:verbose:keymanager:trustmanager:certpath",
            //"vmArgs": "-Djava.security.debug=certpath",
            "env": {
                "KEYCLOAK_ADMIN": "admin",
                "KEYCLOAK_ADMIN_PASSWORD": "admin",
            },
        }
    ]
}
EOF

https://keycloak.127-0-0-1.nip.io:8443/admin/master/console/


(cd certs; openssl pkcs12 -password pass:secret -nodes -in truststore.p12 > ca.pem)
echo Q | openssl s_client -connect localhost:636 -crl_check_all -CAfile certs/ca.pem -CRL certs/crl-server-revoked.pem
echo Q | openssl s_client -connect localhost:636 -crl_check_all -CAfile certs/ca.pem -CRL certs/crl-server-not-revoked.pem
echo Q | openssl s_client -connect localhost:636 -crl_check_all -CAfile certs/ca.pem -CRL certs/crl-unrelated.pem


## Run locally


java -jar quarkus/server/target/lib/quarkus-run.jar --verbose start-dev --hostname=keycloak.127-0-0-1.nip.io --spi-truststore-file-file=$HOME/work/devenvs/keycloak/certs/truststore.p12 --spi-truststore-file-password=secret --spi-userstorage-ldap-crl-file=$HOME/work/devenvs/keycloak/certs/crl-server-revoked.pem^C


## Links


API
server-spi-private/src/main/java/org/keycloak/connections/httpclient/HttpClientProvider.java

Implementation
services/src/main/java/org/keycloak/connections/httpclient/DefaultHttpClientFactory.java
services/src/main/java/org/keycloak/connections/httpclient/HttpClientBuilder.java


https://www.keycloak.org/server/configuration-provider#_configuring_a_default_provider
https://www.keycloak.org/server/outgoinghttp
https://www.keycloak.org/server/keycloak-truststore
