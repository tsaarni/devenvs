def ssl_context(conf, default_ssl_context_factory):
    import ssl
    context = default_ssl_context_factory()
    def sni_callback(socket, server_hostname, context):
        new_context = default_ssl_context_factory()
        new_context.load_cert_chain(conf.certfile, conf.keyfile)
        socket.context = new_context

    context.sni_callback = sni_callback
    return context
