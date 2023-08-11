

https://github.com/keycloak/keycloak/issues/10733
https://github.com/keycloak/keycloak/pull/10831


add following to envs
"KC_SPI_EVENTS_LISTENER_JBOSS_LOGGING_SUCCESS_LEVEL": "info",

or on command line 
--spi-events-listener-jboss-logging-success-level=info



mvn -Pdistribution -DskipTests clean install
mvn clean install -DskipTests -pl services       # compile only services after a change


mvn clean install -f testsuite/integration-arquillian/pom.xml -Dtest=org.keycloak.testsuite.admin.event.AdminEventTest -Dkeycloak.logging.level=debug
mvn clean install -f testsuite/integration-arquillian/pom.xml -Dtest=org.keycloak.testsuite.admin.event.AdminEventTest#createAndDeleteRealm -Dkeycloak.logging.level=debug




# failing

mvn clean install -f testsuite/integration-arquillian/pom.xml -Dtest=org.keycloak.testsuite.admin.FineGrainAdminUnitTest#testCreateRealmCreateClient -Dkeycloak.logging.level=debug
