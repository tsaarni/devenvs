
#############################################################################
#
# Random notes (work-in-progress)
#

https://stackoverflow.com/questions/33117068/use-of-supervisor-in-docker/33119321#33119321
handle SIGTERM from container runtime and distribute it to childs

https://medium.com/@gchudnov/trapping-signals-in-docker-containers-7a57fdda7d86
https://blog.phusion.nl/2015/01/20/docker-and-the-pid-1-zombie-reaping-problem/






#############################################################################
#
# Start kind cluster
#
kind create cluster --config configs/kind-cluster-config.yaml --name envoy


#
# "ad-hoc" build of envoy container
#
cd ~/work/envoy
bazel build -c fastbuild //source/exe:envoy-static
bazel build -c fastbuild //source/restarter:restarter
cp -af bazel-bin/source/exe/envoy-static ~/work/envoy-devenv/docker/envoy/envoy
cp -af bazel-bin/source/restarter/restarter ~/work/envoy-devenv/docker/envoy/restarter

cd ~/work/envoy-devenv
docker build -f docker/envoy/Dockerfile docker/envoy -t envoy:latest

kind load docker-image envoy:latest --name envoy  # upload image to kind cluster


#
# build envoy-control-plane-stub
#
cd ~/work/envoy-devenv
docker build -f docker/envoy-control-plane-stub/Dockerfile docker/envoy-control-plane-stub -t envoy-control-plane-stub:latest
kind load docker-image envoy-control-plane-stub --name envoy


#
# pull backend service image for tests and upload that to cluster
#
docker pull tsaarni/httpbin:latest
kind load docker-image tsaarni/httpbin:latest --name envoy


#
# Deploy
#

# configmap for envoy config
kubectl create configmap envoy-config --from-file=envoy.yaml=envoy-xds-over-tls.yaml

# generate certificates
mkdir -p certs
cfssl genkey -initca configs/cfssl-csr-root-ca-server.json | cfssljson -bare certs/server-root
cfssl genkey -initca configs/cfssl-csr-root-ca-control-plane.json | cfssljson -bare certs/control-plane-root
cfssl gencert -ca certs/server-root.pem -ca-key certs/server-root-key.pem configs/cfssl-csr-endentity-httpbin.json | cfssljson -bare certs/httpbin
cfssl gencert -ca certs/control-plane-root.pem -ca-key certs/control-plane-root-key.pem configs/cfssl-csr-endentity-envoy.json | cfssljson -bare certs/envoy
cfssl gencert -ca certs/control-plane-root.pem -ca-key certs/control-plane-root-key.pem configs/cfssl-csr-endentity-controlplane.json | cfssljson -bare certs/controlplane

# upload secrets
kubectl create secret generic httpbin-certs --from-file=certs/httpbin.pem --from-file=certs/httpbin-key.pem
kubectl create secret generic envoy-certs --from-file=certs/control-plane-root.pem --from-file=certs/envoy.pem --from-file=certs/envoy-key.pem
kubectl create secret generic controlplane-certs --from-file=certs/controlplane.pem --from-file=certs/controlplane-key.pem --from-file=certs/control-plane-root.pem

# deploy
kubectl apply -f deploy-backend-no-tls.yaml
kubectl apply -f deploy-envoy-control-plane-stub.yaml
kubectl apply -f deploy-envoy-with-restarter.yaml



#############################################################################
#
# DEMO
#

# show what is running
#   - httpbin as backend
#   - envoy-control-plane which dynamically configures envoy
#   - httpbin as backend service
kubectl get pod


# show what is running
kubectl exec -it envoy-66b6fc46db-6kbdh bash

# run inside pod: show restarter
ps uaxww

# inside pod: show envoy config and the certificate
cat /etc/envoy/envoy.yaml
openssl x509 -in /run/secrets/certs/envoy.pem -text -noout


# show manifest, secret mount
cat deploy-envoy-with-restarter.yaml


# run traffic
http http://localhost/status/418
http --stream http://localhost/sse

# show envoy logs to observe the inotify watch being triggered
kubectl logs -f -lapp=envoy

# trigger update of envoy cert
cfssl gencert -ca certs/control-plane-root.pem -ca-key certs/control-plane-root-key.pem configs/cfssl-csr-endentity-envoy.json | cfssljson -bare certs/envoy

kubectl create secret generic envoy-certs --from-file=certs/control-plane-root.pem --from-file=certs/envoy.pem --from-file=certs/envoy-key.pem -o yaml --dry-run | kubectl apply -f -
