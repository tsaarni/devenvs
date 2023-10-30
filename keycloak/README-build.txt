
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




*** Debugging


# Debug directly from vscode

mkdir -p .vscode
cp -a ~/work/devenvs/keycloak/configs/launch.json .vscode

##  1. Build with maven
##  2. Start vscode
##  3. Launch the debug session



# Sometimes keycloak does not start, it only prints:
Press [e] to edit command line args (currently 'start-dev --hostname=keycloak.127-0-0-1.nip.io'), [h] for more options>

select "Java: Clean Java language server workspace"



# Remote debug
mvn clean install -f testsuite/integration-arquillian/pom.xml -DforkMode=never -Dmaven.surefire.debug  ...   # attach to port 5005 (not 8000)





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


Resource:


./model/jpa/src/main/resources/default-persistence.xml





2023-10-27 13:37:31,544 ERROR [org.keycloak.services.error.KeycloakErrorHandler] (executor-thread-4) Uncaught server error: org.keycloak.theme.FreeMarkerException: Failed to process template index.ftl
...

Caused by: freemarker.template.TemplateNotFoundException: Template not found for name "index.ftl".


mvn -f js/pom.xml clean install




# when having errors running in vscode, test also via command line:
mvn -f quarkus/server/pom.xml compile quarkus:dev -Dquarkus.args="start-dev"
