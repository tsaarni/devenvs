
https://github.com/keycloak/keycloak/pull/39556


# We did not notice the above (back then unmerged) PR and created own version
https://github.com/keycloak/keycloak/pull/40251


            //"args": "start-dev --verbose --features=\"docker,oid4vc-vci\"",
            "args": "start-dev --verbose",


function get_admin_token() {   http --form POST http://127.0.0.1:8080/realms/master/protocol/openid-connect/token username=admin password=admin grant_type=password client_id=admin-cli | jq -r .access_token; }


http POST http://127.0.0.1:8080/admin/realms/master/client-scopes name=foo1 protocol=openid-connect Authorization:"Bearer $(get_admin_token)"
http POST http://127.0.0.1:8080/admin/realms/master/client-scopes name=foo2 protocol=saml Authorization:"Bearer $(get_admin_token)"
http POST http://127.0.0.1:8080/admin/realms/master/client-scopes name=foo3 protocol=docker-v2 Authorization:"Bearer $(get_admin_token)"
http POST http://127.0.0.1:8080/admin/realms/master/client-scopes name=foo4 protocol=oid4vc Authorization:"Bearer $(get_admin_token)"
http POST http://127.0.0.1:8080/admin/realms/master/client-scopes name=foo5 protocol=invalid Authorization:"Bearer $(get_admin_token)"


# clean up
rm ./quarkus/dist/target/keycloakdb.*
KC_BOOTSTRAP_ADMIN_USERNAME=admin KC_BOOTSTRAP_ADMIN_PASSWORD=admin java -jar quarkus/server/target/lib/quarkus-run.jar start-dev --features=docker --db-url "jdbc:h2:./quarkus/dist/target/keycloakdb;NON_KEYWORDS=VALUE;AUTO_SERVER=TRUE"
