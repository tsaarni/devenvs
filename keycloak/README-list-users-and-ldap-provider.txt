
Following change caused change in behavior


Make LDAP searchForUsersStream consistent with other storages
https://github.com/keycloak/keycloak/issues/17294

> Searching in LDAP works differently than other storages
> ... 
> with an empty params argument, it returns an empty result list. This is inconsistent with other storages, where empty params return all users.

 

PR https://github.com/keycloak/keycloak/pull/19050


￼



TOKEN=$(http --form POST http://localhost:8081/auth/realms/master/protocol/openid-connect/token username=admin password=admin grant_type=password client_id=admin-cli | jq -r .access_token)


http -v GET http://localhost:8081/auth/admin/realms/master/users Authorization:"bearer $TOKEN"

vs

http -v GET "http://localhost:8081/auth/admin/realms/master/users?search=*" Authorization:"bearer $TOKEN"


Eg difference between 21.1.1 vs 22.0.5


- Keycloak 21 returns only local users for /users endpoint without parameters
- Keycloak 21 returns only local users if searching with /users?search=test but also federated users if searching with /users?seach=test*
- Keycloak 22 return all users for /users endpoint without parameters
- both return all users with parameter /users?search=*

Behavior when configured with invalid federation config

- Keycloak 21 fetching /users was working, /users?search=* returned error
- Keycloak 22 always returns error

Note that same applies to other search filters that accept *, not just /users?search=*


Questions:
How to search only local users?
Can that local search be executed while federated config is invalid or if server does not respond?

