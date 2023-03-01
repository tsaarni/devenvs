
# create certs
rm -rf certs
mkdir -p certs
certyaml --destination certs configs/certs.yaml

rm -f certs/truststore.p12
keytool -importcert -storetype PKCS12 -keystore certs/truststore.p12 -storepass secret -noprompt -alias ca -file certs/ca.pem

rm -f certs/admin-keystore.p12
openssl pkcs12 -export -passout pass:secret -noiter -nomaciter -in certs/ldap-admin.pem -inkey certs/ldap-admin-key.pem -out admin-keystore.p12


docker-compose rm -f
docker-compose up openldap ldap-client


sudo nsenter -n -t $(pidof slapd) wireshark -f "port 389 or port 636" -Y ldap -k -o tls.keylog_file:$HOME/work/devenvs/keycloak/output/wireshark-keys.log


sshpass -p user ssh user@localhost -p 2222 -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no "echo Hello world!"
sshpass -p invalid ssh user@localhost -p 2222 -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no "echo Hello world!"

sshpass -p mustchange ssh mustchange@localhost -p 2222 -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no



mvn clean install -DskipTestsuite -DskipExamples -DskipTests

rm -rf target/kc/data/h2/keycloakdb*


# run "Import realm" under vscode debugger
# run "Debug Quarkus" under vscode debugger


# or alternatively run
mvn -f testsuite/utils/pom.xml exec:java -Pkeycloak-server -Djavax.net.ssl.trustStore=$HOME/work/devenvs/keycloak/certs/truststore.p12 -Djavax.net.ssl.trustStorePassword=password -Djavax.net.ssl.javax.net.ssl.trustStoreType="PKCS12" -Djavax.net.ssl.keyStore=$HOME/work/devenvs/keycloak/certs/admin-keystore.p12 -Djavax.net.ssl.keyStorePassword=password -Djavax.net.ssl.javax.net.ssl.keyStoreType="PKCS12" -Dresources -Dkeycloak.migration.action=import -Dkeycloak.migration.provider=dir -Dkeycloak.migration.dir=$HOME/work/devenvs/keycloak/migrations/ldap-federation-simple/
#

http://localhost:8080/
http://localhost:8080/realms/master/account/

http://localhost:8081/auth
http://localhost:8081/auth/realms/master/account/#/



http --form POST http://localhost:8081/auth/realms/master/protocol/openid-connect/token username=mustchange password=mustchange grant_type=password client_id=account-console

###
#
# integration tests
#

# preparation
mvn clean install -DskipTests
mvn clean install -DskipTests -pl util/embedded-ldap/

# compile only changed code
mvn clean install -DskipTests -pl federation/ldap/

# run policy test case
mvn clean install -f testsuite/integration-arquillian/pom.xml -Dtest=org.keycloak.testsuite.federation.ldap.LDAPPasswordPolicyTest -Dkeycloak.logging.level=debug


# failing test case
mvn clean install -f testsuite/integration-arquillian/pom.xml -Dtest=org.keycloak.testsuite.federation.ldap.LDAPProvidersIntegrationTest -Dkeycloak.logging.level=debug
mvn clean install -f testsuite/integration-arquillian/pom.xml -Dtest=org.keycloak.testsuite.federation.ldap.LDAPUserLoginTest -Dkeycloak.logging.level=debug
mvn clean install -f testsuite/integration-arquillian/pom.xml -Dtest=org.keycloak.testsuite.federation.ldap.LDAPUserLoginTest#loginLDAPUserAuthenticationNoneEncryptionStartTLS -Dkeycloak.logging.level=debug

# capture the traffic towards embedded-ldap during test case
patch -p1 < ~/work/devenvs/keycloak/testsuite-tls-secrets-for-wireshark.patch
wireshark -i lo -d tcp.port==10389,ldap -f "port 10389" -Y ldap -k -o tls.keylog_file:/tmp/wireshark-keys.log

wireshark -i lo -d tcp.port==10389,ldap -f "port 10389" -Y ldap.bindResponse_element or ldap.bindRequest_element -k -o tls.keylog_file:/tmp/wireshark-keys.log

# and add new column with field
#   tls.record.version
# to see which bind requests were sent over TLS and which not
