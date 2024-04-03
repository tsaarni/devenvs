
https://github.com/elastic/logstash/issues/16025

OSS environment variable should not be used at all?



Error


[2024-03-28T14:03:21,002][WARN ][logstash.config.source.multilocal] Ignoring the 'pipelines.yml' file because modules or command line options are specified
[2024-03-28T14:03:21,189][INFO ][logstash.geoipdatabasemanagement.manager] database manager is disabled; removing managed databases from disk``
[2024-03-28T14:03:21,192][FATAL][logstash.runner          ] An unexpected error occurred! {:error=>#<Errno::ENOENT: No such file or directory - /home/tsaarni/work/logstash/data/geoip_database_management>, :backtrace=>["org/jruby/RubyDir.java:147:in `initialize'", "org/jruby/RubyClass.java:917:in `new'", "org/jruby/RubyDir.java:492:in `children'", "org/jruby/RubyDir.java:487:in `children'", "/home/tsaarni/work/logstash/x-pack/lib/geoip_database_management/manager.rb:146:in `clean_up_database'", "/home/tsaarni/work/logstash/x-pack/lib/geoip_database_management/manager.rb:64:in `initialize'", "org/jruby/RubyClass.java:897:in `new'", "/home/tsaarni/work/logstash/vendor/jruby/lib/ruby/stdlib/singleton.rb:127:in `block in instance'", "org/jruby/ext/thread/Mutex.java:171:in `synchronize'", "/home/tsaarni/work/logstash/vendor/jruby/lib/ruby/stdlib/singleton.rb:125:in `instance'", "/home/tsaarni/work/logstash/logstash-core/lib/logstash/agent.rb:633:in `initialize_geoip_database_metrics'", "/home/tsaarni/work/logstash/logstash-core/lib/logstash/agent.rb:97:in `initialize'", "org/jruby/RubyClass.java:931:in `new'", "/home/tsaarni/work/logstash/logstash-core/lib/logstash/runner.rb:541:in `create_agent'", "/home/tsaarni/work/logstash/logstash-core/lib/logstash/runner.rb:423:in `execute'", "/home/tsaarni/work/logstash/vendor/bundle/jruby/3.1.0/gems/clamp-1.0.1/lib/clamp/command.rb:68:in `run'", "/home/tsaarni/work/logstash/logstash-core/lib/logstash/runner.rb:288:in `run'", "/home/tsaarni/work/logstash/vendor/bundle/jruby/3.1.0/gems/clamp-1.0.1/lib/clamp/command.rb:133:in `run'", "/home/tsaarni/work/logstash/lib/bootstrap/environment.rb:89:in `<main>'"]}



Potential fix (if OSS should have been used)

diff --git a/logstash-core/lib/logstash/agent.rb b/logstash-core/lib/logstash/agent.rb
index 98b16f371..2e22e82fc 100644
--- a/logstash-core/lib/logstash/agent.rb
+++ b/logstash-core/lib/logstash/agent.rb
@@ -94,7 +94,9 @@ class LogStash::Agent
     initialize_agent_metrics
     initialize_flow_metrics

-    initialize_geoip_database_metrics(metric)
+    unless LogStash::OSS
+      initialize_geoip_database_metrics(metric)
+    end

     @pq_config_validator = LogStash::PersistedQueueConfigValidator.new



