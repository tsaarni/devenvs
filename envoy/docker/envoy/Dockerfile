FROM ubuntu:20.04

RUN apt-get update \
    && apt-get upgrade -y \
    && apt-get install -y ca-certificates python3 \
    && apt-get autoremove -y \
    && apt-get clean \
    && rm -rf /tmp/* /var/tmp/* \
    && rm -rf /var/lib/apt/lists/*

RUN mkdir -p /etc/envoy

ADD envoy /usr/local/bin/envoy
#ADD restarter /usr/local/bin/

EXPOSE 10000

COPY docker-entrypoint.sh /
ENTRYPOINT ["/docker-entrypoint.sh"]
CMD ["envoy"]
