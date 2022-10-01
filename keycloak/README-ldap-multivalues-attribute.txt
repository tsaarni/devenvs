


Cannot add multi valued user attribute when using LDAP federation
https://issues.redhat.com/browse/KEYCLOAK-16575



#
# start new cluster
#

kind delete cluster --name keycloak
kind create cluster --config configs/kind-cluster-config.yaml --name keycloak

# deploy contour
kubectl apply -f https://projectcontour.io/quickstart/contour.yaml

# deploy openldap
docker build docker/openldap/ -t localhost/openldap:latest
kind load docker-image --name keycloak localhost/openldap:latest
kubectl create configmap openldap-config --dry-run=client -o yaml --from-file=templates/database.ldif --from-file=templates/users-and-groups.ldif | kubectl apply -f -
kubectl create secret tls openldap-cert --cert=certs/ldap.pem --key=certs/ldap-key.pem --dry-run=client -o yaml | kubectl apply -f -

# patch tls secret to inject ca.crt
kubectl patch secret openldap-cert --patch-file /dev/stdin <<EOF
data:
  ca.crt: $(cat certs/client-ca.pem | base64 -w 0)
EOF

kubectl apply -f manifests/openldap.yaml

# deploy postgres
kubectl apply -f manifests/postgresql.yaml


# keycloak
kubectl create secret tls keycloak-external --cert=certs/keycloak-server.pem --key=certs/keycloak-server-key.pem --dry-run=client -o yaml | kubectl apply -f -
kubectl apply -f manifests/keycloak-19.yaml


# https://keycloak.127-0-0-121.nip.io/




# capture LDAP traffic
sudo nsenter --target $(pidof slapd) --net wireshark -f  "port 389" --display-filter ldap -k



# Test case


1. In User Federation, add LDAP provider

   vendor: Other
   connection URL:  ldap://openldap
   bind DN:  cn=ldap-admin,ou=users,o=example
   bind credentials: ldap-admin
   edit mode: WRITABLE
   users DN: ou=users,o=example
   User object classes: person, organizationalPerson, user, inetOrgPerson
   Sync registrations: true

   # MISSING ATTRIBUTE IN UI: syncRegistrations true
   # https://github.com/keycloak/keycloak-ui/issues/3241
   # https://github.com/keycloak/keycloak-ui/pull/3339
   #
   # WORKAROUND: switch realm theme temporarily to keycloak (v1)
   # Realm settings / Themes / Admin console theme: keycloak -> reload page

2. In LDAP provider settings, Mappers, Add Mapper

   Name: type
   Mapper type: user-attribute-ldap-mapper
   User model attribute: employeeType
   LDAP attribute: employeeType

3. Add user Joe with employeeType foo##bar via GUI

   Definition of the inetOrgPerson LDAP Object Class
   https://www.ietf.org/rfc/rfc2798.txt

4. Add user Joe

   add attribute "type" with value "foo##bar"

5. Check user in LDAP

kubectl exec deployments/openldap -- ldapsearch -H ldapi:/// -D cn=ldap-admin,ou=users,o=example -w ldap-admin -b ou=users,o=example

1. Create mapping "type" -> "employeeType"
here employeeType is from "inetOrgPerson" schema and it is multivalued attribute (i.e. in LDAP schema it is not set as SINGLE-VALUE)

2. Create user Joe with type foo##bar via GUI

3. Check that user has multiple attributes as "type" when user is fetched from Keycloak

TOKEN=$(http --form POST http://keycloak.127-0-0-121.nip.io/realms/master/protocol/openid-connect/token username=admin password=admin grant_type=password client_id=admin-cli | jq -r .access_token)

http "http://keycloak.127-0-0-121.nip.io//admin/realms/master/users?username=joe" Authorization:"bearer $TOKEN"


HTTP/1.1 200 OK
Cache-Control: no-cache
Connection: keep-alive
Content-Length: 619
Content-Type: application/json
Date: Fri, 01 Jul 2022 11:46:19 GMT
Referrer-Policy: no-referrer
Strict-Transport-Security: max-age=31536000; includeSubDomains
X-Content-Type-Options: nosniff
X-Frame-Options: SAMEORIGIN
X-XSS-Protection: 1; mode=block
[
    {
        "access": {
            "impersonate": true,
            "manage": true,
            "manageGroupMembership": true,
            "mapRoles": true,
            "view": true
        },
        "attributes": {
            "LDAP_ENTRY_DN": [
                "uid=joe,ou=users,o=example"
            ],
            "LDAP_ID": [
                "9959b12e-8d72-103c-8018-35a2d282bd8e"
            ],
            "createTimestamp": [
                "20220701101621Z"
            ],
            "modifyTimestamp": [
                "20220701114331Z"
            ],
            "type": [
                "foo",
                "bar"
            ]
        },
        "createdTimestamp": 1656670581614,
        "disableableCredentialTypes": [],
        "emailVerified": false,
        "enabled": true,
        "federationLink": "4591c981-faa4-4be2-86b3-167ee466d0f9",
        "firstName": "First",
        "id": "039d4b89-dc15-48d2-b3a7-ab4e8c6834b7",
        "lastName": "Last",
        "notBefore": 0,
        "requiredActions": [],
        "totp": false,
        "username": "joe"
    }
]

4. Check that multiple "employeeType" attributes are there when fetched from LDAP

kubectl exec deployments/openldap -- ldapsearch -H ldapi:/// -D cn=ldap-admin,ou=users,o=example -w ldap-admin -b uid=joe,ou=users,o=example


# extended LDIF
#
# LDAPv3
# base <uid=joe,ou=users,o=example> with scope subtree
# filter: (objectclass=*)
# requesting: ALL
#
# joe, users, example
dn: uid=joe,ou=users,o=example
uid: joe
objectClass: inetOrgPerson
objectClass: organizationalPerson
employeeType: foo
employeeType: bar
sn: Last
cn: First
# search result
search: 2
result: 0 Success
# numResponses: 2
# numEntries: 1









# start LDAP server
docker-compose rm -f  # clean previous containers
docker-compose up
docker exec keycloak-openldap-1 ldapsearch -H ldapi:/// -D cn=ldap-admin,ou=users,o=example -w ldap-admin -b cn=config
docker exec keycloak-openldap-1 ldapsearch -H ldapi:/// -D cn=ldap-admin,ou=users,o=example -w ldap-admin -b ou=users,o=example
