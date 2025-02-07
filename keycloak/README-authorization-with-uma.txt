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

### Create new realm with refresh token revoke enabled
http -v POST http://keycloak.127.0.0.1.nip.io:8080/admin/realms \
  Authorization:"Bearer $(get_admin_token)" \
  realm=example-realm \
  enabled:=true \
  revokeRefreshToken:=true \
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

### Create a new user
http -v POST http://keycloak.127.0.0.1.nip.io:8080/admin/realms/example-realm/users \
  Authorization:"Bearer $(get_admin_token)" \
  username=joe \
  enabled:=true \
  email=joe@example.com \
  firstName=Joe \
  lastName=Average \
  emailVerified:=true \
  credentials:='[{"type":"password","value":"joe","temporary":false}]'

### Create protected resource
http -v POST http://keycloak.127.0.0.1.nip.io:8080/admin/realms/example-realm/clients/$CLIENT_ID/authz/resource-server/resource \
  Authorization:"Bearer $(get_admin_token)" \
  name=example-resource \
  type=urn:resource-server:example-resource \
  uris:='["/"]' \
  scopes:='[{"name":"GET"}]'

### Create policy where user joe is allowed to access the resource
http -v POST http://keycloak.127.0.0.1.nip.io:8080/admin/realms/example-realm/clients/$CLIENT_ID/authz/resource-server/policy/user \
  Authorization:"Bearer $(get_admin_token)" \
  name=joe-policy \
  users:='["joe"]'

### Bind policy to the resource
http -v POST http://keycloak.127.0.0.1.nip.io:8080/admin/realms/example-realm/clients/$CLIENT_ID/authz/resource-server/permission/resource \
  Authorization:"Bearer $(get_admin_token)" \
  name=example-resource-permission \
  resources:='["example-resource"]' \
  policies:='["joe-policy"]'


# Check Authorization


function callback_server {
    python3 <<EOF
import socket;from urllib.parse import urlparse,parse_qs
s=socket.socket(); s.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1); s.bind(('localhost',18080)); s.listen(1)
c,_=s.accept(); r=c.recv(1024).decode(); c.send(b'HTTP/1.1 200 OK\r\nContent-Length: 2\r\n\r\nOK'); c.close(); s.close()
print(parse_qs(urlparse(r.split()[1]).query).get('code',[''])[0])
EOF
}

AUTHORIZATION_CODE=$(callback_server)
echo "AUTHORIZATION_CODE=$AUTHORIZATION_CODE"

### Get access token for the user

google-chrome --incognito "http://keycloak.127.0.0.1.nip.io:8080/realms/example-realm/protocol/openid-connect/auth?response_type=code&client_id=example-client&redirect_uri=http://localhost:18080/foo&scope=openid"

http --form POST http://keycloak.127.0.0.1.nip.io:8080/realms/example-realm/protocol/openid-connect/token \
  grant_type=authorization_code \
  code=$AUTHORIZATION_CODE \
  client_id=example-client \
  client_secret=example-secret \
  redirect_uri=http://localhost:18080/foo \
| tee joe-token.json | jq .

### Try refreshing the token

function get_joe_refresh_token() {
  local REFRESH_TOKEN=$(jq -r .refresh_token joe-token.json)
  http --form POST http://keycloak.127.0.0.1.nip.io:8080/realms/example-realm/protocol/openid-connect/token \
    refresh_token=$REFRESH_TOKEN \
    grant_type=refresh_token \
    scope=openid \
    client_id=example-client \
    client_secret=example-secret \
  | tee joe-token.json | jq .
}

get_joe_refresh_token

### Test UMA ticket grant
http -v --form POST http://keycloak.127.0.0.1.nip.io:8080/realms/example-realm/protocol/openid-connect/token \
    grant_type=urn:ietf:params:oauth:grant-type:uma-ticket \
    claim_token=$(jq -r .id_token joe-token.json) \
    claim_token_format=http://openid.net/specs/openid-connect-core-1_0.html#IDToken \
    client_id=example-client \
    client_secret=example-secret \
    audience=example-client \
    permission=example-resource#GET

get_joe_refresh_token

### Test UMA ticket grant with response_mode=decision
http -v --form POST http://keycloak.127.0.0.1.nip.io:8080/realms/example-realm/protocol/openid-connect/token \
    grant_type=urn:ietf:params:oauth:grant-type:uma-ticket \
    claim_token=$(jq -r .id_token joe-token.json) \
    claim_token_format=http://openid.net/specs/openid-connect-core-1_0.html#IDToken \
    client_id=example-client \
    client_secret=example-secret \
    audience=example-client \
    permission=example-resource#GET \
    response_mode=decision

get_joe_refresh_token



# Delete realm
http -v DELETE http://keycloak.127.0.0.1.nip.io:8080/admin/realms/example-realm \
  Authorization:"Bearer $(get_admin_token)"


# Invalid UMA 2.0 claim_token_format value for ID token in the documentation
# https://github.com/keycloak/keycloak/issues/30778
# https://www.keycloak.org/docs/latest/authorization_services/index.html#_service_obtaining_permissions


### Test with some specific Keycloak version

docker run --rm --publish 8080:8080 --env KC_HOSTNAME=keycloak.127.0.0.1.nip.io --env KC_BOOTSTRAP_ADMIN_USERNAME=admin --env KC_BOOTSTRAP_ADMIN_PASSWORD=admin quay.io/keycloak/keycloak:26.0.7 start-dev






### Get access token for the user using password grant
http --form POST http://keycloak.127.0.0.1.nip.io:8080/realms/example-realm/protocol/openid-connect/token \
  grant_type=password \
  username=joe \
  password=joe \
  scope=openid \
  client_id=example-client \
  client_secret=example-secret \
| tee joe-token.json | jq .
