FROM ubuntu:22.04

ENV DEBIAN_FRONTEND noninteractive

RUN apt-get update && \
    apt-get install --yes --no-install-recommends \
        ldap-utils \
        slapd \
        whois
# note: whois package includes mkpasswd which is used in ldif templates

COPY files /

EXPOSE 389

CMD ["/docker-entrypoint.sh"]
