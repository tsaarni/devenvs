https://www.keycloak.org/docs/latest/upgrading/index.html


### Postgres clone

# pg_dump
https://stackoverflow.com/a/1238305/458531

# copy to separate logical db
https://www.postgresql.org/docs/current/manage-ag-templatedbs.html



### Keycloak HA / multi-site

https://www.keycloak.org/high-availability/introduction
https://github.com/keycloak/keycloak/discussions/25269


### Using export / import (instead of database backup)


https://www.keycloak.org/server/importExport


# example of export maybe being incomplete sometimes
Declarative User Profile export
https://github.com/keycloak/keycloak/pull/24147



# See also note https://www.keycloak.org/docs/latest/server_admin/#using-the-admin-console
Export files from the Admin Console are not suitable for backups or data transfer between servers. Only boot-time exports are suitable for backups or data transfer between servers.





### migrations

# automatic or manual (via kc.sh script)
https://www.keycloak.org/docs/latest/upgrading/index.html#_migrate_db


# liquibase rollbacks (not used in keycloak)
https://www.baeldung.com/liquibase-rollback
