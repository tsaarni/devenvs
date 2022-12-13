




Autorebuild during runtime

"Changes detected in configuration. Updating the server image."
"Updating the configuration and installing your custom providers, if any. Please wait."



https://www.baeldung.com/quarkus-bean-discovery-index








During buildtime

[ERROR] Failed to execute goal io.quarkus:quarkus-maven-plugin:2.14.1.Final:build (default) on project keycloak-quarkus-server-app: Failed to build quarkus application: io.quarkus.builder.BuildException: Build failure: Build failed due to errors
[ERROR]         [error]: Build step org.keycloak.quarkus.deployment.KeycloakProcessor#configureKeycloakSessionFactory threw an exception: java.util.ServiceConfigurationError: org.keycloak.models.map.storage.MapStorageProviderFactory: Provider org.keycloak.quarkus.runtime.storage.database.jpa.QuarkusJpaMapStorageProviderFactory could not be instantiated
[ERROR]         at java.base/java.util.ServiceLoader.fail(ServiceLoader.java:582)
[ERROR]         at java.base/java.util.ServiceLoader$ProviderImpl.newInstance(ServiceLoader.java:804)
[ERROR]         at java.base/java.util.ServiceLoader$ProviderImpl.get(ServiceLoader.java:722)
[ERROR]         at java.base/java.util.ServiceLoader$3.next(ServiceLoader.java:1395)
[ERROR]         at org.keycloak.provider.DefaultProviderLoader.load(DefaultProviderLoader.java:60)
[ERROR]         at org.keycloak.provider.ProviderManager.load(ProviderManager.java:94)
[ERROR]         at org.keycloak.quarkus.deployment.KeycloakProcessor.loadFactories(KeycloakProcessor.java:654)
[ERROR]         at org.keycloak.quarkus.deployment.KeycloakProcessor.configureKeycloakSessionFactory(KeycloakProcessor.java:366)
[ERROR]         at java.base/jdk.internal.reflect.NativeMethodAccessorImpl.invoke0(Native Method)
[ERROR]         at java.base/jdk.internal.reflect.NativeMethodAccessorImpl.invoke(NativeMethodAccessorImpl.java:62)
[ERROR]         at java.base/jdk.internal.reflect.DelegatingMethodAccessorImpl.invoke(DelegatingMethodAccessorImpl.java:43)
[ERROR]         at java.base/java.lang.reflect.Method.invoke(Method.java:566)
[ERROR]         at io.quarkus.deployment.ExtensionLoader$3.execute(ExtensionLoader.java:909)
[ERROR]         at io.quarkus.builder.BuildContext.run(BuildContext.java:281)
[ERROR]         at org.jboss.threads.ContextHandler$1.runWith(ContextHandler.java:18)
[ERROR]         at org.jboss.threads.EnhancedQueueExecutor$Task.run(EnhancedQueueExecutor.java:2449)
[ERROR]         at org.jboss.threads.EnhancedQueueExecutor$ThreadBody.run(EnhancedQueueExecutor.java:1478)
[ERROR]         at java.base/java.lang.Thread.run(Thread.java:829)
[ERROR]         at org.jboss.threads.JBossThread.run(JBossThread.java:501)
[ERROR] Caused by: java.lang.ExceptionInInitializerError
[ERROR]         at java.base/jdk.internal.reflect.NativeConstructorAccessorImpl.newInstance0(Native Method)
[ERROR]         at java.base/jdk.internal.reflect.NativeConstructorAccessorImpl.newInstance(NativeConstructorAccessorImpl.java:62)
[ERROR]         at java.base/jdk.internal.reflect.DelegatingConstructorAccessorImpl.newInstance(DelegatingConstructorAccessorImpl.java:45)
[ERROR]         at java.base/java.lang.reflect.Constructor.newInstance(Constructor.java:490)
[ERROR]         at java.base/java.util.ServiceLoader$ProviderImpl.newInstance(ServiceLoader.java:780)
[ERROR]         ... 17 more
[ERROR] Caused by: java.lang.NullPointerException
[ERROR]         at java.base/java.util.HashMap.putMapEntries(HashMap.java:497)
[ERROR]         at java.base/java.util.HashMap.<init>(HashMap.java:486)
[ERROR]         at org.keycloak.models.map.common.DeepCloner$Builder.<init>(DeepCloner.java:115)
[ERROR]         at org.keycloak.models.map.storage.jpa.JpaMapStorageProviderFactory.<clinit>(JpaMapStorageProviderFactory.java:165)
[ERROR]         ... 22 more




JAX-RS filter gets loaded (by quarkus provider loader?)

https://github.com/thomasdarimont/keycloak-project-example/blob/main/keycloak/extensionsx/src/main/java/com/github/thomasdarimont/keycloakx/custom/security/AccessFilter.java
