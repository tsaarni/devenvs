
https://gist.github.com/tsaarni/4297da5db37d5c373b5133bf0d9977ce

# Cipher configuration in Keycloak

## Infinispan & JGroups

Infinispan is utilized by Keycloak as an in-memory cache for storing user sessions and other data,
with JGroups managing commuication between Keycloak cluster instances.

The predefined JGroups transport stacks [1] are set by the Infinispan [2], though user can define custom transport stacks as well [3].
These stack definitions do not inherently include TLS settings; instead, TLS configurations are specified through Keycloak settings [4].
Keycloak's code integrates this configuration with the transport stack setup [5].
Although JGroups offers a method to set the cipher suites [6], Keycloak does not make use of it.

To enable the configuration of cipher suites, changes to the code are required, along with the introduction of a new configuration parameter.

[1] https://www.keycloak.org/server/caching#_transport_stacks
[2] https://github.com/infinispan/infinispan/tree/b7c63784aea1ad719a6fd97d8b9c3284a6364c19/core/src/main/resources/default-configs
[3] https://www.keycloak.org/server/caching#_custom_transport_stacks
[4] https://www.keycloak.org/server/caching#_securing_cache_communication
[5] https://github.com/keycloak/keycloak/blob/3d340d17a4891140219c01315755380b7f7da898/quarkus/runtime/src/main/java/org/keycloak/quarkus/runtime/storage/legacy/infinispan/CacheManagerFactory.java#L370-L394
[6] https://github.com/belaban/JGroups/blob/f565afebec5e0d5d8fb36094753f36335813d14e/src/org/jgroups/util/TLS.java#L39-L40


## pgJDBC

Keycloak uses pgJDBC to establish connections with PostgreSQL databases [1].
This is implemented via Quarkus Datasources, which incorporates Hibernate ORM and the Agroal connection pool [2].
While Keycloak's database URL can be configured to include TLS connection properties, pgJDBC itself does not support cipher suite customization [3].

To facilitate cipher suite configuration, Keycloak needs to implement a custom socket factory for pgJDBC, and a new configuration parameter must be introduced [4].
Other database drivers might have different approach than pgJDBC so this is not generic solution.
Alternatively pgJDBC could be modified to provide connection property for cipher suite selection.

[1] https://www.keycloak.org/server/db
[2] https://quarkus.io/guides/datasource
[3] https://jdbc.postgresql.org/documentation/use/#connection-parameters/
[4] https://jdbc.postgresql.org/documentation/ssl/#custom-sslsocketfactory



## LDAP client

Keycloak utilizes the JDK LDAP client to connect to LDAP servers, and the JDK provides the ability to define custom socket factories for TLS configuration [1].
Keycloak implements a custom LDAP socket factory, but it currently lacks support for customizing cipher suites [2].

To allow cipher suite configurability, Keycloak's LDAP socket factory needs to be extended, and a corresponding configuration parameter should be added.


[1] https://docs.oracle.com/javase/jndi/tutorial/ldap/security/ssl.html
[2] https://github.com/keycloak/keycloak/blob/3d340d17a4891140219c01315755380b7f7da898/services/src/main/java/org/keycloak/truststore/SSLSocketFactory.java


## HTTP Client

Keycloak leverages Apache HttpClient for making HTTP(S) requests to external services.
While Apache HttpClient supports defining cipher suites [1], this feature is not currently used by Keycloak [2].

To enable this functionality, the code must be updated, and a new configuration parameter needs to be exposed.

[1] https://hc.apache.org/httpcomponents-client-4.5.x/current/httpclient/apidocs/org/apache/http/conn/ssl/SSLConnectionSocketFactory.html?is-external=true
[2] https://github.com/keycloak/keycloak/blob/3d340d17a4891140219c01315755380b7f7da898/services/src/main/java/org/keycloak/connections/httpclient/HttpClientBuilder.java#L232-L250

## Vert.x

Keycloak uses Vert.x as an HTTP(S) server through Quarkus.
Keycloak already provides a configuration option for HTTPS cipher suites, which is passed to Vert.x [1].

No further action is required.


[1]  https://www.keycloak.org/server/enabletls
