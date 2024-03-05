

https://groups.google.com/g/keycloak-dev/c/lq_Q_GZPwUI/m/rASBroloAwAJ

### In docker-compose.yml
# check that statement logging is enabled
    command: -c log_statement=all -c log_destination=stderr

# check that database is persisted outside of the container
      - ./tmp:/var/lib/postgresql/data



docker-compose up postgres


### In .vscode/launch.json

# check that keycloak logging is set to debug
"KC_LOG_LEVEL": "debug",

# check that hibernate ORM logging is true
"QUARKUS_HIBERNATE_ORM_LOG_SQL": "true",





# Test admin login
http --form POST http://keycloak.127-0-0-1.nip.io:8080/realms/master/protocol/openid-connect/token username=admin password=admin grant_type=password client_id=admin-cli



TOKEN=$(http --form POST http://keycloak.127-0-0-1.nip.io:8080/realms/master/protocol/openid-connect/token username=admin password=admin grant_type=password client_id=admin-cli | jq -r .access_token)


http -v POST http://keycloak.127-0-0-1.nip.io:8080/admin/realms/master/users Authorization:"bearer $TOKEN" username=foo










#################
#
# Maybe related
#

https://github.com/quarkusio/quarkus/issues/16198


https://smallrye.io/docs/smallrye-fault-tolerance/6.2.0/reference/retry.html
https://vladmihalcea.com/optimistic-locking-retry-with-jpa/


Retry if a connection becomes invalid in some cases
https://hibernate.atlassian.net/jira/software/c/projects/HHH/issues/HHH-7511

> Hibernate ORM can't really do that because the connection might be associated to a running transaction and it's impossible to re-run all the actions that have ran before a failure on a new connection. > Even if it were possible, it might simply be wrong.
> Use connection validation and/or implement a retry mechanism on your application layer in case of "transient errors".




Hibernate transaction can't recover from a DB restart
https://hibernate.atlassian.net/jira/software/c/projects/HHH/issues/HHH-12513



Transaction documentation?
https://groups.google.com/g/keycloak-dev/c/TJswilRe4mw/m/fRi3MdmWCQAJ
