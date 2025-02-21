





diff --git a/apis/projectcontour/v1/httpproxy.go b/apis/projectcontour/v1/httpproxy.go
index 1978f866..922eec9c 100644
--- a/apis/projectcontour/v1/httpproxy.go
+++ b/apis/projectcontour/v1/httpproxy.go
@@ -317,7 +317,7 @@ type VirtualHost struct {
        // The fully qualified domain name of the root of the ingress tree
        // all leaves of the DAG rooted at this object relate to the fqdn.
        //
-       // +kubebuilder:validation:Pattern="^(\\*\\.)?[a-z0-9]([-a-z0-9]*[a-z0-9])?(\\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*$"
+       // +kubebuilder:validation:Pattern="^\\*|(\\*\\.)?[a-z0-9]([-a-z0-9]*[a-z0-9])?(\\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*$"
        Fqdn string `json:"fqdn"`

        // If present the fields describes TLS properties of the virtual



make generate
kubectl apply -f examples/contour/01-crds.yaml


rm -rf certs
mkdir certs
certyaml --destination certs configs/certs.yaml



kubectl apply -f manifests/echoserver-wildcard-fqdn.yaml
kubectl get httpproxies.projectcontour.io





kubectl create secret tls echoserver1 --cert=certs/wildcard-ingress-1.pem --key=certs/wildcard-ingress-1-key.pem --dry-run=client -o yaml | kubectl apply -f -
kubectl create secret tls echoserver2 --cert=certs/wildcard-ingress-2.pem --key=certs/wildcard-ingress-2-key.pem --dry-run=client -o yaml | kubectl apply -f -
kubectl create secret tls echoserver3 --cert=certs/ingress.pem --key=certs/ingress-key.pem --dry-run=client -o yaml | kubectl apply -f -

http --verify=certs/external-root-ca.pem https://foo.127-0-0-101.nip.io
http --verify=certs/external-root-ca.pem https://echoserver.127-0-0-101.nip.io

wireshark -k -i lo port 443


curl --cacert certs/external-root-ca.pem https://foo.127-0-0-101.nip.io
curl --cacert certs/external-root-ca.pem --resolve foo.example.com:443:127.0.0.101 https://foo.example.com

openssl s_client -connect 127.0.0.101:443 -servername foo.example.com -CAfile certs/external-root-ca.pem -showcerts </dev/null 2>/dev/null | openssl x509 -noout -text

openssl s_client -connect 127.0.0.101:443 -CAfile certs/external-root-ca.pem -showcerts </dev/null 2>/dev/null | openssl x509 -noout -text


import ssl
import socket

context = ssl.create_default_context(cafile="certs/external-root-ca.pem")
sock = socket.create_connection(("127.0.0.101", 443))
conn = context.wrap_socket(sock, server_hostname="foo.example.com")
conn.sendall(b"GET / HTTP/1.1\r\nHost: foo.example.com\r\n\r\n")
print(conn.recv(4096).decode().replace("\r\n", "\n"))
conn.close()

context = ssl.create_default_context(cafile="certs/external-root-ca.pem")
context.check_hostname = False
sock = socket.create_connection(("127.0.0.101", 443))
conn = context.wrap_socket(sock)
conn.sendall(b"GET / HTTP/1.1\r\nHost: foo.example.com\r\n\r\n")
print(conn.recv(4096).decode().replace("\r\n", "\n"))
conn.close()


context = ssl.create_default_context(cafile="certs/external-root-ca.pem")
context.check_hostname = False
sock = socket.create_connection(("127.0.0.101", 443))
conn = context.wrap_socket(sock, server_hostname="1.2.3.4")
conn.sendall(b"GET / HTTP/1.1\r\nHost: 1.2.3.4\r\n\r\n")
print(conn.recv(4096).decode().replace("\r\n", "\n"))
conn.close()




# create HTTPProxy with IP address as hostname
cat <<EOF | kubectl apply -f -
apiVersion: projectcontour.io/v1
kind: HTTPProxy
metadata:
  name: echoserver-with-ipaddr
spec:
  virtualhost:
    fqdn: 127.0.0.101
  routes:
    - services:
        - name: echoserver
          port: 80
EOF

http https://127.0.0.101/foo


kubectl create secret tls ingress --cert=certs/ingress.pem --key=certs/ingress-key.pem --dry-run=client -o yaml | kubectl apply -f -

# create HTTPProxy with IP address as hostname
cat <<EOF | kubectl apply -f -
apiVersion: projectcontour.io/v1
kind: HTTPProxy
metadata:
  name: echoserver-with-ipaddr
spec:
  virtualhost:
    fqdn: 127.0.0.101
    #fqdn: echoserver.127-0-0-101.nip.io
    tls:
      secretName: ingress
  routes:
    - services:
        - name: echoserver
          port: 80
EOF

http --verify=certs/external-root-ca.pem https://echoserver.127-0-0-101.nip.io/foo
http --verify=false https://echoserver.127-0-0-101.nip.io/foo
http --verify=false https://127.0.0.101/foo
