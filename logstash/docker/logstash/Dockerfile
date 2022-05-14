# docker build docker/logstash/ -t logstash-oss-tls:8.1.0

FROM docker.elastic.co/logstash/logstash-oss:8.1.0

# add forked versions of the gems
ADD gems/ /usr/share/logstash/vendor/bundle/jruby/2.5.0/gems/

