
export WORKDIR=~/work/devenvs/contour
###export KUBEBUILDER_ASSETS=~/work/contour-devenv/kubebuilder/bin

curl -L -o kubebuilder-bin https://go.kubebuilder.io/dl/latest/$(go env GOOS)/$(go env GOARCH)
chmod +x kubebuilder-bin
mv kubebuilder-bin ~/go/bin/kubebuilder
curl -sSLo envtest-bins.tar.gz "https://storage.googleapis.com/kubebuilder-tools/kubebuilder-tools-1.19.2-$(go env GOOS)-$(go env GOARCH).tar.gz"
rm -rf kubebuilder
tar zxf envtest-bins.tar.gz


# start new cluster
kind delete cluster --name contour
kind create cluster --config $WORKDIR/configs/kind-cluster-config.yaml --name contour


mkdir -p .vscode
cp ~/work/devencs/contour/configs/contour-vscode-launch.json .vscode/launch.json


##############################################################################
#
# Running a container
#

# build container
make container
docker tag ghcr.io/projectcontour/contour:$(git rev-parse --short HEAD) localhost/contour:latest
kind load docker-image localhost/contour:latest --name contour

kubectl apply -f examples/contour

cat <<EOF | kubectl -n projectcontour patch deployment contour --patch-file=/dev/stdin
spec:
  template:
    spec:
      containers:
      - name: contour
        image: localhost/contour:latest
        imagePullPolicy: Never
EOF

cat <<EOF | kubectl -n projectcontour patch daemonset envoy --patch-file=/dev/stdin
spec:
  template:
    spec:
      containers:
      - name: shutdown-manager
        image: localhost/contour:latest
        imagePullPolicy: Never
      initContainers:
      - name: envoy-initconfig
        image: localhost/contour:latest
        imagePullPolicy: Never
EOF


kubectl -n projectcontour patch daemonset envoy -f -




# running old versions
curl https://raw.githubusercontent.com/projectcontour/contour/release-1.0/examples/render/contour.yaml -o contour.yaml
curl https://raw.githubusercontent.com/projectcontour/contour/release-1.1/examples/render/contour.yaml -o contour.yaml
curl https://raw.githubusercontent.com/projectcontour/contour/release-1.7/examples/render/contour.yaml -o contour.yaml


#####################################################################################
#
# Troubleshooting
#

# compile binary
go build -o contour -v github.com/projectcontour/contour/cmd/contour

kubectl -n projectcontour logs deployments/contour
kubectl -n projectcontour logs daemonsets/envoy envoy



# monitor SDS
./contour cli sds --cafile=ca.crt --cert-file=tls.crt --key-file=tls.key

# contour debug api

Debug API https://projectcontour.io/docs/master/troubleshooting/
kubectl -n projectcontour port-forward pod/$(kubectl -n projectcontour get pod -l app=contour -o jsonpath='{.items[0].metadata.name}') 6060
http localhost:6060/debug/dag | dot -Tpng -o dag.png


# create envoy bootstrap
./contour bootstrap tmp/boostrap.yaml --envoy-cert-file=$WORKDIR/certs/envoy.pem --envoy-key-file=$WORKDIR/certs/envoy-key.pem --envoy-cafile=$WORKDIR/certs/internal-root-ca.pem --resources-dir=tmp

# capture grpc traffic
sudo nsenter --target $(pidof contour) --net wireshark -f "port 8001" -k


# contour metrics
kubectl -n projectcontour port-forward pod/$(kubectl -n projectcontour get pod -l app=contour -o jsonpath='{.items[0].metadata.name}') 8000
http localhost:8000/metrics

# Envoy admin interface operations
https://www.envoyproxy.io/docs/envoy/latest/operations/admin

# envoy config_dump
kubectl -n projectcontour port-forward daemonset/envoy 9001
http http://localhost:9001/help
http http://localhost:9001/config_dump?include_eds | jq -C . | less
http http://localhost:9001/config_dump| jq '.configs[].dynamic_active_clusters'
http http://localhost:9001/config_dump| jq '.configs[].dynamic_route_configs'



#####################################################################################
#
# Run contour directly from source directory
#

kubectl apply -f https://projectcontour.io/quickstart/contour.yaml
kubectl apply -f examples/contour

# create endpoints that directs traffic to host, to execute controllers directly from source code without deploying
sed "s/REPLACE_ADDRESS_HERE/$(docker network inspect kind | jq -r '.[0].IPAM.Config[0].Gateway')/" manifests/contour-endpoints-dev.yaml | kubectl apply -f -

cat manifests/contour-endpoints-dev.yaml



kubectl -n projectcontour scale deployment --replicas=0 contour
kubectl -n projectcontour rollout restart daemonset envoy
kubectl apply -f manifests/echoservers-tls.yaml





# generate certificates
certyaml --destination certs configs/certs.yaml


kubectl create secret generic httpbin --dry-run=client -o yaml --from-file=certs/httpbin.pem --from-file=certs/httpbin-key.pem | kubectl apply -f -
kubectl create secret tls echoserver-cert --cert=certs/echoserver.pem --key=certs/echoserver-key.pem --dry-run=client -o yaml | kubectl apply -f -
kubectl create secret generic internal-root-ca --from-file=ca.crt=certs/internal-root-ca.pem --dry-run=client -o yaml | kubectl apply -f -
kubectl create secret tls ingress --cert=certs/ingress.pem --key=certs/ingress-key.pem --dry-run=client -o yaml | kubectl apply -f -
kubectl create secret tls client --cert=certs/client-1.pem --key=certs/client-1-key.pem --dry-run=client -o yaml | kubectl apply -f -
kubectl create secret generic client-root-ca-1 --from-file=cacert.pem=certs/client-root-ca-1.pem --dry-run=client -o yaml | kubectl apply -f -
kubectl -n projectcontour create secret tls envoy-client-cert --cert=certs/envoy.pem --key=certs/envoy-key.pem --dry-run=client -o yaml | kubectl apply -f -
kubectl -n projectcontour create secret tls fallback-cert --cert=certs/fallback.pem --key=certs/fallback-key.pem --dry-run=client -o yaml | kubectl apply -f -
kubectl -n projectcontour create secret generic metrics-cert --from-file=ca.pem=certs/metrics-client-root-ca.pem --from-file=tls.crt=certs/metrics.pem --from-file=tls.key=certs/metrics-key.pem --dry-run=client -o yaml | kubectl apply -f -


kubectl -n projectcontour get secret contourcert -o jsonpath='{..ca\.crt}' | base64 -d > ca.crt
kubectl -n projectcontour get secret contourcert -o jsonpath='{..tls\.crt}' | base64 -d > tls.crt
kubectl -n projectcontour get secret contourcert -o jsonpath='{..tls\.key}' | base64 -d > tls.key

go run github.com/projectcontour/contour/cmd/contour serve --xds-address=0.0.0.0 --xds-port=8001 --envoy-service-http-port=8080 --envoy-service-https-port=8443 --contour-cafile=ca.crt --contour-cert-file=tls.crt --contour-key-file=tls.key --config-path=$WORKDIR/configs/contour.yaml


http http://insecure.127-0-0-101.nip.io/foo
http --verify=certs/external-root-ca.pem https://protected.127-0-0-101.nip.io/foo

kubectl -n projectcontour get secret envoy-client-cert -o yaml | kubectl replace --force -f -


sudo nsenter --target $(pidof envoy) --net wireshark -f "port 443" -k
sudo nsenter --target $(pgrep hypercorn) --net wireshark -f "port 443" -k



#######

kubectl apply -f https://projectcontour.io/quickstart/contour.yaml
kubectl -n projectcontour get daemonsets.apps envoy -o yaml | sed "s/xds-address=contour/xds-address=$(docker network inspect kind | jq -r '.[0].IPAM.Config[0].Gateway')/" | kubectl apply -f -

kubectl -n projectcontour get secret contourcert -o jsonpath='{..ca\.crt}' | base64 -d > ca.crt
kubectl -n projectcontour get secret contourcert -o jsonpath='{..tls\.crt}' | base64 -d > tls.crt
kubectl -n projectcontour get secret contourcert -o jsonpath='{..tls\.key}' | base64 -d > tls.key
go run github.com/projectcontour/contour/cmd/contour serve --xds-address=0.0.0.0 --xds-port=8001 --envoy-service-http-port=8080 --envoy-service-https-port=8443 --contour-cafile=ca.crt --contour-cert-file=tls.crt --contour-key-file=tls.key --debug

kubectl apply -f manifests/echoservers.yaml

http http://httpproxy.127-0-0-101.nip.io/
http http://httpproxy.127-0-0-101.nip.io/subpage
http http://httpproxy.127-0-0-101.nip.io/included/
http http://ingress.127-0-0-101.nip.io/
http http://ingress.127-0-0-101.nip.io/subpage


http http://httpproxy.127-0-0-101.nip.io/ | jq '.path, .pod' ; http http://httpproxy.127-0-0-101.nip.io/subpage | jq '.path, .pod' ; http http://httpproxy.127-0-0-101.nip.io/included/ | jq '.path, .pod'






### metrics

kubectl apply -f https://projectcontour.io/quickstart/contour.yaml


make container
docker tag ghcr.io/projectcontour/contour:$(git rev-parse --short HEAD) localhost/contour:latest
kind load docker-image localhost/contour:latest --name contour


kubectl -n projectcontour create secret generic metrics-cert --from-file=certs/metrics-client-root-ca.pem --from-file=certs/metrics.pem --from-file=certs/metrics-key.pem --dry-run=client -o yaml | kubectl apply -f -

kubectl -n projectcontour create configmap contour --from-file=contour.yaml=configs/contour-metrics-tls-in-cluster.yaml --dry-run=client -o yaml | kubectl apply -f -

cat <<EOF | kubectl -n projectcontour patch deployment contour --patch-file=/dev/stdin
spec:
  template:
    spec:
      containers:
      - name: contour
        image: localhost/contour:latest
        imagePullPolicy: Never
        volumeMounts:
        - mountPath: /metrics-cert
          name: metrics-cert
      volumes:
      - name: metrics-cert
        secret:
          secretName: metrics-cert
EOF

cat <<EOF | kubectl -n projectcontour patch daemonset envoy --patch-file=/dev/stdin
spec:
  template:
    spec:
      containers:
      - name: envoy
        volumeMounts:
        - mountPath: /metrics-cert
          name: metrics-cert
      - name: shutdown-manager
        image: localhost/contour:latest
        imagePullPolicy: Never
      initContainers:
      - name: envoy-initconfig
        image: localhost/contour:latest
        imagePullPolicy: Never
      volumes:
      - name: metrics-cert
        secret:
          secretName: metrics-cert
EOF

kubectl -n projectcontour get pod

#
# Tests
#

kubectl -n projectcontour port-forward deployment/contour 8003:8003
http --verify=certs/internal-root-ca.pem --cert=certs/metrics-client-1.pem --cert-key=certs/metrics-client-1-key.pem https://localhost:8003/metrics   # success
http --verify=certs/internal-root-ca.pem --cert=certs/untrusted-client.pem --cert-key=certs/untrusted-client-key.pem https://localhost:8003/metrics   # failure

kubectl -n projectcontour port-forward daemonset/envoy 8003:8003
http --verify=certs/internal-root-ca.pem --cert=certs/metrics-client-1.pem --cert-key=certs/metrics-client-1-key.pem https://localhost:8003/stats   # success
http --verify=certs/internal-root-ca.pem --cert=certs/untrusted-client.pem --cert-key=certs/untrusted-client-key.pem https://localhost:8003/stats   # failure

kubectl -n projectcontour port-forward deployment/contour 8000:8000
http http://localhost:8000/healthz  # success
http http://localhost:8000/metrics  # failure

kubectl -n projectcontour port-forward daemonset/envoy 8002:8002
http http://localhost:8002/ready  # success
http http://localhost:8002/stats  # failure


kubectl -n projectcontour port-forward deployment/contour 8003:8003
openssl s_client -connect localhost:8003 | openssl x509 -noout -text | grep Validity -A2
sslyze --cert=certs/metrics-client-1.pem --key=certs/metrics-client-1-key.pem localhost:8003

kubectl -n projectcontour port-forward daemonset/envoy 8003:8003
openssl s_client -connect localhost:8003 | openssl x509 -noout -text | grep Validity -A2
sslyze --cert=certs/metrics-client-1.pem --key=certs/metrics-client-1-key.pem localhost:8003


# rotate certs
rm certs/metrics.pem certs/metrics-key.pem
certyaml --destination certs configs/certs.yaml
kubectl -n projectcontour create secret generic metrics-cert --from-file=certs/metrics-client-root-ca.pem --from-file=certs/metrics.pem --from-file=certs/metrics-key.pem --dry-run=client -o yaml | kubectl apply -f -

kubectl -n projectcontour port-forward deployment/contour 8003:8003
openssl s_client -connect localhost:8003 | openssl x509 -noout -text | grep Validity -A2

kubectl -n projectcontour port-forward daemonset/envoy 8003:8003
openssl s_client -connect localhost:8003 | openssl x509 -noout -text | grep Validity -A2




kubectl apply -f manifests/shell.yaml
kubectl cp certs/metrics-client-1.pem shell:.
kubectl cp certs/metrics-client-1-key.pem shell:.
kubectl cp certs/internal-root-ca.pem shell:.
kubectl cp certs/untrusted-client.pem shell:.
kubectl cp certs/untrusted-client-key.pem shell:.

kubectl exec -it shell -- ash

CONTOUR_ADDR=$(kubectl -n projectcontour get pod -lapp=contour -o=jsonpath={..podIP})
ENVOY_ADDR=$(kubectl -n projectcontour get pod -lapp=envoy -o=jsonpath={..podIP})

kubectl exec -it shell -- http --verify=internal-root-ca.pem --cert=metrics-client-1.pem --cert-key=metrics-client-1-key.pem https://$CONTOUR_ADDR:8002/metrics
kubectl exec -it shell -- http --verify=false --cert=metrics-client-1.pem --cert-key=metrics-client-1-key.pem https://$CONTOUR_ADDR:8002/metrics
kubectl exec -it shell -- http --verify=false --cert=untrusted-client.pem --cert-key=untrusted-client-key.pem https://$CONTOUR_ADDR:8002/metrics
kubectl exec -it shell -- http http://$CONTOUR_ADDR:8000/healthz

kubectl exec -it shell -- http --verify=false --cert=metrics-client-1.pem --cert-key=metrics-client-1-key.pem https://$ENVOY_ADDR:8003/stats
kubectl exec -it shell -- http --verify=false --cert=untrusted-client.pem --cert-key=untrusted-client-key.pem https://$ENVOY_ADDR:8003/stats
kubectl exec -it shell -- http http://$ENVOY_ADDR:8002/ready

kubectl exec -it shell -- openssl s_client -connect $ENVOY_ADDR:8002

kubectl -n projectcontour port-forward envoy:8002

http --verify=certs/internal-root-ca.pem --cert=certs/metrics-client-1.pem --cert-key=certs/metrics-client-1-key.pem https://localhost:8002/metrics
http --verify=false --cert=certs/untrusted-client.pem --cert-key=certs/untrusted-client-key.pem https://localhost:8002/metrics

kubectl -n projectcontour port-forward $(kubectl -n projectcontour get pod -lapp=envoy -ojsonpath="{.items[0].metadata.name}") 8003
http --verify=certs/internal-root-ca.pem --cert=certs/metrics-client-1.pem --cert-key=certs/metrics-client-1-key.pem https://localhost:8003/stats


make setup-kind-cluster
make run-e2e

export CONTOUR_E2E_LOCAL_HOST=$(ifconfig | grep inet | grep -v '::' | grep -v 127.0.0.1 | head -n1 | awk '{print $2}')
USE_CONTOUR_CONFIGURATION_CRD=true ginkgo -tags=e2e -mod=readonly -skip-package=upgrade -r -v ./test/e2e/

ginkgo -tags=e2e -mod=readonly -r -v ./test/e2e/infra/
ginkgo -tags=e2e -mod=readonly --fail-fast --until-it-fails -r -v ./test/e2e/infra/

make cleanup-kind

if CurrentSpecReport().Failed {
   os.Exit(1)
}



##############


Install specific old version

https://raw.githubusercontent.com/projectcontour/contour/v1.18.3/examples/render/contour.yaml
