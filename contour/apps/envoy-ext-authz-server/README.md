# envoy-ext-authz-server

Test server for Envoy `ext_authz` (gRPC and HTTP). The client request controls the verdict.
One process is one authz server. Run several processes to get several servers.

## Configuration

Each option is a command line flag and an environment variable (`AUTHZ_` + flag name in uppercase, `-`
replaced by `_`). The command line has priority. Run `--help` for the full list.

| flag | default | meaning |
|---|---|---|
| `--name` | `authz` | server name, used in logs and responses |
| `--grpc-addr`, `--http-addr` | `:9001`, `:9002` | listen addresses. An empty value disables the listener |
| `--grpc-tls-cert`, `--grpc-tls-key`, `--grpc-tls-ca` | | gRPC TLS. A CA enables mTLS |
| `--http-tls-cert`, `--http-tls-key`, `--http-tls-ca` | | HTTP TLS. A CA enables mTLS |
| `--default-status` | `200` | status when the client sends no `status` |
| `--quiet` | `false` | do not log individual requests, see [Logging](#logging) |
| `--metrics-addr` | | Prometheus listen address, serves `/metrics` without authentication. Empty disables, see [Metrics](#metrics) |

## Controlling the verdict

The client sets the verdict with up to four controls. Each control has two equivalent forms:

| control | header | `Authorization` parameter | values | effect |
|---|---|---|---|---|
| status | `x-authz-status: 403` | `status=403` | 200..599 | 200 allows. Any other value denies with that status |
| delay | `x-authz-delay: 2s` | `delay=2s` | Go duration | waits before answering. Tests `responseTimeout` |
| fail | `x-authz-fail: error` | `fail=error` | `error`, `reset` | simulates a server failure (see below) |
| tag | `x-authz-tag: test-42` | `tag=test-42` | max 128 printable characters, no spaces | label to tell requests apart, see [Tag](#tag) |

`Authorization` form: `Authorization: x-authz status=403 delay=2s`. Separate parameters with spaces or
commas. The scheme name is case-insensitive. Other schemes (`Bearer`, `Basic`) are ignored.

- Controls combine. Without controls, `--default-status` applies.
- If both forms set the same control, the header wins.
- An invalid value or unknown parameter is logged and denied with 400.

Responses contain `x-authz-server: <name>`, and `x-authz-tag` if the client sent a tag.
Allow adds `x-authz-verdict: allow`. Deny has the body `denied by <name> (<code>)`.

### Failures

| `fail` | gRPC | HTTP |
|---|---|---|
| `error` | returns `Unavailable`. Envoy uses `status_on_error` or `failOpen` | aborts the connection |
| `reset` | not supported (`Unimplemented`) | aborts the connection |

HTTP aborts the connection because Envoy treats any HTTP status as a denial, not an error.
To test an unreachable server, stop the process.

### Tag

The tag tells requests apart when needed, for example one value per test case or per client. It is part
of the label set of `authz_checks_total`, so every distinct tag creates new time series. Do not change
it for every request (a counter, a timestamp or a random value). With a huge number of requests the number of
series would explode and use all memory. Reuse a small set of tags and use `request_id` in the log to find a single request.

## Metrics

With `--metrics-addr :9100` the server exposes Prometheus metrics at `/metrics`. They are independent
of `--quiet`, so verdicts can be checked without per-request logs. The endpoint has no authentication,
bind it to localhost or the test network only. All `authz_*` metrics have the label `name=<server name>`.

| metric | labels | meaning |
|---|---|---|
| `authz_checks_total` | `protocol`, `verdict`, `code`, `tag`, `host`, `path`, `sni`, `client_cert` | checks handled. `verdict` has the same values as in the log. `code` is empty unless the verdict is `allow` or `deny`. `path` has no query string |
| `authz_controls_total` | `control`, `source` | controls found in checks. `source` is `header` or `authorization` |
| `authz_request_headers_total` | `header` | checks that carried the header (lowercase name, values are not recorded) |
| `authz_context_extensions_total` | `key`, `value` | gRPC context extensions |
| `authz_check_duration_seconds` | `protocol`, `verdict` | histogram of handling time, including `delay` |

Counters only increase and reset when the process restarts. A check is counted before its response is
sent, so a scrape right after a response includes it. Use a distinct `tag` per test case, or compare
values before and after the run. Every distinct combination of label values is a separate time series,
so memory grows if a load test uses a new `tag` or `path` for each request. Values that are only in the
log (`request_id`, header values, `body_bytes`) are not in metrics. Run without `--quiet` to see them.
Metrics show what the server decided. What Envoy did with it is in Envoy's `ext_authz.*` stats.

## Logging

Each request produces one `CHECK` line. For high request rates, set `--quiet` (`AUTHZ_QUIET=true`).
It disables all per-request logging: `CHECK` lines and the warnings about invalid controls or failed
writes. Startup messages and listener errors are still logged. Responses are unchanged.


```
CHECK name=strict tag=test-42 request_id=abc proto=http tls=true sni=... client_cert="CN=envoy" method=GET host=... path=/x ctx=env:staging body_bytes=0 controls="status=\"403\" delay=\"2s\" authorization=[status]" verdict=deny status=403 headers="..."
```

| field | meaning |
|---|---|
| `tag` | client tag. Use it to find the log lines of a test case |
| `request_id` | `x-request-id` from Envoy, if forwarded |
| `controls` | controls the server found in the request, as `name="value"` pairs. Only controls that are set are shown, including `tag`. `authorization=[...]` lists those read from the `Authorization` header. `none (default status)` means none were set. Invalid values are shown as received, with the reason in `error` |
| `verdict`, `status` | outcome and HTTP status. `verdict` is `allow`, `deny`, `error` or `reset`, or `canceled` (client went away during delay), `read-error` (HTTP body read failed), `unsupported-reset` (gRPC `fail=reset`). `status` is 0 unless allow/deny |
| `headers` | all request headers received by the server |
| `ctx` | gRPC context extensions (not available for HTTP) |
| `sni`, `client_cert` | TLS details |

## Header forwarding in Contour

| authz server | `x-authz-*` headers | `Authorization` |
|---|---|---|
| gRPC | forwarded | forwarded |
| HTTP, vhost `authorization` or global | only with `httpSettings.allowedAuthorizationHeaders: [{prefix: x-authz-}]` | forwarded |
| HTTP, `authzProviders` entry | never forwarded. Use `httpSettings.headersToAdd` or a gRPC provider | forwarded |

Use the `Authorization` form when unsure. To see the response headers at the upstream, add
`allowedUpstreamHeaders: [{prefix: x-authz-}]`.

The ExtensionService needs `protocol: h2c` for gRPC (or `h2` with TLS and `validation`), and
`protocol: http/1.1` for HTTP.

## Examples

```shell
make build
runagent run -n auth-default -- ./envoy-ext-authz-server --name default --grpc-addr :9001 --http-addr :9003
runagent run -n auth-strict -e AUTHZ_NAME=strict -e AUTHZ_GRPC_ADDR=:9002 -e AUTHZ_HTTP_ADDR=:9004 \
  -e AUTHZ_DEFAULT_STATUS=403 -- ./envoy-ext-authz-server

curl -si localhost:9003/                                         # allow
curl -si localhost:9003/ -H 'x-authz-status: 403'                # deny
curl -si localhost:9003/ -H 'Authorization: x-authz status=403'  # deny, same as above
curl -si localhost:9003/ -H 'x-authz-delay: 2s'                  # slow answer
curl -si localhost:9003/ -H 'x-authz-fail: error'                # connection aborted
curl -si localhost:9003/ -H 'x-authz-tag: test-42'                # find "tag=test-42" in the log
```

Contour virtual host that uses the server:

```yaml
virtualhost:
  fqdn: echoserver.127-0-0-101.nip.io
  authorization:
    extensionRef: {name: auth-strict}
    responseTimeout: 1s
    httpSettings:
      allowedAuthorizationHeaders: [{prefix: x-authz-}]
      allowedUpstreamHeaders: [{prefix: x-authz-}]
```

## Tests

```shell
go test ./...
```
