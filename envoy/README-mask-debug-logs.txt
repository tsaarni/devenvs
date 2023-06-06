
https://github.com/envoyproxy/envoy/issues/9652

bazel build -c dbg //source/exe:envoy-static


# test envoy
python3 app/server-that-sets-cookies.py
bazel-bin/source/exe/envoy-static -c ~/work/devenvs/envoy/configs/envoy-debug-redact-enabled.yaml --log-level debug
bazel-bin/source/exe/envoy-static -c ~/work/devenvs/envoy/configs/envoy-static-virtualhost.yaml --log-level debug

http http://127.0.0.1:8080/
http "http://127.0.0.1:8080/foo?secret" "Cookie: sessionid=secret; another=secret"
http -a joe:password "http://127.0.0.1:8080/foo?secret" "Cookie: sessionid=secret; another=secret"





# To run under debugger, create .vscode/launch.json entry for running
mkdir -p .vscode
cat > .vscode/launch.json <<EOF
{
  "version": "0.2.0",
  "configurations": [
    {
      "name": "Run envoy under gdb",
      "request": "launch",
      "type": "gdb",
      "target": "\${workspaceFolder}/bazel-bin/source/exe/envoy-static",
      "arguments": "-c /home/tsaarni/work/devenvs/envoy/configs/envoy-static-virtualhost.yaml --log-level debug",
      "cwd": "${workspaceFolder}",
      "valuesFormatting": "disabled"
    }
  ]
}
EOF





Headers are printed here
source/common/http/header_map_impl.cc:444
   #  void HeaderMapImpl::dumpState()




source/common/http/header_utility.cc


[2023-06-05 11:25:52.331][218006][debug][http] [source/common/http/conn_manager_impl.cc:1069] [C0][S12657990578490477947] request headers complete (end_stream=true):
':authority', '127.0.0.1:8080'
':path', '/foo?secret'
':method', 'GET'
'user-agent', 'HTTPie/2.6.0'
'accept-encoding', 'gzip, deflate'
'accept', '*/*'
'connection', 'keep-alive'
'cookie', 'sessionid=secret; another=secret'
'authorization', 'Basic am9lOnBhc3N3b3Jk'

[2023-06-05 11:25:52.331][218006][debug][router] [source/common/router/router.cc:690] [C0][S12657990578490477947] router decoding headers:
':authority', '127.0.0.1:8080'
':path', '/foo?secret'
':method', 'GET'
':scheme', 'http'
'user-agent', 'HTTPie/2.6.0'
'accept-encoding', 'gzip, deflate'
'accept', '*/*'
'cookie', 'sessionid=secret; another=secret'
'authorization', 'Basic am9lOnBhc3N3b3Jk'
'x-forwarded-proto', 'http'
'x-request-id', '28e945f2-50f3-4020-a7a2-588c747fcdcc'
'x-envoy-expected-rq-timeout-ms', '15000'

[2023-06-05 11:25:52.332][218006][debug][http] [source/common/http/conn_manager_impl.cc:1700] [C0][S12657990578490477947] encoding headers via codec (end_stream=false):
':status', '200'
'server', 'envoy'
'date', 'Mon, 05 Jun 2023 08:25:52 GMT'
'content-type', 'text/html'
'set-cookie', 'sessionid=secret'
'set-cookie', 'another=secret'
'x-envoy-upstream-service-time', '0'


-------

[2023-05-24 12:40:04.361][73529][debug][http] [source/common/http/conn_manager_impl.cc:1069] [C0][S11179102941975859115] request headers complete (end_stream=true):
':authority', '127.0.0.1:8080'
':path', '/foo?[redacted]'
':method', 'GET'
'user-agent', 'HTTPie/2.6.0'
'accept-encoding', 'gzip, deflate'
'accept', '*/*'
'connection', 'keep-alive'
'cookie', 'sessionid=[redacted];another=[redacted]'
'authorization', '[redacted]'

[2023-05-24 12:40:04.361][73529][debug][http] [source/common/http/conn_manager_impl.cc:1052] [C0][S11179102941975859115] request end stream
[2023-05-24 12:40:04.361][73529][debug][connection] [./source/common/network/connection_impl.h:98] [C0] current connecting state: false
[2023-05-24 12:40:04.361][73529][debug][router] [source/common/router/router.cc:478] [C0][S11179102941975859115] cluster 'mycluster' match for URL '/foo?secret'
[2023-05-24 12:40:04.361][73529][debug][router] [source/common/router/router.cc:691] [C0][S11179102941975859115] router decoding headers:
':authority', '127.0.0.1:8080'
':path', '/foo?[redacted]'
':method', 'GET'
':scheme', 'http'
'user-agent', 'HTTPie/2.6.0'
'accept-encoding', 'gzip, deflate'
'accept', '*/*'
'cookie', 'sessionid=[redacted];another=[redacted]'
'authorization', '[redacted]'
'x-forwarded-proto', 'http'
'x-request-id', 'a142296b-b4b3-415e-8b9c-5e67cc3e70b7'
'x-envoy-expected-rq-timeout-ms', '15000'







# related issues
# - protobuf
#     https://github.com/envoyproxy/envoy/issues/4757
#     https://github.com/envoyproxy/envoy/pull/9315
# - config dump
#     https://github.com/envoyproxy/envoy/pull/7365
# - accesslog
#     https://github.com/envoyproxy/envoy/issues/7583
#     https://github.com/envoyproxy/envoy/pull/15711



While doing troubleshooting it might be necessary to collect, distribute and store debug level logs.
It introduces a risk of unintended data leakage, where sensitive information is revealed accidentally as part of the log files to people who should not have access to that data.
Most typically, this would happen through authorization header, cookies and query string value, dumped in the log files.

See example of a log entry revealing sensitive information:

':authority', '127.0.0.1:8080'
':path', '/foo?secret'
':method', 'GET'
'user-agent', 'HTTPie/2.6.0'
'accept-encoding', 'gzip, deflate'
'accept', '*/*'
'connection', 'keep-alive'
'cookie', 'sessionid=secret;another=secret'
'authorization', 'Basic am9lOnBhc3N3b3Jk

A new feature is proposed in this issue. Sensitive information is masked by adding [redacted] in place of potentially sensitive values:

':authority', '127.0.0.1:8080'
':path', '/foo?[redacted]'
':method', 'GET'
'user-agent', 'HTTPie/2.6.0'
'accept-encoding', 'gzip, deflate'
'accept', '*/*'
'connection', 'keep-alive'
'cookie', 'sessionid=[redacted];another=[redacted]'
'authorization', '[redacted]'



'x-vault-token', 'hvs.Cjlceb2odgPPYOXJ3IjiJwDZ'
