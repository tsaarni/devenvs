


# build and install org.keycloak modules into maven cache
mvn clean install -DskipTestsuite -DskipExamples -DskipTests

# after main codebase is built, to build the quarkus distribution
mvn -f quarkus/pom.xml clean install -DskipTests



#
# Test server certificate rotation
#

rm -rf certs
mkdir -p certs
certyaml --destination certs configs/certs.yaml
openssl pkcs12 -export -passout pass:secret -noiter -nomaciter -in certs/keycloak-server.pem -inkey certs/keycloak-server-key.pem -out certs/keycloak-server.p12


export WORKDIR=/home/tsaarni/work/devenvs/keycloak

java -jar quarkus/server/target/lib/quarkus-run.jar --verbose start --hostname=keycloak.127-0-0-1.nip.io --https-key-store-file=$WORKDIR/certs/keycloak-server.p12 --https-key-store-password=secret

# check the expiration date of current server certificate
echo Q | openssl s_client -connect keycloak.127-0-0-1.nip.io:8443 | openssl x509 -text -noout

# create new server certificate and keystore
rm certs/keycloak-server*
certyaml --destination certs configs/certs.yaml
openssl pkcs12 -export -passout pass:secret -noiter -nomaciter -in certs/keycloak-server.pem -inkey certs/keycloak-server-key.pem -out certs/keycloak-server.p12

# check the expiration date of current server certificate
echo Q | openssl s_client -connect keycloak.127-0-0-1.nip.io:8443 | openssl x509 -text -noout







### OLDOLDOLD

diff --git a/pom.xml b/pom.xml
index 177bf6bd2c..b3bff9c71a 100644
--- a/pom.xml
+++ b/pom.xml
@@ -41,7 +41,7 @@
         <jboss.snapshots.repo.id>jboss-snapshots-repository</jboss.snapshots.repo.id>
         <jboss.snapshots.repo.url>https://s01.oss.sonatype.org/content/repositories/snapshots/</jboss.snapshots.repo.url>

-        <quarkus.version>2.7.5.Final</quarkus.version>
+        <quarkus.version>999-SNAPSHOT</quarkus.version>

         <!--
         Performing a Wildfly upgrade? Run the:
@@ -82,7 +82,7 @@
         <dom4j.version>2.1.3</dom4j.version>
         <h2.version>1.4.197</h2.version>
         <jakarta.persistence.version>2.2.3</jakarta.persistence.version>
-        <hibernate.core.version>5.3.24.Final</hibernate.core.version>
+        <hibernate.core.version>5.6.8.Final</hibernate.core.version>
         <hibernate.c3p0.version>5.3.24.Final</hibernate.c3p0.version>
         <infinispan.version>12.1.7.Final</infinispan.version>
         <infinispan.protostream.processor.version>4.4.1.Final</infinispan.protostream.processor.version>
@@ -283,6 +283,12 @@
     <dependencyManagement>

         <dependencies>
+<!-- https://mvnrepository.com/artifact/io.vertx/vertx-core -->
+<dependency>
+    <groupId>io.vertx</groupId>
+    <artifactId>vertx-core</artifactId>
+    <version>4.3.1-SNAPSHOT</version>
+</dependency>

             <dependency>
                 <groupId>org.keycloak</groupId>
diff --git a/quarkus/pom.xml b/quarkus/pom.xml
index a2010602e3..c2031469ab 100644
--- a/quarkus/pom.xml
+++ b/quarkus/pom.xml
@@ -40,7 +40,7 @@
         <resteasy.version>4.7.5.Final</resteasy.version>
         <jackson.version>2.13.2</jackson.version>
         <jackson.databind.version>2.13.2.2</jackson.databind.version>
-        <hibernate.core.version>5.6.5.Final</hibernate.core.version>
+        <hibernate.core.version>5.6.8.Final</hibernate.core.version>
         <mysql.driver.version>8.0.28</mysql.driver.version>
         <postgresql.version>42.3.3</postgresql.version>
         <microprofile-metrics-api.version>3.0.1</microprofile-metrics-api.version>
diff --git a/quarkus/runtime/src/test/java/org/keycloak/quarkus/runtime/configuration/test/ConfigurationTest.java b/quarkus/runtime/src/test/java/org/keycloak/quarkus/runtime/configuration/test/ConfigurationTest.java
index 072db49bd7..12e7c05092 100644
--- a/quarkus/runtime/src/test/java/org/keycloak/quarkus/runtime/configuration/test/ConfigurationTest.java
+++ b/quarkus/runtime/src/test/java/org/keycloak/quarkus/runtime/configuration/test/ConfigurationTest.java
@@ -45,7 +45,7 @@ import io.quarkus.runtime.configuration.ConfigUtils;
 import io.smallrye.config.SmallRyeConfigProviderResolver;
 import org.keycloak.quarkus.runtime.Environment;
 import org.keycloak.quarkus.runtime.vault.FilesPlainTextVaultProviderFactory;
-import org.mariadb.jdbc.MySQLDataSource;
+import org.mariadb.jdbc.MariaDbDataSource;
 import org.postgresql.xa.PGXADataSource;

 public class ConfigurationTest {
@@ -325,7 +325,7 @@ public class ConfigurationTest {
         config = createConfig();
         assertEquals("jdbc:mariadb://localhost:3306/keycloak?test=test&test1=test1", config.getConfigValue("quarkus.datasource.jdbc.url").getValue());
         assertEquals(MariaDBDialect.class.getName(), config.getConfigValue("quarkus.hibernate-orm.dialect").getValue());
-        assertEquals(MySQLDataSource.class.getName(), config.getConfigValue("quarkus.datasource.jdbc.driver").getValue());
+        assertEquals(MariaDbDataSource.class.getName(), config.getConfigValue("quarkus.datasource.jdbc.driver").getValue());

         System.setProperty(CLI_ARGS, "--db=postgres");
         config = createConfig();
