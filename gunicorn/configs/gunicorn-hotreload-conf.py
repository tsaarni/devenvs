

def ssl_context(conf, default_ssl_context_factory):
    import os

    # Store the current mtime of the certificate and key files and load them once by calling the default_ssl_context_factory.
    cert_mtime = os.path.getmtime(conf.certfile)
    key_mtime = os.path.getmtime(conf.keyfile)
    context = default_ssl_context_factory()

    def sni_callback(socket, server_hostname, context):
        nonlocal cert_mtime, key_mtime
        try:
            # Check the mtime of the certificate and key files and reload them if they have changed.
            new_cert_mtime = os.path.getmtime(conf.certfile)
            new_key_mtime = os.path.getmtime(conf.keyfile)
            if cert_mtime != new_cert_mtime or key_mtime != new_key_mtime:

                # Reload the certificate and key files by calling the default_ssl_context_factory again.
                # This will throw if we attempt read in the middle of rotation.
                new_context = default_ssl_context_factory()

                cert_mtime, key_mtime = new_cert_mtime, new_key_mtime
                socket.context = new_context

        # Catch exceptions and keep the old context.
        except Exception as e:
            print(f"Error loading certfile or keyfile: {e}")
            return

    context.sni_callback = sni_callback
    return context
