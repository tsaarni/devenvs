




git clone https://github.com/elastic/logstash.git
cd logstash
rbenv install  # installs version defined in logstash/.ruby-version

# activate ruby version
rbenv global jruby-9.3.4.0
ruby --version


export OSS=true
export LOGSTASH_SOURCE=1
export LOGSTASH_PATH=/home/tsaarni/work/logstash


gem install rake
gem install bundler

rake bootstrap




# if getting error
#   Gem::GemNotFoundException: can't find gem rake (>= 0.a) with executable rake
#
# or
#
# Execution failed for task ':bootstrap'.
# > (VersionConflict) Bundler could not find compatible versions for gem "aws-sdk-core":
#    In Gemfile:
#      logstash-output-sns was resolved to 4.0.8, which depends on
#        logstash-mixin-aws (>= 1.0.0) was resolved to 4.0.2, which depends on
#
# then just clean the repo from modified and untracked files and start from the beginning


bin/logstash-plugin install logstash-input-stdin
bin/logstash-plugin install logstash-codec-json
bin/logstash-plugin install logstash-codec-json_lines


rbenv uninstall jruby-9.3.4.0





# local development

cd /home/tsaarni/work/logstash
rake bootstrap
rake plugin:install-default


echo 'gem "logstash-output-syslog", :path => "/home/tsaarni/work/logstash-output-syslog"' >> Gemfile

bin/logstash-plugin install --no-verify
bin/logstash-plugin install --local --no-verify logstash-output-syslog


cp -a /home/tsaarni/work/logstash-output-syslog/lib/logstash/outputs/syslog.rb /home/tsaarni/work/logstash/vendor/bundle/jruby/2.6.0/gems/logstash-output-syslog-3.0.5/lib/logstash/outputs/syslog.rb

bin/logstash -f /home/tsaarni/work/devenvs/logstash/logstash-source-syslog.conf --log.level debug



# run tests

export OSS=true
export LOGSTASH_SOURCE=1
export LOGSTASH_PATH=/home/tsaarni/work/logstash

bundle install
bundle exec rspec spec




## test rsyslog container

docker-compose up rsyslog


# tcp
logger --server localhost --tcp --port 6514 -t myapp -p user.notice "foo"

# tls
ncat --listen 7777 --keep-open | openssl s_client -connect localhost:6515 -CAfile certs/server-ca.pem -cert certs/client.pem -key certs/client-key.pem
logger --server localhost --tcp --port 7777 -t myapp -p user.notice "foo"






# logstash-output-syslog plain output is broken
# also tests fail

https://github.com/logstash-plugins/logstash-output-syslog/issues/51
https://github.com/logstash-plugins/logstash-output-syslog/pull/55
