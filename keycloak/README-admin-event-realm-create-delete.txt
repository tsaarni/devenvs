

https://github.com/keycloak/keycloak/issues/10733
https://github.com/keycloak/keycloak/pull/10831


mvn -Pdistribution -DskipTests clean install
mvn clean install -DskipTests -pl services       # compile only services after a change


mvn clean install -f testsuite/integration-arquillian/pom.xml -Dtest=org.keycloak.testsuite.admin.event.AdminEventTest -Dkeycloak.logging.level=debug
mvn clean install -f testsuite/integration-arquillian/pom.xml -Dtest=org.keycloak.testsuite.admin.event.AdminEventTest#createAndDeleteRealm -Dkeycloak.logging.level=debug




# failing

mvn clean install -f testsuite/integration-arquillian/pom.xml -Dtest=org.keycloak.testsuite.admin.FineGrainAdminUnitTest#testCreateRealmCreateClient -Dkeycloak.logging.level=debug
