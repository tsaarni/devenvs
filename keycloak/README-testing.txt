

*** logging for new junit based e2e test framework

See: keycloak/test-framework/LOGGING.md


cat > .env.test <<EOF
KC_TEST_LOG_LEVEL=INFO
KC_TEST_CONSOLE_COLOR=true
KC_TEST_LOG_CATEGORY__MANAGED_KEYCLOAK__LEVEL=INFO
KC_TEST_LOG_CATEGORY__ORG_KEYCLOAK_VAULT__LEVEL=DEBUG
KC_TEST_LOG_CATEGORY__TESTINFO__LEVEL=DEBUG
KC_TEST_LOG_CATEGORY__ORG_APACHE_HTTP__LEVEL=DEBUG
EOF



mvn -f tests/pom.xml test -Dtest=SMTPConnectionVaultTest

mvn -f tests/pom.xml test -Dtest=ClientVaultTest

