
# ssl_supported_protocols
https://github.com/logstash-plugins/logstash-output-syslog/pull/77


# References
https://github.com/logstash-plugins/logstash-mixin-http_client/blob/9478c9991d5f2f037921eae98a854d6fb7321b9f/lib/logstash/plugin_mixins/http_client.rb#L120
https://github.com/logstash-plugins/logstash-output-tcp/blob/4b2265b4b188140e1c83576f2513e94e5f1734bd/lib/logstash/outputs/tcp.rb#L67
https://github.com/logstash-plugins/logstash-input-http/blob/879eb1472c90c12775cfa8d431d0564bdb2ef177/lib/logstash/inputs/http.rb#L113
https://github.com/logstash-plugins/logstash-input-tcp/blob/ebbfa05977d2838201ac49d35425a1f34e5c1e3d/lib/logstash/inputs/tcp.rb#L133
https://github.com/logstash-plugins/logstash-output-elasticsearch/blob/3ef3c0c85c9e5eb928a92a89ef8661610c4c54cb/lib/logstash/plugin_mixins/elasticsearch/api_configs.rb#L83


# run tests

export OSS=true
export LOGSTASH_SOURCE=1
export LOGSTASH_PATH=~/work/logstash

bundle install
bundle exec rspec spec


##export JAVA_OPTS="-javaagent:$HOME/work/tls-testapp/java/app/build/extract-tls-secrets-4.0.0.jar=wireshark-keys.log"
##wireshark -i lo -k -f "port 9999" -o tls.keylog_file:$HOME/work/logstash-output-syslog/wireshark-keys.log
##bundle exec rspec spec/outputs/syslog_tls_spec.rb --example "with TLS versions"


export JAVA_OPTS="-javaagent:$HOME/work/tls-testapp/java/app/build/extract-tls-secrets-4.0.0.jar=wireshark-keys.log"
wireshark -i lo -k -o tls.keylog_file:$HOME/work/logstash-output-syslog/wireshark-keys.log
bundle exec rspec spec/outputs/syslog_tls_spec.rb --example "with TLS versions"

# use display filter
tls.handshake.type == 1
# then see
#  Extension: supported_versions

