

docker run --rm -it \
  -v "$(pwd)":/envoy \
  -v /etc/passwd:/etc/passwd:ro \
  -w /envoy \
  registry.opensuse.org/opensuse/leap:15.6 /bin/bash

zypper install -y clang19 llvm19-libc++-devel gcc glibc-devel autoconf libtool python3-pip python3-virtualenv lld binutils git curl 

# Install bazel
curl --location --output /usr/local/bin/bazel https://github.com/bazelbuild/bazelisk/releases/latest/download/bazelisk-linux-$([ $(uname -m) = "aarch64" ] && echo "arm64" || echo "amd64")
chmod +x /usr/local/bin/bazel


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





cd ~/work/envoy

docker run --rm -it \
  -v "$(pwd)":/envoy \
  -v /etc/passwd:/etc/passwd:ro \
  -w /envoy \
  registry.opensuse.org/opensuse/leap:15.5 /bin/bash


# Run within opensuse container

# Install compiler toolchain and other tools
zypper install -y clang15 llvm15-devel gcc glibc-devel autoconf libtool python3-pip python3-virtualenv lld binutils git curl xz
zypper install -y clang17 llvm17-libc++-devel gcc glibc-devel autoconf libtool python3-pip python3-virtualenv lld binutils git curl xz

# Install bazel
curl --location --output /usr/local/bin/bazel https://github.com/bazelbuild/bazelisk/releases/latest/download/bazelisk-linux-$([ $(uname -m) = "aarch64" ] && echo "arm64" || echo "amd64")
chmod +x /usr/local/bin/bazel

# SUSE libc++ does not come with static libraries
# Copy the static version from pre-compiled clang+llvm distro

mkdir /tmp/clang-download
cd /tmp/clang-download
CLANG_VERSION=$(clang --version | awk '/^clang version/ {print $3}')
curl --location --output clang+llvm.tar.xz https://github.com/llvm/llvm-project/releases/download/llvmorg-${CLANG_VERSION}/clang+llvm-${CLANG_VERSION}-x86_64-linux-gnu-ubuntu-22.04.tar.xz
tar xvf clang+llvm.tar.xz clang+llvm-${CLANG_VERSION}-x86_64-linux-gnu-ubuntu-22.04/lib/x86_64-unknown-linux-gnu/
cp clang+llvm-${CLANG_VERSION}-x86_64-linux-gnu-ubuntu-22.04/lib/x86_64-unknown-linux-gnu/{libc++.a,libc++abi.a} /usr/lib64/



# Switch to non-root user to compile Envoy.
mkdir /home/tsaarni
chown tsaarni /home/tsaarni
su tsaarni

cd /envoy
echo "build --config=clang" > user.bazelrc
echo "--local_ram_resources=20000" >> user.bazelrc

bazel clean --expunge
bazel build //source/exe:envoy-static




# Install compiler toolchain and build dependencies not bundled with bazel build.
zypper install -y gcc13 gcc13-c++ make git binutils-gold python3

# Set gcc-13 as default version of gcc.
update-alternatives --install /usr/bin/gcc gcc /usr/bin/gcc-13 60
update-alternatives --install /usr/bin/g++ g++ /usr/bin/g++-13 60

# Install Bazel.
curl --location --output /usr/local/bin/bazel https://github.com/bazelbuild/bazelisk/releases/latest/download/bazelisk-linux-$([ $(uname -m) = "aarch64" ] && echo "arm64" || echo "amd64")
chmod +x /usr/local/bin/bazel

# Switch to non-root user to compile Envoy.
mkdir /home/tsaarni
chown tsaarni /home/tsaarni
su tsaarni

# Go to source code directory.
cd /envoy

# Configure the build.
cat > user.bazelrc <<EOF
build --config=gcc
--local_ram_resources=20000
EOF

# Compile Envoy.
bazel clean --expunge
bazel build //source/exe:envoy-static




