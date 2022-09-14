FROM alpine:3.16
RUN apk add -U rsyslog rsyslog-tls
COPY files /
CMD ["/docker-entrypoint.sh"]
