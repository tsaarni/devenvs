require 'socket'

host = 'localhost'
port = 8000

begin
  socket = TCPSocket.new(host, port)
  socket.puts "GET / HTTP/1.1"
  socket.puts "Host: #{host}"
  socket.puts ""
  socket.flush

  while line = socket.gets
    puts line
  end

  socket.close
rescue => e
  warn e
end
