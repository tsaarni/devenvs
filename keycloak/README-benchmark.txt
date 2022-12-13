
https://github.com/keycloak/keycloak-benchmark


################################
#
# Compile dataset module
#

mvn clean install -am -pl dataset
mkdir -p provision/docker-compose/keycloak/providers
cp dataset/target/keycloak-benchmark-dataset-*.jar provision/docker-compose/keycloak/providers


docker-compose -f provision/docker-compose/keycloak/keycloak-quarkus-postgres.yml up


## or restart Keycloak if it was running while the provider was added


# remove persistent volume
docker-compose -f provision/docker-compose/keycloak/keycloak-quarkus-postgres.yml down -v


################################
#
# Dataset REST API
#

https://www.keycloak.org/keycloak-benchmark/dataset-guide/latest/using-provider


# print current status
http http://localhost:8080/realms/master/dataset/status

# create a realm with users, roles, groups, clients
http http://localhost:8080/realms/master/dataset/create-realms?count=1
http "http://localhost:8080/realms/master/dataset/create-realms?count=1&users-per-realm=100000"

# get the last created realm name
http http://localhost:8080/realms/master/dataset/last-realm



# create users
http "http://localhost:8080/realms/master/dataset/create-users?count=100000&realm-roles-per-user=0&groups-per-user=0&client-roles-per-user=0&realm-name=master"

http "http://localhost:8080/realms/master/dataset/create-users?count=100000&realm-name=master"

# create clients
http "http://localhost:8080/realms/master/dataset/create-clients?count=1000&realm-name=master"
