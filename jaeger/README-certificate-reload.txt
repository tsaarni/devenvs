

https://github.com/jaegertracing/jaeger/issues/3968
https://github.com/jaegertracing/jaeger/discussions/4073

https://github.com/jaegertracing/jaeger/pull/2389

# {"level":"warn","ts":1673525507.906441,"caller":"tlscfg/cert_watcher.go:92","msg":"Certificate has been removed, using the last known version","certificate":"/certs/server-ca.pem"}



#
# Create Kubernetes cluster.
#
kind delete cluster --name jaeger   # delete old cluster if any
kind create cluster --name jaeger --config configs/kind-cluster.yaml


# install contour as ingress controller
kubectl apply -f https://projectcontour.io/quickstart/contour.yaml


# Create test certificates and upload them as Secrets.
rm -rf certs
mkdir -p certs
certyaml -d certs configs/certs.yaml

kubectl create secret generic opensearch-certs --dry-run=client -o yaml --from-file=certs/client-ca.pem --from-file=certs/opensearch-server.pem --from-file=certs/opensearch-server-key.pem --from-file=certs/server-ca.pem | kubectl apply -f -
kubectl create secret generic jaeger-certs --dry-run=client -o yaml --from-file=certs/server-ca.pem --from-file=certs/opensearch-admin.pem --from-file=certs/opensearch-admin-key.pem | kubectl apply -f -

# Deploy opensearch.
kubectl apply -f manifests/opensearch.yaml


# Opensearch troubleshooting:
#
#   kubectl logs deployment/opensearch -f
#
# Test that anonymous and basic auth access is disabled and that certificate-based auth is enabled
#   http --verify=certs/server-ca.pem https://opensearch.127-0-0-163.nip.io:9200/
#   http --verify=certs/server-ca.pem --auth admin:admin https://opensearch.127-0-0-163.nip.io:9200/
#   http --verify=certs/server-ca.pem --cert=certs/opensearch-admin.pem --cert-key=certs/opensearch-admin-key.pem https://opensearch.127-0-0-163.nip.io:9200/


# deploy jaeger
kubectl apply -f manifests/jaeger-with-es-storage.yaml

# Jaeger troubleshooting
#
#   kubectl logs deployment/jaeger -f
#
# http://jaeger.127-0-0-163.nip.io/


# deploy hotrod sample app
kubectl apply -f manifests/hotrod.yaml

# generate traces with hotrod
#   http://hotrod.127-0-0-163.nip.io/



########################
#
# Working with Jaeger
#
# https://github.com/jaegertracing/jaeger/blob/main/CONTRIBUTING.md
#

# Install required version of node and yarn
#   cat jaeger-ui/.nvmrc
volta install node@16 yarn@v1.22.19


rm -rf jaeger-ui idl   # remove submodules if any
git submodule update --init --recursive
make install-tools


make test


make build-all-in-one
ls -l ./cmd/all-in-one/all-in-one-linux-amd64
SPAN_STORAGE_TYPE=elasticsearch ./cmd/all-in-one/all-in-one-linux-amd64 --help   # show all options


# run own build on kubernetes
kubectl apply -f manifests/jaeger-dev.yaml

make build-all-in-one
kubectl cp ./cmd/all-in-one/all-in-one-linux-amd64 jaeger:/
kubectl exec -it jaeger -- /all-in-one-linux-amd64 --es.server-urls=https://opensearch:9200 --es.tls.enabled=true --es.tls.ca=/certs/server-ca.pem --es.tls.cert=/certs/opensearch-admin.pem --es.tls.key=/certs/opensearch-admin-key.pem --es.tls.server-name=opensearch --es.tls.skip-host-verify=false

# watch events
kubectl exec -it jaeger -- inotifywait -m /certs/

# test with valid keys
rm certs/opensearch-admin*.pem
certyaml -d certs configs/certs.yaml
kubectl create secret generic jaeger-certs --dry-run=client -o yaml --from-file=certs/server-ca.pem --from-file=certs/opensearch-admin.pem --from-file=certs/opensearch-admin-key.pem | kubectl apply -f -

# test with expired keys
kubectl create secret generic jaeger-certs --dry-run=client -o yaml --from-file=certs/server-ca.pem --from-file=opensearch-admin.pem=certs/expired-opensearch-admin.pem --from-file=opensearch-admin-key.pem=certs/expired-opensearch-admin-key.pem | kubectl apply -f -
