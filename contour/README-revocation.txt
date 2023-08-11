


kubectl -n projectcontour get secret contourcert -o jsonpath='{..ca\.crt}' | base64 -d > ca.crt
kubectl -n projectcontour get secret contourcert -o jsonpath='{..tls\.crt}' | base64 -d > tls.crt
kubectl -n projectcontour get secret contourcert -o jsonpath='{..tls\.key}' | base64 -d > tls.key

go run github.com/projectcontour/contour/cmd/contour serve --xds-address=0.0.0.0 --xds-port=8001 --envoy-service-http-port=8080 --envoy-service-https-port=8443 --contour-cafile=ca.crt --contour-cert-file=tls.crt --contour-key-file=tls.key

certyaml --destination certs configs/certs.yaml

kubectl create secret tls ingress --cert=certs/ingress.pem --key=certs/ingress-key.pem --dry-run=client -o yaml | kubectl apply -f -
kubectl create secret generic client-root-ca-1 --from-file=ca.crt=certs/client-root-ca-1.pem --dry-run=client -o yaml | kubectl apply -f -
kubectl create secret generic client-crl --from-file=crl.pem=certs/client-root-ca-1-crl.pem --dry-run=client -o yaml | kubectl apply -f -

kubectl apply -f manifests/echoserver-mutual-tls-with-revocation.yaml
kubectl get httpproxy echoserver-protected -o yaml

http --ssl tls1.2 --verify=certs/external-root-ca.pem --cert=certs/client-1.pem --cert-key=certs/client-1-key.pem https://protected.127-0-0-101.nip.io # successful
http --ssl tls1.2 --verify=certs/external-root-ca.pem --cert=certs/revoked-client-1.pem --cert-key=certs/revoked-client-1-key.pem https://protected.127-0-0-101.nip.io # revoked
http --ssl tls1.2 --verify=certs/external-root-ca.pem --cert=certs/client-2.pem --cert-key=certs/client-2-key.pem https://protected.127-0-0-101.nip.io # untrusted


wireshark -i lo -f "port 443" -k
