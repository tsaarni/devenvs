FROM ubuntu:18.04
ENV DEBIAN_FRONTEND noninteractive

RUN apt-get update && \
    apt-get install --yes --no-install-recommends \
        ca-certificates \
        curl \
        libnss-sss \
        libpam-sss \
        openssh-server \
        sssd \
        sssd-tools \
        cracklib-runtime

RUN  curl -fsSL https://raw.githubusercontent.com/gdraheim/docker-systemctl-replacement/v1.4/files/docker/systemctl3.py -o /bin/systemctl && \
     chmod +x /bin/systemctl

COPY files /

CMD ["/docker-entrypoint.sh"]
