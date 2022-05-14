
###############################################
#
# Example application with UMA
#

https://quarkus.io/guides/security-keycloak-authorization

git clone https://github.com/quarkusio/quarkus-quickstarts.git
cd quarkus-quickstarts/security-keycloak-authorization-quickstart



import realm  quarkus-quickstarts/security-keycloak-authorization-quickstart/config/quarkus-realm.json


http://keycloak.127-0-0-121.nip.io/auth/realms/master/protocol/openid-connect/auth?client_id=security-admin-console
http://keycloak.127-0-0-121.nip.io/auth/realms/quarkus/protocol/openid-connect/auth?client_id=security-admin-console





diff --git a/security-keycloak-authorization-quickstart/src/main/resources/application.properties b/security-keycloak-authorization-quickstart/src/main/resources/application.properties
index 1c607ba9..9c545c0f 100644
--- a/security-keycloak-authorization-quickstart/src/main/resources/application.properties
+++ b/security-keycloak-authorization-quickstart/src/main/resources/application.properties
@@ -1,5 +1,5 @@
 # Configuration file
-%prod.quarkus.oidc.auth-server-url=https://localhost:8543/auth/realms/quarkus
+quarkus.oidc.auth-server-url=http://keycloak.127-0-0-121.nip.io/auth/realms/quarkus
 quarkus.oidc.client-id=backend-service
 quarkus.oidc.credentials.secret=secret
 quarkus.oidc.tls.verification=none



mvn quarkus:dev

http://localhost:8080/q/dev/io.quarkus.quarkus-oidc/provider


