
https://github.com/keycloak/keycloak/issues/27349
https://github.com/keycloak/keycloak/pull/27351

1. login to http://keycloak.127.0.0.1.nip.io:8080
2. click authentication / policies / otp policy
3. try google authenticator after setting
  - hash algorithm to SHA256 and SHA512
  - number of digits to 8
  - otp type to hotp
4. try to setup authenticator in http://keycloak.127.0.0.1.nip.io:8080/realms/master/account/
5. click account security / signing in / Setup Authenticator Application
6. scan the code in authenticator and try that code works



services/src/main/java/org/keycloak/authentication/otp/GoogleAuthenticatorProvider.java



mvn clean install -f testsuite/integration-arquillian/pom.xml -Dkeycloak.logging.level=debug -Dtest=org.keycloak.testsuite.admin.realm.**
mvn clean install -f testsuite/integration-arquillian/pom.xml -Dkeycloak.logging.level=debug -Dtest=org.keycloak.testsuite.admin.realm.RealmTest#testSupportedOTPApplications

mvn clean install -f testsuite/integration-arquillian/pom.xml -Dkeycloak.logging.level=debug -Dtest=org.keycloak.testsuite.actions.**
mvn clean install -f testsuite/integration-arquillian/pom.xml -Dkeycloak.logging.level=debug -Dtest=org.keycloak.testsuite.actions.AppInitiatedActionTotpSetupTest
mvn clean install -f testsuite/integration-arquillian/pom.xml -Dkeycloak.logging.level=debug -Dtest=org.keycloak.testsuite.actions.RequiredActionTotpSetupTest#setupTotpModifiedPolicy


