# Run server that requires client to send client cert:
#   openssl s_server -accept 8443 -cert certs/server.pem -key certs/server-key.pem -tls1_2 -Verify 1
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
ssl_context.cert_store = cert_store
ssl_context.verify_mode = OpenSSL::SSL::VERIFY_PEER|OpenSSL::SSL::VERIFY_FAIL_IF_NO_PEER_CERT

socket = TCPSocket.new("localhost", "8443")

socket = OpenSSL::SSL::SSLSocket.new(socket, ssl_context)

socket.connect
