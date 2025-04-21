
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


function get_admin_token() {
  http --form POST http://keycloak.127-0-0-1.nip.io:8080/realms/master/protocol/openid-connect/token username=admin password=admin grant_type=password client_id=admin-cli | jq -r .access_token
}


http POST http://keycloak.127-0-0-1.nip.io:8080/admin/realms/master/clients \
    Authorization:"Bearer $(get_admin_token)" \
    name="Test client" \
    clientId="test-client" \
    enabled:=true \
    secret="secret" \
    protocol="openid-connect" \
    serviceAccountsEnabled:=true \
    publicClient:=false

http --form POST http://keycloak.127-0-0-1.nip.io:8080/realms/master/protocol/openid-connect/token \
    client_id="test-client" \
    client_secret="secret" \
    grant_type="client_credentials"





#############
#
# Vault approach




# to see all places where vault is accessed
services/src/main/java/org/keycloak/services/DefaultKeycloakSession.java  (see vault() method)

# all vault access is done through the KeycloakSession interface
session.vault()


# Clients are created
services/src/main/java/org/keycloak/services/resources/admin/ClientsResource.java

# And updated via
services/src/main/java/org/keycloak/services/resources/admin/ClientResource.java


# helper classes for creating and updating clients
services/src/main/java/org/keycloak/services/managers/ClientManager.java
server-spi-private/src/main/java/org/keycloak/models/utils/RepresentationToModel.java   # createClient()



# access client secret
services/src/main/java/org/keycloak/protocol/oidc/OIDCClientSecretConfigWrapper.java

## authenticators
# Validates client based on "client_id" and "client_secret"
services/src/main/java/org/keycloak/authentication/authenticators/client/ClientIdAndSecretAuthenticator.java

# Client authentication based on JWT signed by client secret instead of private key
services/src/main/java/org/keycloak/authentication/authenticators/client/JWTClientSecretAuthenticator.java
