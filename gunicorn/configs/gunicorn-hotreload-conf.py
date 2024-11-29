def ssl_context(conf, default_ssl_context_factory):
    import ssl
    context = default_ssl_context_factory()
    def sni_callback(conn, domain, context):
        context.load_cert_chain(conf.certfile, conf.keyfile)
    context.sni_callback = sni_callback
    return context
