


# start LDAP server
docker-compose rm -f  # clean previous containers
docker-compose up



#




docker exec keycloak-openldap-1 ldapsearch -H ldapi:/// -D cn=ldap-admin,ou=users,o=example -w ldap-admin -b cn=config
docker exec keycloak-openldap-1 ldapsearch -H ldapi:/// -D cn=ldap-admin,ou=users,o=example -w ldap-admin -b ou=users,o=example


TOKEN=$(http --form POST http://localhost:8081/auth/realms/master/protocol/openid-connect/token username=admin password=admin grant_type=password client_id=admin-cli | jq -r .access_token)
http -v GET "http://localhost:8081/auth/admin/realms/master/users?username=joe" Authorization:"bearer $TOKEN"
