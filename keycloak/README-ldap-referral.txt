

https://github.com/keycloak/keycloak/issues/13500







# Compile only javascript part after updating UI code
cd js/apps/admin-ui
mvn install




# add ldap federation

vendor:             other
connection url:     ldap://localhost
bind dn:            cn=ldap-admin,ou=users,o=example
bind credentials:   ldap-admin
edit mode:          read-only
users dn:           ou=nonexisting,o=example
referral:           follow

*** Workaround

adding directory to classpath and creating file like

$ cat /my/properties/jndi.properties
java.naming.referral=follow

works as well
