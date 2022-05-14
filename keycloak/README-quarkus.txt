
export WORKDIR=/home/tsaarni/work/devenvs/keycloak


# build and install org.keycloak modules into maven cache
mvn -f pom.xml clean install -DskipTestsuite -DskipExamples -DskipTests

# after main codebase is built, to build the quarkus distribution
mvn clean install -DskipTests


#########################
#
# certs
#

rm -rf certs
mkdir -p certs
certyaml --destination certs configs/certs.yaml


# renew test
rm certs/* && certyaml --destination certs configs/certs.yaml
openssl s_client -connect keycloak.127-0-0-1.nip.io:8443 | openssl x509 -text -noout



#########################
#
# postgress
#

# postgres requires specific permissions for the certs
mkdir -p certs/pg
cp certs/postgres.pem certs/postgres-key.pem certs/client-ca.pem certs/pg
sudo chown 70 certs/pg/*
sudo chmod 600 certs/pg/*

# postgres requires private key in DER format
openssl pkcs8 -topk8 -inform PEM -outform DER -nocrypt -in certs/postgres-admin-key.pem -out certs/postgres-admin-key2.pem

# run postgres in docker
docker-compose up postgres
docker-compose rm postgres

# test connection
openssl s_client -starttls postgres -cert certs/keycloak.pem -key certs/keycloak-key.pem -CAfile certs/ca.pem -connect localhost:5432


#########################
#
# testing
#


java -jar quarkus/server/target/lib/quarkus-run.jar show-config      # show runtime config
java -jar quarkus/server/target/lib/quarkus-run.jar --verbose        # see exceptions in full

# enable preview features https://www.keycloak.org/server/features
java -jar quarkus/server/target/lib/quarkus-run.jar build --features=admin2

# select db (alternative to --auto-rebuild)
java -jar quarkus/server/target/lib/quarkus-run.jar build --db=dev-mem   # H2 but do not persist at ~/data/h2/
java -jar quarkus/server/target/lib/quarkus-run.jar build --db=postgres  # changes database to postgres
                                                                         # ignore ERROR: Failed to run 'build' command.

# import realm at start  https://www.keycloak.org/server/importExport

# clean realms
rm -rf ~/data/h2/



# build
mvn clean install -DskipTestsuite -DskipExamples -DskipTests

java -jar quarkus/server/target/lib/quarkus-run.jar --verbose start-dev


# run production mode with H2
java -jar quarkus/server/target/lib/quarkus-run.jar --verbose start --hostname=keycloak.127-0-0-1.nip.io --https-certificate-file=$WORKDIR/certs/keycloak-server.pem --https-certificate-key-file=$WORKDIR/certs/keycloak-server-key.pem


# Run with postgres
java -jar quarkus/server/target/lib/quarkus-run.jar --verbose start --auto-build --db=postgres --hostname=keycloak.127-0-0-1.nip.io:8443 --db-username=keycloak --db-password=keycloak --http-enabled=true --https-certificate-file=$WORKDIR/certs/keycloak-server.pem --https-certificate-key-file=$WORKDIR/certs/keycloak-server-key.pem --db-url="jdbc:postgresql://localhost:5432/keycloak?sslcert=$WORKDIR/certs/postgres-admin.pem&sslkey=$WORKDIR/certs/postgres-admin-key2.pem&sslrootcert=$WORKDIR/certs/ca.pem&sslmode=verify-full"





http://keycloak.127-0-0-1.nip.io:8080/q/dev/
http://keycloak.127-0-0-1.nip.io:8080/
https://keycloak.127-0-0-1.nip.io:8443/





TOKEN=$(http --verify false --form POST https://keycloak.127-0-0-1.nip.io:8443/realms/master/protocol/openid-connect/token username=joe password=joe grant_type=password client_id=admin-cli | jq -r .access_token)


http --verify false -v GET https://keycloak.127-0-0-1.nip.io:8443/admin/realms/master/users Authorization:"bearer $TOKEN
http --verify false -v POST https://keycloak.127-0-0-1.nip.io:8443/admin/realms/master/users Authorization:"bearer $TOKEN" username=foo






