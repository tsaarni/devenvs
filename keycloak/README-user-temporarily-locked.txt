https://github.com/keycloak/keycloak/issues/11726
https://github.com/keycloak/keycloak/pull/8432

mkdir -p .vscode
cp -a ~/work/devenvs/keycloak/configs/launch.json .vscode/


mvn clean install -DskipTestsuite -DskipExamples -DskipTests

java -jar quarkus/server/target/lib/quarkus-run.jar --verbose start-dev

mvn -f quarkus/server/pom.xml compile quarkus:dev -Dquarkus.args="start-dev"

mvn -f quarkus/server/pom.xml -Dsuspend=true -Ddebug=8000 compile quarkus:dev -Dquarkus.args="start-dev"


cp ~/work/keycloak-benchmark/dataset/target/keycloak-benchmark-dataset-*.jar quarkus/server/target/providers

rm -rf target/kc
mkdir -p target/kc/providers
cp ~/work/keycloak-benchmark/dataset/target/keycloak-benchmark-dataset-*.jar target/kc/providers


http://keycloak.127-0-0-1.nip.io:8080/


# print current status for dataset
http http://keycloak.127-0-0-1.nip.io:8080/realms/master/dataset/status










Environment.getProviderFiles() (/home/tsaarni/work/keycloak-worktree/fix-user-enabled-status/quarkus/runtime/src/main/java/org/keycloak/quarkus/runtime/Environment.java:161)
KeycloakProcessor.persistBuildTimeProperties(BuildProducer) (/home/tsaarni/work/keycloak-worktree/fix-user-enabled-status/quarkus/deployment/src/main/java/org/keycloak/quarkus/deployment/KeycloakProcessor.java:491)
NativeMethodAccessorImpl.invoke0(Method,Object,Object[])[native method] (/java.base/jdk.internal.reflect/NativeMethodAccessorImpl.class:-1)
NativeMethodAccessorImpl.invoke(Object,Object[]) (/java.base/jdk.internal.reflect/NativeMethodAccessorImpl.class:62)
DelegatingMethodAccessorImpl.invoke(Object,Object[]) (/java.base/jdk.internal.reflect/DelegatingMethodAccessorImpl.class:43)
Method.invoke(Object,Object[]) (/rt.jar/java.lang.reflect/Method.class:566)
ExtensionLoader$3.execute(BuildContext) (/quarkus-core-deployment-2.13.3.Final.jar/io.quarkus.deployment/ExtensionLoader.class:909)
BuildContext.run() (/quarkus-builder-2.13.3.Final.jar/io.quarkus.builder/BuildContext.class:281)
1935751909.run() (Unknown Source:-1)
ContextHandler$1.runWith(Runnable,Object) (/jboss-threads-3.4.3.Final.jar/org.jboss.threads/ContextHandler.class:18)
EnhancedQueueExecutor$Task.run() (/jboss-threads-2.3.3.Final.jar/org.jboss.threads/EnhancedQueueExecutor.class:2449)
EnhancedQueueExecutor$ThreadBody.run() (/jboss-threads-2.3.3.Final.jar/org.jboss.threads/EnhancedQueueExecutor.class:1452)
Thread.run() (/rt.jar/java.lang/Thread.class:829)
JBossThread.run() (/jboss-threads-2.3.3.Final.jar/org.jboss.threads/JBossThread.class:501)



DefaultProviderLoader.loadSpis() (/home/tsaarni/work/keycloak-worktree/fix-user-enabled-status/services/src/main/java/org/keycloak/provider/DefaultProviderLoader.java:45)
ProviderManager.loadSpis() (/home/tsaarni/work/keycloak-worktree/fix-user-enabled-status/services/src/main/java/org/keycloak/provider/ProviderManager.java:86)
KeycloakProcessor.loadFactories(Map) (/home/tsaarni/work/keycloak-worktree/fix-user-enabled-status/quarkus/deployment/src/main/java/org/keycloak/quarkus/deployment/KeycloakProcessor.java:644)
KeycloakProcessor.configureKeycloakSessionFactory(KeycloakRecorder,List) (/home/tsaarni/work/keycloak-worktree/fix-user-enabled-status/quarkus/deployment/src/main/java/org/keycloak/quarkus/deployment/KeycloakProcessor.java:362)
NativeMethodAccessorImpl.invoke0(Method,Object,Object[])[native method] (/java.base/jdk.internal.reflect/NativeMethodAccessorImpl.class:-1)
NativeMethodAccessorImpl.invoke(Object,Object[]) (/java.base/jdk.internal.reflect/NativeMethodAccessorImpl.class:62)
DelegatingMethodAccessorImpl.invoke(Object,Object[]) (/java.base/jdk.internal.reflect/DelegatingMethodAccessorImpl.class:43)
Method.invoke(Object,Object[]) (/rt.jar/java.lang.reflect/Method.class:566)
ExtensionLoader$3.execute(BuildContext) (/quarkus-core-deployment-2.13.3.Final.jar/io.quarkus.deployment/ExtensionLoader.class:909)
BuildContext.run() (/quarkus-builder-2.13.3.Final.jar/io.quarkus.builder/BuildContext.class:281)
860579510.run() (Unknown Source:-1)
ContextHandler$1.runWith(Runnable,Object) (/jboss-threads-3.4.3.Final.jar/org.jboss.threads/ContextHandler.class:18)
EnhancedQueueExecutor$Task.run() (/jboss-threads-2.3.3.Final.jar/org.jboss.threads/EnhancedQueueExecutor.class:2449)
EnhancedQueueExecutor$ThreadBody.run() (/jboss-threads-2.3.3.Final.jar/org.jboss.threads/EnhancedQueueExecutor.class:1478)
Thread.run() (/rt.jar/java.lang/Thread.class:829)
JBossThread.run() (/jboss-threads-2.3.3.Final.jar/org.jboss.threads/JBossThread.class:501)
