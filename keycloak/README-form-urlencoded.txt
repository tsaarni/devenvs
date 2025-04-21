
https://github.com/keycloak/keycloak/issues/16401
https://datatracker.ietf.org/doc/html/rfc6749#section-4.1.3


function get_admin_token() {
  http --form POST http://keycloak.127-0-0-1.nip.io:8080/realms/master/protocol/openid-connect/token username=admin password=admin grant_type=password client_id=admin-cli | jq -r .access_token
}

# create client with % in secret
http -v POST http://keycloak.127-0-0-1.nip.io:8080/admin/realms/master/clients "Authorization: Bearer $(get_admin_token)" clientId=test-client publicClient=false secret="secret-with-percent-%" serviceAccountsEnabled=true
POST /admin/realms/master/clients HTTP/1.1
Accept: application/json, */*;q=0.5
Accept-Encoding: gzip, deflate, br
Authorization: Bearer <token omitted>
Connection: keep-alive
Content-Length: 121
Content-Type: application/json
Host: keycloak.127-0-0-1.nip.io:8080
User-Agent: HTTPie/3.2.2

{
    "clientId": "test-client",
    "publicClient": "false",
    "secret": "secret-with-percent-%",
    "serviceAccountsEnabled": "true"
}


HTTP/1.1 201 Created
Location: http://keycloak.127-0-0-1.nip.io:8080/admin/realms/master/clients/d1d95593-5e2c-4578-a6f4-7225daab37cc
Referrer-Policy: no-referrer
Strict-Transport-Security: max-age=31536000; includeSubDomains
X-Content-Type-Options: nosniff
X-Frame-Options: SAMEORIGIN
X-XSS-Protection: 1; mode=block
content-length: 0


# Succeeds: fetch token with correctly formatted URL encoded form parameters
http -v POST http://keycloak.127-0-0-1.nip.io:8080/realms/master/protocol/openid-connect/token Content-Type:application/x-www-form-urlencoded < <(echo -n "grant_type=client_credentials&client_id=test-client&client_secret=secret-with-percent-%25")

POST /realms/master/protocol/openid-connect/token HTTP/1.1
Accept: application/json, */*;q=0.5
Accept-Encoding: gzip, deflate, br
Connection: keep-alive
Content-Length: 89
Content-Type: application/x-www-form-urlencoded
Host: keycloak.127-0-0-1.nip.io:8080
User-Agent: HTTPie/3.2.2

grant_type=client_credentials&client_id=test-client&client_secret=secret-with-percent-%25


HTTP/1.1 200 OK
Cache-Control: no-store
Content-Type: application/json
Pragma: no-cache
Referrer-Policy: no-referrer
Strict-Transport-Security: max-age=31536000; includeSubDomains
X-Content-Type-Options: nosniff
X-Frame-Options: SAMEORIGIN
X-XSS-Protection: 1; mode=block
content-length: 1391

{
    "access_token": "<token omitted>",
    "expires_in": 60,
    "not-before-policy": 0,
    "refresh_expires_in": 0,
    "scope": "profile email",
    "token_type": "Bearer"
}



# Fails: when we try fetch token withtout URL encoding then request fails and Keycloak logs
# 2025-01-22 12:53:16,913 ERROR [org.keycloak.services.error.KeycloakErrorHandler] (executor-thread-29) Uncaught server error: java.lang.RuntimeException: Failed to decode URL secret-with-percent-% to UTF-8
http -v POST http://keycloak.127-0-0-1.nip.io:8080/realms/master/protocol/openid-connect/token Content-Type:application/x-www-form-urlencoded < <(echo -n "grant_type=client_credentials&client_id=test-client&client_secret=secret-with-percent-%")

POST /realms/master/protocol/openid-connect/token HTTP/1.1
Accept: application/json, */*;q=0.5
Accept-Encoding: gzip, deflate, br
Connection: keep-alive
Content-Length: 87
Content-Type: application/x-www-form-urlencoded
Host: keycloak.127-0-0-1.nip.io:8080
User-Agent: HTTPie/3.2.2

grant_type=client_credentials&client_id=test-client&client_secret=secret-with-percent-%


HTTP/1.1 500 Internal Server Error
Content-Type: application/json
Referrer-Policy: no-referrer
Strict-Transport-Security: max-age=31536000; includeSubDomains
X-Content-Type-Options: nosniff
X-Frame-Options: SAMEORIGIN
X-XSS-Protection: 1; mode=block
content-length: 94

{
    "error": "unknown_error",
    "error_description": "For more on this error consult the server log."
}
