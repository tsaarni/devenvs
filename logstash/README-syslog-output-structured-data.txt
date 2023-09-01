


https://datatracker.ietf.org/doc/html/rfc5424
https://datatracker.ietf.org/doc/html/rfc5424#section-6         # 6. Syslog Message Format
https://datatracker.ietf.org/doc/html/rfc5424#section-6.3.5     # 6.3.5.  Examples
https://datatracker.ietf.org/doc/html/rfc5424#section-6.5       # 6.5.  Examples


# create certs
mkdir -p certs
certyaml -d certs configs/certs.yaml




export OSS=true
export LOGSTASH_SOURCE=1
export LOGSTASH_PATH=$HOME/work/logstash


# build gem
gem build logstash-output-syslog.gemspec


# in ~/work/logstash install the .gem
bin/logstash-plugin install --local ~/work/logstash-output-syslog/logstash-output-syslog-3.0.5.gem

# run logstash
bin/logstash -f ~/work/devenvs/logstash/configs/logstash-source-syslog-w-structured-data.conf --log.level debug


# config example

        #rfc => "rfc5424"
        #structured_data => [
        #    'exampleSDID@32473 iut="3" eventSource="Application" eventID="1011"',
        #    'examplePriority@32473 class="high"'
        #]




# run rsyslog
docker-compose rm -f
docker-compose up rsyslog



# run tests  (but see WORKAROUND below)
bundle exec rspec




*** Unit tests will fail

Errors such as:

  10) LogStash::Outputs::Syslog escape carriage return, newline and newline to \n behaves like syslog output should write expected format
      Failure/Error: @client_socket.write(syslog_msg + "\n")

        #<Double "fake socket"> received :write with unexpected arguments
          expected: (/^<0>(Jan|Feb|Mar|Apr|May|Jun|Jul|Aug|Sep|Oct|Nov|Dec) (0[1-9]|[12][0-9]|3[01]) ([01][0-9]|2[0-3]):([0-5][0-9]):([0-5][0-9]|60) baz LOGSTASH\[-\]: foo\\nbar\\nbaz\n/m)
               got: ("<0>Aug 03 07:41:56 baz LOGSTASH[-]: 2023-08-03T07:41:56.492198020Z baz bar\n")
        Diff:
        @@ -1 +1 @@
        -[/^<0>(Jan|Feb|Mar|Apr|May|Jun|Jul|Aug|Sep|Oct|Nov|Dec) (0[1-9]|[12][0-9]|3[01]) ([01][0-9]|2[0-3]):([0-5][0-9]):([0-5][0-9]|60) baz LOGSTASH\[-\]: foo\\nbar\\nbaz\n/m]



This is because https://github.com/logstash-plugins/logstash-output-syslog/issues/51

WORKAROUND:  To run the tests successfully, apply the change from this PR first

https://github.com/logstash-plugins/logstash-output-syslog/pull/55


patch -p1 <<EOF
diff --git a/lib/logstash/outputs/syslog.rb b/lib/logstash/outputs/syslog.rb
index 55197d9..add50c0 100644
--- a/lib/logstash/outputs/syslog.rb
+++ b/lib/logstash/outputs/syslog.rb
@@ -138,7 +138,7 @@ class LogStash::Outputs::Syslog < LogStash::Outputs::Base
       @ssl_context = setup_ssl
     end

-    if @codec.instance_of? LogStash::Codecs::Plain
+    if @codec.class.to_s == "LogStash::Codecs::Plain"
       if @codec.config["format"].nil?
         @codec = LogStash::Codecs::Plain.new({"format" => @message})
       end
EOF


