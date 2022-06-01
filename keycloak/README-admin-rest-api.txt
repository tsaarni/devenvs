


https://github.com/keycloak/keycloak-documentation/blob/main/server_development/topics/admin-rest-api.adoc
https://www.keycloak.org/docs-api/17.0/rest-api/index.html


change "Access Token Lifespan" from 1 min to 100 days in realm settings

http://localhost:8081/auth/admin/master/console/#/realms/master/token-settings



TOKEN=$(http --form POST http://localhost:8081/auth/realms/master/protocol/openid-connect/token username=admin password=admin grant_type=password client_id=admin-cli | jq -r .access_token)


# get realm
http -v GET http://localhost:8081/auth/admin/realms/master  Authorization:"bearer $TOKEN"

# create realm
http -v POST http://localhost:8081/auth/admin/realms/ Authorization:"bearer $TOKEN" id=my-realm realm=my-realm adminEventsEnabled=true

# delete realm
http -v DELETE http://localhost:8081/auth/admin/realms/my-realm Authorization:"bearer $TOKEN"

# get users
http -v GET http://localhost:8081/auth/admin/realms/master/users Authorization:"bearer $TOKEN"

# create user
http -v POST http://localhost:8081/auth/admin/realms/master/users Authorization:"bearer $TOKEN" username=foo

http -v POST http://localhost:8081/auth/admin/realms/master/users Authorization:"bearer $TOKEN" username=user3 enabled:=true totp:=false emailVerified:=false firstName="" lastName="" email="" credentials:='[{"type": "password", "value": "mypass", "temporary": false}]'

http -v POST http://localhost:8081/auth/admin/realms/master/users Authorization:"bearer $TOKEN" username=ldapuser enabled:=true firstName=Ldap lastName=User attributes:='{"telephoneNumber": ["1", "2", "3"]}'

http -v http://localhost:8081/auth/realms/master/.well-known/openid-configuration
http -v http://localhost:8081/auth/realms/master/protocol/openid-connect/certs


# enable failed login attempt detection
#   Realm settings / Security Defenses / Brute Force Detection
#
# - create user joe
# - do failed login attempt
http --form POST http://localhost:8081/auth/realms/master/protocol/openid-connect/token username=joe password=wrong grant_type=password client_id=test-client

# correct
http --form POST http://localhost:8081/auth/realms/master/protocol/openid-connect/token username=joe password=joe grant_type=password client_id=test-client


# get user
http -v GET "http://localhost:8081/auth/admin/realms/master/users?username=joe" Authorization:"bearer $TOKEN"
http -v GET http://localhost:8081/auth/admin/realms/master/users/c3240bbe-c996-465e-a7d5-e4870f34aebc Authorization:"bearer $TOKEN"


# list user (lists IDs)
http -v GET "http://localhost:8081/auth/admin/realms/master/users" Authorization:"bearer $TOKEN"

# delete user by ide
http -v DELETE http://localhost:8081/auth/admin/realms/master/users/24b39570-bc5b-46a3-a5eb-3d9dc9feb561 Authorization:"bearer $TOKEN"



# get admin events ("save admin events" must be enabled first)
http -v GET http://localhost:8081/auth/admin/realms/master/admin-events Authorization:"bearer $TOKEN"


#####################################
#
# create confidential client
#

# get admin token
TOKEN=$(http --form POST http://keycloak.127-0-0-121.nip.io/auth/realms/master/protocol/openid-connect/token username=admin password=admin grant_type=password client_id=admin-cli | jq -r .access_token)

# create client
http POST http://keycloak.127-0-0-121.nip.io/auth/admin/realms/master/clients Authorization:"bearer $TOKEN" clientId=foo publicClient=false redirectUris:='["http://localhost"]' serviceAccountsEnabled=true secret=mysecret


# get client token using client credentials (requires serviceAccountsEnabled=true for the client)
http --form POST http://keycloak.127-0-0-121.nip.io/auth/realms/master/protocol/openid-connect/token grant_type=client_credentials client_id=foo client_secret=mysecret



# list all clients
http GET http://keycloak.127-0-0-121.nip.io/auth/admin/realms/master/clients Authorization:"bearer $TOKEN"

# list specific client
http GET http://keycloak.127-0-0-121.nip.io/auth/admin/realms/master/clients/29d6a41b-766f-42f6-9cbf-d5b18e476d8a Authorization:"bearer $TOKEN"

# list client secret
http GET http://keycloak.127-0-0-121.nip.io/auth/admin/realms/master/clients/29d6a41b-766f-42f6-9cbf-d5b18e476d8a/client-secret Authorization:"bearer $TOKEN"



#################################
#
# password grant
#

TOKEN=$(http --form POST http://keycloak.127-0-0-121.nip.io/auth/realms/master/protocol/openid-connect/token username=group-admin-user password=secret grant_type=password client_id=admin-cli | jq -r .access_token)

http -v GET http://keycloak.127-0-0-121.nip.io/auth/admin/realms/master/users Authorization:"bearer $TOKEN"



##################################
#
# Authorizatin code flow
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





