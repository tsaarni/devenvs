def app(environ, start_response):
    if environ["PATH_INFO"] == "/stream":
        start_response("200 OK", [
            ("Content-Type", "application/octet-stream")
        ])
        f = open("/dev/urandom", "rb")
        return f
    else:
        data = b"Hello, World!\n"
        start_response("200 OK", [
            ("Content-Type", "text/plain"),
            ("Content-Length", str(len(data)))
        ])
        return iter([data])

