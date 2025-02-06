






# start a new cluster

kind delete cluster --name contour
kind create cluster --config ~/work/devenvs/contour/configs/kind-cluster-config-with-custom-port.yaml --name contour
kubectl apply -f https://projectcontour.io/quickstart/contour-gateway-provisioner.yaml



kubectl apply -f manifests/gateway-provisioner.yaml
kubectl apply -f manifests/echoserver-gatewayapi.yaml


http http://echoserver.127-0-0-101.nip.io






####################
#
# WITH METALLB


# start a new cluster

make install-provisioner-working


kubectl apply -f manifests/gateway-provisioner.yaml
kubectl apply -f manifests/echoserver-gatewayapi.yaml

# address is from metallb ip pool
http http://echoserver-172.20.255.200.nip.io


$ kubectl -n metallb-system get ipaddresspools.metallb.io pool
NAME   AUTO ASSIGN   AVOID BUGGY IPS   ADDRESSES
pool   true          false             ["172.20.255.200-172.20.255.250"]

$ kubectl -n projectcontour get service envoy-mygateway
NAME              TYPE           CLUSTER-IP     EXTERNAL-IP      PORT(S)        AGE
envoy-mygateway   LoadBalancer   10.96.135.70   172.20.255.200   80:30093/TCP   34m





# check the certificates
kubectl -n projectcontour get secrets contourcert-mygateway -o jsonpath='{.data.tls\.crt}' | base64 -d | openssl x509 -text -noout
kubectl -n projectcontour get secrets contourcert-mygateway -o jsonpath='{.data.ca\.crt}' | base64 -d | openssl x509 -text -noout
kubectl -n projectcontour get secrets envoycert-mygateway -o jsonpath='{.data.tls\.crt}' | base64 -d | openssl x509 -text -noout
kubectl -n projectcontour get secrets envoycert-mygateway -o jsonpath='{.data.ca\.crt}' | base64 -d | openssl x509 -text -noout






kubectl apply -f - <<EOF
kind: GatewayClass
apiVersion: gateway.networking.k8s.io/v1
metadata:
  name: contour-with-envoy-deployment
spec:
  controllerName: projectcontour.io/gateway-controller
  parametersRef:
    kind: ContourDeployment
    group: projectcontour.io
    name: contour-with-envoy-deployment-params
    namespace: projectcontour
---
kind: ContourDeployment
apiVersion: projectcontour.io/v1alpha1
metadata:
  namespace: projectcontour
  name: contour-with-envoy-deployment-params
spec:
  envoy:
    workloadType: Deployment
EOF
