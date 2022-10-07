
# build
cd ~/work/jruby-openssl
mvn clean package -Dmaven.test.skip=true


# copy from test branch
cp -a ~/work/jruby-openssl-test/src/main/java/org/jruby/ext/openssl/ src/main/java/org/jruby/ext/

# copy to logstash
cp ~/work/jruby-openssl/lib/jopenssl.jar ./vendor/jruby/lib/ruby/stdlib/jopenssl.jar


# install to logstash jruby environment
gem install --no-document /home/tsaarni/work/jruby-openssl/pkg/jruby-openssl-0.14.1.cr2-SNAPSHOT.gem



# test
irb

require "openssl"
OpenSSL::PKey::read(File.read("/home/tsaarni/work/devenvs/logstash/certs/client-key.pem"))
OpenSSL::PKey::read(File.read("/home/tsaarni/work/devenvs/logstash/certs/server-key.pem"))







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

