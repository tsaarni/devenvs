


# Comment LDAP login test case

https://github.com/keycloak/keycloak/pull/7049/files



# StartTLS hangs with JDK11


~/work/keycloaktest$


mvn install  -Dmaven.surefire.debug -f testsuite/integration-arquillian/tests/base/  -Dtest=org.keycloak.testsuite.federation.ldap.LDAPUserLoginTest#loginLDAPUserCredentialVaultAuthenticationNoneEncryptionStartTLS



MINA hangs


https://issues.apache.org/jira/browse/DIRMINA-1023

https://issues.apache.org/jira/browse/DIRSTUDIO-1091?jql=text%20~%20%22hang%20starttls%22

https://issues.apache.org/jira/browse/DIRMINA-1119?jql=project%20%3D%20DIRMINA%20AND%20resolution%20%3D%20Unresolved%20AND%20text%20~%20%22ssl%22%20ORDER%20BY%20priority%20DESC%2C%20updated%20DESC




# Uplift of ApacheDS

API reference
- https://directory.apache.org/apacheds/gen-docs/2.0.0.AM26/apidocs/
- https://directory.apache.org/apacheds/gen-docs/2.0.0-M24/apidocs/org/apache/directory/server/core/factory/DirectoryServiceFactory.html


example use
https://www.codota.com/web/assistant/code/rs/5c65ce171095a500018102cc#L137

old was based on
https://github.com/kwart/ldap-server/blob/master/src/main/java/org/jboss/test/ldap/InMemoryDirectoryServiceFactory.java


# Certificate rotation with keyloack

Just create tickets for documenting in keycloak admin manual https://www.keycloak.org/docs/latest/server_admin/ ?

- https cert for keycloak itself
- jgroups cert rotation
   - https://issues.redhat.com/browse/WFLY-12164 (depends on the jgroups elytron patch?)
- jdbc client cert
   - https://jdbc.postgresql.org/documentation/head/ssl-client.html
   - https://quarkus.io/guides/security-jdbc   https://github.com/quarkusio/quarkus
- ldap client cert




# Using Elytron with Undertow server

http://darranl.blogspot.com/2017/09/using-wildfly-elytron-with-undertow.html
https://github.com/wildfly-security-incubator/elytron-examples/blob/master/undertow-standalone/src/main/java/org/wildfly/security/examples/HelloWorld.java
https://github.com/wildfly-security/elytron-web
https://github.com/wildfly-security/elytron-web/blob/1.x/undertow/src/main/java/org/wildfly/elytron/web/undertow/server/SecurityContextImpl.java
https://github.com/wildfly-security/wildfly-elytron

patch for using elytron with jgroups
- https://issues.redhat.com/browse/WFLY-12164
- https://github.com/wildfly/wildfly/pull/12517/files


guides (see 4.3.14 Default SSLContext)
- https://docs.wildfly.org/17/WildFly_Elytron_Security.html
- https://developers.redhat.com/blog/2018/04/20/elytron-new-security-framework-wildfly-jboss-eap/
- http://www.mastertheboss.com/jboss-server/jboss-security/complete-tutorial-for-configuring-ssl-https-on-wildfly


generic cert reload for java
- https://tersesystems.com/blog/2018/09/08/keymanagers-and-keystores/
- https://github.com/cloudfoundry/java-buildpack-security-provider/blob/master/src/main/java/org/cloudfoundry/security/FileWatcher.java
- https://github.com/eclipse/jetty.project/blob/ce6e146ac10607b2893961c47287598969f4b4c9/jetty-server/src/test/java/org/eclipse/jetty/server/ssl/SslContextFactoryReloadTest.java



# IPv4 vs IPv6 problem

https://stackoverflow.com/questions/62017891/optimal-ipv4-ipv6-settings-for-wildfly-on-linux
https://developer.jboss.org/thread/280976




# Installing Keycloak with keyloack operator to kind

https://github.com/keycloak/keycloak-operator






# StartTLS with SASL EXTERNAL method 

https://docs.oracle.com/javase/jndi/tutorial/ldap/security/ldap.html







mvn  verify -f testsuite/integration-arquillian/pom.xml -Dtest=org.keycloak.testsuite.federation.ldap.LDAPUserLoginTest#loginLDAPUserAuthenticationSimpleEncryptionStartTLS
mvn  verify -f testsuite/integration-arquillian/pom.xml -Dtest=org.keycloak.testsuite.federation.ldap.**



# run tests with wildfly
mvn clean install -DskipTests
(cd distribution; mvn clean install)

mvn clean verify -f testsuite/integration-arquillian/pom.xml -Pauth-server-wildfly -Dtest="org.keycloak.testsuite.federation.ldap.LDAPUserLoginTest#loginLDAPUserCredentialVaultAuthenticationSimpleEncryptionStartTLS"

mvn clean verify -f testsuite/integration-arquillian/pom.xml -Pauth-server-wildfly -Dtest="org.keycloak.testsuite.federation.ldap.LDAPUserLoginTest#loginLDAPUserCredentialVaultAuthenticationSimpleEncryptionSSL"





Caused by: java.lang.ClassNotFoundException: org.keycloak.truststore.SSLSocketFactory
        at java.base/java.lang.Class.forNameImpl(Native Method)
        at java.base/java.lang.Class.forName(Class.java:412)
        at java.naming/com.sun.jndi.ldap.VersionHelper.loadClass(VersionHelper.java:76)
        at java.naming/com.sun.jndi.ldap.Connection.createSocket(Connection.java:275)
        at java.naming/com.sun.jndi.ldap.Connection.<init>(Connection.java:216)





[ERROR] Failures:
[ERROR]   LDAPGroupMapperTest.lambda$test08_ldapOnlyGroupMappingsRanged$26a8868a$1:729
[ERROR]   LDAPGroupMapperNoImportTest>LDAPGroupMapperTest.lambda$test08_ldapOnlyGroupMappingsRanged$26a8868a$1:729


org.keycloak.testsuite.federation.ldap.LDAPGroupMapperTest
org.keycloak.testsuite.federation.ldap.noimport.LDAPGroupMapperNoImportTest



-------------------------------------------------------------------------------
Test set: org.keycloak.testsuite.federation.ldap.LDAPGroupMapperTest
-------------------------------------------------------------------------------
Tests run: 9, Failures: 1, Errors: 0, Skipped: 0, Time elapsed: 7.804 s <<< FAILURE! - in org.keycloak.testsuite.federation.ldap.LDAPGroupMapperTest
test08_ldapOnlyGroupMappingsRanged(org.keycloak.testsuite.federation.ldap.LDAPGroupMapperTest)  Time elapsed: 0.631 s  <<< FAILURE!
java.lang.AssertionError
        at org.junit.Assert.fail(Assert.java:86)
        at org.junit.Assert.assertTrue(Assert.java:41)
        at org.junit.Assert.assertFalse(Assert.java:64)
        at org.junit.Assert.assertFalse(Assert.java:74)
        at org.keycloak.testsuite.federation.ldap.LDAPGroupMapperTest.lambda$test08_ldapOnlyGroupMappingsRanged$26a8868a$1(LDAPGroupMapperTest.java:729)
        at org.keycloak.testsuite.federation.ldap.LDAPGroupMapperTest$$Lambda$708/0000000000000000.run(Unknown Source)
        at org.keycloak.testsuite.rest.TestingResourceProvider.runOnServer(TestingResourceProvider.java:796)
        at jdk.internal.reflect.GeneratedMethodAccessor856.invoke(Unknown Source)
        at java.base/jdk.internal.reflect.DelegatingMethodAccessorImpl.invoke(DelegatingMethodAccessorImpl.java:43)
        at java.base/java.lang.reflect.Method.invoke(Method.java:566)
        at org.jboss.resteasy.core.MethodInjectorImpl.invoke(MethodInjectorImpl.java:138)
        at org.jboss.resteasy.core.ResourceMethodInvoker.internalInvokeOnTarget(ResourceMethodInvoker.java:535)
        at org.jboss.resteasy.core.ResourceMethodInvoker.invokeOnTargetAfterFilter(ResourceMethodInvoker.java:424)
        at org.jboss.resteasy.core.ResourceMethodInvoker.lambda$invokeOnTarget$0(ResourceMethodInvoker.java:385)
        at org.jboss.resteasy.core.ResourceMethodInvoker$$Lambda$454/0000000000000000.get(Unknown Source)
        at org.jboss.resteasy.core.interception.PreMatchContainerRequestContext.filter(PreMatchContainerRequestContext.java:356)
        at org.jboss.resteasy.core.ResourceMethodInvoker.invokeOnTarget(ResourceMethodInvoker.java:387)
        at org.jboss.resteasy.core.ResourceMethodInvoker.invoke(ResourceMethodInvoker.java:356)
        at org.jboss.resteasy.core.ResourceLocatorInvoker.invokeOnTargetObject(ResourceLocatorInvoker.java:150)
        at org.jboss.resteasy.core.ResourceLocatorInvoker.invoke(ResourceLocatorInvoker.java:104)
        at org.jboss.resteasy.core.SynchronousDispatcher.invoke(SynchronousDispatcher.java:440)
        at org.jboss.resteasy.core.SynchronousDispatcher.lambda$invoke$4(SynchronousDispatcher.java:229)
        at org.jboss.resteasy.core.SynchronousDispatcher$$Lambda$452/0000000000000000.run(Unknown Source)
        at org.jboss.resteasy.core.SynchronousDispatcher.lambda$preprocess$0(SynchronousDispatcher.java:135)
        at org.jboss.resteasy.core.SynchronousDispatcher$$Lambda$453/0000000000000000.get(Unknown Source)
        at org.jboss.resteasy.core.interception.PreMatchContainerRequestContext.filter(PreMatchContainerRequestContext.java:356)
        at org.jboss.resteasy.core.SynchronousDispatcher.preprocess(SynchronousDispatcher.java:138)
        at org.jboss.resteasy.core.SynchronousDispatcher.invoke(SynchronousDispatcher.java:215)
        at org.jboss.resteasy.plugins.server.servlet.ServletContainerDispatcher.service(ServletContainerDispatcher.java:227)
        at org.jboss.resteasy.plugins.server.servlet.HttpServletDispatcher.service(HttpServletDispatcher.java:56)
        at org.jboss.resteasy.plugins.server.servlet.HttpServletDispatcher.service(HttpServletDispatcher.java:51)
        at javax.servlet.http.HttpServlet.service(HttpServlet.java:847)
        at io.undertow.servlet.handlers.ServletHandler.handleRequest(ServletHandler.java:74)
        at io.undertow.servlet.handlers.FilterHandler$FilterChainImpl.doFilter(FilterHandler.java:129)
        at org.keycloak.services.filters.KeycloakSessionServletFilter.doFilter(KeycloakSessionServletFilter.java:91)
        at org.keycloak.testsuite.TestKeycloakSessionServletFilter.doFilter(TestKeycloakSessionServletFilter.java:38)
        at io.undertow.servlet.core.ManagedFilter.doFilter(ManagedFilter.java:61)
        at io.undertow.servlet.handlers.FilterHandler$FilterChainImpl.doFilter(FilterHandler.java:131)
        at io.undertow.servlet.handlers.FilterHandler.handleRequest(FilterHandler.java:84)
        at io.undertow.servlet.handlers.security.ServletSecurityRoleHandler.handleRequest(ServletSecurityRoleHandler.java:62)
        at io.undertow.servlet.handlers.ServletChain$1.handleRequest(ServletChain.java:68)
        at io.undertow.servlet.handlers.ServletDispatchingHandler.handleRequest(ServletDispatchingHandler.java:36)
        at io.undertow.servlet.handlers.RedirectDirHandler.handleRequest(RedirectDirHandler.java:68)
        at io.undertow.servlet.handlers.security.SSLInformationAssociationHandler.handleRequest(SSLInformationAssociationHandler.java:132)
        at io.undertow.servlet.handlers.security.ServletAuthenticationCallHandler.handleRequest(ServletAuthenticationCallHandler.java:57)
        at io.undertow.server.handlers.PredicateHandler.handleRequest(PredicateHandler.java:43)
        at io.undertow.security.handlers.AbstractConfidentialityHandler.handleRequest(AbstractConfidentialityHandler.java:46)
        at io.undertow.servlet.handlers.security.ServletConfidentialityConstraintHandler.handleRequest(ServletConfidentialityConstraintHandler.java:64)
        at io.undertow.security.handlers.AuthenticationMechanismsHandler.handleRequest(AuthenticationMechanismsHandler.java:60)
        at io.undertow.servlet.handlers.security.CachedAuthenticatedSessionHandler.handleRequest(CachedAuthenticatedSessionHandler.java:77)
        at io.undertow.security.handlers.AbstractSecurityContextAssociationHandler.handleRequest(AbstractSecurityContextAssociationHandler.java:43)
        at io.undertow.server.handlers.PredicateHandler.handleRequest(PredicateHandler.java:43)
        at io.undertow.server.handlers.PredicateHandler.handleRequest(PredicateHandler.java:43)
        at io.undertow.servlet.handlers.ServletInitialHandler.handleFirstRequest(ServletInitialHandler.java:269)
        at io.undertow.servlet.handlers.ServletInitialHandler.access$100(ServletInitialHandler.java:78)
        at io.undertow.servlet.handlers.ServletInitialHandler$2.call(ServletInitialHandler.java:133)
        at io.undertow.servlet.handlers.ServletInitialHandler$2.call(ServletInitialHandler.java:130)
        at io.undertow.servlet.core.ServletRequestContextThreadSetupAction$1.call(ServletRequestContextThreadSetupAction.java:48)
        at io.undertow.servlet.core.ContextClassLoaderSetupAction$1.call(ContextClassLoaderSetupAction.java:43)
        at io.undertow.servlet.handlers.ServletInitialHandler.dispatchRequest(ServletInitialHandler.java:249)
        at io.undertow.servlet.handlers.ServletInitialHandler.access$000(ServletInitialHandler.java:78)
        at io.undertow.servlet.handlers.ServletInitialHandler$1.handleRequest(ServletInitialHandler.java:99)
        at io.undertow.server.Connectors.executeRootHandler(Connectors.java:370)
        at io.undertow.server.HttpServerExchange$1.run(HttpServerExchange.java:830)
        at org.jboss.threads.ContextClassLoaderSavingRunnable.run(ContextClassLoaderSavingRunnable.java:35)
        at org.jboss.threads.EnhancedQueueExecutor.safeRun(EnhancedQueueExecutor.java:1982)
        at org.jboss.threads.EnhancedQueueExecutor$ThreadBody.doRunTask(EnhancedQueueExecutor.java:1486)
        at org.jboss.threads.EnhancedQueueExecutor$ThreadBody.run(EnhancedQueueExecutor.java:1377)
        at java.base/java.lang.Thread.run(Thread.java:834)





-------------------------------------------------------------------------------
Test set: org.keycloak.testsuite.federation.ldap.noimport.LDAPGroupMapperNoImportTest
-------------------------------------------------------------------------------
Tests run: 9, Failures: 1, Errors: 0, Skipped: 1, Time elapsed: 4.747 s <<< FAILURE! - in org.keycloak.testsuite.federation.ldap.noimport.LDAPGroupMapperNoImportTest
test08_ldapOnlyGroupMappingsRanged(org.keycloak.testsuite.federation.ldap.noimport.LDAPGroupMapperNoImportTest)  Time elapsed: 0.187 s  <<< FAILURE!
java.lang.AssertionError
        at org.junit.Assert.fail(Assert.java:86)
        at org.junit.Assert.assertTrue(Assert.java:41)
        at org.junit.Assert.assertFalse(Assert.java:64)
        at org.junit.Assert.assertFalse(Assert.java:74)
        at org.keycloak.testsuite.federation.ldap.LDAPGroupMapperTest.lambda$test08_ldapOnlyGroupMappingsRanged$26a8868a$1(LDAPGroupMapperTest.java:729)
        at org.keycloak.testsuite.federation.ldap.LDAPGroupMapperTest$$Lambda$708/0000000000000000.run(Unknown Source)
        at org.keycloak.testsuite.rest.TestingResourceProvider.runOnServer(TestingResourceProvider.java:796)
        at jdk.internal.reflect.GeneratedMethodAccessor856.invoke(Unknown Source)
        at java.base/jdk.internal.reflect.DelegatingMethodAccessorImpl.invoke(DelegatingMethodAccessorImpl.java:43)
        at java.base/java.lang.reflect.Method.invoke(Method.java:566)
        at org.jboss.resteasy.core.MethodInjectorImpl.invoke(MethodInjectorImpl.java:138)
        at org.jboss.resteasy.core.ResourceMethodInvoker.internalInvokeOnTarget(ResourceMethodInvoker.java:535)
        at org.jboss.resteasy.core.ResourceMethodInvoker.invokeOnTargetAfterFilter(ResourceMethodInvoker.java:424)
        at org.jboss.resteasy.core.ResourceMethodInvoker.lambda$invokeOnTarget$0(ResourceMethodInvoker.java:385)
        at org.jboss.resteasy.core.ResourceMethodInvoker$$Lambda$454/0000000000000000.get(Unknown Source)
        at org.jboss.resteasy.core.interception.PreMatchContainerRequestContext.filter(PreMatchContainerRequestContext.java:356)
        at org.jboss.resteasy.core.ResourceMethodInvoker.invokeOnTarget(ResourceMethodInvoker.java:387)
        at org.jboss.resteasy.core.ResourceMethodInvoker.invoke(ResourceMethodInvoker.java:356)
        at org.jboss.resteasy.core.ResourceLocatorInvoker.invokeOnTargetObject(ResourceLocatorInvoker.java:150)
        at org.jboss.resteasy.core.ResourceLocatorInvoker.invoke(ResourceLocatorInvoker.java:104)
        at org.jboss.resteasy.core.SynchronousDispatcher.invoke(SynchronousDispatcher.java:440)
        at org.jboss.resteasy.core.SynchronousDispatcher.lambda$invoke$4(SynchronousDispatcher.java:229)
        at org.jboss.resteasy.core.SynchronousDispatcher$$Lambda$452/0000000000000000.run(Unknown Source)
        at org.jboss.resteasy.core.SynchronousDispatcher.lambda$preprocess$0(SynchronousDispatcher.java:135)
        at org.jboss.resteasy.core.SynchronousDispatcher$$Lambda$453/0000000000000000.get(Unknown Source)
        at org.jboss.resteasy.core.interception.PreMatchContainerRequestContext.filter(PreMatchContainerRequestContext.java:356)
        at org.jboss.resteasy.core.SynchronousDispatcher.preprocess(SynchronousDispatcher.java:138)
        at org.jboss.resteasy.core.SynchronousDispatcher.invoke(SynchronousDispatcher.java:215)
        at org.jboss.resteasy.plugins.server.servlet.ServletContainerDispatcher.service(ServletContainerDispatcher.java:227)
        at org.jboss.resteasy.plugins.server.servlet.HttpServletDispatcher.service(HttpServletDispatcher.java:56)
        at org.jboss.resteasy.plugins.server.servlet.HttpServletDispatcher.service(HttpServletDispatcher.java:51)
        at javax.servlet.http.HttpServlet.service(HttpServlet.java:847)
        at io.undertow.servlet.handlers.ServletHandler.handleRequest(ServletHandler.java:74)
        at io.undertow.servlet.handlers.FilterHandler$FilterChainImpl.doFilter(FilterHandler.java:129)
        at org.keycloak.services.filters.KeycloakSessionServletFilter.doFilter(KeycloakSessionServletFilter.java:91)
        at org.keycloak.testsuite.TestKeycloakSessionServletFilter.doFilter(TestKeycloakSessionServletFilter.java:38)
        at io.undertow.servlet.core.ManagedFilter.doFilter(ManagedFilter.java:61)
        at io.undertow.servlet.handlers.FilterHandler$FilterChainImpl.doFilter(FilterHandler.java:131)
        at io.undertow.servlet.handlers.FilterHandler.handleRequest(FilterHandler.java:84)
        at io.undertow.servlet.handlers.security.ServletSecurityRoleHandler.handleRequest(ServletSecurityRoleHandler.java:62)
        at io.undertow.servlet.handlers.ServletChain$1.handleRequest(ServletChain.java:68)
        at io.undertow.servlet.handlers.ServletDispatchingHandler.handleRequest(ServletDispatchingHandler.java:36)
        at io.undertow.servlet.handlers.RedirectDirHandler.handleRequest(RedirectDirHandler.java:68)
        at io.undertow.servlet.handlers.security.SSLInformationAssociationHandler.handleRequest(SSLInformationAssociationHandler.java:132)
        at io.undertow.servlet.handlers.security.ServletAuthenticationCallHandler.handleRequest(ServletAuthenticationCallHandler.java:57)
        at io.undertow.server.handlers.PredicateHandler.handleRequest(PredicateHandler.java:43)
        at io.undertow.security.handlers.AbstractConfidentialityHandler.handleRequest(AbstractConfidentialityHandler.java:46)
        at io.undertow.servlet.handlers.security.ServletConfidentialityConstraintHandler.handleRequest(ServletConfidentialityConstraintHandler.java:64)
        at io.undertow.security.handlers.AuthenticationMechanismsHandler.handleRequest(AuthenticationMechanismsHandler.java:60)
        at io.undertow.servlet.handlers.security.CachedAuthenticatedSessionHandler.handleRequest(CachedAuthenticatedSessionHandler.java:77)
        at io.undertow.security.handlers.AbstractSecurityContextAssociationHandler.handleRequest(AbstractSecurityContextAssociationHandler.java:43)
        at io.undertow.server.handlers.PredicateHandler.handleRequest(PredicateHandler.java:43)
        at io.undertow.server.handlers.PredicateHandler.handleRequest(PredicateHandler.java:43)
        at io.undertow.servlet.handlers.ServletInitialHandler.handleFirstRequest(ServletInitialHandler.java:269)
        at io.undertow.servlet.handlers.ServletInitialHandler.access$100(ServletInitialHandler.java:78)
        at io.undertow.servlet.handlers.ServletInitialHandler$2.call(ServletInitialHandler.java:133)
        at io.undertow.servlet.handlers.ServletInitialHandler$2.call(ServletInitialHandler.java:130)
        at io.undertow.servlet.core.ServletRequestContextThreadSetupAction$1.call(ServletRequestContextThreadSetupAction.java:48)
        at io.undertow.servlet.core.ContextClassLoaderSetupAction$1.call(ContextClassLoaderSetupAction.java:43)
        at io.undertow.servlet.handlers.ServletInitialHandler.dispatchRequest(ServletInitialHandler.java:249)
        at io.undertow.servlet.handlers.ServletInitialHandler.access$000(ServletInitialHandler.java:78)
        at io.undertow.servlet.handlers.ServletInitialHandler$1.handleRequest(ServletInitialHandler.java:99)
        at io.undertow.server.Connectors.executeRootHandler(Connectors.java:370)
        at io.undertow.server.HttpServerExchange$1.run(HttpServerExchange.java:830)
        at org.jboss.threads.ContextClassLoaderSavingRunnable.run(ContextClassLoaderSavingRunnable.java:35)
        at org.jboss.threads.EnhancedQueueExecutor.safeRun(EnhancedQueueExecutor.java:1982)
        at org.jboss.threads.EnhancedQueueExecutor$ThreadBody.doRunTask(EnhancedQueueExecutor.java:1486)
        at org.jboss.threads.EnhancedQueueExecutor$ThreadBody.run(EnhancedQueueExecutor.java:1377)
        at java.base/java.lang.Thread.run(Thread.java:834)












***********************



mvn  verify -f testsuite/integration-arquillian/pom.xml -Dtest=org.keycloak.testsuite.federation.ldap.LDAPUserLoginTest#loginLDAPUserAuthenticationSASLExternalEncryptionSSL -Djavax.net.ssl.keyStore=$WORKDIR/foo-keystore.p12 -Djavax.net.ssl.keyStorePassword=password -Djavax.net.ssl.javax.net.ssl.keyStoreType="PKCS12" -Dkeycloak.logging.level=debug
mvn  verify -f testsuite/integration-arquillian/pom.xml -Dtest=org.keycloak.testsuite.federation.ldap.LDAPUserLoginTest#loginLDAPUserAuthenticationSASLExternalEncryptionStartTLS -Djavax.net.ssl.keyStore=$WORKDIR/foo-keystore.p12 -Djavax.net.ssl.keyStorePassword=password -Djavax.net.ssl.javax.net.ssl.keyStoreType="PKCS12" -Dkeycloak.logging.level=debug


mvn exec:java -pl util/embedded-ldap/ -Dexec.mainClass=org.keycloak.util.ldap.LDAPEmbeddedServer -DenableSSL=true -DenableStartTLS=true  -Djavax.net.ssl.trustStore=$WORKDIR/truststore.p12 -Djavax.net.ssl.trustStorePassword=password -Djavax.net.ssl.javax.net.ssl.trustStoreType="PKCS12" -DkeystoreFile=$WORKDIR/server-keystore.p12 -DcertificatePassword=password

LD_PRELOAD=./libsslkeylog.so SSLKEYLOGFILE=wireshark-keys.log LDAPSASL_MECH=EXTERNAL LDAPCA_CERT=$PWD/certs/ca.pem LDAPTLS_CERT=certs/foo.pem LDAPTLS_KEY=certs/foo-key.pem ldapsearch -ZZ -H ldap://localhost:10389 -b ou=People,dc=keycloak,dc=org







works:

mvn -f testsuite/utils/pom.xml exec:java -Pkeycloak-server -Dkeycloak.migration.action=import -Dkeycloak.migration.provider=dir -Dkeycloak.migration.dir=$WORKDIR/keycloak/ -Djavax.net.ssl.trustStore=$WORKDIR/truststore.p12 -Djavax.net.ssl.trustStorePassword=password -Djavax.net.ssl.javax.net.ssl.trustStoreType="PKCS12" -Djavax.net.ssl.keyStore=$WORKDIR/foo-keystore.p12 -Djavax.net.ssl.keyStorePassword=password -Djavax.net.ssl.javax.net.ssl.keyStoreType="PKCS12" -Dresources


ldap://localhost:10389
ou=People,dc=keycloak,dc=org


mvn exec:java -pl util/embedded-ldap/ -Dexec.mainClass=org.keycloak.util.ldap.LDAPEmbeddedServer -DenableSSL=true -DenableStartTLS=true -Djavax.net.ssl.trustStore=$WORKDIR/truststore.p12 -Djavax.net.ssl.trustStorePassword=password -Djavax.net.ssl.javax.net.ssl.trustStoreType="PKCS12" -DkeystoreFile=$WORKDIR/server-keystore.p12 -DcertificatePassword=password


export WORKDIR=/home/tsaarni/work/keycloak-devenv
testsuite/integration-arquillian/tests/base/target/kcinit login --config $WORKDIR/configs

bwilson:password




testsuite/integration-arquillian/tests/base/src/test/resources/ldap/users.ldif

dn: uid=admin,ou=People,dc=keycloak,dc=org
objectclass: top
objectclass: person
objectclass: organizationalPerson
objectclass: inetOrgPerson
uid: admin
cn: admin
sn: admin
userCertificate:: MIICwDCCAaigAwIBAgIIFijxpwSjmU0wDQYJKoZIhvcNAQELBQAwDTELMAkGA1UEAxMCY2EwHhcNMjAwODA3MDkxNjA3WhcNMjEwODA3MDkxNjA3WjAfMQ8wDQYDVQQLEwZzeXN0ZW0xDDAKBgNVBAMTA2ZvbzCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBAMahcnlPe+sL5V9CKydSzukcLtkD5wdy/odnGIILDGx9HQv35K4UbtCktpD+ZyXtzolkAdblL8ilU0MgbRRg+fhvIwgQ8YCc6uNQXkD902+Wp1/rFXefPoVVu/9zJkY+ps95eMsUxcKOAr2qGYE30FGfMGO/OlhYw7iX45ZVLQm6P7Cnt2Vv7Fn1jKoIcGZ6wjulOuqunn/fYCNeL3dpWxhgNPerqwGJZm/yHT586zDQIwDYHFHsconH07Ni9Id/J+XfFtpFjU0pqbPiqV94sA2i5AS+rQFxDHuIzRzbilabi5y5x4Dwhs0c7PSOZjBnitwVlg0ouqApQs5R4l74/9UCAwEAAaMSMBAwDgYDVR0PAQH/BAQDAgWgMA0GCSqGSIb3DQEBCwUAA4IBAQAr2n93+DX/hHB3v1bCBaVItnPrlHaS3qiPUl4qrJvLDg4iaxiTW5Irib5VuwkUYgl0XNNWmrm36lePUAvQB4c38qL7ESriPhsQgOVQB7e7tE57nsHuqPJjarZG28TT9c8WSon7Ob5R/LhCR7NSYSrWWXgiHCjhW3VgHUNxSG40QL+Yk6TfjFxH3Pczoj7a8Ywvdth2MCMr/85+YPa9VAiVMhTd08NyBY+Bnpz7xXS1GJ0Egv+vx2UuiyfcojAJYjRTbcO8M8/TuP9x++QmwXVZEFASUQ9dx8auyUchFLihn9eOZm1msdbI/ODFWP8PXilLFRV6M+mirER+/Z0eiBmw

