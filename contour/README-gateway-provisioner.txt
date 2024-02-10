

make install-provisioner-working

kubectl apply -f manifests/gateway-provisioner.yaml
kubectl apply -f manifest/echoserver.yaml
