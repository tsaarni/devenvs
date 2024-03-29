# Run server that requires client to send client cert:
#   openssl s_server -accept 8443 -cert certs/server.pem -key certs/server-key.pem -CAfile certs/client-ca.pem -tls1_2 -Verify 1
#   openssl s_server -accept 8443 -cert certs/server-revoked.pem -key certs/server-revoked-key.pem -CAfile certs/client-ca.pem -tls1_2 -Verify 1
#
# Then run this script.
#

require "openssl"

ssl_context = OpenSSL::SSL::SSLContext.new

#ssl_context.cert = OpenSSL::X509::Certificate.new(File.read("/home/tsaarni/work/devenvs/logstash/certs/client-rsa.pem"))
#ssl_context.key = OpenSSL::PKey::read(File.read("/home/tsaarni/work/devenvs/logstash/certs/client-rsa-key.pem"))
ssl_context.cert = OpenSSL::X509::Certificate.new(File.read("/home/tsaarni/work/devenvs/logstash/certs/client.pem"))
ssl_context.key = OpenSSL::PKey::read(File.read("/home/tsaarni/work/devenvs/logstash/certs/client-key.pem"))

cert_store = OpenSSL::X509::Store.new
cert_store.add_file("/home/tsaarni/work/devenvs/logstash/certs/server-ca.pem")

# copy the behavior of X509_load_crl_file() which supports loading bundles of CRLs
END_TAG = "\n-----END X509 CRL-----\n"
File.read("/home/tsaarni/work/devenvs/logstash/certs/crl.pem").split(END_TAG).each do |crl|
    crl << END_TAG
    cert_store.add_crl(OpenSSL::X509::CRL.new(crl))
end

#cert_store.add_crl(OpenSSL::X509::CRL.new(File.read("/home/tsaarni/work/devenvs/logstash/certs/server-ca-crl.pem")))
#cert_store.add_crl(OpenSSL::X509::CRL.new(File.read("/home/tsaarni/work/devenvs/logstash/certs/client-ca-crl.pem")))
#cert_store.add_crl(OpenSSL::X509::CRL.new(File.read("/home/tsaarni/work/devenvs/logstash/certs/crl.pem"))) # bundle
cert_store.flags = OpenSSL::X509::V_FLAG_CRL_CHECK | OpenSSL::X509::V_FLAG_CRL_CHECK_ALL

ssl_context.cert_store = cert_store
ssl_context.verify_mode = OpenSSL::SSL::VERIFY_PEER|OpenSSL::SSL::VERIFY_FAIL_IF_NO_PEER_CERT
#ssl_context.verify_callback = -> (preverify_ok, ssl_context) {
#    cert_store.verify(ssl_context.current_cert)
#}

socket = TCPSocket.new("localhost", "8443")

socket = OpenSSL::SSL::SSLSocket.new(socket, ssl_context)

socket.connect
socket.puts "Hello World!"
socket.close
