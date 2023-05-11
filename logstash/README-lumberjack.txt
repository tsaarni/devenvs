
rbenv global jruby-9.3.4.0


# create test certs
cd ~/work/devenv/logstash

mkdir -p certs
certyaml -d certs configs/certs.yaml
cat certs/client-ca-crl.pem certs/server-ca-crl.pem > certs/crl.pem


# compile and install ruby-lumberjack
cd ~/work/ruby-lumberjack
bundle install
gem build jls-lumberjack.gemspec





cd ~/work/logstash-output-lumberjack/

export OSS=true
export LOGSTASH_SOURCE=1
export LOGSTASH_PATH=~/work/logstash
bundle install
gem build logstash-output-lumberjack.gemspec





cd ~/work/logstash

bin/logstash-plugin install --no-verify ~/work/logstash-output-lumberjack/logstash-output-lumberjack-3.1.9.gem

# override packages with a local version
####gem install --install-dir vendor/bundle/jruby/2.6.0/ ~/work/jruby-openssl/pkg/jruby-openssl-0.14.1.cr2-SNAPSHOT.gem
gem install --no-document --install-dir vendor/bundle/jruby/2.6.0/ ~/work/ruby-lumberjack/jls-lumberjack-0.0.26.gem


bin/logstash -f ~/work/devenvs/logstash/configs/logstash-source-lumberjack.conf --log.level debug

LS_JAVA_OPTS=-Djruby.openssl.debug=true

docker-compose rm -f
docker-compose up logstash-destination
