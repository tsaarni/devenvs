

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

mvn clean install -DskipTests
(cd distribution; mvn clean install)
mvn install -f testsuite/integration-arquillian/pom.xml -Dtest=org.keycloak.testsuite.federation.ldap.LDAPPasswordPolicyForcedPasswordChange -Dkeycloak.logging.level=debug
