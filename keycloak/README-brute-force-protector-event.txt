


# PR that added the event for permanent lockout
#
# KEYCLOAK-15985 Add Brute Force Detection Lockout Event #8679
# https://github.com/keycloak/keycloak/pull/8679


Add following to debugger environment variables to see Events in the console

"KC_SPI_EVENTS_LISTENER_JBOSS_LOGGING_SUCCESS_LEVEL": "INFO",





# login with wrong password
http --form POST http://keycloak.127.0.0.1.nip.io:8080/realms/master/protocol/openid-connect/token username=joe password=invalid grant_type=password client_id=admin-cli

# login with correct password
http --form POST http://keycloak.127.0.0.1.nip.io:8080/realms/master/protocol/openid-connect/token username=joe password=joe grant_type=password client_id=admin-cli






mvn clean install -f testsuite/integration-arquillian/pom.xml -Dkeycloak.logging.level=debug -Dtest=org.keycloak.testsuite.forms.BruteForceTest#testBrowserInvalidPassword


mvn clean install -f testsuite/integration-arquillian/pom.xml -Dkeycloak.logging.level=debug -Dtest=org.keycloak.testsuite.forms.BruteForceTest#testPermanentLockout
