# Use the official envoy image as a base for out own image.
# It is built on top of Ubuntu 20.04.5 LTS.
FROM envoyproxy/envoy:v1.24.1

# Install some tools that are required by the install scripts that follow.
RUN apt-get update && apt-get install -y curl sudo gpg

# Install Fluent Bit and catatonit according to instructions at:
#   - https://docs.fluentbit.io/manual/installation/linux/ubuntu
#   - https://software.opensuse.org/download/package?package=catatonit&project=devel%3Akubic%3Alibcontainers%3Astable
RUN curl https://raw.githubusercontent.com/fluent/fluent-bit/master/install.sh | sh && \
    echo 'deb http://download.opensuse.org/repositories/devel:/kubic:/libcontainers:/stable/xUbuntu_20.04/ /' | sudo tee /etc/apt/sources.list.d/devel:kubic:libcontainers:stable.list && \
    curl -fsSL https://download.opensuse.org/repositories/devel:kubic:libcontainers:stable/xUbuntu_20.04/Release.key | gpg --dearmor | sudo tee /etc/apt/trusted.gpg.d/devel_kubic_libcontainers_stable.gpg > /dev/null && \
    sudo apt update && \
    sudo apt install catatonit

# Copy entrypoint.sh script to the container.
COPY files .

ENTRYPOINT ["/usr/libexec/catatonit/catatonit", "/docker-entrypoint-with-fluentbit.sh"]
