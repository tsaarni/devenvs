
# Environment to experiment with Envoy Proxy


## Building and debugging

Build Envoy

bazel build -c fastbuild //source/exe:envoy-static


Build with debug symbols

bazel build -c dbg //source/exe:envoy-static
gdb --args bazel-bin/source/envoy-static ...

The binary will be stored in bazel-bin/source/exe/envoy-static


Build and run Envoy test suite within container image, utilizing [Bazel remote cache](https://github.com/buchgr/bazel-remote)

# start bazel cache in one terminal and then build with --remote_http_cache
docker run -v $HOME/.cache/bazel-remote-cache:/data -p 28080:8080 buchgr/bazel-remote-cache

# build release and test
ci/run_envoy_docker.sh "BAZEL_BUILD_EXTRA_OPTIONS='--remote_http_cache=http://127.0.0.1:28080' ./ci/do_ci.sh bazel.release"

# create container
ci/docker_build.sh


Debugging
bazel build -c dbg //source/restarter:restarter
gdb --args bazel-bin/source/restarter/restarter --command bazel-bin/source/exe/envoy-static --watch certs/ -- -c examples/front-proxy/front-envoy.yaml


## Tips and tricks

Using vscode:

# copy vscode settings to a new workspace
cp c_cpp_properties.json $ENVOY_SRC/.vscode/

# generate compile_commands.json
./tools/gen_compilation_database.py


Add -s to bazel command line to see the actual compilation commands

bazel build -c fastbuild -s //source/exe:envoy-static
