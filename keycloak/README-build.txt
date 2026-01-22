
# Build
mvn clean install -DskipTestsuite -DskipExamples -DskipTests

# Parallel build
mvnd clean install -DskipTestsuite -DskipExamples -DskipTests
mvn -T4C clean install -DskipTestsuite -DskipExamples -DskipTests


# distribution
mvn clean install -Pdistribution -DskipTests
ls -l ./quarkus/dist/target/

# Compile just server distribution
###   ls -l quarkus/dist/target/keycloak-*.gz quarkus/dist/target/keycloak-*.zip
mvnd -pl quarkus/deployment,quarkus/dist -am -DskipTests clean install




# Run in dev mode
java -jar quarkus/server/target/lib/quarkus-run.jar start-dev

# Run distro
cd quarkus/dist/target/
tar zxvf keycloak-*.gz
cd keycloak-20.0.2/
bin/kc.sh



*** Build custom docker container



cd /home/tsaarni/work/keycloak-worktree/26.2.10-nordix
mvnd clean install -DskipTestsuite -DskipExamples -DskipTests
cp ./quarkus/dist/target/keycloak-26.2.10.tar.gz quarkus/container
docker build --build-arg KEYCLOAK_DIST=keycloak-26.2.10.tar.gz -f quarkus/container/Dockerfile -t localhost/keycloak:latest quarkus/container



*** Inspect dependencies

mvn dependency:tree -Pdistribution    # Dependency tree
mvn dependency:tree -Pdistribution -Dincludes=jakarta.xml.bind:jakarta.xml.bind-api   # Dedendency on particular package



*** Running unit tests

mvn clean install -DskipTests
(cd distribution; mvn clean install)

mvn clean install -f testsuite/integration-arquillian/pom.xml -Dtest=org.keycloak.testsuite.federation.storage.UserStorageDirtyDeletionUnsyncedImportTest#testMembersWhenCachedUsersRemovedFromBackend -Dkeycloak.logging.level=debug


mvn -Dtest=TestCircle#mytest                                  # run single test case
mvn -Dtest=org.keycloak.testsuite.federation.ldap.*AuthTest   # run matching test classes
mvn -Dtest=org.keycloak.testsuite.federation.ldap.**          # all test cases in package (recursively)



*** Debugging in vscode


# Debug directly from vscode

mkdir -p .vscode
cp -a ~/work/devenvs/keycloak/configs/launch.json  ~/work/devenvs/keycloak/configs/settings.json .vscode

##  1. Build with maven
##  2. Start vscode, wait for build
##  3. run build again without clean
##       mvn install -DskipTestsuite -DskipExamples -DskipTests
##  3. Launch the debug session


*** UI development


1. Add following to launch.json
   "KC_ADMIN_VITE_URL": "http://localhost:5174",
2. remove KC_HOSTNAME from launch.json if it exists
3. start keycloak in vscode debugger
4. run ui under vite
   cd js
   pnpm --filter keycloak-admin-ui run dev
5. access keycloak http://127.0.0.1:8080/



Do you get:  TypeError: crypto.hash is not a function
Workaround:
export VOLTA_FEATURE_PNPM=1
volta install pnpm


*****



# Sometimes keycloak does not start, it only prints:
Press [e] to edit command line args (currently 'start-dev --hostname=keycloak.127-0-0-1.nip.io'), [h] for more options>

select "Java: Clean Java language server workspace"



# Remote debug
mvn clean install -f testsuite/integration-arquillian/pom.xml -DforkMode=never -Dmaven.surefire.debug  ...   # attach to port 5005 (not 8000)







mvn -f js/pom.xml install
mvn -f themes/pom.xml install





2023-10-27 13:10:56,152 ERROR [io.qua.dep.dev.IsolatedDevModeMain] (main) Failed to start quarkus: java.lang.RuntimeException: io.quarkus.builder.BuildException: Build failure: Build failed due to errors
        [error]: Build step org.keycloak.quarkus.deployment.KeycloakProcessor#produceDefaultPersistenceUnit threw an exception: java.lang.NullPointerException: Cannot invoke "java.net.URL.toExternalForm()" because "xmlUrl" is null
        at org.hibernate.jpa.boot.internal.PersistenceXmlParser.parsePersistenceXml(PersistenceXmlParser.java:243)
...
Caused by: io.quarkus.builder.BuildException: Build failure: Build failed due to errors
        [error]: Build step org.keycloak.quarkus.deployment.KeycloakProcessor#produceDefaultPersistenceUnit threw an exception: java.lang.NullPointerException: Cannot invoke "java.net.URL.toExternalForm()" because "xmlUrl" is null
        at org.hibernate.jpa.boot.internal.PersistenceXmlParser.parsePersistenceXml(PersistenceXmlParser.java:243)
...
Caused by: java.lang.NullPointerException: Cannot invoke "java.net.URL.toExternalForm()" because "xmlUrl" is null
        at org.hibernate.jpa.boot.internal.PersistenceXmlParser.parsePersistenceXml(PersistenceXmlParser.java:243)
        ...
        at org.keycloak.quarkus.deployment.KeycloakProcessor.produceDefaultPersistenceUnit(KeycloakProcessor.java:314)


Code:

        if (storage == null) {
            descriptor = PersistenceXmlParser.locateIndividualPersistenceUnit(
                    Thread.currentThread().getContextClassLoader().getResource("default-persistence.xml"));


# Missing resource: ./model/jpa/src/main/resources/default-persistence.xml

mvn -f model/pom.xml install -DskipTests




2023-10-27 13:37:31,544 ERROR [org.keycloak.services.error.KeycloakErrorHandler] (executor-thread-4) Uncaught server error: org.keycloak.theme.FreeMarkerException: Failed to process template index.ftl
...

Caused by: freemarker.template.TemplateNotFoundException: Template not found for name "index.ftl".


mvn -f js/pom.xml install
mvn -f themes/pom.xml install




2023-11-23 13:46:07,333 ERROR [io.qua.dep.dev.IsolatedDevModeMain] (main) Failed to start quarkus: java.lang.RuntimeException: io.quarkus.builder.BuildException: Build failure: Build failed due to errors
	[error]: Build step org.keycloak.quarkus.deployment.CacheBuildSteps#configureInfinispan threw an exception: java.lang.IllegalArgumentException: Option 'configFile' needs to be specified
	at org.keycloak.quarkus.deployment.CacheBuildSteps.configureInfinispan(CacheBuildSteps.java:80)
CacheBuildSteps.java:80
	at java.base/jdk.internal.reflect.NativeMethodAccessorImpl.invoke0(Native Method)
	at java.base/jdk.internal.reflect.NativeMethodAccessorImpl.invoke(NativeMethodAccessorImpl.java:77)
NativeMethodAccessorImpl.java:77
	at java.base/jdk.internal.reflect.DelegatingMethodAccessorImpl.invoke(DelegatingMethodAccessorImpl.java:43)
...



# Missing resource: cache-local.xml
mvn -f quarkus/runtime/pom.xml install -DskipTests



2023-11-23 14:04:51,718 ERROR [io.quarkus.runner.bootstrap.StartupActionImpl] (Quarkus Main Thread) Error running Quarkus: java.lang.reflect.InvocationTargetException
	at java.base/jdk.internal.reflect.NativeMethodAccessorImpl.invoke0(Native Method)
	at java.base/jdk.internal.reflect.NativeMethodAccessorImpl.invoke(NativeMethodAccessorImpl.java:77)
NativeMethodAccessorImpl.java:77
	at java.base/jdk.internal.reflect.DelegatingMethodAccessorImpl.invoke(DelegatingMethodAccessorImpl.java:43)
DelegatingMethodAccessorImpl.java:43
	at java.base/java.lang.reflect.Method.invoke(Method.java:568)
Method.java:568
	at io.quarkus.runner.bootstrap.StartupActionImpl$1.run(StartupActionImpl.java:104)
StartupActionImpl.java:104
	at java.base/java.lang.Thread.run(Thread.java:833)
Thread.java:833
Caused by: java.lang.ExceptionInInitializerError
	at org.keycloak.quarkus.runtime.KeycloakMain.main(KeycloakMain.java:66)
KeycloakMain.java:66
	... 6 more
Caused by: java.lang.NullPointerException: inStream parameter is null
	at java.base/java.util.Objects.requireNonNull(Objects.java:235)
Objects.java:235
	at java.base/java.util.Properties.load(Properties.java:407)
Properties.java:407
	at org.keycloak.common.Version.<clinit>(Version.java:39)


mvn install -f common/pom.xml install -DskipTests



mvn install -f quarkus/pom.xml install -DskipTests





# when having errors running in vscode, test also via command line:
mvn -f quarkus/server/pom.xml compile quarkus:dev -Dquarkus.args="start-dev"

# vs

java -jar quarkus/server/target/lib/quarkus-run.jar start-dev






*** Debug in quarkus dev mode


mkdir -p .vscode
cp -a ~/work/devenvs/keycloak/configs/launch.json .vscode


# compile keycloak

KEYCLOAK_ADMIN=admin KEYCLOAK_ADMIN_PASSWORD=admin mvn -f quarkus/server/pom.xml compile quarkus:dev -Dquarkus.args="start-dev --hostname=keycloak.127.0.0.1.nip.io --db=postgres --db-url=jdbc:postgresql://localhost:5432/keycloak --db-username=keycloak --db-password=keycloak"

# launch vscode
# attach with port 5005







*** Maven problems

To troubleshoot, use following parameters

mkdir -p .vscode
cat > .vscode/settings.json <<EOF
{
    "java.jdt.ls.vmargs": "-Dlog.level=ALL -Djdt.ls.debug=true -Dmaven.plugin.validation=verbose -Xmx16G -Xms100m"
}
EOF



*** Proto file generation errors

[ERROR] Failed to execute goal org.infinispan.protostream:proto-schema-compatibility-maven-plugin:5.0.12.Final:proto-schema-compatibility-check (default) on project keycloak-model-infinispan: Execution default of goal org.infinispan.protostream:proto-schema-compatibility-maven-plugin:5.0.12.Final:proto-schema-compatibility-check failed: IPROTO000039: Incompatible schema changes:
[ERROR] IPROTO000035: Field 'keycloak.State.LOGGING_OUT' number was changed from 1 to 3
[ERROR] IPROTO000035: Field 'keycloak.State.LOGGED_OUT' number was changed from 2 to 1
[ERROR] IPROTO000035: Field 'keycloak.State.LOGGED_OUT_UNCONFIRMED' number was changed from 3 to 2
[ERROR] IPROTO000035: Field 'keycloak.ExecutionStatus.FAILED' number was changed from 0 to 4
[ERROR] IPROTO000035: Field 'keycloak.ExecutionStatus.SUCCESS' number was changed from 1 to 7
[ERROR] IPROTO000035: Field 'keycloak.ExecutionStatus.SETUP_REQUIRED' number was changed from 2 to 5
[ERROR] IPROTO000035: Field 'keycloak.ExecutionStatus.ATTEMPTED' number was changed from 3 to 0
[ERROR] IPROTO000035: Field 'keycloak.ExecutionStatus.SKIPPED' number was changed from 4 to 6
[ERROR] IPROTO000035: Field 'keycloak.ExecutionStatus.CHALLENGED' number was changed from 5 to 1
[ERROR] IPROTO000035: Field 'keycloak.ExecutionStatus.EVALUATED_TRUE' number was changed from 6 to 3
[ERROR] IPROTO000035: Field 'keycloak.ExecutionStatus.EVALUATED_FALSE' number was changed from 7 to 2



Workaround:

rm ./model/infinispan/target/classes/proto/generated/KeycloakModelSchema.proto
mvn -f model/infinispan/pom.xml clean install -DskipTests


*** ConfigError: The project 'keycloak-junit5' is not a valid java project.

The error comes when trying to launch debugger

???? NOT SOLVED  ????  (failed in 26.0.7, worked in newer version)


*** cannot convert from Stream<Object> to Stream<RoleModel>

Caused by: java.lang.Error: Unresolved compilation problems: 
	Type mismatch: cannot convert from Stream<Object> to Stream<RoleModel>
	Type mismatch: cannot convert from Stream<RoleModel> to Stream<Object>
	Type mismatch: cannot convert from Stream<Object> to Stream<RoleModel>
	Type mismatch: cannot convert from Stream<RoleModel> to Stream<Object>
	Type mismatch: cannot convert from Stream<Object> to Stream<RoleModel>
	Type mismatch: cannot convert from Stream<RoleModel> to Stream<Object>
	Type mismatch: cannot convert from Stream<Object> to Stream<RoleModel>
	Type mismatch: cannot convert from Stream<RoleModel> to Stream<Object>

	at org.keycloak.storage.RoleStorageManager.<init>(RoleStorageManager.java:172)
RoleStorageManager.java:172
	at org.keycloak.storage.datastore.DefaultDatastoreProvider.roleStorageManager(DefaultDatastoreProvider.java:93)
DefaultDatastoreProvider.java:93
	at org.keycloak.models.cache.infinispan.RealmCacheSession.getRoleDelegate(RealmCacheSession.java:211)
RealmCacheSession.java:211
	at org.keycloak.models.cache.infinispan.RealmCacheSession.addRealmRole(RealmCacheSession.java:715)
RealmCacheSession.java:715
	at org.keycloak.models.cache.infinispan.RealmCacheSession.addRealmRole(RealmCacheSession.java:710)
RealmCacheSession.java:710
	at org.keycloak.models.jpa.RealmAdapter.addRole(RealmAdapter.java:855)
RealmAdapter.java:855
	at org.keycloak.models.utils.KeycloakModelUtils.setupDefaultRole(KeycloakModelUtils.java:630)
KeycloakModelUtils.java:630
	at org.keycloak.services.managers.RealmManager.createRealm(RealmManager.java:130)
RealmManager.java:130
	at org.keycloak.services.managers.RealmManager.createRealm(RealmManager.java:113)
RealmManager.java:113
	at org.keycloak.services.managers.ApplianceBootstrap.createMasterRealm(ApplianceBootstrap.java:78)
ApplianceBootstrap.java:78
	at org.keycloak.services.resources.KeycloakApplication$2.run(KeycloakApplication.java:144)
KeycloakApplication.java:144
	at org.keycloak.models.utils.KeycloakModelUtils.lambda$1(KeycloakModelUtils.java:335)
KeycloakModelUtils.java:335
	at org.keycloak.models.utils.KeycloakModelUtils.runJobInTransactionWithResult(KeycloakModelUtils.java:449)
KeycloakModelUtils.java:449
	at org.keycloak.models.utils.KeycloakModelUtils.runJobInTransaction(KeycloakModelUtils.java:334)
KeycloakModelUtils.java:334
	at org.keycloak.models.utils.KeycloakModelUtils.runJobInTransaction(KeycloakModelUtils.java:324)
KeycloakModelUtils.java:324
	at org.keycloak.services.resources.KeycloakApplication.bootstrap(KeycloakApplication.java:116)
KeycloakApplication.java:116
	at org.keycloak.services.resources.KeycloakApplication$1.run(KeycloakApplication.java:86)
KeycloakApplication.java:86
	at org.keycloak.models.utils.KeycloakModelUtils.lambda$1(KeycloakModelUtils.java:335)
KeycloakModelUtils.java:335
	at org.keycloak.models.utils.KeycloakModelUtils.runJobInTransactionWithResult(KeycloakModelUtils.java:449)
KeycloakModelUtils.java:449
	at org.keycloak.models.utils.KeycloakModelUtils.runJobInTransaction(KeycloakModelUtils.java:334)
KeycloakModelUtils.java:334
	at org.keycloak.models.utils.KeycloakModelUtils.runJobInTransaction(KeycloakModelUtils.java:324)
KeycloakModelUtils.java:324
	at org.keycloak.services.resources.KeycloakApplication.startup(KeycloakApplication.java:78)
KeycloakApplication.java:78
	at org.keycloak.quarkus.runtime.integration.jaxrs.QuarkusKeycloakApplication.onStartupEvent(QuarkusKeycloakApplication.java:52)
QuarkusKeycloakApplication.java:52
	at org.keycloak.quarkus.runtime.integration.jaxrs.QuarkusKeycloakApplication_Observer_onStartupEvent_GNZ8m5QenZ9h9VNelo7awjUZFDE.notify(Unknown Source)
	at io.quarkus.arc.impl.EventImpl$Notifier.notifyObservers(EventImpl.java:365)
EventImpl.java:365
	at io.quarkus.arc.impl.EventImpl$Notifier.notify(EventImpl.java:347)
EventImpl.java:347
	at io.quarkus.arc.impl.EventImpl.fire(EventImpl.java:81)
EventImpl.java:81
	at io.quarkus.arc.runtime.ArcRecorder.fireLifecycleEvent(ArcRecorder.java:163)
ArcRecorder.java:163
	at io.quarkus.arc.runtime.ArcRecorder.handleLifecycleEvents(ArcRecorder.java:114)
ArcRecorder.java:114
	at io.quarkus.runner.recorded.LifecycleEventsBuildStep$startupEvent1144526294.deploy_0(Unknown Source)
	at io.quarkus.runner.recorded.LifecycleEventsBuildStep$startupEvent1144526294.deploy(Unknown Source)
	... 24 more
￼


workaround


diff --git a/model/storage-private/src/main/java/org/keycloak/storage/RoleStorageManager.java b/model/storage-private/src/main/java/org/keycloak/storage/RoleStorageManager.java
index a3f453a1b4..48b3b60a72 100644
--- a/model/storage-private/src/main/java/org/keycloak/storage/RoleStorageManager.java
+++ b/model/storage-private/src/main/java/org/keycloak/storage/RoleStorageManager.java
@@ -170,7 +170,7 @@ public class RoleStorageManager implements RoleProvider {
     public Stream<RoleModel> searchForRolesStream(RealmModel realm, String search, Integer first, Integer max) {
         Stream<RoleModel> local = localStorage().searchForRolesStream(realm, search, first, max);
         Stream<RoleModel> ext = getEnabledStorageProviders(session, realm, RoleLookupProvider.class)
-                .flatMap(ServicesUtils.timeBound(session,
+                .<RoleModel>flatMap(ServicesUtils.timeBound(session,
                         roleStorageProviderTimeout,
                         p -> ((RoleLookupProvider) p).searchForRolesStream(realm, search, first, max)));

@@ -237,7 +237,7 @@ public class RoleStorageManager implements RoleProvider {
     public Stream<RoleModel> searchForClientRolesStream(ClientModel client, String search, Integer first, Integer max) {
         Stream<RoleModel> local = localStorage().searchForClientRolesStream(client, search, first, max);
         Stream<RoleModel> ext = getEnabledStorageProviders(session, client.getRealm(), RoleLookupProvider.class)
-                .flatMap(ServicesUtils.timeBound(session,
+                .<RoleModel>flatMap(ServicesUtils.timeBound(session,
                         roleStorageProviderTimeout,
                         p -> ((RoleLookupProvider) p).searchForClientRolesStream(client, search, first, max)));

@@ -248,7 +248,7 @@ public class RoleStorageManager implements RoleProvider {
     public Stream<RoleModel> searchForClientRolesStream(RealmModel realm, Stream<String> ids, String search, Integer first, Integer max) {
         Stream<RoleModel> local = localStorage().searchForClientRolesStream(realm, ids, search, first, max);
         Stream<RoleModel> ext = getEnabledStorageProviders(session, realm, RoleLookupProvider.class)
-                .flatMap(ServicesUtils.timeBound(session,
+                .<RoleModel>flatMap(ServicesUtils.timeBound(session,
                         roleStorageProviderTimeout,
                         p -> ((RoleLookupProvider) p).searchForClientRolesStream(realm, ids, search, first, max)));

@@ -259,7 +259,7 @@ public class RoleStorageManager implements RoleProvider {
     public Stream<RoleModel> searchForClientRolesStream(RealmModel realm, String search, Stream<String> excludedIds, Integer first, Integer max) {
         Stream<RoleModel> local = localStorage().searchForClientRolesStream(realm, search, excludedIds, first, max);
         Stream<RoleModel> ext = getEnabledStorageProviders(session, realm, RoleLookupProvider.class)
-                .flatMap(ServicesUtils.timeBound(session,
+                .<RoleModel>flatMap(ServicesUtils.timeBound(session,
                         roleStorageProviderTimeout,
                         p -> ((RoleLookupProvider) p).searchForClientRolesStream(realm, search, excludedIds, first, max)));




*** Vscode java language server cannot show any symbols on testsuite/integration-arquillan/


workaround: remove following line


diff --git a/distribution/pom.xml b/distribution/pom.xml
index d2b8affb0c..f90bce7b1b 100755
--- a/distribution/pom.xml
+++ b/distribution/pom.xml
@@ -43,7 +43,6 @@
     <modules>
         <module>saml-adapters</module>
         <module>galleon-feature-packs</module>
-        <module>licenses-common</module>
         <module>maven-plugins</module>
     </modules>




