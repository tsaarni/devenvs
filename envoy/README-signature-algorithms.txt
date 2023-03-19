



bazel build -c dbg //source/exe:envoy-static             # debug build

tools/vscode/refresh_compdb.sh


# Changing API
#
# Check api/STYLE.md
# Edit the proto files and run

# fix code formatting before commit
./tools/code_format/check_format.py fix
git add api/ generated_api_shadow/










#
# Testing
#
bazel test -c dbg //test/...
bazel test -c dbg //test/common/quic/... --test_output=streamed
bazel test -c dbg //test/common/quic:envoy_quic_proof_source_test --test_output=streamed

bazel test -c dbg //test/common/quic:envoy_quic_proof_source_test --test_output=streamed


# To run under debugger, create vscode launch.json entry for running
# gdb bazel-bin/test/common/quic/envoy_quic_proof_source_test


    {
      "name": "Run test under gdb",
      "request": "launch",
      "type": "gdb",
      "target": "bazel-bin/test/extensions/transport_sockets/tls/ssl_socket_test",
      "arguments": "--gtest_filter=IpVersions/SslSocketTest.SignatureAlgorithms/IPv4_with_sync_cert_validation",
      "cwd": "${workspaceFolder}",
      "valuesFormatting": "disabled"
    }

bazel test -c dbg //test/extensions/transport_sockets/tls:ssl_socket_test --test_output=streamed --test_arg="--gtest_filter=IpVersions/SslSocketTest.SignatureAlgorithms/IPv4_with_sync_cert_validation"


bazel test -c dbg //test/extensions/access_loggers/grpc:tcp_grpc_access_log_integration_test --test_arg="--gtest_filter=IpVersionsCientType/TcpGrpcAccessLogIntegrationTest.SslNotTerminated/IPv4_EnvoyGrpc"



# signature algorithms in boringssl
https://github.com/google/boringssl/blob/8c7aa6bfcd7573d7b904fde6acb4f3652a3ebecc/ssl/extensions.cc#L457-L478


bazel test -c dbg //test/extensions/transport_sockets/tls:ssl_socket_test --test_arg="--gtest_filter=IpVersions/SslSocketTest.SignatureAlgorithms*" --test_output=streamed --cache_test_results=no


bazel test -c dbg //test/extensions/transport_sockets/tls:ssl_socket_test -gtest_filter=SignatureAlgorithms
IpVersions/SslSocketTest.SignatureAlgorithms/IPv6_with_custom_cert_validation



test/extensions/transport_sockets/tls/ssl_socket_test.cc:842: Failure
Expected equality of these values:
  1UL
    Which is: 1
  server_stats_store.counter(options.expectedServerStats()).value()
    Which is: 0
ssl.sigalgs.rsa_pss_rsae_sha256
Stack trace:
  0x215b29f: Envoy::Extensions::TransportSockets::Tls::(anonymous namespace)::testUtilV2()
  0x2181769: Envoy::Extensions::TransportSockets::Tls::SslSocketTest_SignatureAlgorithms_Test::TestBody()
  0x56bf38b: testing::internal::HandleSehExceptionsInMethodIfSupported<>()
  0x56af5fa: testing::internal::HandleExceptionsInMethodIfSupported<>()
  0x569b2b3: testing::Test::Run()
  0x569bc54: testing::TestInfo::Run()
... Google Test internal frames ...

test/extensions/transport_sockets/tls/ssl_socket_test.cc:847: Failure
Expected equality of these values:
  1UL
    Which is: 1
  client_stats_store.counter(options.expectedClientStats()).value()
    Which is: 0
Stack trace:
  0x215b48b: Envoy::Extensions::TransportSockets::Tls::(anonymous namespace)::testUtilV2()
  0x2181769: Envoy::Extensions::TransportSockets::Tls::SslSocketTest_SignatureAlgorithms_Test::TestBody()
  0x56bf38b: testing::internal::HandleSehExceptionsInMethodIfSupported<>()
  0x56af5fa: testing::internal::HandleExceptionsInMethodIfSupported<>()
  0x569b2b3: testing::Test::Run()
  0x569bc54: testing::TestInfo::Run()
... Google Test internal frames ...

test/extensions/transport_sockets/tls/ssl_socket_test.cc:842: Failure
Expected equality of these values:
  1UL
    Which is: 1
  server_stats_store.counter(options.expectedServerStats()).value()
    Which is: 0
ssl.sigalgs.rsa_pss_rsae_sha256
Stack trace:
  0x215b29f: Envoy::Extensions::TransportSockets::Tls::(anonymous namespace)::testUtilV2()
  0x2181790: Envoy::Extensions::TransportSockets::Tls::SslSocketTest_SignatureAlgorithms_Test::TestBody()
  0x56bf38b: testing::internal::HandleSehExceptionsInMethodIfSupported<>()
  0x56af5fa: testing::internal::HandleExceptionsInMethodIfSupported<>()
  0x569b2b3: testing::Test::Run()
  0x569bc54: testing::TestInfo::Run()
... Google Test internal frames ...

test/extensions/transport_sockets/tls/ssl_socket_test.cc:847: Failure
Expected equality of these values:
  1UL
    Which is: 1
  client_stats_store.counter(options.expectedClientStats()).value()
    Which is: 0
Stack trace:
  0x215b48b: Envoy::Extensions::TransportSockets::Tls::(anonymous namespace)::testUtilV2()
  0x2181790: Envoy::Extensions::TransportSockets::Tls::SslSocketTest_SignatureAlgorithms_Test::TestBody()
  0x56bf38b: testing::internal::HandleSehExceptionsInMethodIfSupported<>()
  0x56af5fa: testing::internal::HandleExceptionsInMethodIfSupported<>()
  0x569b2b3: testing::Test::Run()
  0x569bc54: testing::TestInfo::Run()
... Google Test internal frames ...

[external/com_google_absl/absl/flags/internal/flag.cc : 115] RAW: Restore saved value of envoy_reloadable_features_tls_async_cert_validation to: true
[  FAILED  ] IpVersions/SslSocketTest.SignatureAlgorithms/IPv6_with_custom_cert_validation, where GetParam() = (4-byte object <01-00 00-00>, true) (476 ms)
