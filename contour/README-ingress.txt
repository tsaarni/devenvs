
# start new cluster
kind delete cluster --name contour
kind create cluster --config ~/work/devenvs/contour/configs/kind-cluster-config.yaml --name contour


kubectl apply -f https://projectcontour.io/quickstart/contour.yaml


kubectl create secret tls ingress --cert=certs/ingress.pem --key=certs/ingress-key.pem --dry-run=client -o yaml | kubectl apply -f -

kubectl apply -f manifests/echoserver-ingress-tls.yaml


http --verify=certs/external-root-ca.pem https://protected.127-0-0-101.nip.io
http http://protected.127-0-0-101.nip.io

