

https://github.com/projectcalico/calico/issues/7983


# reproduction case
https://github.com/tsaarni/kubernetes-tcp-retransmit-retry-timeout/
git clone git@github.com:tsaarni/kubernetes-tcp-retransmit-retry-timeout.git



# create cluster
kind delete cluster --name echo
kind create cluster --config configs/kind-cluster-config-for-calico.yaml --name echo

kubectl apply -f https://raw.githubusercontent.com/projectcalico/calico/v3.26.1/manifests/calico.yaml
kubectl apply -f https://raw.githubusercontent.com/projectcalico/calico/v3.26.4/manifests/calico.yaml
kubectl apply -f https://raw.githubusercontent.com/projectcalico/calico/v3.27.4/manifests/calico.yaml
kubectl apply -f https://raw.githubusercontent.com/projectcalico/calico/v3.28.0/manifests/calico.yaml
kubectl apply -f https://raw.githubusercontent.com/projectcalico/calico/v3.28.1/manifests/calico.yaml


# install/delete echo client/server test app
kubectl apply -f manifests/echo.yaml
kubectl delete -f manifests/echo.yaml


# logs
kubectl logs deployment/client -f
kubectl logs deployment/server -f

# force delete pods
kubectl delete pod -l app=server --force
kubectl delete pod -l app=client --force


# capture traffic from client pod
sudo nsenter -t $(pgrep -f "echo client") --net wireshark -i any -k






# capture traffic from workers

# get IP addresses (pods + service)
kubectl get pod -o wide
kubectl get service server


# capture traffic from workers
sudo nsenter -t $(docker inspect --format '{{.State.Pid}}' $(kubectl get pod -l app=client -o jsonpath='{.items[0].spec.nodeName}')) --net wireshark -i any -k -o gui.window_title:client-worker -Y "tcp.port==8000"
sudo nsenter -t $(docker inspect --format '{{.State.Pid}}' $(kubectl get pod -l app=server -o jsonpath='{.items[0].spec.nodeName}')) --net wireshark -i any -k -o gui.window_title:server-worker -Y "tcp.port==8000"
￼
# capture traffic from pods
sudo nsenter -t $(pgrep -f "/app/echo client") --net wireshark -i any -k -o gui.window_title:client-pod -Y "tcp.port==8000"   # client
sudo nsenter -t $(pgrep -f "/app/echo --catch-sigterm server") --net wireshark -i any -k -o gui.window_title:server-pod -Y "tcp.port==8000"   # server


# or: -f "port 8000" for capture filter


# list interfaces on pods
kubectl exec $(kubectl get pod -l app=client -o jsonpath='{.items[0].metadata.name}') -- ip addr

# list interfaces on workers
docker exec echo-worker ip addr
docker exec echo-worker2 ip addr


*** Building and troubleshooting Felix



# enable debug
kubectl -n kube-system patch FelixConfiguration default --type=merge  -p '{"spec":{"logSeverityScreen": "Debug"}}'
kubectl -n kube-system get FelixConfiguration default -o yaml

# change calico config
kubectl -n kube-system edit FelixConfiguration default

# restart calico
kubectl -n kube-system rollout restart daemonset calico-node


# check logs on specific worker
kubectl -n kube-system logs $(kubectl -n kube-system get pods --field-selector spec.nodeName=echo-worker -l k8s-app=calico-node -o name)
kubectl -n kube-system logs $(kubectl -n kube-system get pods --field-selector spec.nodeName=echo-worker2 -l k8s-app=calico-node -o name)



make -C node image
kind load docker-image --name echo docker.io/library/node:latest-amd64
kubectl -n kube-system set image daemonset/calico-node calico-node=docker.io/library/node:latest-amd64    # note: the image name from build really is: docker.io/library/node
kubectl -n kube-system rollout restart daemonset calico-node


kubectl get pod -o wide
kubectl -n kube-system get pod -o wide

kubectl delete pod -l app=client --force
kubectl delete pod -l app=server --force

kubectl logs deployment/client -f



kubectl describe pod -l app=client
kubectl describe pod -l app=server

# client
echo "#### routes on client worker $(date)"
docker exec $(kubectl get pod -l app=client -o jsonpath='{.items[0].spec.nodeName}') sh -xc "ip addr; ip route"
echo "#### iptables on client worker $(date)"
docker exec $(kubectl get pod -l app=client -o jsonpath='{.items[0].spec.nodeName}') iptables -L -v -n
echo "#### conntrack on client worker $(date)"
docker exec $(kubectl get pod -l app=client -o jsonpath='{.items[0].spec.nodeName}') conntrack -L

# server
echo "#### routes on server worker $(date)"
docker exec $(kubectl get pod -l app=server -o jsonpath='{.items[0].spec.nodeName}') sh -xc "ip addr; ip route"
echo "#### iptables on server worker $(date)"
docker exec $(kubectl get pod -l app=server -o jsonpath='{.items[0].spec.nodeName}') iptables -L -v -n
echo "#### conntrack on server worker $(date)"
docker exec $(kubectl get pod -l app=server -o jsonpath='{.items[0].spec.nodeName}') conntrack -L


# netstat on server pod
kubectl exec $(kubectl get pod -l app=server -o jsonpath='{.items[0].metadata.name}') -- netstat -tuln



docker exec echo-worker2 sh -c "sysctl -a | grep \\.rp_filter"


for iface in /proc/sys/net/ipv4/conf/*; do echo 1 > $iface/log_martians; done


docker exec $(kubectl get pod -l app=server -o jsonpath='{.items[0].spec.nodeName}') iptables -A INPUT -m conntrack --ctstate INVALID -j LOG --log-prefix "Dropped packet: " --log-level 7
docker exec $(kubectl get pod -l app=server -o jsonpath='{.items[0].spec.nodeName}') iptables -A FORWARD -m conntrack --ctstate INVALID -j LOG --log-prefix "Dropped packet: " --log-level 7
docker exec $(kubectl get pod -l app=server -o jsonpath='{.items[0].spec.nodeName}') iptables -A OUTPUT -m conntrack --ctstate INVALID -j LOG --log-prefix "Dropped packet: " --log-level 7

*** Python client

# TODO: consider also using scapy
#  see example in https://github.com/projectcalico/calico/issues/8882


import socket
s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)

#s.setsockopt(socket.SOL_SOCKET, socket.SO_KEEPALIVE, 1)
#s.setsockopt(socket.IPPROTO_TCP, socket.TCP_KEEPIDLE, 1)
#s.setsockopt(socket.IPPROTO_TCP, socket.TCP_KEEPINTVL, 1)
#s.setsockopt(socket.IPPROTO_TCP, socket.TCP_KEEPCNT, 20)

s.connect(("server", 8000))    # s.connect(("[server IP address]", 8000))
s.sendall(b"Hello, world")
s.recv(1024)



*** BPF mode

# enable BPF mode
kubectl get felixconfigurations default -o yaml
kubectl patch felixconfigurations default --type='json' -p='[{"op": "add", "path": "/spec/bpfEnabled", "value": true}]'    # enable EBFP
kubectl patch felixconfigurations default --type='json' -p='[{"op": "add", "path": "/spec/bpfEnabled", "value": false}]'   # disable EBFP

# then restart client and server

kubectl delete pod -l app=client --force
kubectl delete pod -l app=server --force
