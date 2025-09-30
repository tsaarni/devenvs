




bazel/setup_clang.sh /opt/llvm/


# Configure the build.
cat > user.bazelrc <<EOF
build --config=clang
build --disk_cache=~/.cache/bazel-build-cache
--local_ram_resources=20000
EOF

bazel build -c fastbuild //source/exe:envoy-static



# Generate launch.json for vscode
##  from envoy/tools/vscode/README.md
tools/vscode/generate_debug_config.py //source/exe:envoy-static --args "-c envoy.yaml"




# copy vscode settings
mkdir -p .vscode
cp ~/work/devenvs/envoy/configs/launch.json .vscode/

# setup git hooks
./support/bootstrap

# setup clang
bazel/setup_clang.sh /usr/local/clang+llvm-14.0.0-x86_64-linux-gnu-ubuntu-18.04

# compile with clang (not gcc) and libc++ (not libstdc++)
echo "build --config=libc++" > user.bazelrc
echo "--local_ram_resources=20000" >> user.bazelrc

# generate compile_commands.json
tools/vscode/refresh_compdb.sh




### Install toolchain


# Toolchain instructions are here bazel/README.md
# https://github.com/envoyproxy/envoy/blob/master/bazel/README.md


# get pre-compiled clang
wget https://github.com/llvm/llvm-project/releases/download/llvmorg-14.0.0/clang+llvm-14.0.0-x86_64-linux-gnu-ubuntu-18.04.tar.xz
sudo tar xf clang+llvm-14.0.0-x86_64-linux-gnu-ubuntu-18.04.tar.xz -C/usr/local/


# check list from bazel/README.md
sudo apt-get install \
   autoconf \
   automake \
   cmake \
   curl \
   libtool \
   make \
   ninja-build \
   patch \
   python3-pip \
   unzip \
   virtualenv


# error:
# ./bootstrap0: error while loading shared libraries: libc++.so.1: cannot open shared object file: No such file or directory
# then install also libc++1-14




# fix code formatting before commit
./tools/code_format/check_format.py fix





### Building and debugging

# Build Envoy
bazel build -c fastbuild //source/exe:envoy-static

# Add -s to bazel command line to see the actual compilation commands
bazel build -c fastbuild -s //source/exe:envoy-static


# Build with debug symbols
bazel build -c dbg //source/exe:envoy-static
gdb --args bazel-bin/source/envoy-static ...

# The binary will be stored in bazel-bin/source/exe/envoy-static


## Build and run Envoy test suite within container image, utilizing [Bazel remote cache](https://github.com/buchgr/bazel-remote)

# start bazel cache in one terminal and then build with --remote_http_cache
rm -rf ~/.cache/bazel-remote-cache/ && mkdir -p ~/.cache/bazel-remote-cache/
docker run -v $HOME/.cache/bazel-remote-cache:/data -u $(id -u):$(id -g) -p 28080:8080 quay.io/bazel-remote/bazel-remote --max_size 5

# build release and test
ci/run_envoy_docker.sh "BAZEL_BUILD_EXTRA_OPTIONS='--remote_http_cache=http://127.0.0.1:28080' ./ci/do_ci.sh bazel.release"
ci/run_envoy_docker.sh "BAZEL_BUILD_EXTRA_OPTIONS='--remote_http_cache=http://127.0.0.1:28080' ./ci/do_ci.sh bazel.debug.server_only"

# create container
ci/docker_build.sh
./ci/run_envoy_docker.sh './ci/do_ci.sh release'


Debugging
bazel build -c dbg //source/restarter:restarter
gdb --args bazel-bin/source/restarter/restarter --command bazel-bin/source/exe/envoy-static --watch certs/ -- -c examples/front-proxy/front-envoy.yaml


## Tips and tricks

Using vscode:

# copy vscode settings to a new workspace
cp c_cpp_properties.json $ENVOY_SRC/.vscode/

# generate compile_commands.json
./tools/gen_compilation_database.py


### Running unit tests

Construct the name according to example: /[test directory]:[test target]

For example, following definition in test/common/secret/BUILD

      envoy_cc_test(
         name = "sds_api_test",


Would mean that the test can be run with following command:

bazel test -c dbg //test/common/secret:sds_api_test --test_output=streamed


Use --test_arg to pass arguments to the test binary.
For example, to run only one test case:

bazel test -c dbg //test/common/secret:sds_api_test --test_output=streamed --test_arg="--gtest_filter=CertificateValidationContextSdsRotationApiTestParams/CertificateValidationContextSdsRotationApiTest*"


To avoid caching test results, use --cache_test_results=no

--cache_test_results=no



### Envoy release builds


docker run --volume $HOME/work/devenvs/envoy:/home/tsaarni/work/devenvs/envoy:ro --network host --user $(id -u):$(id -g) --rm envoyproxy/envoy:v1.24-latest --log-level debug -c ~/work/devenvs/envoy/configs/envoy-client-validation-crl-check.yaml

docker run --volume $HOME/work/devenvs/envoy:/home/tsaarni/work/devenvs/envoy:ro --network host --user $(id -u):$(id -g) --rm envoyproxy/envoy:v1.25-latest --log-level debug -c ~/work/devenvs/envoy/configs/envoy-client-validation-crl-check.yaml

docker run --volume $HOME/work/devenvs/envoy:/home/tsaarni/work/devenvs/envoy:ro --network host --user $(id -u):$(id -g) --rm envoyproxy/envoy:v1.26-latest --log-level debug -c ~/work/devenvs/envoy/configs/envoy-client-validation-crl-check.yaml


#### Dev container

# compile
bazel build -c dbg //source/exe:envoy-static && cp -a bazel-bin/source/exe/envoy-static .

# compilation database
tools/gen_compilation_database.py    # reload vscode window afterwards

# fix format
./tools/code_format/check_format.py fix


### Document changes

ci/do_ci.sh docs

xdg-open ./generated/docs/index.html
