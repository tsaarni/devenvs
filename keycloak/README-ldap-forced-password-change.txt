

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


# l: mustchange  p: mustchange

http://keycloak.127.0.0.1.nip.io:8080/
http://keycloak.127.0.0.1.nip.io:8080/realms/master/account/

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

# ldap federation tests
mvn clean install -DskipTests
mvn install -Pauth-server-quarkus -f testsuite/integration-arquillian/pom.xml -Dtest=org.keycloak.testsuite.federation.ldap.** -Dkeycloak.logging.level=debug
mvn install -Pauth-server-quarkus -f testsuite/integration-arquillian/pom.xml -Dtest=org.keycloak.testsuite.admin.UserFederationLdapConnectionTest

# run policy test case only
mvn clean install -Pauth-server-quarkus -f testsuite/integration-arquillian/pom.xml -Dtest=org.keycloak.testsuite.federation.ldap.LDAPPasswordPolicyTest -Dkeycloak.logging.level=debug


### Note:  if you get failures like following, you are missing parameter -Pauth-server-quarkus
##
## [INFO] Results:
## [INFO]
## [ERROR] Failures:
## [ERROR]   LDAPProvidersIntegrationTest.testSyncRegistrationEmailRDNNoDefault:218 expected:<400> but was:<500>
## [ERROR]   LDAPProvidersIntegrationNoImportTest>LDAPProvidersIntegrationTest.testSyncRegistrationEmailRDNNoDefault:218 expected:<400> but was:<500>
## [ERROR] Errors:
## [ERROR]   LDAPAdminRestApiTest.testErrorResponseWhenLdapIsFailing:258 » Processing com.fasterxml.jackson.core.JsonParseException: Unrecognized token 'Not': was expecting (JSON String, Number, Array, Object or token 'null', 'true' or 'false')
 at [Source: REDACTED (`StreamReadFeature.INCLUDE_SOURCE_IN_LOCATION` disabled); line: 1, column: 5]
## [INFO]
## [ERROR] Tests run: 231, Failures: 2, Errors: 1, Skipped: 9
## [INFO]




mvn clean install -f testsuite/integration-arquillian/pom.xml -Pauth-server-quarkus -Dtest=org.keycloak.testsuite.federation.ldap.** -Dkeycloak.logging.level=debug

mvn clean install -f testsuite/integration-arquillian/pom.xml -Pauth-server-quarkus -Dtest=org.keycloak.testsuite.federation.ldap.LDAPProvidersIntegrationTest#testSyncRegistrationEmailRDNNoDefault -Dkeycloak.logging.level=debug

mvn clean install -f testsuite/integration-arquillian/pom.xml -Pauth-server-quarkus -Dtest=org.keycloak.testsuite.admin.UserFederationLdapConnectionTest -Dkeycloak.logging.level=debug

mvn clean install -f testsuite/integration-arquillian/pom.xml  -Dtest=org.keycloak.testsuite.admin.UserFederationLdapConnectionTest#testLdapConnectionMoreServers -Dkeycloak.logging.level=debug


# capture the traffic towards embedded-ldap during test case
patch -p1 < ~/work/devenvs/keycloak/testsuite-tls-secrets-for-wireshark.patch
wireshark -i lo -d tcp.port==10389,ldap -d tcp.port==10636,ldap -f "port 10389 or port 10636" -Y ldap -k -o tls.keylog_file:/tmp/wireshark-keys.log

wireshark -i lo -d tcp.port==10389,ldap -d tcp.port==10636,ldap -f "port 10389 or port 10636" -Y ldap.bindResponse_element or ldap.bindRequest_element -k -o tls.keylog_file:/tmp/wireshark-keys.log

# and add new column with field
#   tls.record.version
# to see which bind requests were sent over TLS and which not






###
[keycloak-dev] Support for password-only sync in user federation
https://lists.jboss.org/pipermail/keycloak-dev/2018-September/011258.html



### TODO
check how MSAD avoid "recursing" when adding required action i.e. writing it back to LDAP





###


apps/create-components.py --server=https://keycloak.127-0-0-121.nip.io/ rest-requests/create-ldap-simple-auth-provider.json
apps/create-components.py --server=https://keycloak.127-0-0-121.nip.io:8080/ rest-requests/create-ldap-simple-auth-provider.json




## Manual test


1. Start OpenLDAP and LDAP client containers

cd ~/work/devenvs/keycloak
docker-compose rm -f
docker-compose up openldap ldap-client


2. Capture traffic to/from OpenLDAP

sudo nsenter --target $(pidof slapd) --net wireshark -i any -k -o tls.keylog_file:output/wireshark-keys.log -Y "ldap"


3. Run Keycloak

mkdir -p .vscode
cp -a ~/work/devenvs/keycloak/configs/launch.json  ~/work/devenvs/keycloak/configs/settings.json .vscode
code .


4. Create LDAP federation provider

cd ~/work/devenvs/keycloak
apps/create-components.py --server=http://127.0.0.1:8080/ rest-requests/create-ldap-ppolicy-provider.json


5. Check that the provider is created

http://keycloak.127.0.0.1.nip.io:8080/admin/master/console/#/master/user-federation


6. Try forced password change (in chrome with dev profile)

http://keycloak.127.0.0.1.nip.io:8080/realms/master/account/

   1. l: mustchange  p: mustchange
   2. Keycloak will prompt "Update password"
   3. Enter new password and submit:  foobarbaz1
   4. Should redirect to account console
   5. Sign out
   6. Sign in again with new password: l: mustchange  p: foobarbaz1
   7. Should redirect to account console


7. Clear Keycloak database

rm -rf ./quarkus/dist/target/keycloakdb*






Compare to SSH login with forced password change

ssh mustchange@localhost -p 2222 -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no
p: mustchange




# To get ppolicy control also to password modify requests

index 22f2af5564..6d539126f2 100644
--- a/federation/ldap/src/main/java/org/keycloak/storage/ldap/idm/store/ldap/LDAPContextManager.java
+++ b/federation/ldap/src/main/java/org/keycloak/storage/ldap/idm/store/ldap/LDAPContextManager.java
@@ -4,6 +4,7 @@ import org.jboss.logging.Logger;
 import org.keycloak.models.KeycloakSession;
 import org.keycloak.models.LDAPConstants;
 import org.keycloak.storage.ldap.LDAPConfig;
+import org.keycloak.storage.ldap.idm.store.ldap.control.PasswordPolicyControlFactory;
 import org.keycloak.tracing.TracingProvider;
 import org.keycloak.truststore.TruststoreProvider;
 import org.keycloak.vault.VaultStringSecret;
@@ -157,6 +158,7 @@ public final class LDAPContextManager implements AutoCloseable {
         HashMap<String, Object> env = new HashMap<>();

         env.put(Context.INITIAL_CONTEXT_FACTORY, ldapConfig.getFactoryName());
+        env.put(LdapContext.CONTROL_FACTORIES, PasswordPolicyControlFactory.class.getName());

         String url = ldapConfig.getConnectionUrl();

diff --git a/federation/ldap/src/main/java/org/keycloak/storage/ldap/idm/store/ldap/LDAPOperationManager.java b/federation/ldap/src/main/java/org/keycloak/storage/ldap/idm/store/ldap/LDAPOperationManager.java
index 043cbe4460..1b94a7941e 100644
--- a/federation/ldap/src/main/java/org/keycloak/storage/ldap/idm/store/ldap/LDAPOperationManager.java
+++ b/federation/ldap/src/main/java/org/keycloak/storage/ldap/idm/store/ldap/LDAPOperationManager.java
@@ -735,6 +735,7 @@ public class LDAPOperationManager {
         try {
             execute(context -> {
                 PasswordModifyRequest modifyRequest = new PasswordModifyRequest(dn.toString(), null, password);
+                context.setRequestControls(getControls());
                 return context.extendedOperation(modifyRequest);
             }, decorator);
         } catch (NamingException e) {
