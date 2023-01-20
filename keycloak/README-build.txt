# Building


# compile distribution
mvn -Pdistribution -DskipTests clean install

ls -l ./quarkus/dist/target/keycloak-*.gz ./quarkus/dist/target/keycloak-*.zip


# OLDOLDOLD
# ./distribution/server-dist/target/keycloak-legacy-19.0.3.tar.gz

cd quarkus/dist/target/
tar zxvf keycloak-*.gz
cd keycloak-20.0.2/
bin/kc.sh




# Inspect dependencies
mvn dependency:tree -Pdistribution    # Dependency tree
mvn dependency:tree -Pdistribution -Dincludes=jakarta.xml.bind:jakarta.xml.bind-api   # Dedendency on particular package




*** Debugging

mkdir -p .vscode
cp -a ~/work/devenvs/keycloak/configs/launch.json .vscode


mvn clean install -f testsuite/integration-arquillian/pom.xml -DforkMode=never -Dmaven.surefire.debug  ...   # attach to port 5005 (not 8000)






#### Troubleshooting

# To access H2 database

--- a/testsuite/utils/src/main/resources/META-INF/keycloak-server.json
+++ b/testsuite/utils/src/main/resources/META-INF/keycloak-server.json
@@ -96,7 +96,7 @@

     "connectionsJpa": {
         "default": {
-            "url": "${keycloak.connectionsJpa.url:jdbc:h2:mem:test;DB_CLOSE_DELAY=-1}",
+            "url": "${keycloak.connectionsJpa.url:jdbc:h2:${jboss.server.data.dir}/test;DB_CLOSE_DELAY=-1;AUTO_SERVER=TRUE;AUTO_SERVER_PORT=9090}",
             "driver": "${keycloak.connectionsJpa.driver:org.h2.Driver}",
             "driverDialect": "${keycloak.connectionsJpa.driverDialect:}",
             "user": "${keycloak.connectionsJpa.user:sa}",



# H2 wrong user name or password

2022-03-08 09:49:47,422 WARN  [io.agroal.pool] (agroal-11) Datasource '<default>': Wrong user name or password [28000-197]
2022-03-08 09:49:47,423 DEBUG [io.agroal.pool] (agroal-11) Cause: : org.h2.jdbc.JdbcSQLException: Wrong user name or password [28000-197]
        at org.h2.message.DbException.getJdbcSQLException(DbException.java:357)
        at org.h2.message.DbException.get(DbException.java:179)
        at org.h2.message.DbException.get(DbException.java:155)
        at org.h2.message.DbException.get(DbException.java:144)
        at org.h2.engine.Engine.validateUserAndPassword(Engine.java:341)
        at org.h2.engine.Engine.createSessionAndValidate(Engine.java:165)
        at org.h2.engine.Engine.createSession(Engine.java:140)
        at org.h2.engine.Engine.createSession(Engine.java:28)
        at org.h2.engine.SessionRemote.connectEmbeddedOrServer(SessionRemote.java:351)
        at org.h2.jdbc.JdbcConnection.<init>(JdbcConnection.java:124)
        at org.h2.jdbc.JdbcConnection.<init>(JdbcConnection.java:103)
        at org.h2.Driver.connect(Driver.java:69)
        at org.h2.jdbcx.JdbcDataSource.getJdbcConnection(JdbcDataSource.java:189)
        at org.h2.jdbcx.JdbcDataSource.getXAConnection(JdbcDataSource.java:352)
        at io.agroal.pool.ConnectionFactory.createConnection(ConnectionFactory.java:216)
        at io.agroal.pool.ConnectionPool$CreateConnectionTask.call(ConnectionPool.java:513)
        at io.agroal.pool.ConnectionPool$CreateConnectionTask.call(ConnectionPool.java:494)
        at java.base/java.util.concurrent.FutureTask.run(FutureTask.java:264)
        at io.agroal.pool.util.PriorityScheduledExecutor.beforeExecute(PriorityScheduledExecutor.java:75)
        at java.base/java.util.concurrent.ThreadPoolExecutor.runWorker(ThreadPoolExecutor.java:1126)
        at java.base/java.util.concurrent.ThreadPoolExecutor$Worker.run(ThreadPoolExecutor.java:628)
        at java.base/java.lang.Thread.run(Thread.java:829)



Solution:

remove old h2 file database ~/data/ or target/kc/data/h2/
rm /home/tsaarni/work/keycloak/target/kc/data/h2/*
