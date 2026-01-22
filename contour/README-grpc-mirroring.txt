# https://github.com/projectcontour/contour/issues/7277


kind delete cluster --name contour
kind create cluster --config ~/work/devenvs/contour/configs/kind-cluster-config.yaml --name contour

kubectl apply -f https://projectcontour.io/quickstart/contour.yaml

# enable externalname support in contour
kubectl -n projectcontour create configmap contour --from-file=contour.yaml=<(cat <<EOF
enableExternalNameService: true
#listener:
#  http2-max-concurrent-streams: 5
#cluster:
#  circuit-breakers:
#    max-connections: 5
EOF
) --dry-run=client -o yaml | kubectl apply -f -


# restart contour pods
kubectl -n projectcontour delete pod -l app=contour --force


kubectl apply -f manifests/echoserver-grpc-mirror.yaml


# capture wireshark dump
sudo nsenter --target $(kindps contour envoy -o json | jq -r .[0].pids[0].pid) --net wireshark -i any -k -Y 'tcp.port == 8080 && grpc'



grpcurl -plaintext -d '{"message": "Hello"}' echoserver.127-0-0-101.nip.io:80 echo.EchoService/Echo

rm -f foo
touch foo
while true; do grpcurl -plaintext -d '{"message": "Hello"}' echoserver.127-0-0-101.nip.io:80 echo.EchoService/Echo | jq -r .remoteAddr >>foo; done
$ sort -u foo | wc -l
20


while true; do http  echoserver.127-0-0-101.nip.io/metrics|grep connec; done

kubectl port-forward deployment/echoserver-mirror 8080
while true; do http localhost:8080/metrics|grep connec; done



kubectl delete pod -l app=echoserver --force
kubectl delete pod -l app=echoserver-mirror --force

kubectl -n projectcontour delete pod -l app=envoy --force




## https test

http --verify=certs/external-root-ca.pem https://echoserver.127-0-0-101.nip.io

sudo nsenter --target $(kindps contour echoserver-mirror -o json | jq -r .[0].pids[0].pid) --net wireshark -i any -k -Y 'tcp.port == 8443' -o tls.keylog_file:/proc/$(kindps contour echoserver-mirror -o json | jq -r .[0].pids[0].pid)/root/tmp/wireshark-keys.log
