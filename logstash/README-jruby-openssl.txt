
# build
cd ~/work/jruby-openssl
mvn clean package -Dmaven.test.skip=true



cp -a ~/work/jruby-openssl-test/src/main/java/org/jruby/ext/openssl/ src/main/java/org/jruby/ext/



# install to logstash jruby environment
gem install /home/tsaarni/work/jruby-openssl/pkg/jruby-openssl-0.14.1.cr2-SNAPSHOT.gem



# test
irb

require "openssl"
OpenSSL::PKey::read(File.read("/home/tsaarni/work/devenvs/logstash/certs/client-key.pem"))
OpenSSL::PKey::read(File.read("/home/tsaarni/work/devenvs/logstash/certs/server-key.pem"))
