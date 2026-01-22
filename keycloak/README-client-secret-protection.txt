
### Add support for looking up client secrets via Vault SPI
### https://github.com/keycloak/keycloak/pull/39650



### Client secret: hashing of stored secrets
### https://github.com/keycloak/keycloak/discussions/8455

### Impacts
### https://gist.github.com/tsaarni/be561c8ea8bdcdb2ebe02d30937eb7e0


### Old jira
### https://github.com/keycloak/keycloak/discussions/8455


# realm should have an option to configure client secret hashing algorithm from:

No hashing
sha-256
sha-512

## From Argon2 password provider implementation
## https://github.com/keycloak/keycloak/pull/28031/files


* argon2:: Argon2 (recommended for non-FIPS deployments)
￼* pbkdf2-sha512:: PBKDF2 with SHA512 (default, recommended for FIPS deployments)
￼* pbkdf2-sha256:: PBKDF2 with SHA256
￼

crypto/default/src/main/java/org/keycloak/crypto/hash/Argon2PasswordHashProvider.java
server-spi/src/main/java/org/keycloak/credential/hash/PasswordHashProvider.java
server-spi-private/src/main/java/org/keycloak/credential/hash/Pbkdf2PasswordHashProvider.java


server-spi/src/main/java/org/keycloak/credential/CredentialModel.java
model/jpa/src/main/java/org/keycloak/models/jpa/entities/CredentialEntity.java


server-spi/src/main/java/org/keycloak/models/ClientModel.java

server-spi-private/src/main/java/org/keycloak/models/ClientSecretConstants.java

model/jpa/src/main/java/org/keycloak/models/jpa/ClientAdapter.java
model/jpa/src/main/java/org/keycloak/models/jpa/entities/ClientEntity.java

services/src/main/java/org/keycloak/services/resources/admin/ClientResource.java
services/src/main/java/org/keycloak/protocol/oidc/OIDCClientSecretConfigWrapper.java





Client secrets are currently stored in clear plain-text in the database.

We've got a strong requirement to secure client secrets at rest.

Can it be possible to either :

allow to obtain client secrets from an external vault
support client secret encryption as it is for user password



Concept for hashed client secrets and other custom secret types
A client can have a secret of type "plain" (current secret), "hashed" (secret hashed like user password) and "none" (for public clients)
Additional secret types can be added with a SPI, so that secrets could be stored at an external vault.

Every secret implementation is responsible for generating and storing secrets as well as secret verification for authentication. An operation to get the secret is optional (e.g. type "plain" can do this, type "hashed" can not):

An implementation provides the following operations:

- SecretValue generate(secretId)
- boolean verify(secretId, candidateSecretValue) throws SecretFormatException
- SecretValue getSecret(secretId) // optional, could be in a sub-interface which is implemented for type "plain"
- deleteSecret(secretId)

If getSecret is not supported, admin console shows the client secret only one time directly after creation, so that it can be copied. In this case secret creation isn't done automatically when client is created. It is an own action after client creation (Could be the existing "regenerate secret" action).

New realm settings allow to set the allowed client secret types and the default client secret type. This allows to get a realm with more secure secret handling by default.




*** Create client


mkdir -p .vscode
cp -a ~/work/devenvs/keycloak/configs/launch.json  ~/work/devenvs/keycloak/configs/settings.json .vscode


Start parameters in .vscode/launch.json

            "args": "start-dev --verbose --vault=file --vault-dir=${workspaceFolder}/quarkus/dist/target/vault-secrets --log-level=INFO,org.keycloak.vault:debug",


Create secret

mkdir -p quarkus/dist/target/vault-secrets
echo -n "my-client-secret" > quarkus/dist/target/vault-secrets/second_client-secret


First realm

{
  "realm": "first",
  "displayName": "First Keycloak Realm",
  "enabled": true,
  "identityProviders": [
    {
      "alias": "keycloak-oidc",
      "displayName": "Second Keycloak Realm (federated)",
      "providerId": "keycloak-oidc",
      "enabled": true,
      "config": {
        "tokenUrl": "http://keycloak.127.0.0.1.nip.io:8080/realms/second/protocol/openid-connect/token",
        "authorizationUrl": "http://keycloak.127.0.0.1.nip.io:8080/realms/second/protocol/openid-connect/auth",
        "clientId": "federator",
        "clientSecret": "my-client-secret"
      }
    }
  ]
}

Second realm

{
  "realm": "second",
  "displayName": "Second Keycloak Realm",
  "enabled": true,
  "clients": [
    {
      "clientId": "federator",
      "enabled": true,
      "redirectUris": ["http://keycloak.127.0.0.1.nip.io:8080/*"],
      "webOrigins": ["*"],
      "publicClient": false,
      "protocol": "openid-connect",
      "serviceAccountsEnabled": true,
      "secret": "${vault.client-secret}"
    }
  ],
  "users": [
    {
      "username": "joe",
      "email": "joe@example.com",
      "firstName": "Joe",
      "lastName": "Average",
      "enabled": true,
      "credentials": [
        {
          "type": "password",
          "value": "joe",
          "temporary": false
        }
      ]
    }
  ]
}



# login with IdP federation
http://keycloak.127.0.0.1.nip.io:8080/realms/first/account/


wireshark -i lo -k -Y 'http' -f 'tcp port 8080'



function get_admin_token() {
  http --form POST http://keycloak.127-0-0-1.nip.io:8080/realms/master/protocol/openid-connect/token username=admin password=admin grant_type=password client_id=admin-cli | jq -r .access_token
}


http POST http://keycloak.127-0-0-1.nip.io:8080/admin/realms/master/clients \
    Authorization:"Bearer $(get_admin_token)" \
    name="Test client" \
    clientId="test-client" \
    enabled:=true \
    secret="\${vault.client-secret}" \
    protocol="openid-connect" \
    serviceAccountsEnabled:=true \
    publicClient:=false

http http://keycloak.127-0-0-1.nip.io:8080/admin/realms/master/clients Authorization:"Bearer $(get_admin_token)"

http --form POST http://keycloak.127-0-0-1.nip.io:8080/realms/master/protocol/openid-connect/token \
    client_id="test-client" \
    client_secret="my-client-secret" \
    grant_type="client_credentials"

#############
#
# Client secret rotation and client secret policy
#

https://github.com/keycloak/keycloak/discussions/9156
https://github.com/keycloak/keycloak-community/blob/main/design/client-secret-rotation.md

--features=client-secret-rotation



#############
#
# Vault SPI approach
#



# to see all places where vault is accessed
services/src/main/java/org/keycloak/services/DefaultKeycloakSession.java  (see callers of vault() method)

# all vault access is done through the KeycloakSession interface
session.vault()


# Clients are created
services/src/main/java/org/keycloak/services/resources/admin/ClientsResource.java

# And updated via
services/src/main/java/org/keycloak/services/resources/admin/ClientResource.java


# helper classes for creating and updating clients
services/src/main/java/org/keycloak/services/managers/ClientManager.java
server-spi-private/src/main/java/org/keycloak/models/utils/RepresentationToModel.java   # createClient()


# Rotation
services/src/main/java/org/keycloak/services/clientpolicy/executor/ClientSecretRotationExecutor.java


# access client secret
services/src/main/java/org/keycloak/protocol/oidc/OIDCClientSecretConfigWrapper.java

## authenticators
# Validates client based on "client_id" and "client_secret", calls client.validateSecret() and wrapper.validateRotatedSecret()
services/src/main/java/org/keycloak/authentication/authenticators/client/ClientIdAndSecretAuthenticator.java   # authenticateClient()

# Client authentication based on JWT signed by client secret instead of private key
services/src/main/java/org/keycloak/authentication/authenticators/client/JWTClientSecretAuthenticator.java


# Validation happens in
model/jpa/src/main/java/org/keycloak/models/jpa/ClientAdapter.java                      # validateSecret()
model/infinispan/src/main/java/org/keycloak/models/cache/infinispan/ClientAdapter.java  # validateSecret()
services/src/main/java/org/keycloak/protocol/oidc/OIDCClientSecretConfigWrapper.java    # validateRotatedSecret()








##

services/src/main/java/org/keycloak/vault/AbstractVaultProvider.java
services/src/main/java/org/keycloak/vault/FilesPlainTextVaultProviderFactory.java
services/src/main/java/org/keycloak/vault/FilesPlainTextVaultProvider.java
services/src/main/java/org/keycloak/vault/FilesKeystoreVaultProvider.java
services/src/main/java/org/keycloak/vault/FilesKeystoreVaultProviderFactory.java

services/src/main/resources/META-INF/services/org.keycloak.vault.VaultProviderFactory


Remove Hashicorp Support
https://github.com/keycloak/keycloak/issues/9144
https://github.com/keycloak/keycloak/discussions/16446

HashiCorp Vault provider for keycloak
https://github.com/InseeFrLab/keycloak-hashicorp-vault-ext






mvn clean install -f testsuite/integration-arquillian/pom.xml  -Dtest=org.keycloak.testsuite.client.ClientSecretRotationTest



######


swap container

# Compile only Keycloak distribution
mvn clean install -DskipTests -DskipTestsuite

KEYCLOAK_DIST=$(basename quarkus/dist/target/keycloak-*.tar.gz)
KEYCLOAK_VERSION=$(echo $KEYCLOAK_DIST | grep -oP 'keycloak-\K\d+\.\d+\.\d+')

# Copy the tar.gz package to container build directory
cp ./quarkus/dist/target/$KEYCLOAK_DIST quarkus/container/

# Create container image
(cd quarkus/container/; docker build --build-arg KEYCLOAK_DIST=$KEYCLOAK_DIST -t localhost/keycloak:$KEYCLOAK_VERSION .)

kind load docker-image --name secrets-provider localhost/keycloak:999.0.0
kubectl set image statefulset/keycloak keycloak=localhost/keycloak:999.0.0



# enable kubernetes auth, write secret
./mvnw verify -DskipEnvSetup





function get_admin_token() {
  http --form POST http://127.0.0.127:8080/realms/master/protocol/openid-connect/token username=admin password=admin grant_type=password client_id=admin-cli | jq -r .access_token
}

http POST http://127.0.0.127:8080/admin/realms/first/clients \
    Authorization:"Bearer $(get_admin_token)" \
    name="Test client" \
    clientId="test-client" \
    enabled:=true \
    secret="\${vault.client-secret}" \
    protocol="openid-connect" \
    serviceAccountsEnabled:=true \
    publicClient:=false

# view that client has the secret reference
http http://127.0.0.127:8080/admin/realms/first/clients Authorization:"Bearer $(get_admin_token)"


http --form POST http://127.0.0.127:8080/realms/first/protocol/openid-connect/token \
    client_id="test-client" \
    client_secret="my-client-secret" \
    grant_type="client_credentials"


kubectl logs keycloak-0 -f

./mvnw clean install -DskipTests && kubectl delete pod keycloak-0 --force && kubectl wait --for=condition=Running pod/keycloak-0 --timeout=60s && kubectl logs keycloak-0 -f

# list secrets in vault
http GET  http://127.0.0.127:8080/admin/realms/first/secrets-manager/ Authorization:"Bearer $(get_admin_token)"

# create a new secret
http POST http://127.0.0.127:8080/admin/realms/first/secrets-manager/my-client-secret Authorization:"Bearer $(get_admin_token)"

# view that secret
http GET http://127.0.0.127:8080/admin/realms/first/secrets-manager/my-client-secret Authorization:"Bearer $(get_admin_token)"

# delete the secret
http DELETE http://127.0.0.127:8080/admin/realms/first/secrets-manager/my-client-secret Authorization:"Bearer $(get_admin_token)"


http LIST http://127.0.0.127:8200/v1/secretv1/keycloak/first/ X-Vault-Token:"my-root-token"



sudo nsenter --target $(kindps secrets-provider openbao --output json | jq -r '.[0].pids[0].pid') --net -- wireshark -f "port 8200" -k -Y http -i any










1. Create secret for --vault=file SPI

mkdir -p quarkus/dist/target/vault-secrets
echo -n "my-client-secret" > quarkus/dist/target/vault-secrets/master_client-secret

2. Start Keycloak with following command line in .vscode/launch.json

            "args": "start-dev --verbose --vault=file --vault-dir=${workspaceFolder}/quarkus/dist/target/vault-secrets",
            "env": {
                "KC_BOOTSTRAP_ADMIN_USERNAME": "admin",
                "KC_BOOTSTRAP_ADMIN_PASSWORD": "admin",
                "KC_ADMIN_VITE_URL": "http://localhost:5174",
            },

# Start dev admin ui with nodejs on the background
cd js
pnpm --filter keycloak-admin-ui run dev


3. Create client with vault reference

function get_admin_token() {
  http --form POST http://127.0.0.1:8080/realms/master/protocol/openid-connect/token username=admin password=admin grant_type=password client_id=admin-cli | jq -r .access_token
}

http POST http://127.0.0.1:8080/admin/realms/master/clients \
    Authorization:"Bearer $(get_admin_token)" \
    name="Test client" \
    clientId="test-client" \
    enabled:=true \
    secret="\${vault.client-secret}" \
    protocol="openid-connect" \
    serviceAccountsEnabled:=true \
    publicClient:=false

4. View that client has the secret reference ${{vault.client-secret}}

http http://127.0.0.1:8080/admin/realms/master/clients?clientId=test-client Authorization:"Bearer $(get_admin_token)"


5. Verify that token can be fetched with client_secret while client secret is in vault

http --form POST http://127.0.0.1:8080/realms/master/protocol/openid-connect/token \
    client_id="test-client" \
    client_secret="my-client-secret" \
    grant_type="client_credentials"


6. Get the OIDC client adapter configuration and check they have the actual secret "my-client-secret"

CLIENT_ID=$(http http://127.0.0.1:8080/admin/realms/master/clients?clientId=test-client Authorization:"Bearer $(get_admin_token)" | jq -r '.[0].id')
http http://127.0.0.1:8080/admin/realms/master/clients/$CLIENT_ID/installation/providers/keycloak-oidc-keycloak-json Authorization:"Bearer $(get_admin_token)"
http http://127.0.0.1:8080/admin/realms/master/clients/$CLIENT_ID/installation/providers/keycloak-oidc-jboss-subsystem Authorization:"Bearer $(get_admin_token)"
http http://127.0.0.1:8080/admin/realms/master/clients/$CLIENT_ID/installation/providers/keycloak-oidc-jboss-subsystem-cli Authorization:"Bearer $(get_admin_token)"



7. Check that the web ui looks correct

- Navigate to  http://127.0.0.1:8080
- Log in admin:admin
- Navigate to the "Clients" section
- Click on "Test client"
- Click on the "Credentials" tab
- Click on the eye button next to the client secret to verify it is correctly referenced as ${vault.client-secret}
- Click "Client Secret (?)" help button and validate it has "This field is able to obtain its value from vault, use ${vault.ID} format."
- Edit the client secret field to "new-client-secret"
- Save the changes
- Click "Action" and select "Download adapter configuration" from the dropdown
- Check that the configuration file contains the updated client secret "new-client-secret"




## Custom admin UI extension
## https://github.com/keycloak/keycloak/discussions/24805

# org.keycloak.services.ui.extend.UiPageSpi
# org.keycloak.services.ui.extend.UiTabSpi

# https://github.com/keycloak/keycloak-quickstarts/blob/latest/extension/extend-admin-console-spi/src/main/java/org/keycloak/admin/ui/AdminUiPage.java

## Allow overriding storage behavior for Admin UI extensions with UiTabProviderFactory and UiPageProviderFactory
## https://github.com/keycloak/keycloak/issues/28931







#### Run test cases

# uses new junit based e2e test framework

mvn -f tests/pom.xml test -Dtest=ClientVaultTest



