

export OSS=true
export LOGSTASH_SOURCE=1
export LOGSTASH_PATH=$HOME/work/logstash


# build gem
gem build logstash-output-syslog.gemspec


# in ~/work/logstash install the .gem
bin/logstash-plugin install --local ~/work/logstash-output-syslog/logstash-output-syslog-3.0.5.gem

# run logstash
bin/logstash -f ~/work/devenvs/logstash/configs/logstash-source-syslog-socketopts.conf



# run rsyslog
docker-compose rm -f
docker-compose up rsyslog
