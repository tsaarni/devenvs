


rm -rf certs
mkdir -p certs
certyaml -d certs configs/certs.yaml
chmod +r certs/*



docker compose -f docker-compose-auth-filters.yaml  up


http --verify certs/server-ca.pem https://keycloak.127-0-0-15.nip.io/realms/envoy/.well-known/openid-configuration


    "authorization_endpoint": "https://keycloak.127-0-0-15.nip.io/realms/envoy/protocol/openid-connect/auth",
    "token_endpoint": "https://keycloak.127-0-0-15.nip.io/realms/envoy/protocol/openid-connect/token",
    "jwks_uri": "https://keycloak.127-0-0-15.nip.io/realms/envoy/protocol/openid-connect/certs",

http --verify certs/server-ca.pem https://echoserver.127-0-0-15.nip.io/
https://echoserver.127-0-0-15.nip.io/



wireshark -i lo -f "port 443" -k -Y tls -o tls.keylog_file:/tmp/envoy-wireshark-keys.log


sudo nsenter --target $(docker inspect --format '{{.State.Pid}}' envoy-keycloak-1) --net -- wireshark -i any -k -f "port 8080"
