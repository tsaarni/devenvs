from http.server import BaseHTTPRequestHandler, ThreadingHTTPServer


class HttpHandler(BaseHTTPRequestHandler):
    protocol_version = 'HTTP/1.1'

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
