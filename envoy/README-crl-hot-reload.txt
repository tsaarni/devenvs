
tools/gen_compilation_database.py


# compile inside devcontainer
bazel build -c dbg //source/exe:envoy-static && cp -af bazel-bin/source/exe/envoy-static .

# run outside devcontainer
./envoy-static --log-level debug -c ~/work/devenvs/envoy/configs/envoy-client-validation-crl-check.yaml


# run using released container image
docker run --volume $HOME/work/devenvs/envoy:/home/tsaarni/work/devenvs/envoy:ro --network host --user $(id -u):$(id -g) --rm envoyproxy/envoy:v1.27-latest --log-level debug -c ~/work/devenvs/envoy/configs/envoy-client-validation-crl-check.yaml


# run backend service
python3 -m http.server --bind 127.0.0.1 8081


bazel-bin/source/exe/envoy-static --log-level debug -c ~/work/devenvs/envoy/configs/envoy-client-validation-crl-check.yaml


# generate certs and two CRLS: one with client cert revoked, one with client cert not revoked
mkdir -p certs
rm certs/*
(cd app/generate-certs-and-crl; go run main.go)
cp certs/crl-client-not-revoked.pem certs/client-ca-crl.pem

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






###########################
#
# Test inside kubernetes
#


# start new cluster
kind delete cluster --name envoy
kind create cluster --config configs/kind-cluster-config.yaml --name envoy

# generate certs and two CRLS: one with client cert revoked, one with client cert not revoked
mkdir -p certs
rm certs/*
(cd app/generate-certs-and-crl; go run main.go)

# deploy envoy + backend service
kubectl apply -f manifests/envoy-with-crl.yaml

# create secrets
kubectl create secret generic envoy-certs --from-file=envoy.pem=certs/envoy.pem --from-file=envoy-key.pem=certs/envoy-key.pem --from-file=client-ca.pem=certs/client-ca.pem  --dry-run=client -o yaml | kubectl apply -f -
kubectl create secret generic envoy-crl --from-file=client-ca-crl.pem=certs/crl-client-not-revoked.pem --dry-run=client -o yaml | kubectl apply -f -

# make successful request with not revoked client cert
http --verify certs/server-ca.pem --cert certs/client.pem --cert-key certs/client-key.pem https://protected.127-0-0-135.nip.io

# update CRL
kubectl create secret generic envoy-crl --from-file=client-ca-crl.pem=certs/crl-client-revoked.pem --dry-run=client -o yaml | kubectl apply -f -

# make unsuccessful request with revoked client cert
http --verify certs/server-ca.pem --cert certs/client.pem --cert-key certs/client-key.pem https://protected.127-0-0-135.nip.io
