
# Create cluster
kind delete cluster --name permissions
kind create cluster --name permissions

kind create cluster --name permissions --config configs/kind-cluster-with-version.conf

# build test app
docker build --tag dirwatcher:latest docker/dirwatcher/
kind load docker-image dirwatcher:latest --name permissions

# deploy test app
kubectl apply -f manifests/secret.yaml
kubectl delete -f manifests/dirwatcher-nonroot.yaml --force   # delete previous pod
kubectl apply -f manifests/dirwatcher-nonroot.yaml

# update secret
while true; do kubectl create secret generic mysecret --from-file=password=/proc/sys/kernel/random/uuid --dry-run=client -o yaml | kubectl apply -f -; sleep 30; done

# read testapp logs
kubectl logs -f dirwatcher-nonroot

kubectl exec dirwatcher-nonroot ls -laR /secret /var/run/secrets/kubernetes.io/serviceaccount


###########
#
# tips when working with kubelet code in kind
#

# run unit tests
make check test
make test WHAT=./pkg/volume GOFLAGS=-v

# build kubelet
make WHAT=cmd/kubelet
ls -l ./_output/local/bin/linux/amd64/kubelet

# replace kubelet version with own build
docker exec permissions-control-plane systemctl stop kubelet
docker cp ./_output/local/bin/linux/amd64/kubelet permissions-control-plane:/usr/bin/kubelet
docker exec permissions-control-plane systemctl start kubelet

# trigger service account token update manually
docker exec permissions-control-plane systemctl restart kubelet

# read logs
docker exec -it permissions-control-plane bash
sed -i.bak 's/verbosity: 0/verbosity: 4/' /var/lib/kubelet/config.yaml
systemctl restart kubelet
journalctl -u kubelet -f

# check syscalls
sudo strace -f -p $(pidof /usr/bin/kubelet) --trace=file 2>trace.log
