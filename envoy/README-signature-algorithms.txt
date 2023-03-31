
bazel build -c dbg //source/exe:envoy-static  # debug build
tools/vscode/refresh_compdb.sh                # generate compile_commands.json for vscode


# test by running envoy

mkdir -p certs
certyaml -d certs configs/certs.yaml

python3 -m http.server --bind 127.0.0.1 8081   # test server

rm /tmp/envoy-wireshark-keys.log
bazel-bin/source/exe/envoy-static -c ~/work/devenvs/envoy/configs/envoy-signature-algorithms.yaml --log-level debug

http --cert certs/client.pem --cert-key certs/client-key.pem --verify certs/server-ca.pem https://localhost:8443/

# check the signature algorithms in the decrypted server hello
wireshark -i lo -f "port 8443" -Y tls -o tls.keylog_file:/tmp/envoy-wireshark-keys.log -k


# openssl and sslyze does not show the offered signature algorithms
openssl s_client -cert certs/client.pem -key certs/client-key.pem -CAfile certs/server-ca.pem -connect localhost:8443
sslyze localhost:8443



# Changing API
#
# Check api/STYLE.md
# Edit the proto files and run

# fix code formatting before commit
./tools/code_format/check_format.py fix
git add api/ generated_api_shadow/



# Generate documentation
docs/build.sh
xdg-open generated/docs/index.html


# add changelog entry to: changelogs/current.yaml







#
# Testing
#
bazel test -c dbg //test/...
bazel test -c dbg //test/extensions/transport_sockets/tls/...
bazel test -c dbg //test/common/quic/... --test_output=streamed
bazel test -c dbg //test/common/quic:envoy_quic_proof_source_test --test_output=streamed

bazel test -c dbg //test/common/quic:envoy_quic_proof_source_test --test_output=streamed

bazel test -c dbg //test/extensions/transport_sockets/tls:ssl_socket_test --test_arg="--gtest_filter=IpVersions/SslSocketTest.SetSignatureAlgorithms*" --test_output=streamed --cache_test_results=no



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
