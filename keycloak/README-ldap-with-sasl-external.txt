


# Add support for SASL EXTERNAL authentication for LDAP federation
# https://github.com/keycloak/keycloak/issues/11725

# KEYCLOAK-14055 Add SASL EXTERNAL authentication for LDAP federation
# https://github.com/keycloak/keycloak/pull/7365


# PR tests

mvnd clean install -DskipExamples -DskipTests
mvn -T4C clean install -DskipExamples -DskipTests


./mvnw test -Pauth-server-quarkus -pl testsuite/integration-arquillian/tests/base -Dkeycloak.logging.level=debug -Dtest=org.keycloak.testsuite.federation.ldap.**
./mvnw test -Pauth-server-quarkus -pl testsuite/integration-arquillian/tests/base -Dkeycloak.logging.level=debug -Dtest=org.keycloak.testsuite.admin.UserFederationLdapConnectionTest
./mvnw test -Pauth-server-quarkus -Dtest=org.keycloak.keystore.TestReloadingKeystore


###########


export WORKDIR=~/work/devenvs/keycloak


## Preparation

# Create certificates and keys for LDAP server and clients
rm -rf certs
mkdir -p certs
certyaml --destination certs configs/certs.yaml   # generate certificates and keys


# Create truststore and keystore for Keycloak and LDAP server
# trusted cert for keycloak
rm -f certs/truststore.p12
keytool -importcert -storetype PKCS12 -keystore certs/truststore.p12 -storepass secret -noprompt -alias ca -file certs/ca.pem

# ldap client cert for keycloak
openssl pkcs12 -export -passout pass:secret -noiter -nomaciter -in certs/ldap-admin.pem -inkey certs/ldap-admin-key.pem -out admin-keystore.p12

# optional: keycloak embedded ldap test server cert
openssl pkcs12 -export -passout pass:password -noiter -nomaciter -in certs/server.pem -inkey certs/server-key.pem -out server-keystore.p12




## Build and run test services
# Run LDAP server (OpenLDAP) and LDAP client (SSH and SSSD)

docker-compose up
docker-compose rm -f  # clean previous containers


## Test that LDAP is up and running
# dump configuration
docker exec keycloak-devenv-openldap-1 ldapsearch -H ldapi:/// -Y EXTERNAL -b cn=config

# dump users and groups
docker exec keycloak-devenv-openldap-1 slapcat -F /data/config

# connect externally with TLS

# password authentication
LDAPTLS_CACERT=$WORKDIR/certs/ca.pem LDAPTLS_CERT=$WORKDIR/certs/ldap-admin.pem LDAPTLS_KEY=$WORKDIR/certs/ldap-admin-key.pem ldapsearch -H ldap://ldap.127-0-0-1.nip.io -D cn=ldap-admin,ou=users,o=example -w ldap-admin -b ou=users,o=example

# StartTLS with SASL EXTERNAL
LDAPSASL_MECH=EXTERNAL LDAPTLS_CACERT=$WORKDIR/certs/ca.pem LDAPTLS_CERT=$WORKDIR/certs/ldap-admin.pem LDAPTLS_KEY=$WORKDIR/certs/ldap-admin-key.pem ldapsearch -ZZ -H ldap://ldap.127-0-0-1.nip.io -b cn=config

# LDAPS with SASL EXTERNAL
LDAPSASL_MECH=EXTERNAL LDAPTLS_CACERT=$WORKDIR/certs/ca.pem LDAPTLS_CERT=$WORKDIR/certs/ldap-admin.pem LDAPTLS_KEY=$WORKDIR/certs/ldap-admin-key.pem ldapsearch -H ldaps://ldap.127-0-0-1.nip.io -b cn=config

# List users
LDAPSASL_MECH=EXTERNAL LDAPTLS_CACERT=$WORKDIR/certs/ca.pem LDAPTLS_CERT=$WORKDIR/certs/ldap-admin.pem LDAPTLS_KEY=$WORKDIR/certs/ldap-admin-key.pem ldapsearch -H ldaps://localhost:636  -b ou=users,o=example


# dump user and groups by using `ldap-admin`
ldapsearch -D cn=ldap-admin,ou=users,o=example -w ldap-admin -b ou=users,o=example

# list user
ldapsearch -H ldaps://localhost:636  -b ou=users,o=example "(&(uid=user)(objectclass=inetOrgPerson)(objectclass=organizationalPerson))" -s one

# test bind (by changing password)
ldappasswd -ZZ -D cn=user,ou=users,o=example -w user -s user



# The client configuration is read from `ldaprc` in home directory by default.
# For parameters seem https://www.openldap.org/software/man.cgi?query=ldap.conf



## Running Keycloak

Clone keycloak repository, build and install to local maven repo at `~/.m2/repository`

```bash
mvn install -DskipTests
mvn clean install -DskipTests  # or alternatively: clean build
```

After editing part of the code, build and install only single module in
multi-module maven project (to save time)

```bash
mvn install -DskipTests -pl federation/ldap/
```

Run Keycloak with embedded undertow server

```bash
# -Dresources will trigger test server to read theme directly from themes directory
mvn -f testsuite/utils/pom.xml exec:java -Pkeycloak-server -Dkeycloak.migration.action=import -Dkeycloak.migration.provider=dir -Dkeycloak.migration.dir=$WORKDIR/migrations/ldap-federation/ -Djavax.net.ssl.trustStore=$WORKDIR/truststore.p12 -Djavax.net.ssl.trustStorePassword=password -Djavax.net.ssl.javax.net.ssl.trustStoreType="PKCS12" -Djavax.net.ssl.keyStore=$WORKDIR/admin-keystore.p12 -Djavax.net.ssl.keyStorePassword=password -Djavax.net.ssl.javax.net.ssl.keyStoreType="PKCS12" -Dresources
```

To login to keycloak using command line use `kcinit` from keycloak

```bash
./testsuite/integration-arquillian/tests/base/target/containers/keycloak-client-tools/bin/kcadm.sh config credentials --server http://localhost:8081/auth --realm master --client test-client --secret "" --config /dev/null --user user --password user
```


When starting keycloak above, the realm configuration is imported from
$WORKDIR/keycloak
(see [here](https://github.com/keycloak/keycloak-documentation/blob/master/server_admin/topics/export-import.adoc))

To re-export realm after doing configuration changes:

1. Choose Export
    - Export groups and roles: on
    - Export clients: on
2. In `realm-export.json` find `bindCredential` and fill in LDAP bind password


To run in debugger, use `mvnDebug` instead of `mvn`.

```bash
mkdir -p .vscode
cp $WORKDIR/configs/launch.json .vscode/
```

Or to create release distribution packages

```
mvn install -DskipTests -Pdistribution
ls distribution/server-dist/target/keycloak*.tar.gz
```


## Capturing LDAP traffic

To debug LDAP traffic, first get the interface name from LDAP server container
and then run wireshark.

sudo nsenter -n -t $(pidof slapd) wireshark -f "port 389 or port 636" -k -o tls.keylog_file:$WORKDIR/output/wireshark-keys.log


The LDAP server container uses openssl wrapper ([see here](docker/openldap/sslkeylog/))
that dumps TLS pre-master secrets to `output/wireshark-keys.log`.
This enables debugging LDAP over TLS.



To capture from Keycloak:

```
git clone https://github.com/neykov/extract-tls-secrets.git
cd extract-tls-secrets
mvn package
```

```diff
diff --git a/pom.xml b/pom.xml
index 1d01b52028..dea2243e56 100755
--- a/pom.xml
+++ b/pom.xml
@@ -1527,7 +1527,7 @@
                     <artifactId>maven-surefire-plugin</artifactId>
                     <configuration>
                         <forkMode>once</forkMode>
-                        <argLine>-Djava.awt.headless=true ${surefire.memory.settings}</argLine>
+                        <argLine>-Djava.awt.headless=true ${surefire.memory.settings} -javaagent:/home/tsaarni/packages/extract-tls-secrets/target/extract-tls-secrets-4.1.0-SNAPSHOT.jar=/home/tsaarni/work/keycloak-devenv/wireshark-keys.log</argLine>
                         <runOrder>alphabetical</runOrder>
                     </configuration>
                 </plugin>
```


## Using LDAP client

Run follownig to login to LDAP client container using LDAP user account

```bash
sshpass -p user ssh user@localhost -p 2222 -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no "echo Hello world!"
```


## Running unittest


# First build and install
mvn clean install -DskipTests
(cd distribution; mvn clean install)


# deploy EmbeddedLDAPServer (dependency of test suite)
mvn clean install -DskipTests -pl util/embedded-ldap/

# run with remote debugger
#  attach to 5005
mvn verify -DforkMode=never -Dmaven.surefire.debug -f testsuite/integration-arquillian/pom.xml -Dtest=org.keycloak.testsuite.federation.ldap.LDAPAnonymousBindTest -Dkeycloak.logging.level=debug

# run without debugger
mvn install -f testsuite/integration-arquillian/pom.xml -Dtest=org.keycloak.testsuite.federation.ldap.LDAPAnonymousBindTest -Dkeycloak.logging.level=debug
mvn install -f testsuite/integration-arquillian/pom.xml -Dtest=org.keycloak.testsuite.federation.ldap.*AuthTest -Dkeycloak.logging.level=debug

mvn install -f testsuite/integration-arquillian/pom.xml -Dkeycloak.logging.level=debug -Dtest="org.keycloak.testsuite.federation.ldap.LDAPUserLoginTest#loginLDAPUserAuthenticationSASLExternalEncryptionSSL"


mvn install -f testsuite/integration-arquillian/pom.xml -Dtest=org.keycloak.testsuite.federation.ldap.**
mvn install -f testsuite/integration-arquillian/pom.xml -Dtest=org.keycloak.testsuite.admin.UserFederationLdapConnectionTest

# model tests
cd testsuite/model
mvn clean install -Plegacy-jpa-federation+ldap -Dtest=org.keycloak.testsuite.model.UserModelTest#testAddDirtyRemoveFederationUser -Dkeycloak.logging.level=debug



# run with quarkus
mvn clean install -e -Pauth-server-quarkus -f testsuite/integration-arquillian/pom.xml -Dkeycloak.logging.level=debug -Dtest="org.keycloak.testsuite.federation.ldap.LDAPUserLoginTest#loginLDAPUserAuthenticationSASLExternalEncryptionSSL"

mvn clean install -e -Pauth-server-wildfly -f testsuite/integration-arquillian/pom.xml -Dkeycloak.logging.level=debug -Dtest="org.keycloak.testsuite.federation.ldap.LDAPUserLoginTest#loginLDAPUserAuthenticationSASLExternalEncryptionSSL"


to remote debug unittests on quarkus: -Dauth.server.debug=true -Dauth.server.debug.port=5005



# run only LDAPEmbeddedServer
mvn exec:java -pl util/embedded-ldap/ -Dexec.mainClass=org.keycloak.util.ldap.LDAPEmbeddedServer
mvn exec:java -pl util/embedded-ldap/ -Dexec.mainClass=org.keycloak.util.ldap.LDAPEmbeddedServer -DenableSSL=true -DenableStartTLS=true -Djavax.net.ssl.trustStore=$WORKDIR/truststore.p12 -Djavax.net.ssl.trustStorePassword=password -Djavax.net.ssl.javax.net.ssl.trustStoreType="PKCS12" -DkeystoreFile=$WORKDIR/server-keystore.p12 -DcertificatePassword=password



ldapsearch -H ldap://localhost:10389 -D uid=admin,ou=system -w secret -b ou=People,dc=keycloak,dc=org

# anonymous bind
ldapsearch -H ldap://localhost:10389 -x -b ou=People,dc=keycloak,dc=org

# ldaps
LDAPTLS_REQCERT=never ldapsearch -H ldaps://localhost:10636 -x -D uid=admin,ou=system -w secret -b ou=People,dc=keycloak,dc=org
LDAPTLS_CACERT=/path
LDAPTLS_CERT=/path LDAPTLS_KEY=/path LDAPTLS_CACERT=/path

# starttls
LDAPTLS_REQCERT=never ldapsearch -H ldap://localhost:10389 -ZZ -D uid=admin,ou=system -w secret -b ou=People,dc=keycloak,dc=org

LDAPTLS_REQCERT=never LDAPTLS_CERT=certs/ldap-admin.pem LDAPTLS_KEY=certs/ldap-admin-key.pem ldapsearch -H ldap://localhost:10389 -ZZ -Y EXTERNAL -b ou=People,dc=keycloak,dc=org

LD_PRELOAD=./libsslkeylog.so SSLKEYLOGFILE=wireshark-keys.log LDAPSASL_MECH=EXTERNAL LDAPCA_CERT=$PWD/certs/ca.pem LDAPTLS_CERT=certs/foo.pem LDAPTLS_KEY=certs/foo-key.pem ldapsearch -ZZ -H ldap://localhost:10389 -b ou=People,dc=keycloak,dc=org




### Connect to H2 database

testsuite/utils/src/main/resources/META-INF/keycloak-server.json
-            "url": "${keycloak.connectionsJpa.url:jdbc:h2:mem:test;DB_CLOSE_DELAY=-1}",
+            "url": "${keycloak.connectionsJpa.url:jdbc:h2:${jboss.server.data.dir}/test;DB_CLOSE_DELAY=-1;AUTO_SERVER=TRUE;AUTO_SERVER_PORT=9090}",
             "driver": "${keycloak.connectionsJpa.driver:org.h2.Driver}",


connect by using

java -cp /home/tsaarni/.m2/repository/com/h2database/h2/1.4.197/h2-1.4.197.jar org.h2.tools.Console -url "jdbc:h2:/tmp/keycloak-server-4036973945339407039/data/test;AUTO_SERVER=TRUE" -user SA




### Using REST API


https://github.com/keycloak/keycloak-documentation/blob/master/server_development/topics/admin-rest-api.adoc
https://www.keycloak.org/docs-api/10.0/rest-api/index.html


change "Access Token Lifespan" from 1 min to 100 days in realm settings

http://localhost:8081/auth/admin/master/console/#/realms/master/token-settings



TOKEN=$(http --form POST http://localhost:8081/auth/realms/master/protocol/openid-connect/token username=admin password=admin grant_type=password client_id=admin-cli | jq -r .access_token)

http -v GET http://localhost:8081/auth/admin/realms/master  Authorization:"bearer $TOKEN"

http -v GET http://localhost:8081/auth/admin/realms/master/users Authorization:"bearer $TOKEN"

# create user
http -v POST http://localhost:8081/auth/admin/realms/master/users Authorization:"bearer $TOKEN" username=foo

http -v POST http://localhost:8081/auth/admin/realms/master/users Authorization:"bearer $TOKEN" username=user3 enabled:=true totp:=false emailVerified:=false firstName="" lastName="" email="" credentials:='[{"type": "password", "value": "mypass", "temporary": false}]'

http -v POST http://localhost:8081/auth/admin/realms/master/users Authorization:"bearer $TOKEN" username=ldapuser enabled:=true firstName=Ldap lastName=User attributes:='{"telephoneNumber": ["1", "2", "3"]}'


http -v http://localhost:8081/auth/realms/master/.well-known/openid-configuration
http -v http://localhost:8081/auth/realms/master/protocol/openid-connect/certs


# enable failed login attempt detection
#   Realm settings / Security Defenses / Brute Force Detection
#
# - create user joe
# - do failed login attempt
http --form POST http://localhost:8081/auth/realms/master/protocol/openid-connect/token username=joe password=wrong grant_type=password client_id=test-client

# correct
http --form POST http://localhost:8081/auth/realms/master/protocol/openid-connect/token username=joe password=joe grant_type=password client_id=test-client


# get user
http -v GET "http://localhost:8081/auth/admin/realms/master/users?username=joe" Authorization:"bearer $TOKEN"
http -v GET http://localhost:8081/auth/admin/realms/master/users/c3240bbe-c996-465e-a7d5-e4870f34aebc Authorization:"bearer $TOKEN"













# undertow
TOKEN=$(http --form POST http://localhost:8081/auth/realms/master/protocol/openid-connect/token username=admin password=admin grant_type=password client_id=admin-cli | jq -r .access_token)
http -v POST http://localhost:8081/auth/admin/realms/master/components Authorization:"bearer $TOKEN" < rest-requests/create-ldap-starttls-provider.json
http -v "http://localhost:8081/auth/admin/realms/master/components?parent=master&type=org.keycloak.storage.UserStorageProvider" Authorization:"bearer $TOKEN"



# quarkus
TOKEN=$(http --form POST http://localhost:8080/realms/master/protocol/openid-connect/token username=admin password=admin grant_type=password client_id=admin-cli | jq -r .access_token)
http -v POST http://localhost:8080/admin/realms/master/components Authorization:"bearer $TOKEN" < rest-requests/create-ldap-starttls-provider.json
http -v POST http://localhost:8080/admin/realms/master/components Authorization:"bearer $TOKEN" < rest-requests/create-ldaps-provider.json
http -v "http://localhost:8080/admin/realms/master/components?parent=master&type=org.keycloak.storage.UserStorageProvider" Authorization:"bearer $TOKEN"






#
# Test LDAP client certificate rotation
#

docker-compose rm -f  # clean previous containers
docker-compose up

# run wireshark to check the client certificate from keycloak
sudo nsenter -n -t $(pidof slapd) wireshark -f "port 389 or port 636" -k -o tls.keylog_file:$WORKDIR/output/wireshark-keys.log



# generate certificates and keystores
rm -rf certs
mkdir -p certs
certyaml --destination certs configs/certs.yaml
keytool -importcert -storetype PKCS12 -keystore certs/truststore.p12 -storepass secret -noprompt -alias ca -file certs/ca.pem

openssl pkcs12 -export -passout pass:secret -noiter -nomaciter -in certs/ldap-admin.pem -inkey certs/ldap-admin-key.pem -out certs/ldap-admin.p12



# run keycloak with keystore and truststore SPI

rm ~/data/h2/*  # remove database

export KEYCLOAK_ADMIN=admin
export KEYCLOAK_ADMIN_PASSWORD=admin

java -jar quarkus/server/target/lib/quarkus-run.jar --verbose start --hostname-strict-https=false --http-enabled=true --hostname=keycloak.127-0-0-1.nip.io --spi-keystore-default-ldap-keystore-file=$WORKDIR/certs/ldap-admin.p12 --spi-keystore-default-ldap-keystore-password=secret --spi-truststore-file-file=$WORKDIR/certs/truststore.p12 --spi-truststore-file-password=secret

# create provider for LDAP federation
TOKEN=$(http --form POST http://keycloak.127-0-0-1.nip.io:8080/realms/master/protocol/openid-connect/token username=admin password=admin grant_type=password client_id=admin-cli | jq -r .access_token)

http -v POST http://keycloak.127-0-0-1.nip.io:8080/admin/realms/master/components Authorization:"bearer $TOKEN" < rest-requests/create-ldap-starttls-provider.json

# observe certificate details on wireshark

# test LDAP authentication
http -v POST http://keycloak.127-0-0-1.nip.io:8080/admin/realms/master/testLDAPConnection  Authorization:"bearer $TOKEN" < rest-requests/test-ldap-authentication.json

# create user
http -v POST http://keycloak.127-0-0-1.nip.io:8080/admin/realms/master/users Authorization:"bearer $TOKEN" username=user2 enabled:=true totp:=false emailVerified:=false firstName="" lastName="" email="" credentials:='[{"type": "password", "value": "mypass", "temporary": false}]'

# check that user was created
ldapsearch -H ldap://localhost:389 -D cn=ldap-admin,ou=users,o=example -w ldap-admin  -b ou=users,o=example

# regenerate ldap-admin certificate and keystore
rm certs/ldap-admin*
certyaml --destination certs configs/certs.yaml
openssl pkcs12 -export -passout pass:secret -noiter -nomaciter -in certs/ldap-admin.pem -inkey certs/ldap-admin-key.pem -out certs/ldap-admin.p12


# do above again and expiration date on client certificate
















