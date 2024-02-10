

https://github.com/keycloak/keycloak/issues/25993
https://github.com/keycloak/keycloak/pull/26096



kubectl apply -f manifests/keycloak-21.yaml
kubectl apply -f manifests/keycloak-22.yaml


# Make sure that SQL statements are logged
#   manifest/postgresql.yaml

    containers:
        - name: postgresql
          args:
            - "-c"
            - "log_statement=all"



# Change access token expiration to 1 day
TOKEN=$(http --form POST http://keycloak.127-0-0-121.nip.io/realms/master/protocol/openid-connect/token username=admin password=admin grant_type=password client_id=admin-cli | jq -r .access_token)
http -v PUT  http://keycloak.127-0-0-121.nip.io/admin/realms/master Authorization:"bearer $TOKEN" ssoSessionIdleTimeout:=86400 accessTokenLifespan:=86400

TOKEN=$(http --form POST http://localhost:8080/realms/master/protocol/openid-connect/token username=admin password=admin grant_type=password client_id=admin-cli | jq -r .access_token)
http -v PUT  http://localhost:8080/admin/realms/master Authorization:"bearer $TOKEN" ssoSessionIdleTimeout:=86400 accessTokenLifespan:=86400




http -v PUT  http://keycloak.127-0-0-121.nip.io/admin/realms/master Authorization:"bearer $TOKEN" bruteForceProtected=true defaultSignatureAlgorithm=ES384
http -v PUT  http://keycloak.127-0-0-121.nip.io/admin/realms/master Authorization:"bearer $TOKEN" bruteForceProtected=false defaultSignatureAlgorithm=ES256



http http://keycloak.127-0-0-121.nip.io/admin/realms/master Authorization:"bearer $TOKEN" > input.json
jq '.bruteForceProtected = true | .defaultSignatureAlgorithm = "ES384"' input.json > output.json
http -v PUT  http://keycloak.127-0-0-121.nip.io/admin/realms/master Authorization:"bearer $TOKEN" < output.json

while true; do http -v PUT  http://keycloak.    /admin/realms/master Authorization:"bearer $TOKEN" bruteForceProtected=true defaultSignatureAlgorithm=ES384 | grep "HTTP/1.1 4"; done
while true; do http -v PUT  http://keycloak.127-0-0-121.nip.io/admin/realms/master Authorization:"bearer $TOKEN" < input.json | grep "HTTP/1.1 4"; done



# send infinite number of realm update requests to keycloak

import requests

#URL = 'http://keycloak.127-0-0-121.nip.io'
URL = "http://localhost:8080"

TOKEN = requests.post(
    f"{URL}/realms/master/protocol/openid-connect/token",
    data={'username': 'admin', 'password': 'admin', 'grant_type': 'password', 'client_id': 'admin-cli'}
).json()['access_token']

num_requests = 1
while True:
    res = requests.put(
        f"{URL}/admin/realms/master",
        headers={'Authorization': f"bearer {TOKEN}"},
        #json={'bruteForceProtected': True, 'defaultSignatureAlgorithm': 'ES384', 'browserFlow': 'browser2'}
        json={'bruteForceProtected': True, 'defaultSignatureAlgorithm': 'ES384'}
        #json={"sslRequired":"NONE", 'defaultSignatureAlgorithm': 'ES384'}
        #json={'browserFlow': 'browser2'}
    )
    if res.status_code == 204:
        print(f'Attempt {num_requests} successful')
    else:
        print(f'Error at attempt {num_requests}: {res.status_code} {res.reason}: {res.text}')
    num_requests += 1



# send infinite number of realm import requests to keycloak

import requests

#URL = 'http://keycloak.127-0-0-121.nip.io'
URL = "http://localhost:8080"

TOKEN = requests.post(
    f"{URL}/realms/master/protocol/openid-connect/token",
    data={'username': 'admin', 'password': 'admin', 'grant_type': 'password', 'client_id': 'admin-cli'}
).json()['access_token']

REALM_DATA = requests.get(
    f"{URL}/admin/realms/master",
    headers={'Authorization': f"bearer {TOKEN}"}).text

num_requests = 1
while True:
num_requests = 1
while True:
    res = requests.put(
        f"{URL}/admin/realms/master",
        headers={'Authorization': f"bearer {TOKEN}"},
        #json={'bruteForceProtected': True, 'defaultSignatureAlgorithm': 'ES384', 'browserFlow': 'browser2'}
        #json={'bruteForceProtected': True, 'defaultSignatureAlgorithm': 'ES384'}
        json={"sslRequired":"external"}
        #json={'browserFlow': 'browser2'}
    )
    if res.status_code == 204:
        print(f'Attempt {num_requests} successful')
    else:
        print(f'Error at attempt {num_requests}: {res.status_code} {res.reason}: {res.text}')
    num_requests += 1

    res = requests.put(
        f"{URL}/admin/realms/master",
        headers={'Authorization': f"bearer {TOKEN}",
                 'Content-Type': 'application/json'},
        data=REALM_DATA)
    if res.status_code == 204:
        print(f'Attempt {num_requests} successful')
    else:
        print(f'Error at attempt {num_requests}: {res.status_code} {res.reason}: {res.text}')
    num_requests += 1



kubectl logs postgres-0 -f|grep deadlock


# Reset persistent data in postgres

kubectl scale statefulset keycloak --replicas 0
kubectl scale statefulset postgres --replicas 0
kubectl delete persistentvolumeclaims data-postgres-0
kubectl scale statefulset postgres --replicas 1
kubectl scale statefulset keycloak --replicas 1





*** Unit tests

# Failing test
mvn clean install -f testsuite/integration-arquillian/pom.xml -Dtest=org.keycloak.testsuite.federation.storage.UserStorageDirtyDeletionUnsyncedImportTest#testMembersWhenCachedUsersRemovedFromBackend -Dkeycloak.logging.level=debug
