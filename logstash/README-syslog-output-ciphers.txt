
# Add support for ssl_cipher_suites
# https://github.com/logstash-plugins/logstash-output-syslog/issues/74

rbenv install jruby-9.4.9.0
rbenv global jruby-9.4.9.0




./gradlew installDevelopmentGems
./gradlew bootstrap


# run tests

export OSS=true
export LOGSTASH_SOURCE=1
export LOGSTASH_PATH=~/work/logstash

bundle install
bundle exec rspec spec



export JAVA_OPTS="-javaagent:$HOME/work/tls-testapp/java/app/build/extract-tls-secrets-4.0.0.jar=wireshark-keys.log"
wireshark -i lo -k -f "port 9999" -o tls.keylog_file:$HOME/work/logstash-output-syslog/wireshark-keys.log
bundle exec rspec spec/outputs/syslog_tls_spec.rb --example "with custom ciphers"





rm -rf certs
mkdir -p certs
certyaml -d certs configs/certs.yaml

# build gem
rm *.gem
gem build logstash-output-syslog.gemspec


# in ~/work/logstash install the .gem
bin/logstash-plugin install --local ~/work/logstash-output-syslog/logstash-output-syslog-3.0.6.gem

# run logstash
bin/logstash -f ~/work/devenvs/logstash/configs/logstash-source-syslog-ciphers.conf



# run rsyslog
docker compose rm -f
docker compose build rsyslog
docker compose up rsyslog



export JAVA_OPTS="-javaagent:$HOME/work/tls-testapp/java/app/build/extract-tls-secrets-4.0.0.jar=wireshark-keys.log"
wireshark -i lo -k -f "port 6515" -o tls.keylog_file:$HOME/work/logstash/wireshark-keys.log
