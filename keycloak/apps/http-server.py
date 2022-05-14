#!/bin/env python3
#
# Web server for receiving authorization code from Keycloak
# as part of authorization code flow
#

from http.server import BaseHTTPRequestHandler, HTTPServer

server_address = ('', 8000)

class Handler(BaseHTTPRequestHandler):
    def do_GET(self):
        self.send_response(200)
        self.send_header("Content-type", "text/html")
        self.end_headers()

s = HTTPServer(server_address, Handler)
print('Listening: {}:{}'.format(server_address[0], str(server_address[1])))
s.serve_forever()
