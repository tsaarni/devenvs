




cd ~/work/envoy
mkdir -p .vscode
cp ~/work/devenvs/envoy/configs/{settings.json,launch.json} .vscode/



diff --git a/bazel/envoy_internal.bzl b/bazel/envoy_internal.bzl
index 8282fb3e4f..91301e3106 100644
--- a/bazel/envoy_internal.bzl
+++ b/bazel/envoy_internal.bzl
@@ -16,6 +16,7 @@ def envoy_copts(repository, test = False):
         "-Wvla",
         "-Wno-deprecated-declarations",
         "-Wreturn-type",
+        "-Wno-nullability-completeness",
     ]
     _repo = repo_label(repository)


https://github.com/bazel-contrib/toolchains_llvm/blob/master/toolchain/distributions/github.jsonc

diff --git a/bazel/toolchains.bzl b/bazel/toolchains.bzl
index aa296862d7..61630327c2 100644
--- a/bazel/toolchains.bzl
+++ b/bazel/toolchains.bzl
@@ -13,7 +13,7 @@ def envoy_toolchains():
     )
     llvm_toolchain(
         name = "llvm_toolchain",
-        llvm_version = "18.1.8",
+        llvm_version = "21.1.8",
         cxx_standard = {"": "c++20"},
         sysroot = {
             "linux-x86_64": "@sysroot_linux_amd64//:sysroot",



# Configure the build.
cat > user.bazelrc <<EOF
build --config=clang
build --disk_cache=~/.cache/envoy-bazel
build --local_resources=cpu=HOST_CPUS*0.5
build --local_resources=memory=HOST_RAM*0.3
EOF

# llvm 21.1.0
cat > user.bazelrc <<EOF
build --config=clang
build --disk_cache=~/.cache/envoy-bazel
build --local_resources=cpu=HOST_CPUS*0.3
build --local_resources=memory=HOST_RAM*0.3
build --copt=-Wno-nullability-completeness
EOF


cat > user.bazelrc <<EOF
build --config=clang
build --disk_cache=~/.cache/envoy-bazel
build --local_resources=cpu=HOST_CPUS*0.5
build --local_resources=memory=HOST_RAM*0.3
build --override_repository=envoy_build_config=%workspace%/envoy_build_config
EOF


### The binary will be stored in bazel-bin/source/exe/envoy-static

bazel build -c fastbuild //source/exe:envoy-static     # Build binary
bazel build -c fastbuild -s //source/exe:envoy-static  # See the actual compilation commands

bazel build -c dbg //source/exe:envoy-static           # Include debug symbols
gdb --args bazel-bin/source/envoy-static ...           # Run with gdb debugger


nice -n 10 bazel build -c dbg //source/exe:envoy-static --//bazel:http3=False
## workaround for compilation errors
#  --//bazel:http3=False


# build with gperftools tcmalloc
bazel build -c fastbuild //source/exe:envoy-static --define tcmalloc=gperftools





# Generate compilation database
cp .vscode/tasks_example.json .vscode/tasks.json
# run: "Tasks: Run Task" and select "Refresh Compilation Database Exclude Contrib"

or

ENVOY_GEN_COMPDB_OPTIONS="--vscode --include_headers --exclude_contrib" ./ci/do_ci.sh refresh_compdb



# Generate launch.json for vscode
##  from envoy/tools/vscode/README.md
./tools/vscode/generate_debug_config.py //source/exe:envoy-static --args "-c envoy.yaml"



# setup git hooks
./support/bootstrap



# fix code formatting before commit
bazel run //tools/code_format:check_format -- fix
bazel run //tools/code_format:check_format -- fix source/common/stats/thread_local_store.cc

# Uncommitted changes (default)
./tools/local_fix_format.sh

# Changes since main branch
./tools/local_fix_format.sh -main

# Specific files
./tools/local_fix_format.sh source/common/stats/thread_local_store.cc







## Build and run Envoy test suite within container image, utilizing [Bazel remote cache](https://github.com/buchgr/bazel-remote)

# start bazel cache in one terminal and then build with --remote_http_cache
rm -rf ~/.cache/bazel-remote-cache/ && mkdir -p ~/.cache/bazel-remote-cache/
docker run -v $HOME/.cache/bazel-remote-cache:/data -u $(id -u):$(id -g) -p 28080:8080 quay.io/bazel-remote/bazel-remote --max_size 5

# build release and test
ci/run_envoy_docker.sh "BAZEL_BUILD_EXTRA_OPTIONS='--remote_http_cache=http://127.0.0.1:28080' ./ci/do_ci.sh bazel.release"
ci/run_envoy_docker.sh "BAZEL_BUILD_EXTRA_OPTIONS='--remote_http_cache=http://127.0.0.1:28080' ./ci/do_ci.sh bazel.debug.server_only"

# create container 

## 1. compile (/tmp/envoy-docker-build/  ~/.cache/envoy-bazel)
./ci/run_envoy_docker.sh './ci/do_ci.sh release'
./ci/run_envoy_docker.sh './ci/do_ci.sh release.server_only'

## 2. build image  (THIS DOES NOT WORK 2026-05-30)
sudo rm -rf /tmp/bazel-shared
ENVOY_DOCKER_IN_DOCKER=1 ENVOY_DOCKER_BUILD_DIR=/tmp/envoy-docker-build ./ci/run_envoy_docker.sh './ci/do_ci.sh docker'

## 2. build image  (hack)
mkdir -p linux/amd64
cp /tmp/envoy-docker-build/envoy/x64/bin/release.tar.zst linux/amd64/
CI_SHA1=$(git rev-parse HEAD~1) DOCKER_PLATFORM=linux/amd64 DOCKER_LOAD_IMAGES=1 ENVOY_OCI_DIR=/tmp/envoy-oci ./distribution/docker/build.sh
rm -rf linux/


or

diff --git a/ci/docker-entrypoint-extra.sh b/ci/docker-entrypoint-extra.sh
index 891f79b402..e360542ceb 100755
--- a/ci/docker-entrypoint-extra.sh
+++ b/ci/docker-entrypoint-extra.sh
@@ -7,7 +7,7 @@ fi
 SOCKET_GID=$(stat -c '%g' /var/run/docker.sock)
 if ! getent group "$SOCKET_GID" > /dev/null; then
     # Create a group with the same GID as the socket
-    groupmod -g "$SOCKET_GID" docker \
-        || (delgroup docker && addgroup -g "$SOCKET_GID" docker)
-    gpasswd -a envoybuild docker
+    groupadd -g "$SOCKET_GID" docker
 fi
+DOCKER_GROUP=$(getent group "$SOCKET_GID" | cut -d: -f1)
+usermod -aG "$DOCKER_GROUP" envoybuild







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


### Document changes

cd docs
bazel run --@envoy//tools/tarball:target=//:html @envoy//tools/tarball:unpack "$PWD"/generated/docs/
cd generated/docs
python3 -m http.server 8000
xdg-open http://localhost:8000

