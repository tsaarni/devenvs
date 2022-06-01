

# copy sample application to vertx repo
cp -a files/src/ ~/work/vert.x/


Choose debugger in vscode, launch App







# compile on command line to install vert.x to local maven repo for quarkus to use
mvn install -DskipTests=true
mvn clean install -DskipTests=true


cp -a target/vertx-core-4.3.0-SNAPSHOT.jar ~/work/keycloak/quarkus/server/target/lib/lib/main/io.vertx.vertx-core-4.3.0-SNAPSHOT.jar


# run unit tests
mvn test -Dtest=Http1xTLSTest
mvn test -Dtest=Http2TLSTest
mvn test -Dtest=Http2TLSTest#testSNISubjectAltenativeNameCNMatch1


mkdir -p certs

rm -f certs/*
certyaml -d certs configs/certs.yaml

# import certs to p12 and combine into single p12 file with multiple server certs
openssl pkcs12 -export -passout pass:secret -noiter -nomaciter -in certs/server.pem -inkey certs/server-key.pem -out certs/server.p12  -name server.127-0-0-1.nip.io
openssl pkcs12 -export -passout pass:secret -noiter -nomaciter -in certs/server2.pem -inkey certs/server2-key.pem -out certs/server2.p12 -name server2.127-0-0-1.nip.io
openssl pkcs12 -export -passout pass:secret -noiter -nomaciter -in certs/wildcard.pem -inkey certs/wildcard-key.pem -out certs/wildcard.p12 -name *.server2.127-0-0-1.nip.io
openssl pkcs12 -export -passout pass:secret -noiter -nomaciter -in certs/no-dns-name.pem -inkey certs/no-dns-name-key.pem -out certs/no-dns-name.p12 -name foo

for s in certs/server.p12 certs/server2.p12 certs/wildcard.p12 certs/no-dns-name.p12; do keytool -importkeystore -srckeystore $s -srcstoretype pkcs12 -srcstorepass secret -destkeystore certs/keystore.p12 -deststoretype pkcs12 -deststorepass secret; done

for s in certs/server.p12 certs/server2.p12 certs/wildcard.p12 certs/no-dns-name.p12; do keytool -importkeystore -srckeystore $s -srcstoretype pkcs12 -srcstorepass secret -destkeystore certs/keystore.jks -deststoretype jks -deststorepass secret; done

keytool -list -v -keystore certs/keystore.p12 -storepass secret
keytool -list -v -keystore certs/keystore.jks -storepass secret




http --verify certs/ca.pem https://localhost:8443/
http --verify certs/ca.pem https://server.127-0-0-1.nip.io:8443/
http --verify certs/ca.pem https://server2.127-0-0-1.nip.io:8443/

openssl s_client -connect localhost:8443 | openssl x509 -text -noout # returns the first certificate loaded: server.127-0-0-1.nip.io
openssl s_client -connect localhost:8443 -servername server.127-0-0-1.nip.io | openssl x509 -text -noout
openssl s_client -connect localhost:8443 -servername server2.127-0-0-1.nip.io | openssl x509 -text -noout
openssl s_client -connect localhost:8443 -servername wildcard.server2.127-0-0-1.nip.io | openssl x509 -text -noout
openssl s_client -connect localhost:8443 -servername wildcard.127.0.0.1.nip.io | openssl x509 -text -noout
openssl s_client -connect localhost:8443 -servername not-matching | openssl x509 -text -noout

openssl s_client -tls1_3 -connect localhost:8443  -requestCAfile certs/some-other-ca.pem


apps/sni-tester.sh localhost:8443 "" server.127-0-0-1.nip.io server2.127-0-0-1.nip.io wildcard.server2.127-0-0-1.nip.io wildcard.127.0.0.1.nip.io not-matching
apps/sni-tester.sh localhost:8443 "" host1 host2.com wildcard.host3.com host4.com www.host4.com host5.com wildcard.host5.com localhost not-matching


# convert jks to p12
keytool -importkeystore -srckeystore /home/tsaarni/work/vert.x/src/test/resources/tls/sni-keystore.jks -destkeystore certs/sni-keystore.p12 -srcstoretype JKS -deststoretype PKCS12 -srcstorepass wibble -deststorepass wibble

keytool -list -v -keystore /home/tsaarni/work/vert.x/src/test/resources/tls/sni-keystore.jks -storepass wibble
keytool -list -v -keystore certs/sni-keystore.p12 -storepass wibble



###########################
#
# Code layout
#


https://github.com/eclipse-vertx/vert.x/issues/3780



API for configuration options
https://vertx.io/docs/apidocs/index.html?io/vertx/core/net/

# Interfaces
src/main/java/io/vertx/core/net/KeyCertOptions.java         # provide KeyManagerFactory, function that returns KeyManager per hostname,
                                                            # and convert KeyManager to KeyCertOptions
src/main/java/io/vertx/core/net/TrustOptions.java           # provide TrustManagerFactory, function that returns TrustManager per hostname
                                                            # and convert TrustManager to TrustOptions

# Concrete
src/main/java/io/vertx/core/net/PemKeyCertOptions.java      # pem, file or Buffer (not for trusted certs)

src/main/java/io/vertx/core/net/KeyStoreOptionsBase.java    # base class for KeyStore based implementations
                                                            # implements both KeyCertOptions and TrustOptions
src/main/java/io/vertx/core/net/KeyStoreOptions.java        # both jks and pkcs12, file or Buffer
src/main/java/io/vertx/core/net/JksOptions.java             # convenience: sets jks keytype
src/main/java/io/vertx/core/net/PfxOptions.java             # convenience: sets pkcs12 keytype

src/main/java/io/vertx/core/net/KeyManagerFactoryOptions.java  # user provides keymanager


# The above options are set to e.g. HttpServerOptions -> NetServerOptions -> TCPSSLOptions
src/main/java/io/vertx/core/net/TCPSSLOptions.java by calling


  .setKeyCertOptions(KeyCertOptions options)

  .setKeyStoreOptions(JksOptions options)
  .setPfxKeyCertOptions(PfxOptions options)
  .setPemKeyCertOptions(PemKeyCertOptions options)

  .setTrustOptions(TrustOptions options)
  .setTrustStoreOptions(JksOptions options)
  .setPemTrustOptions(PemTrustOptions options)
  .setPfxTrustOptions(PfxOptions options)


# TCPSSLOptions will be then used by TCPServerBase to construct SSLHelper
src/main/java/io/vertx/core/net/impl/TCPServerBase.java    # constructs SSLHelper

src/main/java/io/vertx/core/net/impl/SSLHelper.java        # uses the options


src/main/java/io/vertx/core/net/impl/KeyStoreHelper.java   # creates the keymanagers using the options





# Key Managers and Key Stores
https://tersesystems.com/blog/2018/09/08/keymanagers-and-keystores/





# How to Run Blocking Code in Vert.x
https://dzone.com/articles/how-to-run-blocking-code-in-vertx

# vertx copyInternal (non-bloking)
https://github.com/eclipse-vertx/vert.x/blob/96c360b02f75df7ba6e491c3ac04953ea942c085/src/main/java/io/vertx/core/file/impl/FileSystemImpl.java#L625

# NOTE: the above works only for short blocking code
#   "Long blocking operations should use a dedicated thread managed by the application"
https://vertx.io/docs/vertx-core/java/#blocking_code



# filesystem watchers

# netflix
https://github.com/ReactiveX/RxJavaFileUtils/blob/master/src/main/java/rx/fileutils/FileSystemWatcher.java
https://github.com/ReactiveX/RxJavaFileUtils/blob/master/src/test/java/rx/fileutils/FileSystemEventOnSubscribeTest.java

https://github.com/helmbold/rxfilewatcher/blob/master/src/main/java/de/helmbold/rxfilewatcher/PathObservables.java

https://github.com/alexvictoor/netty-livereload/blob/master/src/main/java/com/github/alexvictoor/livereload/FileSystemWatcher.java


# reactive programming in quarkus
https://quarkus.io/guides/mutiny-primer

https://github.com/smallrye/smallrye-mutiny-vertx-bindings
https://smallrye.io/smallrye-mutiny/





# Force updated version in keycloak

diff --git a/pom.xml b/pom.xml
index d7e5e94e72..5a963add6b 100644
--- a/pom.xml
+++ b/pom.xml
@@ -285,6 +285,12 @@
     <dependencyManagement>

         <dependencies>
+<!-- https://mvnrepository.com/artifact/io.vertx/vertx-core -->
+<dependency>
+    <groupId>io.vertx</groupId>
+    <artifactId>vertx-core</artifactId>
+    <version>4.3.0-SNAPSHOT</version>
+</dependency>

             <dependency>
                 <groupId>org.keycloak</groupId>










Reasons for NewSunX509

- key type for selection logic

https://hg.openjdk.java.net/jdk/jdk/file/ee1d592a9f53/src/java.base/share/classes/sun/security/ssl/X509KeyManagerImpl.java#l699
