

https://github.com/projectcontour/contour/pull/4990
https://github.com/projectcontour/contour/pull/5008

kubectl apply -f manifests/echoserver-ipfilter.yaml

http http://blocked.127-0-0-101.nip.io



kubectl -n projectcontour port-forward daemonset/envoy 9001
http http://localhost:9001/config_dump| jq '.configs[].dynamic_route_configs'
