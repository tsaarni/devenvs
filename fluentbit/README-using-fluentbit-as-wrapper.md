# Using Fluent Bit as a wrapper for Envoy to format and forward logs

This example shows how to use Fluent Bit as a wrapper for Envoy to format and forward logs.
It can be useful when you want to convert application specific log format to a normalized format.


## Configuration

First create a container that includes both Envoy and Fluent Bit.
See the [Dockerfile](docker/envoy-with-fluentbit/Dockerfile) for details.

```bash
docker build docker/envoy-with-fluentbit/ -t envoy-with-fluentbit:latest
```

The container has an [entrypoint script](docker/envoy-with-fluentbit/files/docker-entrypoint-with-fluentbit.sh) that starts Envoy and Fluent Bit together.
It will start `fluentbit` in a subshell and set up forwarding for `stdout` and `stderr` to the Fluent Bit process.

```bash
exec > >(/opt/fluent-bit/bin/fluent-bit --quiet --config=/configs/fluentbit-envoy.conf ; kill -SIGTERM $$ ) 2>&1
exec /usr/libexec/catatonit/catatonit -g /usr/local/bin/envoy -- --config-path /etc/envoy/envoy-httpbingo-config.yaml
```

Note that if `fluentbit` exits e.g. due to crash, the script executes `kill` to force exit of the parent process.
Otherwise Envoy would continue to execute but we would lose rest of the logs.

[Catatonit](https://github.com/openSUSE/catatonit) will execute `envoy` and handle signals.
It is started with `exec` to make `catatonit` process the new PID 1, replacing the bash process which is executing the script.

Fluent Bit configuration files are mounted on `/configs` directory.
The main configuration file [`fluentbit-envoy.conf`](configs/fluentbit-envoy.conf) contains the following.

```ini
[SERVICE]
    Parsers_File /configs/fluentbit-envoy-parsers.conf
    Log_Level off

[INPUT]
    Name stdin
    Parser envoy

[FILTER]
    Name modify
    match *
    condition key_value_matches msg \sDC\s
    add disconnect true

[OUTPUT]
    Name stdout
    Format json_lines
    json_date_key false
```

It uses [Standard Input](https://docs.fluentbit.io/manual/pipeline/inputs/standard-input) as source and parses the logs using the parser defined in [`fluentbit-envoy-parsers.conf`](configs/fluentbit-envoy-parsers.conf).
[Standard output](https://docs.fluentbit.io/manual/pipeline/outputs/standard-output) sink is used to forward the logs to the standard output of the container in JSON format.
The example also demonstrates optional processing of the HTTP access logs by using the [modify filter](https://docs.fluentbit.io/manual/pipeline/filters/modify).
It adds a new field `disconnect` if the HTTP access log event contains [`DC`](https://www.envoyproxy.io/docs/envoy/latest/configuration/observability/access_log/usage#config-access-log-format-response-flags).
The `DC` flag indicates abrupt disconnection of the client before the HTTP response was delivered.

The parser configuration in [`fluentbit-envoy-parsers.conf`](configs/fluentbit-envoy-parsers.conf) is as follows.

```ini
[PARSER]
    name   envoy
    format regex
    regex  ^\[(?<timestamp>[^\]]*)\](\[(?<pid>\d+)\])?(\[(?<severity>\w+)\])?(?<msg>.*)$
```

The regular expression is written to parse following log lines.

```
[2023-01-04 09:08:23.950][8][info][main] [source/server/server.cc:390] initializing epoch 0 (base id=0, hot restart version=11.104)
[2023-01-04T09:08:45.976Z] "GET /get HTTP/1.1" 200 - 0 569 225 225 "-" "HTTPie/1.0.3" "69f185c9-7b1f-4508-af27-8b4a1c82f7c8" "httpbingo.org" "77.83.142.42:80"
[2023-01-04T09:08:51.304Z] "GET /delay/5 HTTP/1.1" 0 DC 0 0 372 - "-" "HTTPie/1.0.3" "33e366d4-ef22-4d5c-bddc-21cb4469a20d" "httpbingo.org" "77.83.142.42:80"
```

It will split the logs into the following JSON fields:

* `timestamp`
* `pid`
* `severity`
* `msg`


In addition to the Fluent Bit configuration, the container includes Envoy configuration
See the [envoy-httpbingo-config.yaml](docker/envoy-with-fluentbit/files/etc/envoy/envoy-httpbingo-config.yaml) for details.
With the configuration, Envoy proxies for [httpbingo.org](https://httpbingo.org/) - an online service that provides a set of HTTP endpoints for testing HTTP clients.


## Testing

Run Envoy and Fluent Bit together.

```bash
docker run --volume $PWD/configs:/configs:ro -p 10000:10000 -it --rm --name envoy envoy-with-fluentbit:latest
```

Observe that Envoy debug logs printed to standard output of the container in JSON format:

```json
{"timestamp":"2023-01-04 10:41:02.783","pid":"10","severity":"info","msg":"[main] [source/server/server.cc:390] initializing epoch 0 (base id=0, hot restart version=11.104)"}
{"timestamp":"2023-01-04 10:41:02.783","pid":"10","severity":"info","msg":"[main] [source/server/server.cc:392] statically linked extensions:"}
{"timestamp":"2023-01-04 10:41:02.783","pid":"10","severity":"info","msg":"[main] [source/server/server.cc:394]   envoy.udp_packet_writer: envoy.udp_packet_writer.default, envoy.udp_packet_writer.gso"}
{"timestamp":"2023-01-04 10:41:02.783","pid":"10","severity":"info","msg":"[main] [source/server/server.cc:394]   envoy.matching.network.custom_matchers: envoy.matching.custom_matchers.trie_matcher"}
```

Make a request to Envoy to proxy the request to httpbingo.org `/get` endpoint, which simply returns the GET data.
Envoy is exposed at port 10000.


```bash
curl http://localhost:10000/get
```

Observe that also HTTP access logs are parsed by Fluent Bit:

```json
{"timestamp":"2023-01-04T10:41:18.495Z","msg":" \"GET /get HTTP/1.1\" 200 - 0 570 234 230 \"-\" \"HTTPie/1.0.3\" \"b71d1d1b-432d-42c2-83c6-275339c47eea\" \"httpbingo.org\" \"77.83.142.42:80\""}
```

Make a request to `/delay/5` endpoint, which delays the response for 5 seconds.
Press CTRL+C to trigger client disconnect before the response is received.

```bash
curl http://localhost:10000/delay/5
```

Observe that the logs contains `DC` flag in `msg` field and that Fluent Bit added a new field `disconnect`:

```json
{"timestamp":"2023-01-04T10:41:21.862Z","msg":" \"GET /delay/5 HTTP/1.1\" 0 DC 0 0 508 - \"-\" \"HTTPie/1.0.3\" \"b17b02fe-c89a-44ba-8c8c-02646ed8a4fc\" \"httpbingo.org\" \"77.83.142.42:80\"","disconnect":"true"}
```

List the processes inside the container.

```bash
$ docker exec envoy ps -ef

UID          PID    PPID  C STIME TTY          TIME CMD
root           1       0  0 11:04 pts/0    00:00:00 /bin/sh -c /docker-entrypoint-with-fluent.sh
root           7       1  0 11:04 pts/0    00:00:00 /usr/libexec/catatonit/catatonit -g /usr/local/bin/envoy -- --config-path /etc/envoy/envoy-httpbingo-config.yaml
root           8       7  0 11:04 pts/0    00:00:00 bash /docker-entrypoint-with-fluent.sh
root           9       8  0 11:04 pts/0    00:00:00 /opt/fluent-bit/bin/fluent-bit --quiet --config=/configs/fluentbit-envoy.conf
root          10       7  0 11:04 pts/0    00:00:01 /usr/local/bin/envoy --config-path /etc/envoy/envoy-httpbingo-config.yaml
root          70       0  0 11:06 ?        00:00:00 ps -ef
```
