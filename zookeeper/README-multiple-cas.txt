
# Zookeeper doesn't support multiple ca into truststore
# https://issues.apache.org/jira/browse/ZOOKEEPER-4990


# Loading multiple trusted certificates with identical subject names from a PEM bundle fails
# https://issues.apache.org/jira/browse/ZOOKEEPER-4992

# ZOOKEEPER-4992: Avoid overriding same-subject certs in PEM trust store
# https://github.com/apache/zookeeper/pull/2336



# run unit tests
mvn clean install -DskipTests=true

mvn checkstyle:check  # coding style check
mvn test              # run all tests


# run individual tests in vscode   (<<<<IMPORTANT FOR TESTING)
# but compile with 
#   mvn install -DskipTests=true
# before running individual tests in vscode


ssl.keyStore.location=/var/lib/zookeeper/secrets/server/zk-server-keystore.jks
ssl.quorum.keyStore.password=xxxxxxxxxxxxxxxxx
ssl.quorum.trustStore.password=xxxxxxxxxxxxxxxxx
ssl.quorum.keyStore.location=/var/lib/zookeeper/secrets/server/zk-server-keystore.jks
ssl.quorum.trustStore.location=/var/lib/zookeeper/secrets/server/zk-server-truststore.jks
ssl.trustStore.password=xxxxxxxxxxxxxxxxx
ssl.keyStore.password=xxxxxxxxxxxxxxxxx


zkcli.sh




# generate certs
rm -rf certs
mkdir -p certs
certyaml -d certs configs/certs.yaml

# Create PEM bundles
cat certs/cert1.pem certs/cert1-key.pem > certs/cert1-combined.pem   # create bundle with cert+key
cat certs/cert2.pem certs/cert2-key.pem > certs/cert2-combined.pem
cat certs/clicert.pem certs/clicert-key.pem > certs/clicert-combined.pem
cat certs/clicert2.pem certs/clicert2-key.pem > certs/clicert2-combined.pem

cat certs/cacert.pem certs/cacert2.pem > certs/server-ca-bundle.pem
cat certs/clientcacert.pem certs/clientcacert2.pem > certs/client-ca-bundle.pem

# create jks truststore with multiple ca certs for server to verify clients
keytool -importcert -file certs/clientcacert.pem -alias ca1 -keystore certs/client-ca-truststore.jks -storepass my-password -noprompt
keytool -importcert -file certs/clientcacert2.pem -alias ca2 -keystore certs/client-ca-truststore.jks -storepass my-password -noprompt

# create pkcs12 keystore for client
openssl pkcs12 -export -passout pass:my-password -noiter -nomaciter -in certs/clicert.pem -inkey certs/clicert-key.pem -out certs/clicert-keystore.p12
openssl pkcs12 -export -passout pass:my-password -noiter -nomaciter -in certs/clicert2.pem -inkey certs/clicert2-key.pem -out certs/clicert2-keystore.p12

# create jks truststore with multiple ca certs for server to verify clients
keytool -importcert -file certs/cacert.pem -alias ca1 -keystore certs/server-ca-truststore.jks -storepass my-password -noprompt
keytool -importcert -file certs/cacert2.pem -alias ca2 -keystore certs/server-ca-truststore.jks -storepass my-password -noprompt


# list certs in jks truststore
keytool -list -keystore certs/client-ca-truststore.jks -storepass my-password
keytool -list -keystore certs/server-ca-truststore.jks -storepass my-password


cd ~/work/zookeeper
bin/zkServer.sh --config ~/work/devenvs/zookeeper/configs start-foreground

export CLIENT_JVMFLAGS="
-Dzookeeper.clientCnxnSocket=org.apache.zookeeper.ClientCnxnSocketNetty
-Dzookeeper.client.secure=true
-Dzookeeper.ssl.keyStore.location=$HOME/work/devenvs/zookeeper/certs/clicert-keystore.p12
-Dzookeeper.ssl.keyStore.password=my-password
-Dzookeeper.ssl.trustStore.location=$HOME/work/devenvs/zookeeper/certs/cacert.pem
"


bin/zkCli.sh -server localhost:2281


export CLIENT_JVMFLAGS="
-Dzookeeper.clientCnxnSocket=org.apache.zookeeper.ClientCnxnSocketNetty
-Dzookeeper.client.secure=true
-Dzookeeper.ssl.keyStore.location=$HOME/work/devenvs/zookeeper/certs/clicert2-keystore.p12
-Dzookeeper.ssl.keyStore.password=my-password
-Dzookeeper.ssl.trustStore.location=$HOME/work/devenvs/zookeeper/certs/cacert2.pem
"

bin/zkCli.sh -server localhost:2281


export CLIENT_JVMFLAGS="
-Dzookeeper.clientCnxnSocket=org.apache.zookeeper.ClientCnxnSocketNetty
-Dzookeeper.client.secure=true
-Dzookeeper.ssl.keyStore.location=$HOME/work/devenvs/zookeeper/certs/clicert-keystore.p12
-Dzookeeper.ssl.keyStore.password=my-password
-Dzookeeper.ssl.trustStore.location=$HOME/work/devenvs/zookeeper/certs/server-ca-bundle.pem
"

bin/zkCli.sh -server localhost:2281



export CLIENT_JVMFLAGS="
-Dzookeeper.clientCnxnSocket=org.apache.zookeeper.ClientCnxnSocketNetty
-Dzookeeper.client.secure=true
-Dzookeeper.ssl.keyStore.location=$HOME/work/devenvs/zookeeper/certs/clicert-keystore.p12
-Dzookeeper.ssl.keyStore.password=my-password
-Dzookeeper.ssl.trustStore.location=$HOME/work/devenvs/zookeeper/certs/server-ca-truststore.jks
-Dzookeeper.ssl.trustStore.password=my-password
"

bin/zkCli.sh -server localhost:2281

# Test commands to run in zkCli.sh to verify it is connected and authenticated

ls /
create /testnode "hello"
get /testnode
delete /testnode
