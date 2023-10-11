





To run logstash with custom-built jruby, change the path to jruby.jar


diff --git a/bin/logstash.lib.sh b/bin/logstash.lib.sh
index 3c4da9070..fe2b14a5b 100755
--- a/bin/logstash.lib.sh
+++ b/bin/logstash.lib.sh
@@ -87,7 +87,8 @@ setup_classpath() {
 }

 # set up CLASSPATH once (we start more than one java process)
-CLASSPATH="${JRUBY_HOME}/lib/jruby.jar"
+#CLASSPATH="${JRUBY_HOME}/lib/jruby.jar"
+CLASSPATH="/home/tsaarni/package/jruby/lib/jruby.jar"
 CLASSPATH="$(setup_classpath $CLASSPATH $LOGSTASH_JARS)"

 setup_java() {











diff --git a/lib/logstash/outputs/syslog.rb b/lib/logstash/outputs/syslog.rb
index d3ae6ca..18bbf2c 100644
--- a/lib/logstash/outputs/syslog.rb
+++ b/lib/logstash/outputs/syslog.rb
@@ -54,6 +54,10 @@ class LogStash::Outputs::Syslog < LogStash::Outputs::Base
     "debug",
   ]

+  # constants for Linux socket options, which are missing from jruby socket implementation.
+  IP_TOS = 1
+  IPV6_TCLASS = 64
+
   # syslog server address to connect to
   config :host, :validate => :string, :required => true

@@ -128,6 +132,10 @@ class LogStash::Outputs::Syslog < LogStash::Outputs::Base
   # RFC5424 structured data.
   config :structured_data, :validate => :string, :default => ""

+  # socket options for syslog socket.
+  #
+  config :socket_options, :validate => :array
+
   def register
     @client_socket = nil

@@ -229,6 +237,16 @@ class LogStash::Outputs::Syslog < LogStash::Outputs::Base
         end
       end
     end
+    if !@socket_options.empty?
+      require "socket"
+      if @socket_options.has_key?("tos")
+        puts "Setting TOS to #{@socket_options["tos"]}"
+        socket.setsockopt(Socket::IPPROTO_IP, IP_TOS, @socket_options["tos"])
+      end
+      if @socket_options.has_key?("traffic_class")
+        socket.setsockopt(Socket::IPPROTO_IPV6, IPV6_TCLASS, @socket_options["traffic_class"])
+      end
+    end
     socket
   end




Following jruby patch does not work 

a/core/src/main/java/org/jruby/ext/socket/RubyBasicSocket.java b/core/src/main/java/org/jruby/ext/socket/RubyBasicSocket.java
index 44a2d06c8e..a7348ab8b2 100644
--- a/core/src/main/java/org/jruby/ext/socket/RubyBasicSocket.java
+++ b/core/src/main/java/org/jruby/ext/socket/RubyBasicSocket.java
@@ -518,6 +519,16 @@ public class RubyBasicSocket extends RubyIO {
                     case IPPROTO_HOPOPTS: // these both have value 0 on several platforms
                         if (MulticastStateManager.IP_ADD_MEMBERSHIP == intOpt) {
                             joinMulticastGroup(val);
+                        } else {
+                            ByteBuffer buf = ByteBuffer.allocate(4);
+                            buf.order(ByteOrder.nativeOrder());
+                            //flipBuffer(buf.putInt(val.convertToInteger().getIntValue()));
+                            buf.putInt(64);
+                            System.err.println("setsockopt: level:" + intLevel + " opt:" + intOpt + " buf:" + buf.get(0) + buf.get(1) + buf.get(2) + buf.get(3));
+                            int ret = SOCKOPT.setsockopt(fd.realFileno, intLevel, intOpt, buf, buf.remaining());
+                            if (ret != 0) {
+                                throw runtime.newErrnoEINVALError(SOCKOPT.strerror(ret));
+                            }
                         }

                         break;
