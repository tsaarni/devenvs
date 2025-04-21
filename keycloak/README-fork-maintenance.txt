


Possible ideas to refactor the code to remove merge conflicts:


federation/ldap/src/main/java/org/keycloak/storage/ldap/idm/store/ldap/LDAPContextManager.java

Implement own TrustProvider class that overrides provider.getSSLSocketFactory()
It would return our custom LDAPSSLSocketFactory.

The initialization of LDAPSSLSocketFactory would need to be moved from LDAPContextManager.java to avoid changes there
Could it be done in the constructor of our custom TrustProvider class?



federation/ldap/src/main/java/org/keycloak/storage/ldap/idm/store/ldap/LDAPOperationManager.java

Same as above





quarkus/runtime/src/main/java/org/keycloak/quarkus/runtime/configuration/mappers/HttpPropertyMappers.java

Remove keystore spi use for HTTPS certificates and start to use Quarkus periodic reload functionality uplifted in Keycloak 26.2.0

https://www.keycloak.org/docs/latest/release_notes/index.html#option-to-reload-trust-and-key-material-for-the-management-interface
https://www.keycloak.org/server/management-interface#_tls_support
https://quarkus.io/guides/http-reference#reloading-the-certificates

https://github.com/keycloak/keycloak/pull/32715
https://github.com/quarkusio/quarkus/pull/38608


https-management-certificates-reload-period defaults to 1h



quarkus/runtime/src/main/java/org/keycloak/quarkus/runtime/configuration/mappers/HttpPropertyMappers.java

Remove all changes since they are not needed after moving to Quarkus periodic HTTPS reload
