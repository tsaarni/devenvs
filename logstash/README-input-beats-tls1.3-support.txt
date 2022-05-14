



# Use jruby from same release track as logstash
rbenv install jruby-9.2.20.1
rbenv global jruby-9.2.20.1



# logstash source must be available
git clone https://github.com/elastic/logstash.git

export OSS=true
export LOGSTASH_SOURCE=1
export LOGSTASH_PATH=$PWD
rake bootstrap

# install logstash-core and api gems to gem cache
gem build logstash-core/logstash-core.gemspec 
gem build logstash-core-plugin-api/logstash-core-plugin-api.gemspec 
gem install logstash-core-8.2.0-java.gem
gem install logstash-core-plugin-api-2.1.16-java.gem






# Set env variables that point out local copy of logstash 
export LOGSTASH_PATH=/home/tsaarni/work/logstash
export LOGSTASH_SOURCE=1

# download deps
bundle install
bundle exec rake test:integration:setup

# compile java
rake vendor


# run unit test suite
bundle exec rspec

# run integration test suite
bundle exec rspec spec --tag integration -fd 

# run only test with string in description
bundle exec rspec spec --tag integration -fd -e "minimum protocol version"  
DEBUG=1 bundle exec rspec spec --tag integration -fd -e "minimum protocol version"


# run java unittests
gradle test
gradle test --tests org.logstash.netty.SslContextBuilderTest







#
# Testing with upstream container by overwriting some individual files
#


rm -rf certs
mkdir -p certs
certyaml --destination certs configs/certs.yaml


# compile the plugin and overwrite files e upstream container
mkdir -p docker/logstash/gems/logstash-input-beats-6.2.6-java/vendor/jar-dependencies/org/logstash/beats/logstash-input-beats/6.2.6
mkdir -p docker/logstash/gems/logstash-input-beats-6.2.6-java/lib/logstash/inputs/beats
cp -a SOURCE docker/logstash/gems/logstash-input-beats-6.2.6-java/vendor/jar-dependencies/org/logstash/beats/logstash-input-beats/6.2.6/logstash-input-beats-6.2.6.jar
cp -a SOURCE docker/logstash/gems/logstash-input-beats-6.2.6-java/lib/logstash/inputs/beats/tls.rb
 




docker-compose up

openssl s_client -cert certs/client.pem --key certs/client-key.pem -CAfile certs/server-ca.pem -connect localhost:12345 -tls1_3
sslyze --cert certs/client.pem  --key certs/client-key.pem  localhost:12345

