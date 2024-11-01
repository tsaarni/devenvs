



get_admin_token() {
  http --form POST http://keycloak.127-0-0-121.nip.io/realms/master/protocol/openid-connect/token username=admin password=admin grant_type=password client_id=admin-cli | jq -r .access_token
}


# create realm for testing
http POST http://keycloak.127-0-0-121.nip.io/admin/realms "Authorization: Bearer $(get_admin_token)" realm=test enabled:=true


# create user and set password
http POST http://keycloak.127-0-0-121.nip.io/admin/realms/test/users "Authorization: Bearer $(get_admin_token)" username=joe email=joe@example.com enabled:=true firstName=Joe lastName=Doe attributes:='{"attr1":"val1","attr2":"val2"}' credentials:='[{"type":"password","value":"joe","temporary":false}]'

# get user id
USER_ID=$(http GET http://keycloak.127-0-0-121.nip.io/admin/realms/test/users "Authorization: Bearer $(get_admin_token)" | jq -r '.[] | select(.username == "joe") | .id')

# get user attributes
http GET http://keycloak.127-0-0-121.nip.io/admin/realms/test/users/$USER_ID "Authorization: Bearer $(get_admin_token)"

# update user attributes
http PUT http://keycloak.127-0-0-121.nip.io/admin/realms/test/users/$USER_ID "Authorization: Bearer $(get_admin_token)" firstName=John lastName=Smith

# add custom attribute
http PUT http://keycloak.127-0-0-121.nip.io/admin/realms/test/users/$USER_ID "Authorization: Bearer $(get_admin_token)" attributes:='{"attr3":"val3"}'


# delete user
http DELETE http://keycloak.127-0-0-121.nip.io/admin/realms/test/users/$USER_ID "Authorization: Bearer $(get_admin_token)"


# delete realm
http DELETE http://keycloak.127-0-0-121.nip.io/admin/realms/test "Authorization: Bearer $(get_admin_token)"
