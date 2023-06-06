#!/bin/env python3
from http.server import BaseHTTPRequestHandler, HTTPServer

class MyServer(BaseHTTPRequestHandler):
    def do_GET(self):
        self.send_response(200)
        self.send_header('Content-type', 'text/html')
        self.send_header('Set-Cookie', 'sessionid=secret')
        self.send_header('Set-Cookie', 'another=secret')
        self.end_headers()
        self.wfile.write(b'Hello, World!')

port = 8081
httpd = HTTPServer(('localhost', port), MyServer)
print(f'Starting httpd on port {port}...')
httpd.serve_forever()
