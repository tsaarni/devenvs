


https://github.com/keycloak/keycloak-documentation/blob/main/server_development/topics/admin-rest-api.adoc
https://www.keycloak.org/docs-api/17.0/rest-api/index.html


change "Access Token Lifespan" from 1 min to 100 days in realm settings

http://localhost:8081/auth/admin/master/console/#/realms/master/token-settings




function get_admin_token() {
  http --form POST http://keycloak.127-0-0-1.nip.io:8080/realms/master/protocol/openid-connect/token username=admin password=admin grant_type=password client_id=admin-cli | jq -r .access_token
}



http -v http://keycloak.127-0-0-1.nip.io:8080/realms/master/.well-known/openid-configuration
http -v http://keycloak.127-0-0-1.nip.io:8080/realms/master/protocol/openid-connect/certs


# get realm
http -v GET http://keycloak.127-0-0-1.nip.io:8080/admin/realms/master  Authorization:"bearer $(get_admin_token)"

# create realm
http -v POST http://keycloak.127-0-0-1.nip.io:8080/admin/realms/ Authorization:"bearer $(get_admin_token)" id=my-realm realm=my-realm adminEventsEnabled=true

# delete realm
http -v DELETE http://keycloak.127-0-0-1.nip.io:8080/admin/realms/my-realm Authorization:"bearer $(get_admin_token)"

# get users
http -v GET http://keycloak.127-0-0-1.nip.io:8080/admin/realms/master/users Authorization:"bearer $(get_admin_token)"

# create user
http -v POST http://keycloak.127-0-0-1.nip.io:8080/admin/realms/master/users Authorization:"bearer $(get_admin_token)" username=joe enabled:=true credentials:='[{"type": "password", "value": "joe", "temporary": false}]'

http -v POST http://keycloak.127-0-0-1.nip.io:8080/admin/realms/master/users Authorization:"bearer $(get_admin_token)" username=user3 enabled:=true totp:=false emailVerified:=false firstName="" lastName="" email="" credentials:='[{"type": "password", "value": "mypass", "temporary": false}]'

http -v POST http://keycloak.127-0-0-1.nip.io:8080/admin/realms/master/users Authorization:"bearer $(get_admin_token)" username=ldapuser enabled:=true firstName=Ldap lastName=User attributes:='{"telephoneNumber": ["1", "2", "3"]}'


# enable failed login attempt detection
#   Realm settings / Security Defenses / Brute Force Detection
#
# - create user joe
# - do failed login attempt

http -v POST http://keycloak.127-0-0-1.nip.io:8080/admin/realms/ Authorization:"bearer $(get_admin_token)" id=brute-force-realm realm=brute-force-realm enabled=true bruteForceProtected=true
http -v POST http://keycloak.127-0-0-1.nip.io:8080/admin/realms/brute-force-realm/users Authorization:"bearer $(get_admin_token)" username=joe enabled:=true credentials:='[{"type": "password", "value": "joe", "temporary": false}]'

http --form POST http://keycloak.127-0-0-1.nip.io:8080/realms/brute-force-realm/protocol/openid-connect/token username=joe password=wrong grant_type=password client_id=test-client

# correct
http --form POST http://keycloak.127-0-0-1.nip.io:8080/realms/brute-force-realm/protocol/openid-connect/token username=joe password=joe grant_type=password client_id=test-client




# get user
http -v GET "http://keycloak.127-0-0-1.nip.io:8080/admin/realms/master/users?username=joe" Authorization:"bearer $(get_admin_token)"

id=$(http GET "http://keycloak.127-0-0-1.nip.io:8080/admin/realms/master/users?username=joe" Authorization:"bearer $(get_admin_token)" | jq -r '.[0].id')
http -v GET http://keycloak.127-0-0-1.nip.io:8080/admin/realms/master/users/$id Authorization:"bearer $(get_admin_token)"


# list user (lists IDs)
http -v GET "http://keycloak.127-0-0-1.nip.io:8080/admin/realms/master/users" Authorization:"bearer $(get_admin_token)"

# delete user by id
http -v DELETE http://keycloak.127-0-0-1.nip.io:8080/admin/realms/master/users/$id Authorization:"bearer $(get_admin_token)"



# get admin events ("save admin events" must be enabled first)
http -v GET http://keycloak.127-0-0-1.nip.io:8080/admin/realms/master/admin-events Authorization:"bearer $(get_admin_token)"


#####################################
#
# create confidential client
#

# Create realm
http -v POST http://keycloak.127-0-0-1.nip.io:8080/admin/realms/ Authorization:"bearer $(get_admin_token)" id=my-realm realm=my-realm adminEventsEnabled=true enabled=true

# Create client
http POST http://keycloak.127-0-0-121.nip.io:8080/admin/realms/my-realm/clients \
  Authorization:"bearer $(get_admin_token)" \
  clientId=foo \
  publicClient=false \
  redirectUris[]=http://localhost \
  serviceAccountsEnabled=true \
  secret=mysecret


# Create client
http POST http://keycloak.127-0-0-121.nip.io:8080/admin/realms/my-realm/clients \
  Authorization:"bearer $(get_admin_token)" \
  clientId=bar \
  publicClient=false \
  redirectUris[]=http://localhost \
  serviceAccountsEnabled=true \
  secret="\${vault.mysecret}"


# Get client token using client credentials (requires serviceAccountsEnabled=true for the client)
http --form POST http://keycloak.127-0-0-121.nip.io:8080/realms/my-realm/protocol/openid-connect/token grant_type=client_credentials client_id=foo client_secret=mysecret

# List all clients in the realm
http GET http://keycloak.127-0-0-121.nip.io:8080/admin/realms/my-realm/clients Authorization:"bearer $(get_admin_token)"

# Get client configuration by name or by id
http GET "http://keycloak.127-0-0-121.nip.io:8080/admin/realms/my-realm/clients?clientId=foo" Authorization:"bearer $(get_admin_token)"

http GET http://keycloak.127-0-0-121.nip.io:8080/admin/realms/my-realm/clients/2f214329-2fc7-4cac-b79d-3907138f9887 Authorization:"bearer $(get_admin_token)"

# list client secret
http GET http://keycloak.127-0-0-121.nip.io:8080/admin/realms/my-realm/clients/2f214329-2fc7-4cac-b79d-3907138f9887/client-secret Authorization:"bearer $(get_admin_token)"


# Delete client
http -v DELETE http://keycloak.127-0-0-121.nip.io:8080/admin/realms/my-realm/clients/2f214329-2fc7-4cac-b79d-3907138f9887 Authorization:"bearer $(get_admin_token)"



# Delete realm
http -v DELETE http://keycloak.127-0-0-1.nip.io:8080/admin/realms/my-realm Authorization:"bearer $(get_admin_token)"


###########################################
#
# Create LDAP federation
#

# Create realm
http -v POST http://keycloak.127-0-0-1.nip.io:8080/admin/realms/ Authorization:"bearer $(get_admin_token)" id=my-realm realm=my-realm adminEventsEnabled=true enabled=true

http -v POST http://keycloak.127-0-0-121.nip.io:8080/admin/realms/my-realm/components \
  Authorization:"bearer $(get_admin_token)" \
  id=my-ldap \
  name="my ldap" \
  providerId=ldap \
  providerType="org.keycloak.storage.UserStorageProvider" \
  config[connectionUrl][]=ldap://ldap.example.com \
  config[usersDn][]=ou=scientists,dc=example,dc=com \
  config[bindDn][]=cn=admin,dc=example,dc=com \
  config[bindCredential][]=mypassword \
  config[editMode][]=READ_ONLY

# Get LDAP configuration
http -v GET http://keycloak.127-0-0-121.nip.io:8080/admin/realms/my-realm/components/my-ldap Authorization:"bearer $(get_admin_token)"

# Delete realm
http -v DELETE http://keycloak.127-0-0-1.nip.io:8080/admin/realms/my-realm Authorization:"bearer $(get_admin_token)"


###########################################
#
# Configure IDP brokering
#

# Create realm
http -v POST http://keycloak.127-0-0-1.nip.io:8080/admin/realms/ Authorization:"bearer $(get_admin_token)" id=my-realm realm=my-realm adminEventsEnabled=true enabled=true

http -v POST http://keycloak.127-0-0-121.nip.io:8080/admin/realms/my-realm/identity-provider/instances \
  Authorization:"bearer $(get_admin_token)" \
  alias=oidc-keycloak \
  providerId=oidc \
  config[clientId]=my-client-id \
  config[clientSecret]=my-secret \
  config[authorizationUrl]=https://another-keycloak:8443/realms/other-realm/protocol/openid-connect/auth \
  config[tokenUrl]=https://another-keycloak:8443/realms/other-realm/protocol/openid-connect/token \
  config[userInfoUrl]=https://another-keycloak:8443/realms/other-realm/protocol/openid-connect/userinfo \
  config[jwksUrl]=https://another-keycloak:8443/realms/other-realm/protocol/openid-connect/certs \
  config[logoutUrl]=https://another-keycloak:8443/realms/other-realm/protocol/openid-connect/logout \
  config[issuer]=https://another-keycloak:8443/realms/other-realm \
  config[redirectUri]=https://keycloak.127-0-0-121.nip.io:8443/auth/realms/my-realm/broker/oidc-keycloak/endpoint

# Get IDP configuration
http -v GET http://keycloak.127-0-0-121.nip.io:8080/admin/realms/my-realm/identity-provider/instances/oidc-keycloak Authorization:"bearer $(get_admin_token)"

# Delete realm
http -v DELETE http://keycloak.127-0-0-1.nip.io:8080/admin/realms/my-realm Authorization:"bearer $(get_admin_token)"






###########################################
#
# export configuration for a realm
#

http -v POST http://keycloak.127-0-0-121.nip.io:8080/admin/realms/my-realm/partial-export Authorization:"bearer $(get_admin_token)"


##################################
#
# Authorization code flow
#

Add http://localhost:8000 to the "Valid Redirection URIs" field in "security-admin-console" client


# generate PKCE code_challenge and code_verifier
https://referbruv.com/utilities/pkce-generator-online


# run web server
apps/http-server.py


# authenticate
google-chrome -incognito 'http://keycloak.127-0-0-121.nip.io/auth/realms/master/protocol/openid-connect/auth?client_id=security-admin-console&redirect_uri=http://localhost:8000&response_mode=query&response_type=code&scope=openid&code_challenge=fnnumE0tZuqJLPYuqbSdRaPfPGl-RekkHyfJTciE69I&code_challenge_method=S256'

# check the code=NNN from web server logs and add that to token request
http -v --form post http://keycloak.127-0-0-121.nip.io/auth/realms/master/protocol/openid-connect/token grant_type=authorization_code client_id=security-admin-console redirect_uri=http://localhost:8000 code_verifier=Aw802ZCahS-nXenJgzcU5S1aCGdlbNED_zpyiDe1Y0g code=NNN


# paste access token to see the claims
apps/jwt-decode.py


# Use refresh token flow to get new access token
http -v --form post http://keycloak.127-0-0-121.nip.io/auth/realms/master/protocol/openid-connect/token grant_type=refresh_token client_id=security-admin-console refresh_token=NNN

# Use Keycloak's whoami endpoint to get user permissions
http -v http://keycloak.127-0-0-121.nip.io/auth/admin/master/console/whoami Authorization:"bearer $TOKEN"
