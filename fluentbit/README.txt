


### docker run --volume=$PWD:/conf:ro -it fluent/fluent-bit:1.9.9-debug /bin/bash


docker run --volume=$PWD:/conf:ro --publish 127.0.0.1:5170:5170 -it fluent/fluent-bit:1.9.9-debug /fluent-bit/bin/fluent-bit --config /conf/fluent-bit.conf

echo '{ "message": "Hello world!" }' > /dev/tcp/localhost/5170
