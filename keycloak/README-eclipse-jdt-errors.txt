****

https://github.com/keycloak/keycloak/issues/13145
https://github.com/keycloak/keycloak/pull/13147

Error 1: Cannot refer to 'this' nor 'super' while explicitly invoking a constructor
Error 2: The method myMethod(Map<String,String>) in the type App is not applicable for the arguments (Map<Object,Object>)



****


services/src/main/java/org/keycloak/services/resources/admin/info/ServerInfoAdminResource.java

        // True - asymmetric algorithms, false - symmetric algorithms
        Map<Boolean, List<String>> algorithms = session.getAllProviders(ClientSignatureVerifierProvider.class).stream()
                        .collect(
                                Collectors.toMap(
                                        ClientSignatureVerifierProvider::isAsymmetricAlgorithm,
                                        clientSignatureVerifier -> Collections.singletonList(clientSignatureVerifier.getAlgorithm()),
                                        (l1, l2) -> listCombiner(l1, l2)
                                                .stream()
                                                .sorted()
                                                .collect(Collectors.toList()),
                                        HashMap::new
                                )
                        );



Error:
Type mismatch: cannot convert from Map<Boolean,Object> to Map<Boolean,List<String>>


Fix

diff --git a/services/src/main/java/org/keycloak/services/resources/admin/info/ServerInfoAdminResource.java b/services/src/main/java/org/keycloak/services/resources/admin/info/ServerInfoAdminResource.java
index 9b66c7f8be..8168f5987a 100755
--- a/services/src/main/java/org/keycloak/services/resources/admin/info/ServerInfoAdminResource.java
+++ b/services/src/main/java/org/keycloak/services/resources/admin/info/ServerInfoAdminResource.java
@@ -113,7 +113,7 @@ public class ServerInfoAdminResource {
                                 Collectors.toMap(
                                         ClientSignatureVerifierProvider::isAsymmetricAlgorithm,
                                         clientSignatureVerifier -> Collections.singletonList(clientSignatureVerifier.getAlgorithm()),
-                                        (l1, l2) -> listCombiner(l1, l2)
+                                        (l1, l2) -> ServerInfoAdminResource.<String>listCombiner(l1, l2)
                                                 .stream()
                                                 .sorted()
                                                 .collect(Collectors.toList()),


