
# https://github.com/elastic/logstash/issues/16173

# Logstash 8.11.3's distribution from Elastic is bundled with Jruby 9.4.5.0 and Adoptium's JDK 17.0.9p9, but since you have built from source there are a number of additional variables at play.



./gradlew installDevelopmentGems

export OSS=false    # https://github.com/elastic/logstash/issues/16025

bin/logstash -f ~/work/devenvs/logstash/configs/logstash-destination-tcp.conf --config.reload.automatic

# or with debug
bin/logstash -f ~/work/devenvs/logstash/configs/logstash-destination-tcp.conf --log.level debug --config.reload.automatic



# NOTE: this does not reproduce the fault!!! new reproduction steps are needed
# Test sequence

## 1. test that tcp input in port 5000 is working
echo "Hello, World!" > /dev/tcp/localhost/5000

## 2. Run this in another terminal to keep port 5001 allocated
nc -l 5001

## 3. Change the logstash tcp input from 5000 to 5001 and wait for logstash to reload config and fail binding to 5001
sed -i 's/5000/5001/g' ~/work/devenvs/logstash/configs/logstash-destination-tcp.conf

## 4. Change the config back
sed -i 's/5001/5000/g' ~/work/devenvs/logstash/configs/logstash-destination-tcp.conf

## 5. Check that logstash recovered and is again successfully listening on port 5000
echo "Hello, World!" > /dev/tcp/localhost/5000
