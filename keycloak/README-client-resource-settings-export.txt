
Fix the client's authorization setting import error
https://github.com/keycloak/keycloak/pull/28317


commit
8fb6d43e073471ce01583760c8c8582062fad953



        if (System.getenv("EXPORT_WITH_RESOURCE_ID") != null) {










export KC_DB_URL="jdbc:h2:./quarkus/dist/target/keycloakdb;NON_KEYWORDS=VALUE;AUTO_SERVER=TRUE"
export KC_HOSTNAME="keycloak.127.0.0.1.nip.io"
export KC_BOOTSTRAP_ADMIN_USERNAME="admin"
export KC_BOOTSTRAP_ADMIN_PASSWORD="admin"

java -jar quarkus/server/target/lib/quarkus-run.jar start-dev --verbose


RESOURCE_SERVER_EXPORT_WITH_ID=true java -jar quarkus/server/target/lib/quarkus-run.jar start-dev --verbose


###############################################
#
# Setup Keycloak Authorization Server and Protected Resource
#

function get_admin_token() {
  http --form POST http://keycloak.127.0.0.1.nip.io:8080/realms/master/protocol/openid-connect/token \
    username=admin \
    password=admin \
    grant_type=password \
    client_id=admin-cli \
  | jq -r .access_token
}


### Create new realm
http -v POST http://keycloak.127.0.0.1.nip.io:8080/admin/realms \
  Authorization:"Bearer $(get_admin_token)" \
  realm=example-realm \
  enabled:=true \
  refreshTokenMaxReuse:=0


### Create a new client
http -v POST http://keycloak.127.0.0.1.nip.io:8080/admin/realms/example-realm/clients \
  Authorization:"Bearer $(get_admin_token)" \
  clientId=example-client \
  publicClient=false \
  secret=example-secret \
  directAccessGrantsEnabled=true \
  rootUrl=http://localhost:18080 \
  redirectUris:='["http://localhost:18080/*"]' \
  authorizationServicesEnabled=true \
  serviceAccountsEnabled=true

### Get client id
CLIENT_ID=$(http GET http://keycloak.127.0.0.1.nip.io:8080/admin/realms/example-realm/clients \
  Authorization:"Bearer $(get_admin_token)" \
| jq -r '.[] | select(.clientId=="example-client") | .id')
echo "CLIENT_ID=$CLIENT_ID"



### Create protected resources
http -v POST http://keycloak.127.0.0.1.nip.io:8080/admin/realms/example-realm/clients/$CLIENT_ID/authz/resource-server/resource \
  Authorization:"Bearer $(get_admin_token)" \
  name=example-resource \
  type=urn:resource-server:example-resource \
  uris:='["/"]' \
  scopes:='[{"name":"GET"}]'

http -v POST http://keycloak.127.0.0.1.nip.io:8080/admin/realms/example-realm/clients/$CLIENT_ID/authz/resource-server/resource \
  Authorization:"Bearer $(get_admin_token)" \
  name=example-resource-no-permission \
  type=urn:resource-server:example-resource \
  uris:='["/"]' \
  scopes:='[{"name":"POST"}]'

# Export the resource server settings
http -v GET http://keycloak.127-0-0-1.nip.io:8080/admin/realms/example-realm/clients/$CLIENT_ID/authz/resource-server/settings \
  Authorization:"Bearer $(get_admin_token)"


# Check that "id" and "_id" are not present when running without the "EXPORT_WITH_RESOURCE_ID" environment variable

# Check that "id" and "_id" are present when running with the "EXPORT_WITH_RESOURCE_ID" environment variable
