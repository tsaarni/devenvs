

# start new cluster
kind delete cluster --name keycloak
kind create cluster --config configs/kind-cluster-config.yaml --name keycloak

kubectl apply -f https://projectcontour.io/quickstart/contour.yaml
kubectl apply -f manifests/postgresql.yaml
kubectl apply -f manifests/keycloak-15.yaml
kubectl apply -f manifests/keycloak-18.yaml
kubectl apply -f manifests/keycloak-19.yaml
kubectl apply -f manifests/keycloak-20.yaml
kubectl apply -f manifests/keycloak-21.yaml


http://keycloak.127-0-0-121.nip.io/
http://keycloak.127-0-0-121.nip.io/auth/
http://keycloak.127-0-0-121.nip.io/realms/master/account



mkdir -p certs
certyaml -d certs configs/certs.yaml

kubectl create secret tls keycloak-external --cert=certs/keycloak-server.pem --key=certs/keycloak-server-key.pem --dry-run=client -o yaml | kubectl apply -f -
kubectl create secret tls keycloak-internal --cert=certs/keycloak-internal.pem --key=certs/keycloak-internal-key.pem --dry-run=client -o yaml | kubectl apply -f -
kubectl create secret generic internal-ca --from-file=ca.crt=certs/internal-ca.pem --dry-run=client -o yaml | kubectl apply -f -



# logs
kubectl logs statefulset/keycloak


# openldap

docker build docker/openldap/ -t localhost/openldap:latest
kind load docker-image --name keycloak localhost/openldap:latest
kubectl create configmap openldap-config --dry-run=client -o yaml --from-file=templates/database.ldif --from-file=templates/users-and-groups.ldif | kubectl apply -f -
kubectl create secret tls openldap-cert --cert=certs/ldap.pem --key=certs/ldap-key.pem --dry-run=client -o yaml | kubectl apply -f -

# patch tls secret to inject ca.crt
kubectl patch secret openldap-cert --patch-file /dev/stdin <<EOF
data:
  ca.crt: $(cat certs/client-ca.pem | base64 -w 0)
EOF

kubectl create configmap keycloak-config --dry-run -o yaml --from-file=configs/master-realm.json | kubectl apply -f -


docker build docker/keycloak/ -t localhost/keycloak:latest
kind load docker-image --name keycloak localhost/keycloak:latest


kubectl apply -f manifests

https://keycloak.127-0-0-121.nip.io/auth/


### for reference
# helm repo add codecentric https://codecentric.github.io/helm-charts
# cd helm
# helm fetch --untar codecentric/keycloak



### CERT-MANAGER ALTERNATIVE
### deploy cert-manager and certificates
# kubectl apply -f https://github.com/jetstack/cert-manager/releases/download/v1.9.1/cert-manager.yaml
# kubectl apply -f manifests/certificates.yaml
# kubectl get secret keycloak-external -o jsonpath="{..ca\.crt}" > ca.pem





kubectl exec -it keycloak-0 -- /opt/jboss/keycloak/bin/jboss-cli.sh --connect


# download and unpack dependencies
mkdir -p downloads
curl -L https://downloads.jboss.org/keycloak/10.0.0/keycloak-10.0.0.tar.gz > downloads/keycloak.tar.gz
curl -L https://repo1.maven.org/maven2/org/postgresql/postgresql/42.2.5/postgresql-42.2.5.jar > downloads/postgres-jdbc.jar

rm -rf docker/keycloak/files
mkdir -p docker/keycloak/files/opt/jboss
tar xf downloads/keycloak.tar.gz -Cdocker/keycloak/files/opt/jboss
mv docker/keycloak/files/opt/jboss/keycloak* docker/keycloak/files/opt/jboss/keycloak

mkdir -p docker/keycloak/files/opt/jboss/keycloak/modules/system/layers/base/org/postgresql/jdbc/main
cp -a downloads/postgres-jdbc.jar docker/keycloak/files/opt/jboss/keycloak/modules/system/layers/base/org/postgresql/jdbc/main
cp docker/keycloak/tools/databases/postgres/module.xml docker/keycloak/files/opt/jboss/keycloak/modules/system/layers/base/org/postgresql/jdbc/main

# build keycloak container
docker build docker/keycloak/ -t localhost/keycloak:latest
kind load docker-image --name keycloak localhost/keycloak:latest



# create distribution tarball from source
mvn -Pdistribution -pl distribution/server-dist -am -Dmaven.test.skip clean install
distribution/server-dist/target/keycloak-*.tar.gz







sudo nsenter --target $(pidof slapd) --net wireshark -f  "port 389" -k
sudo nsenter --target $(pgrep -f org.jboss.as.standalone | sed -n 1p) --net wireshark -k
sudo nsenter --target $(pgrep -f org.jboss.as.standalone | sed -n 2p) --net wireshark -k


kubectl get secret keycloakcert -o jsonpath="{..ca\.crt}" | base64 -d > ca.crt
keytool -importcert -storetype PKCS12 -keystore truststore.p12 -storepass secret -noprompt -alias ca -file ca.crt
keytool -importcert -storetype PKCS12 -keystore truststore-new.p12 -storepass secret -noprompt -alias ca -file ca.crt
kubectl create secret generic cacert --dry-run -o yaml --from-file=truststore-new.p12 | kubectl apply -f -






mvn clean install -DskipTests=true
mvn clean install -Pdistribution


kubectl create -f https://raw.githubusercontent.com/keycloak/keycloak-quickstarts/latest/kubernetes-examples/keycloak.yaml

kubectl apply -f https://projectcontour.io/quickstart/contour.yaml
kubectl apply -f manifests/postgresql.yaml
kubectl apply -f manifests/keycloak-x.yaml


docker run --rm -it -e POSTGRES_USER=keycloak -e POSTGRES_PASSWORD=keycloak -e POSTGRES_DB=keycloak -p 5432:5432 docker.io/postgres:14-alpine
http://localhost:8080/q/dev/
http://localhost:8080/













TOKEN=$(http --form POST http://keycloak.127-0-0-121.nip.io/auth/realms/master/protocol/openid-connect/token username=joe password=joe grant_type=password client_id=admin-cli | jq -r .access_token)


http -v GET http://keycloak.127-0-0-121.nip.io/auth/admin/realms/master/users Authorization:"bearer $TOKEN"
http -v POST http://keycloak.127-0-0-121.nip.io/auth/admin/realms/master/users Authorization:"bearer $TOKEN" username=foo





### Login page remains empty

check that logs say "Strict HTTPS: false"


2023-05-24 09:10:15,121 INFO  [org.keycloak.quarkus.runtime.hostname.DefaultHostnameProvider] (main) Hostname settings: Base URL: <unset>, Hostname: <request>, Strict HTTPS: false, Path: <request>, Strict BackChannel: false, Admin URL: <unset>, Admin: <request>, Port: -1, Proxied: true
