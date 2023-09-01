

https://github.com/keycloak/keycloak/issues/14523
https://github.com/keycloak/keycloak/pull/15253

https://datatracker.ietf.org/doc/html/draft-behera-ldap-password-policy-11
https://datatracker.ietf.org/doc/html/draft-behera-ldap-password-policy-11#section-9.1

   *  bindResponse.resultCode = success (0),
      passwordPolicyResponse.error = changeAfterReset (2): The user is
      binding for the first time after the password administrator set
      the password.  In this scenario, the client SHOULD prompt the user
      to change his password immediately.




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



# Start with clean database
rm -rf target/kc/data/h2/keycloakdb*

mkdir -p .vscode
cp ~/work/devenvs/keycloak/configs/launch.json .vscode/launch.json

# run "Debug Quarkus" under vscode debugger


# or alternatively run
mvn -f testsuite/utils/pom.xml exec:java -Pkeycloak-server -Djavax.net.ssl.trustStore=$HOME/work/devenvs/keycloak/certs/truststore.p12 -Djavax.net.ssl.trustStorePassword=password -Djavax.net.ssl.javax.net.ssl.trustStoreType="PKCS12" -Djavax.net.ssl.keyStore=$HOME/work/devenvs/keycloak/certs/admin-keystore.p12 -Djavax.net.ssl.keyStorePassword=password -Djavax.net.ssl.javax.net.ssl.keyStoreType="PKCS12" -Dresources -Dkeycloak.migration.action=import -Dkeycloak.migration.provider=dir -Dkeycloak.migration.dir=$HOME/work/devenvs/keycloak/migrations/ldap-federation-simple/


# create ldap simple auth provider
apps/create-components.py rest-requests/create-ldap-simple-auth-provider.json
apps/create-components.py --server=https://keycloak.127-0-0-121.nip.io/ rest-requests/create-ldap-simple-auth-provider.json



# l: mustchange  p: mustchange

http://localhost:8080/
http://localhost:8080/realms/master/account/

http://localhost:8081/auth
http://localhost:8081/auth/realms/master/account/#/


# Automated script that logs in to check if password change is required
python3 -m venv .venv
source .venv/bin/activate
pip install selenium webdrivermanager
webdrivermanager chrome

apps/login-with-forced-password-change.py


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






###
[keycloak-dev] Support for password-only sync in user federation
https://lists.jboss.org/pipermail/keycloak-dev/2018-September/011258.html



### TODO
check how MSAD avoid "recursing" when adding required action i.e. writing it back to LDAP





###
apps/create-components.py rest-requests/create-ldap-simple-auth-provider.json
