
# https://github.com/keycloak/keycloak/issues/9609

# Keycloak 18.0.2



kind delete cluster --name keycloak
kind create cluster --config configs/kind-cluster-config.yaml --name keycloak

kubectl apply -f https://projectcontour.io/quickstart/contour.yaml
kubectl apply -f manifests/postgresql.yaml
kubectl apply -f manifests/keycloak-18.yaml

kubectl create secret tls keycloak-external --cert=certs/keycloak-server.pem --key=certs/keycloak-server-key.pem --dry-run=client -o yaml | kubectl apply -f -


# openldap

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






http://keycloak.127-0-0-121.nip.io/
http://pgadmin.127-0-0-121.nip.io



apps/create-components.py --server=http://keycloak.127-0-0-121.nip.io/ rest-requests/create-ldap-simple-auth-provider-kubernetes.json



Create group mapper

1. Select type:      group-ldap-mapper
2. Name:             my-group-mapper
3. LDAP Groups DN:   ou=groups,o=example
4. Mode:             IMPORT


or

Create role mapper

1. Select type:      role-ldap-mapper
2. Name:             my-role-mapper
3. LDAP Roles DN:    ou=groups,o=example

alternatively, using memberof

Create group mapper

1. Select type:                   group-ldap-mapper
2. Name:                          my-group-mapper
3. LDAP Groups DN:                ou=groups,o=example
4. User Groups Retrieve Strategy: MEMBEROF_ATTRIBUTE





# Run LDAP queries within the pod

kubectl exec -it deployment/openldap -- bash

# List users and group memberships
ldapsearch -H ldap://localhost -w ldap-admin -D "cn=ldap-admin,ou=users,o=example" -b "ou=users,o=example" memberOf
ldapsearch -H ldap://localhost -w ldap-admin -D "cn=ldap-admin,ou=users,o=example" -b "ou=groups,o=example" member


# create users group
ldapadd -H ldap://localhost -w ldap-admin -D "cn=ldap-admin,ou=users,o=example" <<EOF
dn: cn=users,ou=groups,o=example
objectClass: groupOfNames
cn: users
member: cn=user1,ou=users,o=example
member: cn=user2,ou=users,o=example
EOF



# remove user from group
ldapmodify -H ldap://localhost -w ldap-admin -D "cn=ldap-admin,ou=users,o=example" <<EOF
dn: cn=users,ou=groups,o=example
changetype: modify
delete: member
member: cn=user2,ou=users,o=example
EOF

# add user to group
ldapmodify -H ldap://localhost -w ldap-admin -D "cn=ldap-admin,ou=users,o=example" <<EOF
dn: cn=users,ou=groups,o=example
changetype: modify
add: member
member: cn=user2,ou=users,o=example
EOF



sudo nsenter -n -t $(pidof slapd) wireshark -f "port 389 or port 636" -Y ldap -k -o tls.keylog_file:$WORKDIR/output/wireshark-keys.log




# Reconfigure OpenLDAP

kubectl create configmap openldap-config --dry-run=client -o yaml --from-file=templates/database.ldif --from-file=templates/users-and-groups.ldif | kubectl apply -f -
kubectl delete pod -l app=openldap --force
