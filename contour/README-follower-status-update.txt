
# Non-leader contour controller pod memory keeps increasing until OOM #6860￼
https://github.com/projectcontour/contour/issues/6860


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
