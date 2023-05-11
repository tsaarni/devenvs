

rbenv install jruby-9.3.4.0
rbenv global jruby-9.3.4.0

rbenv uninstall -f jruby-9.3.4.0



mkdir -p certs
certyaml -d certs configs/certs.yaml



export OSS=true
export LOGSTASH_SOURCE=1
export LOGSTASH_PATH=$HOME/work/logstash





bin/logstash-plugin install logstash-input-stdin
bin/logstash-plugin install logstash-codec-json
bin/logstash-plugin install logstash-codec-json_lines


cd ~/work/logstash-output-syslog
gem build logstash-output-syslog.gemspec



# local development

cd ~/work/logstash
./gradlew installDevelopmentGems
rake bootstrap
rake plugin:install-default


bin/logstash-plugin install --local ~/work/logstash-output-syslog/logstash-output-syslog-3.0.5.gem

# First remove duplicate line from Gemfile 
#    gem "logstash-output-syslog"
# and then add
#    gem "logstash-output-syslog", :path => "/home/tsaarni/work/logstash-output-syslog"


# or to update just syslog plugin
bin/logstash-plugin install --no-verify logstash-output-syslog


#  or the default version
#    bin/logstash-plugin install logstash-output-syslog

cp -a ~/work/logstash-output-syslog/lib/logstash/outputs/syslog.rb ~/work/logstash/vendor/bundle/jruby/2.6.0/gems/logstash-output-syslog-3.0.5/lib/logstash/outputs/syslog.rb








bin/logstash -f ~/work/devenvs/logstash/configs/logstash-source-syslog.conf --log.level debug



# run tests

export OSS=true
export LOGSTASH_SOURCE=1
export LOGSTASH_PATH=/home/tsaarni/work/logstash

bundle install
bundle exec rspec spec




## test rsyslog container

docker-compose rm -f
docker-compose up rsyslog


# tcp
logger --server localhost --tcp --port 6514 -t myapp -p user.notice "foo"

# tls
ncat --listen 7777 --keep-open | openssl s_client -connect localhost:6515 -CAfile certs/server-ca.pem -cert certs/client.pem -key certs/client-key.pem
logger --server localhost --tcp --port 7777 -t myapp -p user.notice "foo"

# tls + crl
ncat --listen 7777 --keep-open | openssl s_client -connect localhost:6515 -CAfile certs/server-ca.pem -cert certs/client.pem -key certs/client-key.pem -CRL certs/server-ca-crl.pem -crl_check_all -verify_return_error
logger --server localhost --tcp --port 7777 -t myapp -p user.notice "foo"





# logstash-output-syslog plain output is broken
# also tests fail

https://github.com/logstash-plugins/logstash-output-syslog/issues/51
https://github.com/logstash-plugins/logstash-output-syslog/pull/55





*** Permission denied errors when installing?

   > /home/tsaarni/work/logstash/vendor/jruby/lib/ruby/gems/shared/gems/minitest-5.11.3/.autotest (Permission denied)
   > /home/tsaarni/work/logstash/vendor/jruby/lib/ruby/gems/shared/gems/minitest-5.11.3/History.rdoc (Permission denied)


rm -rf /home/tsaarni/work/logstash/vendor/jruby/lib/ruby/gems/shared/gems/minitest-5.11.3/

