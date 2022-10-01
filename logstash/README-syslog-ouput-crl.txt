

mkdir -p certs
certyaml -d certs configs/certs.yaml



export OSS=true
export LOGSTASH_SOURCE=1
export LOGSTASH_PATH=/home/tsaarni/work/logstash





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

docker-compose rm -f
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
