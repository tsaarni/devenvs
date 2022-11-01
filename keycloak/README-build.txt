

# compile distribution
mvn -Pdistribution -DskipTests clean install

# ./distribution/server-dist/target/keycloak-legacy-19.0.3.tar.gz
# ./quarkus/dist/target/keycloak-19.0.3.tar.gz
#
#or
#
# ./quarkus/dist/target/keycloak-999-SNAPSHOT.tar.gz



# Dependency tree
mvn dependency:tree -Pdistribution

# Dedendency on particular package
mvn dependency:tree -Pdistribution -Dincludes=jakarta.xml.bind:jakarta.xml.bind-api





export WORKDIR=/home/tsaarni/work/devenvs/keycloak
mkdir -p .vscode
cp $WORKDIR/configs/launch.json .vscode/


