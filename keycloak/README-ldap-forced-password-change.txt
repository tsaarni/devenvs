

export WORKDIR=/home/tsaarni/work/devenvs/keycloak

certyaml --destination certs configs/certs.yaml
docker-compose rm -f
docker-compose up openldap ldap-client


sudo nsenter -n -t $(pidof slapd) wireshark -f "port 389 or port 636" -Y ldap -k -o tls.keylog_file:$WORKDIR/output/wireshark-keys.log


sshpass -p user ssh user@localhost -p 2222 -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no "echo Hello world!"
sshpass -p invalid ssh user@localhost -p 2222 -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no "echo Hello world!"

sshpass -p mustchange ssh mustchange@localhost -p 2222 -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no



mvn clean install -DskipTestsuite -DskipExamples -DskipTests

rm -rf target/kc/data/h2/keycloakdb*

# run "Import realm" under debugger
# run "Debug Quarkus" under debugger



http://localhost:8080/
http://localhost:8080/realms/master/account/




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
