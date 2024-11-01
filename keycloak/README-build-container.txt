
### To build custom container use following procedure

# Compile only Keycloak distribution
mvn clean install -DskipTests -DskipTestsuite

KEYCLOAK_DIST=$(basename quarkus/dist/target/keycloak-*.tar.gz)
KEYCLOAK_VERSION=$(echo $KEYCLOAK_DIST | grep -oP 'keycloak-\K\d+\.\d+\.\d+')

# Copy the tar.gz package to container build directory
cp ./quarkus/dist/target/$KEYCLOAK_DIST quarkus/container/

# Create container image
(cd quarkus/container/; docker build --build-arg KEYCLOAK_DIST=$KEYCLOAK_DIST -t localhost/keycloak:$KEYCLOAK_VERSION .)

