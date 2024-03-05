

https://github.com/keycloak/keycloak/issues/13500
https://github.com/keycloak/keycloak/pull/24852







# Compile only javascript part after updating UI code
cd js/apps/admin-ui
mvn install





#
# Testing
#

docker-compose up postgres openldap

# launch keycloak via debugger
# Access console at http://keycloak.127.0.0.1.nip.io:8080/




# Add ldap federation in Admin GUI

vendor:             Other
connection url:     ldap://localhost
bind dn:            cn=ldap-admin,ou=users,o=example
bind credentials:   ldap-admin
edit mode:          READ_ONLY
users dn:           ou=nonexisting,o=example
referral:           follow



# Add LDAP federation via Admin REST API
TOKEN=$(http --form POST http://keycloak.127-0-0-121.nip.io/realms/master/protocol/openid-connect/token username=admin password=admin grant_type=password client_id=admin-cli | jq -r .access_token)
POST http://keycloak.127.0.0.1.nip.io:8080/admin/realms/master/components Authorization:"bearer $TOKEN"

{
  "config": {
    "enabled": [
      "true"
    ],
    "vendor": [
      "other"
    ],
    "connectionUrl": [
      "ldap://localhost"
    ],
    "connectionTimeout": [
      ""
    ],
    "bindDn": [
      "cn=ldap-admin,ou=users,o=example"
    ],
    "bindCredential": [
      "ldap-admin"
    ],
    "startTls": [
      "false"
    ],
    "useTruststoreSpi": [
      "always"
    ],
    "connectionPooling": [
      "false"
    ],
    "authType": [
      "simple"
    ],
    "usersDn": [
      "ou=nonexisting,o=example"
    ],
    "usernameLDAPAttribute": [
      "uid"
    ],
    "rdnLDAPAttribute": [
      "uid"
    ],
    "uuidLDAPAttribute": [
      "entryUUID"
    ],
    "userObjectClasses": [
      "inetOrgPerson, organizationalPerson"
    ],
    "customUserSearchFilter": [
      ""
    ],
    "readTimeout": [
      ""
    ],
    "editMode": [
      "READ_ONLY"
    ],
    "searchScope": [
      ""
    ],
    "pagination": [
      "false"
    ],
    "referral": [
      "follow"
    ],
    "batchSizeForSync": [
      ""
    ],
    "importEnabled": [
      "true"
    ],
    "syncRegistrations": [
      "true"
    ],
    "allowKerberosAuthentication": [
      "false"
    ],
    "useKerberosForPasswordAuthentication": [
      "false"
    ],
    "cachePolicy": [
      "DEFAULT"
    ],
    "usePasswordModifyExtendedOp": [
      "false"
    ],
    "validatePasswordPolicy": [
      "false"
    ],
    "trustEmail": [
      "false"
    ],
    "krbPrincipalAttribute": [
      "krb5PrincipalName"
    ],
    "changedSyncPeriod": [
      "-1"
    ],
    "fullSyncPeriod": [
      "-1"
    ]
  },
  "providerId": "ldap",
  "providerType": "org.keycloak.storage.UserStorageProvider",
  "parentId": "ee82be7a-d618-41d0-8894-58225b3b3612",
  "name": "ldap"
}



*** Workaround

adding directory to classpath and creating file like

$ cat /my/properties/jndi.properties
java.naming.referral=follow

works as well
