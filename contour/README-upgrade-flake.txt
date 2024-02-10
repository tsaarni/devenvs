

kind delete cluster --name contour
kind create cluster --config ~/work/devenvs/contour/configs/kind-cluster-config.yaml --name contour




make container
docker tag ghcr.io/projectcontour/contour:$(git rev-parse --short HEAD) localhost/contour:main
kind load docker-image localhost/contour:latest --name contour





kubectl apply -f https://projectcontour.io/quickstart/v1.25.0/contour-gateway-provisioner.yaml
kubectl rollout status -n gateway-system deployment/gateway-api-admission-server
kubectl apply -f contour-gateway-provisioner.yaml
kubectl rollout status -n gateway-system deployment/gateway-api-admission-server



cat <<EOF | kubectl apply -f -
kind: GatewayClass
apiVersion: gateway.networking.k8s.io/v1beta1
metadata:
  name: example
spec:
  controllerName: projectcontour.io/gateway-controller
EOF



kubectl delete gatewayclass example
kubectl delete namespaces gateway-system projectcontour
kubectl delete clusterrolebinding contour-gateway-provisioner gateway-api-admission
kubectl delete clusterroles contour-gateway-provisioner gateway-api-admission
kubectl delete ValidatingWebhookConfiguration gateway-api-admission



#
docker container inspect contour-control-plane | grep -i cpu
docker update --cpus=".5" contour-control-plane contour-worker

CONTOUR_UPGRADE_FROM_VERSION=v1.25.0 CONTOUR_E2E_IMAGE=localhost/contour:main go run github.com/onsi/ginkgo/v2/ginkgo -tags=e2e -mod=readonly -randomize-all -poll-progress-after=300s -vv ./test/e2e/upgrade
CONTOUR_UPGRADE_FROM_VERSION=v1.25.0 CONTOUR_E2E_IMAGE=ghcr.io/projectcontour/contour:main go run github.com/onsi/ginkgo/v2/ginkgo -tags=e2e -mod=readonly -randomize-all -poll-progress-after=300s -vv ./test/e2e/upgrade



########

./test/scripts/cleanup.sh
./test/scripts/make-kind-cluster.sh

make container
docker tag ghcr.io/projectcontour/contour:$(git rev-parse --short HEAD) ghcr.io/projectcontour/contour:main
kind load docker-image ghcr.io/projectcontour/contour:main --name contour-e2e

docker update --cpus=".4" contour-e2e-control-plane contour-e2e-worker

CONTOUR_UPGRADE_FROM_VERSION=v1.25.0 CONTOUR_E2E_IMAGE=ghcr.io/projectcontour/contour:main go run github.com/onsi/ginkgo/v2/ginkgo -tags=e2e -mod=readonly -randomize-all -poll-progress-after=300s -vv ./test/e2e/upgrade 2>&1 | ts '%H:%M:%.S'

#
