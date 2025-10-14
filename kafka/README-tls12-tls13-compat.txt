https://issues.apache.org/jira/browse/ZOOKEEPER-4987



cd ~/work/devenvs/zookeeper

# create truststore for zookeeper server
keytool -importcert \
  -file certs/cacert.pem \
  -keystore certs/truststore.p12 \
  -storetype PKCS12 \
  -alias CARoot \
  -storepass password \
  -noprompt


# create keystore for kafka zookeeper client

openssl pkcs12 -export \
  -in certs/clicert.pem \
  -inkey certs/clicert-key.pem \
  -certfile certs/cacert.pem \
  -out certs/kafka-client-keystore.p12 \
  -name kafka-client \
  -password pass:password


cd ~/work/kafka

diff --git a/gradle/dependencies.gradle b/gradle/dependencies.gradle
index 4e6f83fade..c7880ab08d 100644
--- a/gradle/dependencies.gradle
+++ b/gradle/dependencies.gradle
@@ -159,7 +159,8 @@ versions += [
   snappy: "1.1.10.5",
   spotbugs: "4.8.6",
   zinc: "1.9.2",
-  zookeeper: "3.8.4",
+//  zookeeper: "3.8.4",
+  zookeeper: "3.9.2",
   // When updating the zstd version, please do as well in docker/native/native-image-configs/resource-config.json
   // Also make sure the compression levels in org.apache.kafka.common.record.CompressionType are still valid
   zstd: "1.5.6-4",



./gradlew jar



bin/kafka-server-start.sh ~/work/devenvs/kafka/configs/server.properties



wireshark -i lo -k --display-filter "tcp.port==2281"
