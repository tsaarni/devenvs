
# import ssl
# import threading

# t = threading.local()
# t.context = None

# def ssl_context(conf):
#     if t.context is None:
#         t.context = ssl.SSLContext(conf.ssl_version)
#         t.context.load_cert_chain(certfile=conf.certfile, keyfile=conf.keyfile)
#         t.context.verify_mode = conf.cert_reqs
#     return t.context

# def ssl_context(conf):
#     import ssl

#     def set_defaults(context):
#         context.verify_mode = conf.cert_reqs
#         context.minimum_version = ssl.TLSVersion.TLSv1_3
#         if conf.ciphers:
#             context.set_ciphers(conf.ciphers)
#         if conf.ca_certs:
#             context.load_verify_locations(cafile=conf.ca_certs)

#     # Return different server certificate depending which hostname the client
#     # used to connect the server. Requires Python 3.7 or later.
#     def sni_callback(socket, server_hostname, context):
#         if server_hostname == "foo.127.0.0.1.nip.io":
#             new_context = ssl.SSLContext()
#             new_context.load_cert_chain(certfile="foo.pem", keyfile="foo-key.pem")
#             set_defaults(new_context)
#             socket.context = new_context

#     context = ssl.SSLContext(conf.ssl_version)
#     context.sni_callback = sni_callback
#     set_defaults(context)

#     # Load fallback certificate that will be returned when there is no match
#     # or client did not set TLS SNI (server_hostname == None)
#     context.load_cert_chain(certfile=conf.certfile, keyfile=conf.keyfile)

#     return context



# def ssl_context(conf, default_ssl_context_factory):
#     import ssl

#     # The default SSLContext returned by the factory function is initialized
#     # with the TLS parameters from config, including server certificate, private
#     # key and other parameters.
#     context = default_ssl_context_factory()

#     # The SSLContext can be further customized, for example to enforce minimum
#     # TLS version.
#     context.minimum_version = ssl.TLSVersion.TLSv1_3

#     # Server can also return different server certificate depending which
#     # hostname the client uses. Requires Python 3.7 or later.
#     def sni_callback(socket, server_hostname, context):
#         if server_hostname == "foo.127.0.0.1.nip.io":
#             new_context = default_ssl_context_factory()
#             new_context.load_cert_chain(certfile="foo.pem", keyfile="foo-key.pem")
#             socket.context = new_context

#     context.sni_callback = sni_callback

#     return context


# def ssl_context(conf, default_context_factory):
#     import ssl
#     context = default_context_factory()
#     context.minimum_version = ssl.TLSVersion.TLSv1_3
#     return context

