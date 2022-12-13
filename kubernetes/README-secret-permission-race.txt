
# Create cluster
kind delete cluster --name permissions
kind create cluster --name permissions

# build test app
docker build --tag dirwatcher:latest docker/dirwatcher/
kind load docker-image dirwatcher:latest --name permissions

# deploy test app
kubectl apply -f manifests/secret.yaml
kubectl delete -f manifests/dirwatcher-nonroot.yaml --force
kubectl apply -f manifests/dirwatcher-nonroot.yaml

# update secret
while true; do kubectl create secret generic mysecret --from-file=password=/proc/sys/kernel/random/uuid --dry-run=client -o yaml | kubectl apply -f -; sleep 1; done

# read testapp logs
kubectl logs -f dirwatcher-nonroot


make WHAT=cmd/kubelet
ls -l ./_output/local/bin/linux/amd64/kubelet

docker exec permissions-control-plane systemctl stop kubelet
docker cp ./_output/local/bin/linux/amd64/kubelet permissions-control-plane:/usr/bin/kubelet
docker exec permissions-control-plane systemctl start kubelet




docker exec -it permissions-control-plane bash
journalctl -u kubelet -f



sudo strace -f -p $(pidof /usr/bin/kubelet) --trace=file 2>trace.log
