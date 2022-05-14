


export WORKDIR=~/work/contour-devenv

# copy vscode settings
mkdir .vscode
cp $WORKDIR/configs/launch.json .vscode/


# Toolchain instructions from here
https://github.com/envoyproxy/envoy/blob/master/bazel/README.md

# compile with clang (not gcc)
echo "build --config=clang" >> user.bazelrc
echo "--local_ram_resources=10000" >> user.bazelrc

# install and setup clang
#   https://github.com/llvm/llvm-project/releases
wget https://github.com/llvm/llvm-project/releases/download/llvmorg-11.0.0/clang+llvm-11.0.0-x86_64-linux-gnu-ubuntu-20.04.tar.xz
sudo tar xf clang+llvm-11.0.0-x86_64-linux-gnu-ubuntu-20.04.tar.xz -C/usr/local/
bazel/setup_clang.sh /usr/local/clang+llvm-11.0.0-x86_64-linux-gnu-ubuntu-20.04/

# setup git hooks
./support/bootstrap

# generate compile_commands.json
#./tools/gen_compilation_database.py
tools/vscode/refresh_compdb.sh


# fix code formatting before commit
./tools/code_format/check_format.py fix


# start new cluster
kind delete cluster --name contour
kind create cluster --config configs/kind-cluster-config.yaml --name contour

# generate certificates
certyaml --destination certs configs/certs.yaml


#############################################################################
#
# BUILDING
#

# Local development builds (first see above for clang install)

bazel build -c fastbuild //source/exe:envoy-static
bazel build -c fastbuild -s //source/exe:envoy-static    # use -s to see compile commands
bazel build -c fastbuild //source/restarter:restarter

bazel build -c dbg //source/exe:envoy-static             # debug build

# run directly on command line
bazel-bin/source/exe/envoy-static -c bootstrap-config.yaml --service-node mynode --service-cluster -mycluster --log-level info

gdb --args bazel-bin/source/exe/envoy-static -c bootstrap-config.yaml --service-node mynode --service-cluster -mycluster --log-level info            # run under debugger
gdbserver localhost:9999 bazel-bin/source/envoy-static -c bootstrap-config.yaml --service-node mynode --service-cluster -mycluster --log-level info  # debug remotely

gdb
target remote localhost:9999

# package as docker image
cp -af bazel-bin/source/exe/envoy-static $WORKDIR/docker/envoy/envoy
###cp -af bazel-bin/source/restarter/restarter $WORKDIR/docker/envoy/restarter
docker build -f $WORKDIR/docker/envoy/Dockerfile $WORKDIR/docker/envoy -t envoy:latest
kind load docker-image envoy:latest --name contour  # upload image to kind cluster





# updating docs
export ENVOY_SRCDIR=$PWD
docs/build.sh
SPHINX_SKIP_CONFIG_VALIDATION=true docs/build.sh
google-chrome generated/docs/index.html




# Release builds

# start bazel cache in one terminal and then build with --remote_http_cache
mkdir -p ~/.cache/bazel-remote-cache
docker run -v $HOME/.cache/bazel-remote-cache:/data -u $(id -u):$(id -g) -p 28080:8080 buchgr/bazel-remote-cache

# build release
ci/run_envoy_docker.sh "BAZEL_BUILD_EXTRA_OPTIONS='--remote_http_cache=http://127.0.0.1:28080' ./ci/do_ci.sh bazel.release.server_only"  # without test
ci/run_envoy_docker.sh "BAZEL_BUILD_EXTRA_OPTIONS='--remote_http_cache=http://127.0.0.1:28080' ./ci/do_ci.sh bazel.release"              # with test

# build devel
./ci/run_envoy_docker.sh "BAZEL_BUILD_EXTRA_OPTIONS='--remote_http_cache=http://127.0.0.1:28080' ./ci/do_ci.sh bazel.dev"


docker build -f ci/Dockerfile-envoy -t envoy .   # only ubuntu image
ci/docker_build.sh                               # all images


# run tests without IPv6
./ci/run_envoy_docker.sh './ci/do_ci.sh bazel.dev //test/... --test_env=ENVOY_IP_TEST_VERSIONS=v4only --test_verbose_timeout_warnings '



#############################################################################
#
# Testing
#


# Build and load services for testing
docker pull tsaarni/httpbin:latest && kind load docker-image tsaarni/httpbin:latest --name contour

docker build -f docker/envoy-control-plane-stub/Dockerfile docker/envoy-control-plane-stub -t envoy-control-plane-stub:latest && kind load docker-image envoy-control-plane-stub --name contour






kubectl create configmap envoy-config --dry-run -o yaml --from-file=envoy.yaml=configs/envoy-xds-over-tls-path-source.yaml --from-file=configs/envoy-sds-auth-secret-tls-certicate.yaml --from-file=configs/envoy-sds-auth-secret-validation-context.yaml | kubectl apply -f -

kubectl create secret generic controlplane --dry-run -o yaml --from-file=certs/controlplane.pem --from-file=certs/controlplane-key.pem --from-file=certs/internal-root-ca.pem | kubectl apply -f -

kubectl create secret generic envoy --dry-run -o yaml --from-file=certs/envoy.pem --from-file=certs/envoy-key.pem --from-file=certs/internal-root-ca.pem | kubectl apply -f -

kubectl apply -f manifests/backend-httpbin-no-tls-no-ingress.yaml
kubectl apply -f manifests/deploy-control-plane-stub.yaml
kubectl apply -f manifests/deploy-envoy-xds-certificate-rotation.yaml




kubectl create configmap envoy-config --dry-run -o yaml --from-file=envoy.yaml=configs/envoy-xds-over-tls-path-source.yaml --from-file=configs/envoy-sds-auth-secret-tls-certicate.yaml --from-file=configs/envoy-sds-auth-secret-validation-context.yaml | kubectl apply -f -

http http://host1.127-0-0-101.nip.io/status/418
http --stream http://host1.127-0-0-101.nip.io/sse




#############################################################################
#
# Debug
#

docker build -f docker/envoy-debug/Dockerfile docker/envoy-debug/ -t envoy-debug:latest && kind load docker-image envoy-debug:latest --name contour
kubectl apply -f manifests/deploy-envoy-debug.yaml


bazel build -c dbg //source/exe:envoy-static

tar cf - bazel-bin/source/exe/envoy-static | kubectl exec -i $(kubectl get pod -l app=envoy -o jsonpath='{.items[0].metadata.name}') -- tar xvvf -
kubectl exec -it $(kubectl get pod -l app=envoy -o jsonpath='{.items[0].metadata.name}') -- bazel-bin/source/exe/envoy-static -c /etc/envoy/envoy.yaml --service-cluster mycluster --service-node envoy --log-level debug


kubectl port-forward $(kubectl get pod -l app=envoy -o jsonpath='{.items[0].metadata.name}') 9999
kubectl exec -it $(kubectl get pod -l app=envoy -o jsonpath='{.items[0].metadata.name}') -- gdbserver localhost:9999 bazel-bin/source/exe/envoy-static -c /etc/envoy/envoy.yaml --service-cluster mycluster --service-node envoy --log-level debug
gdb -iex "set sysroot ." -iex "target remote localhost:9999"


kubectl delete pod -lapp=control-plane --now && kubectl logs -f -lapp=control-plane

kubectl exec -it $(kubectl get pod -l app=envoy -o jsonpath='{.items[0].metadata.name}') bash



kubectl exec -it $(kubectl get pod -l app=envoy -o jsonpath='{.items[0].metadata.name}') -- lldb-server gdbserver "*:9999" -- /bazel-bin/source/exe/envoy-static -c /etc/envoy/envoy.yaml --service-cluster mycluster --service-node envoy --log-level debug

PYTHONPATH=/usr/lib/llvm-9/lib/python3.7/site-packages lldb bazel-bin/source/exe/envoy-static

sudo nsenter --target $(pidof envoy-static) --net wireshark -f "port 8080" -k



kubectl -n projectcontour port-forward envoy-fg4qz 9001
http http://localhost:9001/config_dump



#############################################################################
#
# Cert rotation
#

docker build -f docker/envoy-control-plane-stub/Dockerfile docker/envoy-control-plane-stub -t envoy-control-plane-stub:latest && kind load docker-image envoy-control-plane-stub --name contour
docker pull tsaarni/httpbin:latest && kind load docker-image tsaarni/httpbin:latest --name contour
docker build -f docker/envoy-debug/Dockerfile docker/envoy-debug/ -t envoy-debug:latest && kind load docker-image envoy-debug:latest --name contour

kubectl create configmap envoy-config --dry-run -o yaml --from-file=envoy.yaml=configs/envoy-xds-over-tls-path-source.yaml --from-file=configs/envoy-sds-auth-secret-tls-certicate.yaml --from-file=configs/envoy-sds-auth-secret-validation-context.yaml | kubectl apply -f -
kubectl create secret generic controlplane --dry-run -o yaml --from-file=certs/controlplane.pem --from-file=certs/controlplane-key.pem --from-file=certs/internal-root-ca.pem | kubectl apply -f -
kubectl create secret generic envoy --dry-run -o yaml --from-file=certs/envoy.pem --from-file=certs/envoy-key.pem --from-file=certs/internal-root-ca.pem | kubectl apply -f -

kubectl apply -k manifests/envoy-certificate-rotation/

while inotifywait -qre close_write bazel-bin/source/exe; do
  tar cf - bazel-bin/source/exe/envoy-static | kubectl exec -i $(kubectl get pod -l app=envoy -o jsonpath='{.items[0].metadata.name}') -- tar xvvf -
done

bazel build -c dbg //source/exe:envoy-static
bazel build -c fastbuild //source/exe:envoy-static


kubectl exec -it $(kubectl get pod -l app=envoy -o jsonpath='{.items[0].metadata.name}') -- bazel-bin/source/exe/envoy-static -c /etc/envoy/envoy.yaml --service-cluster mycluster --service-node envoy --log-level debug


http http://host1.127-0-0-101.nip.io/status/418
http --stream http://host1.127-0-0-101.nip.io/sse


# re-generate envoy cert
rm certs/envoy*pem
certyaml --destination certs configs/certs.yaml
kubectl create secret generic envoy --dry-run -o yaml --from-file=certs/envoy.pem --from-file=certs/envoy-key.pem --from-file=certs/internal-root-ca.pem | kubectl apply -f -


sudo nsenter --target $(pidof envoy-static) --net wireshark -f "port 8080" -k

# trigger new TLS connection
kubectl delete pod -lapp=control-plane


# unittest with bazel
bazel test -c dbg //test/extensions/formatter/req_without_query:req_without_query_test --test_output=streamed --test_arg="-l trace"
bazel test -c dbg //test/extensions/formatter/regex_substitute:regex_substitute_test --test_output=streamed --test_arg="-l trace"
bazel test -c dbg //test/extensions/formatter/regex_substitute:regex_substitute_test --test_output=streamed --test_arg="--gtest_filter=RegexSubstituteTest.TestStripQueryString"
bazel test -c dbg //test/common/filesystem:watcher_impl_test --test_output=streamed --test_arg="-l trace"
bazel test -c dbg //test/... --test_output=streamed  # run all unittests
bazel test -c dbg --config=clang-asan //test/...     # address sanitizer
bazel test -c dbg --config=clang-tsan //test/...     # thread sanitizer






docker run --rm --publish 8081:3000 gcr.io/k8s-staging-ingressconformance/echoserver:v20201006-42d00bd
bazel-bin/source/exe/envoy-static -c $WORKDIR/configs/envoy-static-virtualhost.yaml
http "http://localhost:8080/foo?supersecret=password"




###########################################################################
#
# Changing API
#


Check api/STYLE.md


Edit the proto files and run

./tools/proto_format/proto_format.sh fix
git add api/ generated_api_shadow/



############################################################################
#
# Changing docs
#

docs/build.sh
xdg-open generated/docs/index.html






