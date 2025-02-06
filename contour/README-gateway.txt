

# Gateway manual provisioning



kind delete cluster --name contour
kind create cluster --config ~/work/devenvs/contour/configs/kind-cluster-config-with-custom-port.yaml --name contour



kubectl apply -f https://projectcontour.io/quickstart/contour-gateway.yaml


# Note
# This might fail for the first time
#
# resource mapping not found for name: "example" namespace: "" from "https://projectcontour.io/quickstart/contour-gateway.yaml": no matches for kind "GatewayClass" in version "gateway.networking.k8s.io/v1"
ensure CRDs are installed first
#
# due to a race between creating resources and installing CRDs
# repeat the kubectl apply command to fix the issue.


# Add custom port to envoy daemonset

kubectl -n projectcontour patch daemonset envoy --type='json' -p='[{"op": "add", "path": "/spec/template/spec/containers/1/ports/-", "value": {"containerPort": 1234, "hostPort": 1234, "name": "custom-protocol", "protocol": "TCP"}}]'


# Add custom tcp port to gateway

kubectl apply -f - <<EOF
kind: GatewayClass
apiVersion: gateway.networking.k8s.io/v1
metadata:
  name: example
spec:
  controllerName: projectcontour.io/gateway-controller
  parametersRef:
    group: projectcontour.io
    kind: ContourDeployment
    namespace: projectcontour
    name: contour
---
kind: ContourDeployment
apiVersion: projectcontour.io/v1alpha1
metadata:
  namespace: projectcontour
  name: contour
---
kind: Gateway
apiVersion: gateway.networking.k8s.io/v1
metadata:
  name: contour
  namespace: projectcontour
spec:
  gatewayClassName: example
  listeners:
    - name: http
      protocol: HTTP
      port: 80
      allowedRoutes:
        namespaces:
          from: All
    - name: custom-protocol
      protocol: TCP
      port: 1234
      allowedRoutes:
        namespaces:
          from: All
EOF


kubectl -n projectcontour describe gatewayclass example

kubectl apply -f manifests/echoserver-gatewayapi-tcproute.yaml



http http://echoserver.127-0-0-101.nip.io:1234
