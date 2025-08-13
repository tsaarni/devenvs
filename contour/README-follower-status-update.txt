
# Non-leader contour controller pod memory keeps increasing until OOM #6860￼
https://github.com/projectcontour/contour/issues/6860



kubectl -n projectcontour scale deployment --replicas=1 contour
kubectl -n projectcontour rollout restart daemonset envoy





kubectl apply -f manifests/echoserver.yaml


# simulate update of the status of the envoy service with a load balancer IP
function update_envoy_service_ip() {
  NEW_IP=$(printf "1.2.%d.%d" $((RANDOM%256)) $((RANDOM%256)))
  kubectl -n projectcontour patch service envoy --type=merge --subresource=status --patch '"status": {"loadBalancer":{"ingress":[{"ip":"'"$NEW_IP"'"}]}}'
  kubectl -n projectcontour get service envoy
}

update_envoy_service_ip


# Check that the status was updated
kubectl get httpproxy echoserver -o yaml



# check who is the leader (pod name in HOLDER)
kubectl -n projectcontour get leases

# delete the leader to force a new leader election
kubectl -n projectcontour delete leases leader-elect


kubectl -n projectcontour patch lease leader-elect --type=merge -p '{"spec":{"holderIdentity":"<new-pod-name>"}}'






kubectl apply -f https://projectcontour.io/quickstart/contour.yaml
make container
docker tag ghcr.io/projectcontour/contour:$(git rev-parse --short HEAD) localhost/contour:latest
kind load docker-image localhost/contour:latest --name contour
kubectl -n projectcontour set image deployment/contour contour=localhost/contour:latest


kubectl -n projectcontour patch service envoy --type=merge --subresource=status --patch '"status": {"loadBalancer":{"ingress":[{"ip":"1.2.3.4"}]}}'
kubectl get httpproxy echoserver -o jsonpath='{.status.loadBalancer.ingress[*].ip}'
kubectl -n projectcontour scale deployment --replicas=0 contour
