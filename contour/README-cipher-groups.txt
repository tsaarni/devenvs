https://github.com/projectcontour/contour/issues/6380
https://github.com/projectcontour/contour/pull/6461



go run github.com/projectcontour/contour/cmd/contour serve --xds-address=0.0.0.0 --xds-port=8001 --envoy-service-http-port=8080 --envoy-service-https-port=8443 --contour-cafile=ca.crt --contour-cert-file=tls.crt --contour-key-file=tls.key --config-path=$HOME/work/devenvs/contour/configs/contour-config-cipher-group.yaml



kubectl create secret tls ingress --cert=certs/ingress.pem --key=certs/ingress-key.pem --dry-run=client -o yaml | kubectl apply -f -
kubectl apply -f manifests/echoserver-tls.yaml


http --verify=certs/external-root-ca.pem https://protected.127-0-0-101.nip.io


sslyze protected.127-0-0-101.nip.io:443
