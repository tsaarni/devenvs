

kubectl create secret tls ingress --cert=certs/ingress.pem --key=certs/ingress-key.pem --dry-run=client -o yaml | kubectl apply -f -
kubectl create secret generic internal-root-ca --from-file=ca.crt=certs/internal-root-ca.pem --dry-run=client -o yaml | kubectl apply -f -
kubectl create secret tls echoserver-cert --cert=certs/echoserver.pem --key=certs/echoserver-key.pem --dry-run=client -o yaml | kubectl apply -f -
kubectl -n projectcontour create secret tls envoy-client-cert --cert=certs/envoy.pem --key=certs/envoy-key.pem --dry-run=client -o yaml | kubectl apply -f -
kubectl -n projectcontour create secret tls untrusted-client --cert=certs/untrusted-client.pem --key=certs/untrusted-client-key.pem --dry-run=client -o yaml | kubectl apply -f -



go run github.com/projectcontour/contour/cmd/contour serve --xds-address=0.0.0.0 --xds-port=8001 --envoy-service-http-port=8080 --envoy-service-https-port=8443 --contour-cafile=ca.crt --contour-cert-file=tls.crt --contour-key-file=tls.key --config-path=$HOME/work/devenvs/contour/configs/contour-config-with-upstream-tls.yaml

kubectl apply -f manifests/echoserver-with-upstream-tls.yaml



# http
http --verify=certs/external-root-ca.pem https://protected.127-0-0-101.nip.io


# passthrough
http --verify=certs/internal-root-ca.pem --cert=certs/envoy.pem --cert-key=certs/envoy-key.pem https://passthrough.127-0-0-101.nip.io
http --verify=certs/internal-root-ca.pem --cert=certs/untrusted-client.pem --cert-key=certs/untrusted-client-key.pem https://passthrough.127-0-0-101.nip.io


# re-encrypt
http --verify=certs/external-root-ca.pem https://protected2.127-0-0-101.nip.io
http --verify=certs/external-root-ca.pem --cert=certs/envoy.pem --cert-key=certs/envoy-key.pem https://protected2.127-0-0-101.nip.io


## TLS client

ipython3

import ssl
import socket
context = ssl.SSLContext(ssl.PROTOCOL_TLS_CLIENT)
context.load_verify_locations('certs/internal-root-ca.pem')
context.load_cert_chain(certfile="certs/client-1.pem", keyfile="certs/client-1-key.pem")
sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
ssl_sock = context.wrap_socket(sock, server_hostname="passthrough.127-0-0-101.nip.io")
ssl_sock.connect(("passthrough.127-0-0-101.nip.io", 443))


ssl_sock.sendall(b"a")   # send a byte
ssl_sock.recv(1)         # receive a byte


import ssl
import socket
context = ssl.SSLContext(ssl.PROTOCOL_TLS_CLIENT)
context.load_verify_locations('certs/external-root-ca.pem')
context.load_cert_chain(certfile="certs/client-1.pem", keyfile="certs/client-1-key.pem")
sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
ssl_sock = context.wrap_socket(sock, server_hostname="protected2.127-0-0-101.nip.io")
ssl_sock.connect(("protected2.127-0-0-101.nip.io", 443))


## Capture and decrypt

cd ~/work/ingress-controller-conformance/images/echoserver

patch -p1 <<EOF
diff --git a/images/echoserver/echoserver.go b/images/echoserver/echoserver.go
index 60c0c98..63cf4ef 100644
--- a/images/echoserver/echoserver.go
+++ b/images/echoserver/echoserver.go
@@ -201,6 +201,11 @@ func processError(w http.ResponseWriter, err error, code int) {

 func listenAndServeTLS(addr string, serverCert string, serverPrivKey string, clientCA string, handler http.Handler) error {
        var config tls.Config
+       f, err := os.OpenFile("/tmp/wireshark-keys.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
+       if err != nil {
+               return fmt.Errorf("failed to open wireshark-keys.log: %v", err)
+       }
+       config.KeyLogWriter = f

        // Optionally enable client certificate validation when client CA certificates are given.
        if clientCA != "" {
EOF

docker buildx build --load --platform=linux/amd64 -t local/echoserver:0.0 .
kind load docker-image local/echoserver:0.0 --name contour


sudo nsenter -t $(pgrep -f "echoserver") --net wireshark -f "port 8443" -k -o tls.keylog_file:/proc/$(pgrep -f "echoserver")/root/tmp/wireshark-keys.log
