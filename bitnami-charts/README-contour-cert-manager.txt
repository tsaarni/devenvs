

https://github.com/bitnami/charts/pull/29416

cd ~/work/bitnami-charts/bitnami/cert-manager/
helm dependency build
helm install certmanager . --set installCRDs=true


cd ~/work/bitnami-charts/bitnami/contour/
helm dependency build


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




# multiple releases in single namespace

helm install contour1 . --set contour.ingressClass.name=contour1,useCertManager=true
helm install contour2 . --set contour.ingressClass.name=contour2,contour.manageCRDs=false,useCertManager=true




kubectl apply -f - <<EOF
apiVersion: projectcontour.io/v1
kind: HTTPProxy
metadata:
  name: httpproxy1
spec:
  virtualhost:
    fqdn: example.com
  routes:
    - services:
        - name: echoserver
          port: 80
  ingressClassName: contour1
EOF


kubectl apply -f - <<EOF
apiVersion: projectcontour.io/v1
kind: HTTPProxy
metadata:
  name: httpproxy2
spec:
  virtualhost:
    fqdn: example.com
  routes:
    - services:
        - name: echoserver
          port: 80
  ingressClassName: contour2
EOF
