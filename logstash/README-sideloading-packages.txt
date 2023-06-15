
### Preparations

# compile jruby-openssl
cd ~/work/jruby-openssl
./mvnw package -Dmaven.test.skip=true        # Compiled .gem will be in:  pkg/*.gem

# you can set version in: lib/jopenssl/version.rb e.g.
  VERSION = '0.14.2'


# compile logstash-output-syslog
cd ~/work/logstash-output-syslog
gem build logstash-output-syslog.gemspec     # *.gem will be in the top directory


# generate test certs and CRLs
#    note: see https://github.com/tsaarni/certyaml
#          go to releases page for pre-compiled binaries
cd ~/work/devenvs/logstash
mkdir -p certs
certyaml -d certs/ configs/certs.yaml



### Running fixed jruby-openssl inside official upstream logstash OSS container


# on host environment, copy (pre-compiled) gems to make them available in the container
cd ~/work/logstash
mkdir -p tmp
cp ~/work/jruby-openssl/pkg/jruby-openssl-0.14.2.gem tmp/
cp ~/work/logstash-output-syslog/logstash-output-syslog-3.0.5.gem tmp/

# also copy test certs to make them available in the container
cp ~/work/jruby-openssl/src/test/ruby/x509/ec-ca.crt tmp/
cp ~/work/jruby-openssl/src/test/ruby/x509/ec-ca.crl tmp/


# run upstream logstash container, mount logstash repo and certs into the environment
docker run --rm -it --user=$(id -u):$(id -g) --volume=/home/tsaarni/work/logstash:/source:ro --volume=/home/tsaarni/work/devenvs/logstash/certs:/certs:ro --network=host docker.elastic.co/logstash/logstash-oss:8.8.1 /bin/bash


# inside container, install fixed jruby-openssl and logstash-output-syslog

export PATH=$PATH:/usr/share/logstash/jdk/bin/:/usr/share/logstash/vendor/jruby/bin/

gem install --install-dir /usr/share/logstash/vendor/bundle/jruby/2.6.0/ /source/tmp/jruby-openssl-0.14.2.gem
logstash-plugin install --local /source/tmp/logstash-output-syslog-3.0.5.gem

# inspect the versions of jruby-openssl
find /usr/share/logstash -name *jruby-openssl* -type f
gem list
gem info jruby-openssl


# For logstash to pick up new version of jruby-openssl, check that the version in Gemfile matches with the installed version
#   Note: use a hack to edit the file from outside the container, because there is no editor inside the container
vi /proc/<PID_OF_PROCESS_INSIDE_CONTAINER>/root/usr/share/logstash/Gemfile.lock



### Test CRL validation via ruby interactive shell

# run ruby interactive shell
irb

# paste following to irb shell (should run without errors)
require "openssl"
crl = OpenSSL::X509::CRL.new(File.read(File.expand_path('/source/tmp/ec-ca.crl', __FILE__)))
ca = OpenSSL::X509::Certificate.new(File.read(File.expand_path('/source/tmp/ec-ca.crt', __FILE__)))
crl.verify(ca.public_key)




### Test Logstash syslog output with TLS and CRL validation

# inside container, create config file
cat <<EOF > logstash.conf
input {
    stdin {
        id => "myapp"
        codec => "line"
    }
}

output {
    syslog {
        host => "localhost"

        port => 6515
        protocol => "ssl-tcp"
        ssl_cacert => "/certs/server-ca.pem"
        ssl_cert => "/certs/client.pem"
        ssl_key => "/certs/client-key.pem"
        ssl_verify => true
        ssl_crl => "/certs/server-ca-crl.pem"
        ssl_crl_check_all => true
    }
}
EOF


# run logstash
bin/logstash -f logstash.conf

# wait for logstash to start
# and write something to the console (i.e. logstash stdin) to trigger log even to syslog output

# to start simulated syslog server on host, run openssl s_server (good certs)
#    Note: observe on openssl output that logstash will send logs to the server
openssl s_server -accept 6515 -cert certs/server.pem -key certs/server-key.pem -CAfile certs/client-ca.pem

# or run the server with revoked certs
openssl s_server -accept 6515 -cert certs/server-revoked.pem -key certs/server-revoked-key.pem -CAfile certs/client-ca.pem


# Note that for logstash inside container to be able to connect to openssl s_server (on host), you need to run docker run with "--network=host"




#### Random notes


# logstash build command
./gradlew assembleOssTarDistribution

# output artifacts will be generated under build/ dir like this:
build/logstash-oss-8.9.0-SNAPSHOT-linux-x86_64.tar.gz

# the dockerfile for the official logstash container is generated from this template:
docker/templates/Dockerfile.j2

# the source files for generated Gemfile
Gemfile.template
rakelib/plugins-metadata.json




# To see default gem load paths inside container
#   e.g. /usr/share/logstash/.local/share/gem/jruby/2.6.0:/usr/share/logstash/vendor/jruby/lib/ruby/gems/shared
export PATH=$PATH:/usr/share/logstash/jdk/bin/:/usr/share/logstash/vendor/jruby/bin/
gem env
