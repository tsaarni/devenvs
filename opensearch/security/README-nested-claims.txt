


https://github.com/opensearch-project/security/issues/551
https://github.com/opensearch-project/security/pull/1122



######################
#
# setup environment
#

# create certificates
rm -rf certs
mkdir certs
certyaml -d certs/ configs/certs.yaml

# start opensearch, opensearch-dashboards and keycloak
docker-compose rm -f
docker-compose up


# access opensearch-dashboards
http://opensearch-dashboards.127-0-0-1.nip.io:5601


# access keycloak admin ui
http://keycloak.127-0-0-1.nip.io:8080



######################
#
# work with the security plugin source code
#

# compile
cd ~/work/opensearch-security
./gradlew clean assemble

ls -l build/distributions/

# copy the jar so that opensearch container can pick it up
cp build/distributions/opensearch-security-3.0.0.0-SNAPSHOT.zip ~/work/devenvs/opensearch/security/



# over at the devenvs directory (where the docker-compose.yml is located) run the following command to restart opensearch
docker-compose restart opensearch





######################
#
# troubleshooting
#

# test if opensearch is up and running
http --verify=certs/server-ca.pem --cert=certs/opensearch-admin.pem --cert-key=certs/opensearch-admin-key.pem https://localhost:9200/_cluster/health

# fetch well-known configuration from Keycloak
http http://localhost:8080/realms/master/.well-known/openid-configuration

# fetch admin token from keycloak
http --form POST http://localhost:8080/realms/master/protocol/openid-connect/token username=admin password=admin grant_type=password client_id=admin-cli

# capture http traffic (including issued JWT access token) at keycloak container side
sudo nsenter -t $(docker inspect --format '{{ .State.Pid }}' security-keycloak-1) --net wireshark -i any -k -f "tcp port 8080" -Y "http"

# decode jwt by running the following command and pasting the access token on the console
jq -R 'split(".") | .[1] | @base64d | fromjson'


{
  "exp": 1727504492,
  "iat": 1727504192,
  "auth_time": 1727504192,
  "jti": "079b5fe3-aec4-48a3-8f07-279cc1cd9d6a",
  "iss": "http://keycloak.127-0-0-1.nip.io:8080/realms/opensearch",
  "sub": "b855ab13-c0c1-4413-a6a9-ec8ca2d499b7",
  "typ": "Bearer",
  "azp": "opensearch-dashboards",
  "sid": "3f5d91b4-808f-4dc4-a362-128ebff03e28",
  "acr": "1",
  "allowed-origins": [
    "*"
  ],
  "realm_access": {
    "roles": [
      "demo"
    ]
  },
  "scope": "openid address profile phone email",
  "email_verified": false,
  "address": {},
  "name": "Demo User",
  "preferred_username": "demo",
  "given_name": "Demo",
  "family_name": "User",
  "email": "demo@example.com"
}










Author: Andy Lin <ndylin@amazon.com>
Date:   Fri Apr 16 09:28:34 2021 -0700

    Re-implement RolesUtil and add more tests

commit 3b903a6aa18f2ad641746f18889c88481cf236ab
Author: Andy Lin <ndylin@amazon.com>
Date:   Wed Apr 14 16:59:15 2021 -0700

    Rename Roles to RolesUtil, add comments to RolesUtil, add more unit tests

commit 65064e8cfee59ba43f6d9fd8ce8f78417b999a48
Author: Andy Lin <ndylin@amazon.com>
Date:   Wed Apr 14 09:48:55 2021 -0700

    Add rolespath to JWT authentication



https://github.com/opensearch-project/security/pull/3262

After RFC 6901 was introduced and the implementation was added to Jackson,
there is no need to keep the com.jayway.jsonpath:json-path library in our source code,
so we can replace current validation with Jackson's JsonPointer class.



# docs
https://opensearch.org/docs/latest/security/authentication-backends/jwt/
https://opensearch.org/docs/latest/security/authentication-backends/openid-connect/
