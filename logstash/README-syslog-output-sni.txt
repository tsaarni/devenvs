
# create certs
mkdir -p certs
certyaml -d certs configs/certs.yaml



export OSS=true
export LOGSTASH_SOURCE=1
export LOGSTASH_PATH=$HOME/work/logstash


# build gem
gem build logstash-output-syslog.gemspec


# in ~/work/logstash install the .gem
bin/logstash-plugin install --local ~/work/logstash-output-syslog/logstash-output-syslog-3.0.5.gem

# run logstash
bin/logstash -f ~/work/devenvs/logstash/configs/logstash-source-syslog.conf --log.level debug


# run rsyslog
docker-compose rm -f
docker-compose up rsyslog


# observe server_name extension with wireshark




#### Notes

# Jruby OpenSSL will set SNI if @hostname attribute is set in the socket
https://github.com/jruby/jruby-openssl/blob/a4f2d52bffca44f2a02ee9627b895fb17a74285d/src/main/java/org/jruby/ext/openssl/SSLSocket.java#L212-L214
