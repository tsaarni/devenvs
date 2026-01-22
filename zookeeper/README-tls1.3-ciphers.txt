# Zookeeper 3.7.0 : The client supported protocol versions [TLSv1.3] are not accepted by server preferences
# https://issues.apache.org/jira/browse/ZOOKEEPER-4415
# https://github.com/apache/zookeeper/pull/1919




# run unit tests
mvn clean install -DskipTests=true

mvn test   # run all

# to run only X509UtilTest use vscode (for some reason -Dtest=X509UtilTest does not work)





# generate certs
rm -rf certs
mkdir -p certs
certyaml -d certs configs/certs.yaml
cat certs/cert1.pem certs/cert1-key.pem > certs/cert1-combined.pem   # create bundle with cert+key


# run server under vscode
cd ~/work/zookeeper
mkdir -p .vscode
cp ~/work/devenvs/zookeeper/launch.json .vscode



# Run from command line
mvn clean install -DskipTests=true
bin/zkServer.sh --config ~/work/devenvs/zookeeper/configs start-foreground



# test connection
openssl s_client --connect localhost:2281 --cert certs/clicert.pem --key certs/clicert-key.pem --CAfile certs/cacert.pem -tls1_3
openssl s_client --connect localhost:2281 --cert certs/clicert.pem --key certs/clicert-key.pem --CAfile certs/cacert.pem  -min_protocol TLSv1.2 -max_protocol TLSv1.3

sslyze --cert certs/clicert.pem --key certs/clicert-key.pem localhost:2281




###
### Misc
###


# default ciphers
zookeeper-server/src/main/java/org/apache/zookeeper/common/X509Util.java

# SSLContext creation
zookeeper-server/src/main/java/org/apache/zookeeper/common/SSLContextAndOptions.java  createNettyJdkSslContext()




# Available ciphers vs JDK versions

https://docs.oracle.com/javase/8/docs/technotes/guides/security/StandardNames.html#ciphersuites
https://docs.oracle.com/en/java/javase/11/docs/specs/security/standard-names.html#jsse-cipher-suite-names
https://docs.oracle.com/en/java/javase/17/docs/specs/security/standard-names.html#jsse-cipher-suite-names








## zookeeper client fails to fallback to tls1.2 when tls1.3 ciphers are not correct / zookeeper client fails to fallback to tls1.3 when tls1.2 ciphers are not correct
## https://issues.apache.org/jira/projects/ZOOKEEPER/issues/ZOOKEEPER-4987


The following Kafka settings work for the Zookeeper client when Zookeeper server was configured with the default enabled protocols (TLSv1.3 and TLSv1.2) but server is restricted to only TLSv1.2 ciphers:

zookeeper.ssl.client.enable=true
zookeeper.ssl.protocol = TLSv1.3
zookeeper.ssl.enabled.protocols = TLSv1.2, TLSv1.3
zookeeper.clientCnxnSocket=org.apache.zookeeper.ClientCnxnSocketNetty

￼
