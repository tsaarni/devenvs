
kubectl apply -f examples/contour/01-crds.yaml 


kubectl apply -f manifests/shell.yaml
kubectl exec -it shell -- ash







cat >server.py <<EOF
from http.server import BaseHTTPRequestHandler, ThreadingHTTPServer


class HttpHandler(BaseHTTPRequestHandler):
    protocol_version = 'HTTP/1.1'
    #timeout = 10
    timeout = 999999

    def do_GET(self):
        response = bytes("{'response': 'Hello world!'}", 'UTF-8')
        self.send_response(200)
        self.send_header('Content-type', 'application/json')
        self.send_header("Content-Length", str(len(response)))
        self.end_headers()
        self.wfile.write(response)


def run(server_class=ThreadingHTTPServer, handler_class=HttpHandler):
    httpd = server_class(('', 8000), handler_class)
    httpd.serve_forever()


if __name__ == '__main__':
    run()

EOF



python3 server.py


sudo nsenter -t $(pgrep -f "python3 server.py") --net wireshark -f "port 8000" -k



http http://shell.127-0-0-101.nip.io



python3

import requests
s = requests.Session()
s.get('http://shell.127-0-0-101.nip.io')
