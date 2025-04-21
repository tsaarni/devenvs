
# https://quarkus.io/guides/datasource#quarkus-agroal_quarkus-datasource-jdbc-max-lifetime


version: "3.8"
services:
  postgres:
    image: docker.io/postgres:17-alpine
    command: -c log_statement=all -c log_destination=stderr
    ports:
      - "5432:5432"
    environment:
      - POSTGRES_USER=keycloak
      - POSTGRES_PASSWORD=keycloak
      - POSTGRES_DB=keycloak

  keycloak:
    # https://quay.io/repository/keycloak/keycloak?tab=tags
    image: quay.io/keycloak/keycloak:26.1.4
    command:
      - start-dev
    ports:
      - "8080:8080"
    environment:
      - KC_BOOTSTRAP_ADMIN_USERNAME=admin
      - KC_BOOTSTRAP_ADMIN_PASSWORD=admin
      - KC_DB=postgres
      - KC_DB_URL=jdbc:postgresql://postgres/keycloak
      - KC_DB_USERNAME=keycloak
      - KC_DB_PASSWORD=keycloak
      - QUARKUS_DATASOURCE_JDBC_MAX_LIFETIME=1M




docker compose up

sudo nsenter --target $(pgrep -f quarkus) --net wireshark -i any -k -f "port 5432" --display-filter "tcp.flags.syn == 1 || tcp.flags.fin == 1 || tcp.flags.reset == 1"
