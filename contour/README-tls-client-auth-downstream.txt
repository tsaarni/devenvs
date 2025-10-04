
Manual end-to-end test procedure for external client certificate validation with Contour and Envoy

https://gist.github.com/tsaarni/4bc8d5e03bc180437a1aa639




rm -rf certs
mkdir certs
certyaml --destination certs configs/certs.yaml



kubectl create secret tls ingress --cert=certs/ingress.pem --key=certs/ingress-key.pem --dry-run=client -o yaml | kubectl apply -f -
kubectl create secret tls echoserver-cert --cert=certs/echoserver.pem --key=certs/echoserver-key.pem --dry-run=client -o yaml | kubectl apply -f -
kubectl create secret generic internal-root-ca --from-file=ca.crt=certs/internal-root-ca.pem --dry-run=client -o yaml | kubectl apply -f -
kubectl -n projectcontour create secret tls envoy-client-cert --cert=certs/envoy.pem --key=certs/envoy-key.pem --dry-run=client -o yaml | kubectl apply -f -


kubectl -n projectcontour create configmap contour --from-file=contour.yaml=configs/contour-config-with-upstream-tls.yaml --dry-run=client -o yaml | kubectl -n projectcontour apply -f -


kubectl -n projectcontour rollout restart deployment contour

kubectl apply -f manifests/echoserver-backend-tls.yaml



http http://protected.127-0-0-101.nip.io



rm certs/envoy*.pem
certyaml --destination certs configs/certs.yaml
kubectl -n projectcontour create secret tls envoy-client-cert --cert=certs/envoy.pem --key=certs/envoy-key.pem --dry-run=client -o yaml | kubectl apply -f -
