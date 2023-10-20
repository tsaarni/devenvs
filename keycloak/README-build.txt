
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





# Remote debug
mvn clean install -f testsuite/integration-arquillian/pom.xml -DforkMode=never -Dmaven.surefire.debug  ...   # attach to port 5005 (not 8000)
