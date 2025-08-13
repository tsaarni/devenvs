

docker run --rm -it \
  -v "$(pwd)":/envoy \
  -v /etc/passwd:/etc/passwd:ro \
  -w /envoy \
  registry.opensuse.org/opensuse/leap:15.6 /bin/bash

zypper install -y \
  autoconf \
  curl \
  less \
  libtool \
  patch \
  python3-pip \
  unzip \
  python3-virtualenv \
  clang11 \
  lld \
  libstdc++6 \
  libstdc++-devel \
  gcc10-c++ \
  wget \
  git \
  vim

wget -O /usr/local/bin/bazel https://github.com/bazelbuild/bazelisk/releases/latest/download/bazelisk-linux-$([ $(uname -m) = "aarch64" ] && echo "arm64" || echo "amd64")
chmod +x /usr/local/bin/bazel

mkdir /home/tsaarni
chown tsaarni /home/tsaarni
su tsaarni

echo "build --config=clang" > user.bazelrc
echo "--local_ram_resources=20000" >> user.bazelrc

bazel clean
bazel build //source/exe:envoy-static
