
# Distributed tracing, OTL / OpenTelemetry

# https://www.keycloak.org/observability/tracing
# https://quarkus.io/guides/opentelemetry


# https://quarkus.io/blog/quarkus-3-13-0-released/
# https://quarkus.io/guides/tls-registry-reference

docker compose -f docker-compose-tracing.yaml rm --force --stop
docker compose -f docker-compose-tracing.yaml up



# Check jaeger UI
http://localhost:16686/

# Searcg for operations: POST /realms/{realm}/protocol/{protocol}/token
http://localhost:16686/search?operation=POST%20%2Frealms%2F%7Brealm%7D%2Fprotocol%2F%7Bprotocol%7D%2Ftoken&service=keycloak


# Keycloak attribute mapper
quarkus/runtime/src/main/java/org/keycloak/quarkus/runtime/configuration/mappers/TracingPropertyMappers.java

quarkus/config-api/src/main/java/org/keycloak/config/TracingOptions.java


TRACING_ENABLED           -> quarkus.otel.enabled
TRACING_ENDPOINT          -> quarkus.otel.exporter.otlp.traces.endpoint
TRACING_SERVICE_NAME      -> quarkus.otel.service.name
TRACING_RESOURCE_ATTRIBUTES -> quarkus.otel.resource.attributes
   "k8s.namespace.name=default,k8s.pod.name=keycloak-pod,k8s.container.name=keycloak-container"
TRACING_PROTOCOL          -> quarkus.otel.exporter.otlp.traces.protocol
   grpc (default), http/protobuf
   jaeger exporter is deprecated https://opentelemetry.io/blog/2022/jaeger-native-otlp/

TRACING_SAMPLER_TYPE      -> quarkus.otel.traces.sampler
   always_on, always_off, traceidratio (default), parentbased_always_on, parentbased_always_off, parentbased_traceidratio

   jaeger exporter is deprecated to parentbased_jaeger_remote is not relevant

TRACING_SAMPLER_RATIO     -> quarkus.otel.traces.sampler.arg
TRACING_COMPRESSION       -> quarkus.otel.exporter.otlp.traces.compression
TRACING_JDBC_ENABLED      -> quarkus.datasource.jdbc.telemetry


Additional headers
quarkus.otel.exporter.otlp.headers =



# Kubernetes attributes
# https://opentelemetry.io/docs/specs/semconv/attributes-registry/k8s/#kubernetes-attributes


Quarkus attributes:

quarkus.otel.exporter.otlp.headers




function get_admin_token() {
  http --form POST http://keycloak.127.0.0.1.nip.io:8080/realms/master/protocol/openid-connect/token \
    username=admin \
    password=admin \
    grant_type=password \
    client_id=admin-cli \
  | jq -r .access_token
}
