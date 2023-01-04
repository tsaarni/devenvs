# Using Fluent Bit as a wrapper for Envoy to format and forward logs

This example shows how to use Fluent Bit as a wrapper for Envoy to format and forward logs.
It can be useful when you want to convert application specific log format to a normalized format.


## Configuration

First create a container that includes both Envoy and Fluent Bit.

```bash
docker build docker/envoy-with-fluentbit/ -t envoy-with-fluentbit:latest
```

See the [Dockerfile](docker/envoy-with-fluentbit/Dockerfile) for details.
[Catatonit](https://github.com/openSUSE/catatonit) is used as entrypoint (PID 1).
It executes [entrypoint shell script](docker/envoy-with-fluentbit/files/docker-entrypoint-with-fluentbit.sh) that starts Envoy and Fluent Bit together.
`fluent-bit` is started in a subshell and `stdout` and `stderr` are forwarded to the Fluent Bit process.

```bash
exec > >(/opt/fluent-bit/bin/fluent-bit --quiet --config=/configs/fluentbit-envoy.conf ; kill -SIGTERM $$ ) 2>&1
exec /usr/local/bin/envoy --config-path /etc/envoy/envoy-httpbingo-config.yaml
```

Note that if `fluent-bit` exits e.g. due to crash, the script executes `kill` to force exit of the parent process.
Otherwise Envoy would continue to execute but we would lose the rest of the logs.
Envoy is started with `exec`, replacing the bash process which is executing the entrypoint script.

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

It uses [Standard Input](https://docs.fluentbit.io/manual/pipeline/inputs/standard-input) source and parses the logs using the parser defined in [`fluentbit-envoy-parsers.conf`](configs/fluentbit-envoy-parsers.conf).
[Standard output](https://docs.fluentbit.io/manual/pipeline/outputs/standard-output) sink is used to forward the logs to the standard output of the container in JSON format.
It could also stream the logs to a remote destination but for simplicity we are using the standard output only.
The example demonstrates optional processing of the HTTP access logs by using the [modify filter](https://docs.fluentbit.io/manual/pipeline/filters/modify).
It adds a new field `disconnect` if the HTTP access log event contains letters [`DC`](https://www.envoyproxy.io/docs/envoy/latest/configuration/observability/access_log/usage#config-access-log-format-response-flags).
The `DC` flag in Envoy HTTP access log indicates abrupt disconnection of the client before the HTTP response was delivered.

The parser configuration in [`fluentbit-envoy-parsers.conf`](configs/fluentbit-envoy-parsers.conf) is as follows.

```ini
[PARSER]
    name   envoy
    format regex
    regex  ^\[(?<timestamp>[^\]]*)\](\[(?<pid>\d+)\])?(\[(?<severity>\w+)\])?(?<msg>.*)$
```

The regular expression is written to parse log lines that look like following.

```
[2023-01-04 09:08:23.950][7][info][main] [source/server/server.cc:390] initializing epoch 0 (base id=0, hot restart version=11.104)
[2023-01-04T09:08:45.976Z] "GET /get HTTP/1.1" 200 - 0 569 225 225 "-" "curl/7.68.0" "69f185c9-7b1f-4508-af27-8b4a1c82f7c8" "httpbingo.org" "77.83.142.42:80"
[2023-01-04T09:08:51.304Z] "GET /delay/5 HTTP/1.1" 0 DC 0 0 372 - "-" "curl/7.68.0" "33e366d4-ef22-4d5c-bddc-21cb4469a20d" "httpbingo.org" "77.83.142.42:80"
```

The first log line is a debug log from Envoy.
The second and third log lines are HTTP access logs from Envoy.

The parser will split the logs into the following JSON fields:

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
docker run --volume $PWD/configs:/configs:ro -p 10000:10000 --rm --name envoy envoy-with-fluentbit:latest
```

Observe that Envoy debug logs printed to standard output of the container in JSON format:

```json
{"timestamp":"2023-01-04 10:41:02.783","pid":"7","severity":"info","msg":"[main] [source/server/server.cc:390] initializing epoch 0 (base id=0, hot restart version=11.104)"}
{"timestamp":"2023-01-04 10:41:02.783","pid":"7","severity":"info","msg":"[main] [source/server/server.cc:392] statically linked extensions:"}
{"timestamp":"2023-01-04 10:41:02.783","pid":"7","severity":"info","msg":"[main] [source/server/server.cc:394]   envoy.udp_packet_writer: envoy.udp_packet_writer.default, envoy.udp_packet_writer.gso"}
{"timestamp":"2023-01-04 10:41:02.783","pid":"7","severity":"info","msg":"[main] [source/server/server.cc:394]   envoy.matching.network.custom_matchers: envoy.matching.custom_matchers.trie_matcher"}
```

Make a request to Envoy proxy.
Envoy is exposed at port 10000.

```bash
curl http://localhost:10000/get
```

It will forward the request to [httpbingo.org](https://httpbingo.org/) `/get` endpoint, which simply echoes the `GET` HTTP request data back as response.
Observe that also HTTP access logs are parsed by Fluent Bit:

```json
{"timestamp":"2023-01-04T10:41:18.495Z","msg":" \"GET /get HTTP/1.1\" 200 - 0 570 234 230 \"-\" \"curl/7.68.0\" \"b71d1d1b-432d-42c2-83c6-275339c47eea\" \"httpbingo.org\" \"77.83.142.42:80\""}
```

Make a request to `/delay/5` endpoint.

```bash
curl http://localhost:10000/delay/5
```

The server will delay the response for 5 seconds.
Press CTRL+C to interrupt `curl` before 5 seconds has passed.
That will cause abrupt disconnection of the client connection before Envoy had the opportunity to deliver HTTP response.
Observe that the logs contains `DC` flag in `msg` field and that Fluent Bit added a new JSON field `disconnect: true`:

```json
{"timestamp":"2023-01-04T10:41:21.862Z","msg":" \"GET /delay/5 HTTP/1.1\" 0 DC 0 0 508 - \"-\" \"curl/7.68.0\" \"b17b02fe-c89a-44ba-8c8c-02646ed8a4fc\" \"httpbingo.org\" \"77.83.142.42:80\"","disconnect":"true"}
```

List the processes inside the container.

```bash
$ docker exec envoy ps -ef

UID          PID    PPID  C STIME TTY          TIME CMD
root           1       0  0 12:43 ?        00:00:00 /usr/libexec/catatonit/catatonit /docker-entrypoint-with-fluentbit.sh
root           7       1  1 12:43 ?        00:00:00 /usr/local/bin/envoy --config-path /etc/envoy/envoy-httpbingo-config.yaml
root           8       7  0 12:43 ?        00:00:00 bash /docker-entrypoint-with-fluentbit.sh
root           9       8  0 12:43 ?        00:00:00 /opt/fluent-bit/bin/fluent-bit --quiet --config=/configs/fluentbit-envoy.conf
```

As discussed, `catatonit` is the initial command run inside the container (PID 1).
Envoy has replaced the entrypoint script process which was `bash` before the `exec /usr/local/bin/envoy` command was called (PID 7).
There is a `bash` process (PID 8) for the subshell with `fluent-bit` running as a child (PID 9).
The reason for the subshell to remain running is to execute `kill -SIGTERM` on envoy process if `fluent-bit` (PID 9) would exit for any reason.


## Other thoughts

Sometimes more complicated logic is required.
For example, the timestamp format in Envoy debug log might need to be converted into "zulu time" or the severity field might need to be changed to uppercase letters.
For such cases, it is possible to use Fluent Bit [Lua filter](https://docs.fluentbit.io/manual/pipeline/filters/lua) to implement the logic.
