

The original architecture document for Vault SPI was "Secure Credentials Store - Vault SPI (Part I)" ([link](https://github.com/keycloak/keycloak-community/blob/main/design/secure-credentials-store.md)).
There was no part II.


## Offloading token signing to KMS


Stackoverflow question "Keycloak: HSM Provider for Realm Keys" ([link](https://stackoverflow.com/questions/76929725/keycloak-hsm-provider-for-realm-keys)) discusses how to do either:

* offload either token signing to HSM
* fetch keys from HSM and cache them in memory

The response has simple example code that shows the latter, via Key SPIs `KeyProvider`.



Github discussion "Optional outsourcing of encryption and signature (with e.g. KMS)" ([link](https://github.com/keycloak/keycloak/discussions/22264)) is related to the same but did not seem to proceed further.









## References

Following is a collection of related references



* Config attribute encryption "Securing credentials/passwords not possible with Quarkus distribution" ([link](https://github.com/keycloak/keycloak/issues/11089)).
  It was further discussed "Elytron replacement in Quarkus distribution of Keycloak" ([link](https://github.com/keycloak/keycloak/discussions/19281)).
  Issue was solved by "SmallRye Keystore" ([link](https://github.com/keycloak/keycloak/pull/20375)) enabled by "Secret Keys Handlers" in smallrye-config ([link](https://github.com/smallrye/smallrye-config/pull/833)). Example for passing encrypted database password:
  `db-password=${jasypt::ENC(wqp8zDeiCQ5JaFvwDtoAcr2WMLdlD0rjwvo8Rh0thG5qyTQVGxwJjBIiW26y0dtU)}`






# Original design document
#
#


private static final Pattern pattern = Pattern.compile("^\\$\\{vault\\.(.+?)}$");

services/src/main/java/org/keycloak/vault/DefaultVaultTranscriber.java

# https://www.keycloak.org/server/vault
${vault.<realmname>_<secretname>}




Secrets

* Client secrets  (server generated, client configured)
* Access tokens from IdPs in the case of IdP brokering, when "Store Token" is enabled
  https://www.keycloak.org/docs/25.0.1/server_admin/#retrieving-external-idp-tokens
* Realm keys
* OTP shared secrets
* LDAP bind password
* SMTP password








External references

Kubernetes KMS v2
https://kubernetes.io/blog/2022/09/09/kms-v2-improvements/
https://github.com/kubernetes/enhancements/tree/master/keps/sig-auth/3299-kms-v2-improvements
https://kubernetes.io/docs/tasks/administer-cluster/kms-provider/
