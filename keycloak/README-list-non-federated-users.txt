
￼



TOKEN=$(http --form POST http://localhost:8081/auth/realms/master/protocol/openid-connect/token username=admin password=admin grant_type=password client_id=admin-cli | jq -r .access_token)


http -v GET http://localhost:8081/auth/admin/realms/master/users Authorization:"bearer $TOKEN"

vs

http -v GET "http://localhost:8081/auth/admin/realms/master/users?username=*" Authorization:"bearer $TOKEN"


- 21.1.1 return only local users for /users endpoint without parameters
- 22.0.5 return all users for /users endpoint without parameters
- both return all users with parameter /users?username=*
- with invalid federation config: 21.1.1 with /users was working, /users?username=* returns error
- with invalid federation config: 22.0.5 always returns error
- 21.1.1 returns only local users if searching with /users?username=test but also federated users if searching with /users?username=test*


How to search only local users?
How to search only for local users if federation config is broken?

