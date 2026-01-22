
### To build custom container use following procedure

# Compile only Keycloak distribution
mvn clean install -DskipTests -DskipTestsuite

KEYCLOAK_DIST=$(basename quarkus/dist/target/keycloak-*.tar.gz)
KEYCLOAK_VERSION=$(echo $KEYCLOAK_DIST | grep -oP 'keycloak-\K\d+\.\d+\.\d+')

# Copy the tar.gz package to container build directory
cp ./quarkus/dist/target/$KEYCLOAK_DIST quarkus/container/

# Create container image
(cd quarkus/container/; docker build --build-arg KEYCLOAK_DIST=$KEYCLOAK_DIST -t localhost/keycloak:$KEYCLOAK_VERSION .)


### Run custom Keycloak container

docker run \
  --network host \
  -e KC_BOOTSTRAP_ADMIN_USERNAME=admin \
  -e KC_BOOTSTRAP_ADMIN_PASSWORD=admin \
  -e KC_HTTP_ENABLED=true \
  -e KC_HOSTNAME_STRICT=false \
  localhost/keycloak:26.2.10 \
  start --log-level=INFO

http://keycloak.127.0.0.1.nip.io:8080/

# or upstream images:

quay.io/keycloak/keycloak:26.1.1


## Running together with postgres and openldap docker-compose.yml

docker compose up
docker-compose rm -f
