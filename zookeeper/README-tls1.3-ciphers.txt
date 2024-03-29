

https://issues.apache.org/jira/browse/ZOOKEEPER-4415



# run unit tests
mvn clean install -DskipTests=true

mvn test   # run all

# to run only X509UtilTest use vscode (for some reason -Dtest=X509UtilTest does not work)





# generate certs
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






