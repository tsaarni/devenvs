FROM sssd-sssd-build-image:latest

ENV DEBIAN_FRONTEND=noninteractive

ARG USERNAME=ubuntu

RUN echo $USERNAME ALL=\(root\) NOPASSWD:ALL > /etc/sudoers.d/$USERNAME && \
    chmod 0440 /etc/sudoers.d/$USERNAME

RUN curl -fsSL https://get.docker.com -o get-docker.sh && sh get-docker.sh && \
    usermod -aG docker $USERNAME

USER $USERNAME

CMD ["/usr/bin/dumb-init", "--", "sleep", "infinity"]
