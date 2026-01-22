

# AD objectGUID-based user lookups fail in Keycloak 26.2.x (works in 26.0.7)
https://github.com/keycloak/keycloak/issues/43324


#
# https://github.com/dkoudela/active-directory-to-openldap



# Run OpenLDAP server

docker build --tag openldap-msad:latest .

docker run --rm -p 1389:1389 -p 1636:1636 openldap-msad:latest



# Create user


ldapadd -x -H ldap://localhost:1389 -D "cn=dsadmin,dc=example,dc=test" -w 'password1' <<EOF
dn: CN=myuser,ou=users,dc=example,dc=test
changetype: add
objectClass: user
objectClass: organizationalPerson
objectClass: person
cn: My
sn: User
sAMAccountName: myuser
userAccountControl: 512
objectGUID:: 7c9Z7yK7Q0m9k8F8b0pR1A==
EOF


# list user
ldapsearch -x -H ldap://localhost:1389 -D "cn=dsadmin,dc=example,dc=test" -w 'password1' -b "ou=users,dc=example,dc=test" "(sAMAccountName=myuser)"


# remove user
ldapdelete -x -H ldap://localhost:1389 -D "cn=dsadmin,dc=example,dc=test" -w 'password1' "cn=myuser,ou=users,dc=example,dc=test"




######################
#
# Keycloak
#


function get_admin_token() {
  http --form POST http://keycloak.127-0-0-1.nip.io:8080/realms/master/protocol/openid-connect/token username=admin password=admin grant_type=password client_id=admin-cli | jq -r .access_token
}



# Create LDAP provider

http -v POST http://keycloak.127-0-0-1.nip.io:8080/admin/realms/master/components \
  Authorization:"Bearer $(get_admin_token)" \
  id=my-ldap \
  name="my ldap" \
  providerId=ldap \
  providerType="org.keycloak.storage.UserStorageProvider" \
  config[krbPrincipalAttribute][]=userPrincipalName \
  config[pagination][]=true \
  config[startTls][]=false \
  config[searchScope][]=2 \
  config[useTruststoreSpi][]=ldapsOnly \
  config[maxLifespan][]=0 \
  config[connectionPooling][]=true \
  config[usersDn][]=ou=users,dc=example,dc=test \
  config[cachePolicy][]=NO_CACHE \
  config[priority][]=0 \
  config[importEnabled][]=false \
  config[userObjectClasses][]='person, organizationalPerson, user' \
  config[enabled][]=true \
  config[bindDn][]=cn=alice,ou=users,dc=example,dc=test \
  config[usernameLDAPAttribute][]=sAMAccountName \
  config[bindCredential][]=password11 \
  config[rdnLDAPAttribute][]=cn \
  config[vendor][]=ad \
  config[readTimeout][]=30000 \
  config[editMode][]=READ_ONLY \
  config[uuidLDAPAttribute][]=objectGUID \
  config[connectionUrl][]=ldap://localhost:1389 \
  config[authType][]=simple \
  config[connectionTimeout][]=10000



# Lookup user by username
http -v GET http://keycloak.127-0-0-1.nip.io:8080/admin/realms/master/users?username=myuser Authorization:"Bearer $(get_admin_token)"

# Lookup credentials by user ID
http -v GET http://keycloak.127-0-0-1.nip.io:8080/admin/realms/master/users/f:my-ldap:myuser/credentials Authorization:"Bearer $(get_admin_token)"


# capture traffic with wiresarhk
sudo nsenter --target $(pgrep -f slapd) --net wireshark -i any -k -Y "ldap"  -o "ldap.tcp.port:1389"




LDAPMessage searchResEntry(28) "cn=myuser,ou=users,dc=example,dc=test" [1 result]
    messageID: 28
    protocolOp: searchResEntry (4)
        searchResEntry
            objectName: cn=myuser,ou=users,dc=example,dc=test
            attributes: 6 items
                PartialAttributeList item objectClass
                    type: objectClass
                    vals: 3 items
                        AttributeValue: user
                        AttributeValue: organizationalPerson
                        AttributeValue: person
                PartialAttributeList item cn
                    type: cn
                    vals: 2 items
                        AttributeValue: My
                        AttributeValue: myuser
                PartialAttributeList item sn
                    type: sn
                    vals: 1 item
                        AttributeValue: User
                PartialAttributeList item sAMAccountName
                    type: sAMAccountName
                    vals: 1 item
                        AttributeValue: myuser
                PartialAttributeList item userAccountControl
                    type: userAccountControl
                    vals: 1 item
                        AttributeValue: 512
                PartialAttributeList item objectGUID
                    type: objectGUID
                    vals: 1 item
                        GUID: ef59cfed-bb22-4943-bd93-c17c6f4a51d4
    [Response To: 4]
    [Time: 244.477 microseconds]




LDAPMessage searchRequest(30) "ou=users,dc=example,dc=test" wholeSubtree
    messageID: 30
    protocolOp: searchRequest (3)
        searchRequest
            baseObject: ou=users,dc=example,dc=test
            scope: wholeSubtree (2)
            derefAliases: derefAlways (3)
            sizeLimit: 0
            timeLimit: 0
            typesOnly: False
            Filter: (objectGUID=ed:cf:59:ef:22:bb:43:49:bd:93:c1:7c:6f:4a:51:d4)
                filter: equalityMatch (3)
                    equalityMatch
                        attributeDesc: objectGUID
                        assertionValue: ed:cf:59:ef:22:bb:43:49:bd:93:c1:7c:6f:4a:51:d4
            attributes: 11 items
                AttributeDescription: whenChanged
                AttributeDescription: whenCreated
                AttributeDescription: mail
                AttributeDescription: sAMAccountName
                AttributeDescription: objectGUID
                AttributeDescription: sn
                AttributeDescription: cn
                AttributeDescription: objectclass
                AttributeDescription: userPrincipalName
                AttributeDescription: userAccountControl
                AttributeDescription: pwdLastSet
    [Response In: 11]
    controls: 1 item
        Control
            controlType: 2.16.840.1.113730.3.4.2 (Manage DSA IT LDAPv3 control)




https://github.com/keycloak/keycloak/pull/24771/files#diff-303dc3651881a87425f86de0778df2cca8527933f1a7bc21e6da72ca78f2c28f



ldapsearch \
  -x -H ldap://localhost:1389 -D "cn=dsadmin,dc=example,dc=test" -w 'password1'  \
  -b "ou=users,dc=example,dc=test" \
  -s sub \
  -a always \
  -z 0 \
  -l 0 \
  -M \
  "(objectGUID=\ed\cf\59\ef\22\bb\43\49\bd\93\c1\7c\6f\4a\51\d4)" \
  whenChanged whenCreated mail sAMAccountName objectGUID sn cn objectclass userPrincipalName userAccountControl pwdLastSet





federation/ldap/src/main/java/org/keycloak/storage/ldap/idm/store/ldap/LDAPOperationManager.java:400

    public SearchResult lookupById(final LdapName baseDN, final String id, final Collection<String> returningAttributes) {



################


docker run --rm --name keycloak \
  --network host \
  -e KC_HOSTNAME=keycloak.127.0.0.1.nip.io \
  -e KC_BOOTSTRAP_ADMIN_USERNAME=admin \
  -e KC_BOOTSTRAP_ADMIN_PASSWORD=admin \
  quay.io/keycloak/keycloak:26.0.7 \
  start-dev --verbose


REST request:

$ http -v GET http://keycloak.127-0-0-1.nip.io:8080/admin/realms/master/users/f:my-ldap:myuser/credentials Authorization:"Bea
rer $(get_admin_token)"


LDAP capture:


Plain Text
No.	Time	Source	Destination	Protocol	Length	Stream index	Info
4	2025-12-16 10:38:07,689220969	172.17.0.1	172.17.0.2	LDAP	128	0	bindRequest(1) "cn=alice,ou=users,dc=example,dc=test" simple
6	2025-12-16 10:38:07,689561743	172.17.0.2	172.17.0.1	LDAP	82	0	bindResponse(1) success
8	2025-12-16 10:38:07,690736371	172.17.0.1	172.17.0.2	LDAP	396	0	searchRequest(2) "ou=users,dc=example,dc=test" wholeSubtree
9	2025-12-16 10:38:07,691057689	172.17.0.2	172.17.0.1	LDAP	297	0	searchResEntry(2) "cn=myuser,ou=users,dc=example,dc=test"
10	2025-12-16 10:38:07,691085495	172.17.0.2	172.17.0.1	LDAP	82	0	searchResDone(2) success  [1 result]



Search request in detail:


Plain Text
LDAPMessage searchRequest(2) "ou=users,dc=example,dc=test" wholeSubtree
    messageID: 2
    protocolOp: searchRequest (3)
        searchRequest
            baseObject: ou=users,dc=example,dc=test
            scope: wholeSubtree (2)
            derefAliases: derefAlways (3)
            sizeLimit: 0
            timeLimit: 0
            typesOnly: False
            Filter: (&(&(&(&(sAMAccountName=myuser))(objectclass=person))(objectclass=organizationalPerson))(objectclass=user))
            attributes: 11 items
                AttributeDescription: whenChanged
                AttributeDescription: whenCreated
                AttributeDescription: mail
                AttributeDescription: sAMAccountName
                AttributeDescription: objectGUID
                AttributeDescription: sn
                AttributeDescription: cn
                AttributeDescription: objectclass
                AttributeDescription: userPrincipalName
                AttributeDescription: userAccountControl
                AttributeDescription: pwdLastSet
    [Response In: 9]
    controls: 1 item





The REST response contains empty list of credentials:

HTTP/1.1 200 OK
Cache-Control: no-cache
Content-Type: application/json;charset=UTF-8
Referrer-Policy: no-referrer
Strict-Transport-Security: max-age=31536000; includeSubDomains
X-Content-Type-Options: nosniff
X-Frame-Options: SAMEORIGIN
X-XSS-Protection: 1; mode=block
content-length: 2

[]



 ############



keycloak 26.2.10-nordix

same REST request as before

LDAP capture:


Plain Text
No.	Time	Source	Destination	Protocol	Length	Stream index	Info
12	2025-12-16 10:45:26,028987624	172.17.0.1	172.17.0.2	LDAP	128	0	bindRequest(1) "cn=alice,ou=users,dc=example,dc=test" simple
14	2025-12-16 10:45:26,029298653	172.17.0.2	172.17.0.1	LDAP	82	0	bindResponse(1) success
16	2025-12-16 10:45:26,030206947	172.17.0.1	172.17.0.2	LDAP	396	0	searchRequest(2) "ou=users,dc=example,dc=test" wholeSubtree
17	2025-12-16 10:45:26,030512184	172.17.0.2	172.17.0.1	LDAP	297	0	searchResEntry(2) "cn=myuser,ou=users,dc=example,dc=test"
18	2025-12-16 10:45:26,030527620	172.17.0.2	172.17.0.1	LDAP	82	0	searchResDone(2) success  [2 results]
20	2025-12-16 10:45:26,046543107	172.17.0.1	172.17.0.2	LDAP	128	0	bindRequest(3) "cn=alice,ou=users,dc=example,dc=test" simple
21	2025-12-16 10:45:26,046785205	172.17.0.2	172.17.0.1	LDAP	82	0	bindResponse(3) success
22	2025-12-16 10:45:26,047043973	172.17.0.1	172.17.0.2	LDAP	317	0	searchRequest(4) "ou=users,dc=example,dc=test" wholeSubtree
23	2025-12-16 10:45:26,047366316	172.17.0.2	172.17.0.1	LDAP	297	0	searchResEntry(4) "cn=myuser,ou=users,dc=example,dc=test"
24	2025-12-16 10:45:26,047389471	172.17.0.2	172.17.0.1	LDAP	82	0	searchResDone(4) success  [2 results]


Search requests in detail:

1st search


Plain Text
LDAPMessage searchRequest(2) "ou=users,dc=example,dc=test" wholeSubtree
    messageID: 2
    protocolOp: searchRequest (3)
        searchRequest
            baseObject: ou=users,dc=example,dc=test
            scope: wholeSubtree (2)
            derefAliases: derefAlways (3)
            sizeLimit: 0
            timeLimit: 0
            typesOnly: False
            Filter: (&(&(&(&(sAMAccountName=myuser))(objectclass=person))(objectclass=organizationalPerson))(objectclass=user))
            attributes: 11 items
                AttributeDescription: whenChanged
                AttributeDescription: whenCreated
                AttributeDescription: mail
                AttributeDescription: sAMAccountName
                AttributeDescription: objectGUID
                AttributeDescription: sn
                AttributeDescription: cn
                AttributeDescription: objectclass
                AttributeDescription: userPrincipalName
                AttributeDescription: userAccountControl
                AttributeDescription: pwdLastSet
    [Response In: 17]
    controls: 1 item


2nd request (new)


Plain Text
LDAPMessage searchRequest(4) "ou=users,dc=example,dc=test" wholeSubtree
    messageID: 4
    protocolOp: searchRequest (3)
        searchRequest
            baseObject: ou=users,dc=example,dc=test
            scope: wholeSubtree (2)
            derefAliases: derefAlways (3)
            sizeLimit: 0
            timeLimit: 0
            typesOnly: False
            Filter: (objectGUID=ed:cf:59:ef:22:bb:43:49:bd:93:c1:7c:6f:4a:51:d4)
            attributes: 11 items
                AttributeDescription: whenChanged
                AttributeDescription: whenCreated
                AttributeDescription: mail
                AttributeDescription: sAMAccountName
                AttributeDescription: objectGUID
                AttributeDescription: sn
                AttributeDescription: cn
                AttributeDescription: objectclass
                AttributeDescription: userPrincipalName
                AttributeDescription: userAccountControl
                AttributeDescription: pwdLastSet
    [Response In: 23]
    controls: 1 item



NOTE:  in this example I had modified the LDAP server schema by adding EQUALITY and the server DOES RESPOND with results to second request.

If i did not do that, the response to second request is empty (=no user matched) and Keycloak will throw exception because user with the objectGUID should exist


REST response contains one credential:


HTTP/1.1 200 OK
Cache-Control: no-cache
Content-Type: application/json;charset=UTF-8
Referrer-Policy: no-referrer
Strict-Transport-Security: max-age=31536000; includeSubDomains
X-Content-Type-Options: nosniff
X-Frame-Options: SAMEORIGIN
content-length: 65

[
    {
        "createdDate": -1,
        "federationLink": "my-ldap",
        "type": "password"
    }
]



#####



docker run --rm --name keycloak \
  --network host \
  -e KC_HOSTNAME=keycloak.127.0.0.1.nip.io \
  -e KC_BOOTSTRAP_ADMIN_USERNAME=admin \
  -e KC_BOOTSTRAP_ADMIN_PASSWORD=admin \
  quay.io/keycloak/keycloak:26.4.7 \
  start-dev --verbose


uses same LDAP filter as 26.2.x

(objectGUID=ed:cf:59:ef:22:bb:43:49:bd:93:c1:7c:6f:4a:51:d4)
