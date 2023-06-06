
tools/gen_compilation_database.py


# compile inside devcontainer
bazel build -c dbg //source/exe:envoy-static && cp -af bazel-bin/source/exe/envoy-static .

# run outside devcontainer
./envoy-static --log-level debug -c ~/work/devenvs/envoy/configs/envoy-client-validation-crl-check.yaml


python3 -m http.server --bind 127.0.0.1 8081

bazel-bin/source/exe/envoy-static --log-level debug -c ~/work/devenvs/envoy/configs/envoy-client-validation-crl-check.yaml


# generate certs and two CRLS: one with client cert revoked, one with client cert not revoked
mkdir -p certs
rm certs/*
(cd app/generate-certs-and-crl; go run main.go)

# not revoked
cp certs/crl-client-not-revoked.pem certs/crl.pem
mv certs/crl.pem certs/client-ca-crl.pem
http --verify certs/server-ca.pem --cert certs/client.pem --cert-key certs/client-key.pem https://localhost:8443

# revoked
cp certs/crl-client-revoked.pem certs/crl.pem
mv certs/crl.pem certs/client-ca-crl.pem
http --verify certs/server-ca.pem --cert certs/client.pem --cert-key certs/client-key.pem https://localhost:8443




# run tests inside devcontainer
bazel test -c dbg //test/common/secret:sds_api_test --test_output=streamed
bazel test -c dbg //test/common/secret:sds_api_test --test_output=streamed --test_arg="--gtest_filter=CertificateValidationContextSdsRotationApiTestParams/CertificateValidationContextSdsRotationApiTest*"
