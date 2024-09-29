

https://github.com/bitnami/charts/pull/29416


helm install my-release . --set useCertManager=true,envoy.useHostPort.http=true,envoy.useHostPort.https=true
helm upgrade my-release . --set useCertManager=true,envoy.useHostPort.http=true,envoy.useHostPort.https=true
helm uninstall my-release
 
 
kubectl logs daemonsets/my-release-contour-envoy -c envoy
 
 
kubectl delete secret contourcert
kubectl delete secret envoycert
kubectl delete secret my-release-contour-contour-crt
kubectl delete secret my-release-contour-envoy-crt
kubectl delete secret my-release-contour-ca-crt
 
 
 
kubectl apply -f https://raw.githubusercontent.com/tsaarni/devenvs/main/contour/manifests/echoserver.yaml
http http://echoserver.127-0-0-101.nip.io
 
kubectl get secret my-release-contour-envoy-crt -o jsonpath='{.data.tls\.crt}' | base64 -d | openssl x509 -text -noout
kubectl get secret my-release-contour-contour-crt -o jsonpath='{.data.tls\.crt}' | base64 -d | openssl x509 -text -noout
kubectl get secret my-release-contour-root-ca -o jsonpath='{.data.tls\.crt}' | base64 -d | openssl x509 -text -noout
 
 
 
cat <<EOF | kubectl apply -f -
apiVersion: cert-manager.io/v1
kind: Issuer
metadata:
    name: myselfsigned
spec:
    selfSigned: {}
---
apiVersion: cert-manager.io/v1
kind: Certificate
metadata:
  name: mycert
spec:
  secretName: mycert
  commonName: mycert
  dnsNames:
    - mycert
  issuerRef:
    name: myselfsigned
    kind: Issuer
EOF
