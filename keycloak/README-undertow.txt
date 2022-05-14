#
# Running Keycloak using embedded Undertow server (without application server)
#


export WORKDIR=/home/tsaarni/work/devenvs/keycloak


# using vscode
mkdir -p .vscode
cp $WORKDIR/configs/launch.json .vscode/


mvn clean install -DskipTests  # or alternatively: clean build


1. Choose debug or run (CTRL+F5) with "Debug KeycloakServer" selected
2. Click proceed on error "Build failed, do you want to continue?"





# using maven

mvn -f testsuite/utils/pom.xml exec:java -Pkeycloak-server -Dkeycloak.migration.action=import -Dkeycloak.migration.provider=dir -Dkeycloak.migration.dir=$WORKDIR/migrations/ldap-federation/ -Djavax.net.ssl.trustStore=$WORKDIR/truststore.p12 -Djavax.net.ssl.trustStorePassword=password -Djavax.net.ssl.javax.net.ssl.trustStoreType="PKCS12" -Djavax.net.ssl.keyStore=$WORKDIR/admin-keystore.p12 -Djavax.net.ssl.keyStorePassword=password -Djavax.net.ssl.javax.net.ssl.keyStoreType="PKCS12" -Dresources



mvnDebug -f testsuite/utils/pom.xml exec:java -Pkeycloak-server -Dkeycloak.migration.action=import -Dkeycloak.migration.provider=dir -Dkeycloak.migration.dir=$WORKDIR/migrations/ldap-federation/ -Djavax.net.ssl.trustStore=$WORKDIR/truststore.p12 -Djavax.net.ssl.trustStorePassword=password -Djavax.net.ssl.javax.net.ssl.trustStoreType="PKCS12" -Djavax.net.ssl.keyStore=$WORKDIR/admin-keystore.p12 -Djavax.net.ssl.keyStorePassword=password -Djavax.net.ssl.javax.net.ssl.keyStoreType="PKCS12" -Dresources











TOKEN=$(http --form POST http://localhost:8081/auth/realms/master/protocol/openid-connect/token username=admin password=admin grant_type=password client_id=admin-cli | jq -r .access_token)

http -v POST http://localhost:8081/auth/admin/realms/ Authorization:"bearer $TOKEN" id=my-realm realm=my-realm
http -v DELETE http://localhost:8081/auth/admin/realms/my-realm Authorization:"bearer $TOKEN"
http -v GET http://localhost:8081/auth/admin/realms/master/admin-events Authorization:"bearer $TOKEN"














### Connect to H2 database



To access H2 database

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



connect by using

java -cp /home/tsaarni/.m2/repository/com/h2database/h2/1.4.197/h2-1.4.197.jar org.h2.tools.Console -url "jdbc:h2:$(ls -dt /tmp/keycloak-server-* | head -1)/data/test;AUTO_SERVER=TRUE" -user SA

browser will open automatically
http://127.0.1.1:41337/frame.jsp?jsessionid=c134b43d61365d4b6214845902e2e6f6



### H2 wrong user name or password

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

remove old h2 file database ~/data/
