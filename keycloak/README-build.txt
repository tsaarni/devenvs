
# Build
mvn clean install -DskipTestsuite -DskipExamples -DskipTests

# Parallel build
mvnd clean install -DskipTestsuite -DskipExamples -DskipTests
mvn -T4C clean install -DskipTestsuite -DskipExamples -DskipTests



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
