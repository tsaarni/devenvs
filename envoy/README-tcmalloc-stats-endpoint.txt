

# admin: add tcmalloc stats endpoint
https://github.com/envoyproxy/envoy/pull/41376




bazel build -c fastbuild //source/exe:envoy-static
bazel build -c fastbuild //source/exe:envoy-static --define tcmalloc=gperftools

bazel test -c fastbuild //test/server/admin:server_info_handler_test
bazel test -c fastbuild //test/server/admin:server_info_handler_test --define=tcmalloc=gperftools
bazel test -c fastbuild //test/server/admin:server_info_handler_test --define=tcmalloc=disabled
