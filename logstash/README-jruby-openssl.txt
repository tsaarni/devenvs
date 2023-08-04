
# build
cd ~/work/jruby-openssl
mvn clean package -Dmaven.test.skip=true


# copy from test branch
cp -a ~/work/jruby-openssl-test/src/main/java/org/jruby/ext/openssl/ src/main/java/org/jruby/ext/

# copy to logstash
##cp ~/work/jruby-openssl/lib/jopenssl.jar ./vendor/jruby/lib/ruby/stdlib/jopenssl.jar
cp ~/work/jruby-openssl/lib/jopenssl.jar ./vendor/bundle/jruby/3.1.0/gems/jruby-openssl-0.14.2-java/lib/jopenssl.jar

# install to logstash jruby environment
gem install --no-document /home/tsaarni/work/jruby-openssl/pkg/jruby-openssl-0.14.1.cr2-SNAPSHOT.gem
gem install --no-document --install-dir vendor/bundle/jruby/2.6.0/ ~/work/jruby-openssl/pkg/jruby-openssl-0.14.*


# test
irb

require "openssl"
OpenSSL::PKey::read(File.read("/home/tsaarni/work/devenvs/logstash/certs/client-key.pem"))
OpenSSL::PKey::read(File.read("/home/tsaarni/work/devenvs/logstash/certs/server-key.pem"))



openssl req -x509 -nodes -days 3650 -newkey ec -pkeyopt ec_paramgen_curve:prime256v1 -subj "/CN=www.example.com" -keyout certs/key.pem -out certs/cert.pem
openssl req -x509 -nodes -days 3650 -newkey ec -pkeyopt ec_paramgen_curve:secp384r1 -subj "/CN=www.example.com" -keyout certs/key.pem -out certs/cert.pem
openssl req -x509 -nodes -days 3650 -newkey ec -pkeyopt ec_paramgen_curve:secp521r1 -subj "/CN=www.example.com" -keyout certs/key.pem -out certs/cert.pem

OpenSSL::PKey::read(File.read("/home/tsaarni/work/devenvs/logstash/certs/key.pem"))


# unittest
# from: .github/workflows/ci-test.yml

jruby -S bundle install
jruby -rbundler/setup -S rake test_prepare
jruby -rbundler/setup -S rake test

jruby -Ilib:src/test/ruby src/test/ruby/ec/test_ec.rb

# to recompile (running test again will pick the newly compiled version automatically)
mvn clean package -Dmaven.test.skip=true






mvn clean package -Dmaven.test.skip=true
gem install --no-document /home/tsaarni/work/jruby-openssl/pkg/jruby-openssl-0.14.1.cr2-SNAPSHOT.gem

irb

require "openssl"
key = OpenSSL::PKey::EC.new(File.read("/home/tsaarni/work/jruby-openssl/src/test/ruby/ec/private_key.pem"))
data = 'abcd'
digest = OpenSSL::Digest::SHA256.new
sig = key.sign(digest, data)
key.verify(digest, sig, data)




# convert PKCS#8 to encrypted
openssl pkcs8 -in private_key_pkcs8.pem -topk8 -out private_key_pkcs8_enc.pem
openssl pkcs8 -in private_key_pkcs8.pem -topk8 -out private_key_pkcs8_enc.pem -v1 PBE-MD5-DES  # PBES1






### build failure

[INFO] BUILD FAILURE

[ERROR] Plugin org.codehaus.mojo:build-helper-maven-plugin:1.9 or one of its dependencies could not be resolved: Failed to read artifact descriptor for org.codehaus.mojo:build-helper-maven-plugin:jar:1.9: 1 problem was encountered while building the effective model
[ERROR] [FATAL] Non-parseable POM /home/tsaarni/.m2/repository/org/codehaus/mojo/build-helper-maven-plugin/1.9/build-helper-maven-plugin-1.9.pom: UTF-8 BOM plus xml decl of ISO-8859-1 is incompatible (position: START_DOCUMENT seen <?xml version="1.0" encoding="ISO-8859-1"... @1:42)  @ line 1, column 42


# WORKAROUND
# change to old maven:

sdk use maven 3.8.7    # changes only in this terminal
