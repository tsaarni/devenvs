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

## Controlling the verdict

The client sets the verdict with up to four controls. Each control has two equivalent forms:

| control | header | `Authorization` parameter | values | effect |
|---|---|---|---|---|
| status | `x-authz-status: 403` | `status=403` | 200..599 | 200 allows. Any other value denies with that status |
| delay | `x-authz-delay: 2s` | `delay=2s` | Go duration | waits before answering. Tests `responseTimeout` |
| fail | `x-authz-fail: error` | `fail=error` | `error`, `reset` | simulates a server failure (see below) |
| id | `x-authz-id: test-42` | `id=test-42` | max 128 printable characters, no spaces | correlation id, see [Logging](#logging) |

`Authorization` form: `Authorization: x-authz status=403 delay=2s`. Separate parameters with spaces or
commas. The scheme name is case-insensitive. Other schemes (`Bearer`, `Basic`) are ignored.

- Controls combine. Without controls, `--default-status` applies.
- If both forms set the same control, the header wins.
- An invalid value or unknown parameter is logged and denied with 400.

Responses contain `x-authz-server: <name>`, and `x-authz-id` if the client sent an id.
Allow adds `x-authz-verdict: allow`. Deny has the body `denied by <name> (<code>)`.

### Failures

| `fail` | gRPC | HTTP |
|---|---|---|
| `error` | returns `Unavailable`. Envoy uses `status_on_error` or `failOpen` | aborts the connection |
| `reset` | not supported (`Unimplemented`) | aborts the connection |

HTTP aborts the connection because Envoy treats any HTTP status as a denial, not an error.
To test an unreachable server, stop the process.

## Logging

Each request produces one `CHECK` line:

```
CHECK name=strict id=test-42 request_id=abc proto=http tls=true sni=... client_cert="CN=envoy" method=GET host=... path=/x ctx=env:staging body_bytes=0 controls="status=\"403\" delay=\"2s\" authorization=[status]" verdict=deny status=403 headers="..."
```

| field | meaning |
|---|---|
| `id` | client correlation id. Use it to find the log line of a request |
| `request_id` | `x-request-id` from Envoy, if forwarded |
| `controls` | controls the server found in the request, as `name="value"` pairs. Only controls that are set are shown, including `id`. `authorization=[...]` lists those read from the `Authorization` header. `none (default status)` means none were set. Invalid values are shown as received, with the reason in `error` |
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
curl -si localhost:9003/ -H 'x-authz-id: test-42'                # find "id=test-42" in the log
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
