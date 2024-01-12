



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



http -v PUT  http://keycloak.127-0-0-121.nip.io/admin/realms/master Authorization:"bearer $TOKEN" bruteForceProtected=true defaultSignatureAlgorithm=ES384
http -v PUT  http://keycloak.127-0-0-121.nip.io/admin/realms/master Authorization:"bearer $TOKEN" bruteForceProtected=false defaultSignatureAlgorithm=ES256



http http://keycloak.127-0-0-121.nip.io/admin/realms/master Authorization:"bearer $TOKEN" > input.json
jq '.bruteForceProtected = true | .defaultSignatureAlgorithm = "ES384"' input.json > output.json
http -v PUT  http://keycloak.127-0-0-121.nip.io/admin/realms/master Authorization:"bearer $TOKEN" < output.json

while true; do http -v PUT  http://keycloak.127-0-0-121.nip.io/admin/realms/master Authorization:"bearer $TOKEN" bruteForceProtected=true defaultSignatureAlgorithm=ES384 | grep "HTTP/1.1 4"; done
while true; do http -v PUT  http://keycloak.127-0-0-121.nip.io/admin/realms/master Authorization:"bearer $TOKEN" < input.json | grep "HTTP/1.1 4"; done



# send infinite number of realm update requests to keycloak

import os
import requests

TOKEN = requests.post(
    'http://keycloak.127-0-0-121.nip.io/realms/master/protocol/openid-connect/token',
    data={'username': 'admin', 'password': 'admin', 'grant_type': 'password', 'client_id': 'admin-cli'}
).json()['access_token']

num_requests = 1
while True:
    res = requests.put(
        'http://keycloak.127-0-0-121.nip.io/admin/realms/master',
        headers={'Authorization': 'bearer ' + TOKEN},
        json={'bruteForceProtected': True, 'defaultSignatureAlgorithm': 'ES384'}
    )
    if res.status_code == 204:
        print(f'Attempt {num_requests} successful')
    else:
        print(f'Error at attempt {num_requests}: {res.status_code} {res.text}')
    num_requests += 1



# send infinite number of realm import requests to keycloak

import os
import requests

TOKEN = requests.post(
    'http://keycloak.127-0-0-121.nip.io/realms/master/protocol/openid-connect/token',
    data={'username': 'admin', 'password': 'admin', 'grant_type': 'password', 'client_id': 'admin-cli'}
).json()['access_token']

REALM_DATA = open('input.json', 'r').read()

num_requests = 1
while True:
    res = requests.put(
        'http://keycloak.127-0-0-121.nip.io/admin/realms/master',
        headers={'Authorization': 'bearer ' + TOKEN,
                 'Content-Type': 'application/json'},
        data=REALM_DATA)
    if res.status_code == 204:
        print(f'Attempt {num_requests} successful')
    else:
        print(f'Error at attempt {num_requests}: {res.status_code} {res.text}')
    num_requests += 1



kubectl logs postgres-0 -f|grep deadlock


# Reset persistent data in postgres

kubectl scale statefulset keycloak --replicas 0
kubectl scale statefulset postgres --replicas 0
kubectl delete persistentvolumeclaims data-postgres-0
kubectl scale statefulset postgres --replicas 1
kubectl scale statefulset keycloak --replicas 1
