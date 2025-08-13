
cd ~/work/devenvs/keycloak
docker compose up



                "KC_DB_URL": "jdbc:postgresql://localhost:5432/keycloak?connectTimeout=10",
//                "KC_DB_URL": "jdbc:postgresql://localhost:5432/keycloak",


function get_admin_token() {
  http --form POST http://localhost:8080/realms/master/protocol/openid-connect/token username=admin password=admin grant_type=password client_id=admin-cli | jq -r .access_token
}


# List users in Keycloak
time http GET http://localhost:8080/admin/realms/master/users "Authorization:Bearer $(get_admin_token)"

# Create a new user in Keycloak
time http POST http://localhost:8080/admin/realms/master/users "Authorization:Bearer $(get_admin_token)" username=foo1



# Enter into network namespace of postgres container and block all packets between Keycloak and PostgreSQL
sudo nsenter -t $(docker inspect -f '{{.State.Pid}}' keycloak-postgres-1) --net -- tc qdisc add dev eth0 root netem loss 100%

# Allow packets again
sudo nsenter -t $(docker inspect -f '{{.State.Pid}}' keycloak-postgres-1) --net -- tc qdisc del dev eth0 root

lsof -iTCP:5432

docker exec -it keycloak-postgres-1 bash
ps -ef | awk '/postgres: keycloak keycloak/ {print $1}' | xargs kill -9
docker restart keycloak-postgres-1


wireshark -k -i loopback --display-filter "tcp.port == 5432"

wireshark -k -i loopback --display-filter "tcp.port == 5432 && (tcp.flags.syn == 1 || tcp.flags.ack == 1)"









######



cd ~/work/pgjdbc
###git checkout REL42.7.7
./gradlew build -x test && cp ./pgjdbc/build/libs/postgresql-42.7.7-SNAPSHOT-all.jar /home/tsaarni/.m2/repository/org/postgresql/postgresql/42.7.7/postgresql-42.7.7.jar





	... 17 more
Running the server in development mode. DO NOT use this configuration in production.
2025-08-01 17:01:31,220 INFO  [org.keycloak.url.HostnameV2ProviderFactory] (Quarkus Main Thread) If hostname is specified, hostname-strict is effectively ignored
2025-08-01 17:01:32,481 WARN  [io.smallrye.config] (Quarkus Main Thread) SRCFG01008: The value default has been converted by a Boolean Converter to "false"
2025-08-01 17:01:32,482 WARN  [io.smallrye.config] (Quarkus Main Thread) SRCFG01008: The value default has been converted by a Boolean Converter to "false"
2025-08-01 17:01:32,725 WARNING [org.postgresql.Driver] (JPA Startup Thread) Processing option [connectTimeout=10]
2025-08-01 17:01:32,757 WARNING [org.postgresql.Driver] (agroal-11) Processing option [connectTimeout=10]
2025-08-01 17:01:32,767 WARNING [org.postgresql.core.v3.ConnectionFactoryImpl] (agroal-11) ConnectTimeout is 10000 milliseconds
2025-08-01 17:01:32,771 WARNING [org.postgresql.ssl.MakeSSL] (agroal-11) Connecting to localhost/127.0.0.1:5432 with timeout 10,000 ms
2025-08-01 17:01:32,775 WARNING [org.postgresql.ssl.MakeSSL] (agroal-11) PGStream Setting network timeout to 5,000 milliseconds
2025-08-01 17:01:32,776 WARNING [org.postgresql.ssl.MakeSSL] (agroal-11) PGStream Setting network timeout to 0 milliseconds
2025-08-01 17:01:33,922 INFO  [org.keycloak.spi.infinispan.impl.embedded.JGroupsConfigurator] (Quarkus Main Thread) JGroups JDBC_PING discovery enabled.
2025-08-01 17:01:33,931 INFO  [org.keycloak.spi.infinispan.impl.embedded.JGroupsConfigurator] (Quarkus Main Thread) JGroups Encryption enabled (mTLS).
2025-08-01 17:01:34,040 INFO  [org.keycloak.jgroups.certificates.CertificateReloadManager] (Quarkus Main Thread) Starting JGroups certificate reload manager
2025-08-01 17:01:34,111 INFO  [org.infinispan.CONTAINER] (Quarkus Main Thread) ISPN000556: Starting user marshaller 'org.infinispan.commons.marshall.ImmutableProtoStreamMarshaller'
2025-08-01 17:01:34,350 INFO  [org.keycloak.connections.infinispan.DefaultInfinispanConnectionProviderFactory] (Quarkus Main Thread) Node name: node_386656, Site name: null
2025-08-01 17:01:34,591 WARNING [org.postgresql.Driver] (agroal-11) Processing option [connectTimeout=10]
2025-08-01 17:01:34,591 WARNING [org.postgresql.core.v3.ConnectionFactoryImpl] (agroal-11) ConnectTimeout is 10000 milliseconds
2025-08-01 17:01:34,591 WARNING [org.postgresql.ssl.MakeSSL] (agroal-11) Connecting to localhost/127.0.0.1:5432 with timeout 10,000 ms
2025-08-01 17:01:34,591 WARNING [org.postgresql.ssl.MakeSSL] (agroal-11) PGStream Setting network timeout to 5,000 milliseconds
2025-08-01 17:01:34,592 WARNING [org.postgresql.ssl.MakeSSL] (agroal-11) PGStream Setting network timeout to 0 milliseconds
2025-08-01 17:01:34,950 INFO  [io.quarkus] (Quarkus Main Thread) Keycloak 1 on JVM (powered by Quarkus 3.24.5) started in 3.982s. Listening on: http://0.0.0.0:8080
2025-08-01 17:01:34,951 INFO  [io.quarkus] (Quarkus Main Thread) Profile dev activated.
2025-08-01 17:01:34,951 INFO  [io.quarkus] (Quarkus Main Thread) Installed features: [agroal, cdi, hibernate-orm, jdbc-h2, jdbc-mariadb, jdbc-mssql, jdbc-mysql, jdbc-oracle, jdbc-postgresql, keycloak, micrometer, narayana-jta, opentelemetry, reactive-routes, rest, rest-jackson, smallrye-context-propagation, smallrye-health, vertx]





Running the server in development mode. DO NOT use this configuration in production.
2025-08-01 17:02:25,279 INFO  [org.keycloak.url.HostnameV2ProviderFactory] (Quarkus Main Thread) If hostname is specified, hostname-strict is effectively ignored
2025-08-01 17:02:26,159 WARN  [io.smallrye.config] (Quarkus Main Thread) SRCFG01008: The value default has been converted by a Boolean Converter to "false"
2025-08-01 17:02:26,160 WARN  [io.smallrye.config] (Quarkus Main Thread) SRCFG01008: The value default has been converted by a Boolean Converter to "false"
2025-08-01 17:02:26,371 WARNING [org.postgresql.Driver] (JPA Startup Thread) Processing option []
2025-08-01 17:02:26,403 WARNING [org.postgresql.Driver] (agroal-11) Processing option []
2025-08-01 17:02:26,416 WARNING [org.postgresql.core.v3.ConnectionFactoryImpl] (agroal-11) ConnectTimeout is 10000 milliseconds
2025-08-01 17:02:26,421 WARNING [org.postgresql.ssl.MakeSSL] (agroal-11) Connecting to localhost/127.0.0.1:5432 with timeout 10,000 ms
2025-08-01 17:02:26,426 WARNING [org.postgresql.ssl.MakeSSL] (agroal-11) PGStream Setting network timeout to 5,000 milliseconds
2025-08-01 17:02:26,427 WARNING [org.postgresql.ssl.MakeSSL] (agroal-11) PGStream Setting network timeout to 0 milliseconds
2025-08-01 17:02:27,534 INFO  [org.keycloak.spi.infinispan.impl.embedded.JGroupsConfigurator] (Quarkus Main Thread) JGroups JDBC_PING discovery enabled.
2025-08-01 17:02:27,543 INFO  [org.keycloak.spi.infinispan.impl.embedded.JGroupsConfigurator] (Quarkus Main Thread) JGroups Encryption enabled (mTLS).
2025-08-01 17:02:27,680 INFO  [org.keycloak.jgroups.certificates.CertificateReloadManager] (Quarkus Main Thread) Starting JGroups certificate reload manager
2025-08-01 17:02:27,704 INFO  [org.infinispan.CONTAINER] (Quarkus Main Thread) ISPN000556: Starting user marshaller 'org.infinispan.commons.marshall.ImmutableProtoStreamMarshaller'
2025-08-01 17:02:27,958 INFO  [org.keycloak.connections.infinispan.DefaultInfinispanConnectionProviderFactory] (Quarkus Main Thread) Node name: node_458358, Site name: null
2025-08-01 17:02:28,236 WARNING [org.postgresql.Driver] (agroal-11) Processing option []
2025-08-01 17:02:28,237 WARNING [org.postgresql.core.v3.ConnectionFactoryImpl] (agroal-11) ConnectTimeout is 10000 milliseconds
2025-08-01 17:02:28,237 WARNING [org.postgresql.ssl.MakeSSL] (agroal-11) Connecting to localhost/127.0.0.1:5432 with timeout 10,000 ms
2025-08-01 17:02:28,238 WARNING [org.postgresql.ssl.MakeSSL] (agroal-11) PGStream Setting network timeout to 5,000 milliseconds
2025-08-01 17:02:28,238 WARNING [org.postgresql.ssl.MakeSSL] (agroal-11) PGStream Setting network timeout to 0 milliseconds
2025-08-01 17:02:28,561 INFO  [io.quarkus] (Quarkus Main Thread) Keycloak 1 on JVM (powered by Quarkus 3.24.5) started in 3.803s. Listening on: http://0.0.0.0:8080
2025-08-01 17:02:28,561 INFO  [io.quarkus] (Quarkus Main Thread) Profile dev activated.
2025-08-01 17:02:28,561 INFO  [io.quarkus] (Quarkus Main Thread) Installed features: [agroal, cdi, hibernate-orm, jdbc-h2, jdbc-mariadb, jdbc-mssql, jdbc-mysql, jdbc-oracle, jdbc-postgresql, keycloak, micrometer, narayana-jta, opentelemetry, reactive-routes, rest, rest-jackson, smallrye-context-propagation, smallrye-health, vertx]
