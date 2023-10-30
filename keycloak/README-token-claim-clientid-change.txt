

https://github.com/keycloak/keycloak/pull/16359


# patch to put legacy clientId back
https://github.com/Nordix/keycloak/tree/sa-with-legacy-clientid



https://www.keycloak.org/docs/latest/server_admin/#_service_accounts

1. Create client from admin console
2. Enable Client authentication for the client
3. Select "Service account roles" from authentication flow
4. Click save
5. Select "Credentials" tab to check the client id and secret


URL=http://keycloak.127-0-0-1.nip.io:8080


# get the client secret from client credentials tab
CREDENTIALS=myclient:OQMQTd2LXkHwThMiuNHEKyr34bSjDkOn

http -a $CREDENTIALS -f POST $URL/realms/master/protocol/openid-connect/token grant_type=client_credentials | jq -r '.access_token' | apps/jwt-decode.py



# create client with legacy "clientID" claim
TOKEN=$(http --form POST $URL/realms/master/protocol/openid-connect/token username=admin password=admin grant_type=password client_id=admin-cli | jq -r .access_token)
http -v POST $URL/admin/realms/master/clients Authorization:"bearer $TOKEN" clientId=myclient name="my client" publicClient=false serviceAccountsEnabled=true protocolMappers:='[{"name": "Client ID (Legacy)", "protocol": "openid-connect", "protocolMapper": "oidc-usersessionmodel-note-mapper", "config":{ "user.session.note": "clientId", "id.token.claim": "true", "access.token.claim": "true", "introspection.token.claim": "true", "claim.name": "foobar", "jsonType.label": "String" }}]'
