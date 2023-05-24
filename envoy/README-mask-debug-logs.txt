
https://github.com/envoyproxy/envoy/issues/9652

bazel build -c dbg //source/exe:envoy-static


# test envoy
python3 -m http.server --bind 127.0.0.1 8081
bazel-bin/source/exe/envoy-static -c ~/work/devenvs/envoy/configs/envoy-static-virtualhost.yaml --log-level debug
bazel-bin/source/exe/envoy-static -c ~/work/devenvs/envoy/configs/envoy-debug-redact-disabled.yaml --log-level debug

http http://127.0.0.1:8080/
http "http://127.0.0.1:8080/foo?secret" "Cookie:sessionid=secret;another=secret"
http -a joe:password "http://127.0.0.1:8080/foo?secret" "Cookie:sessionid=secret;another=secret"





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



[2023-05-24 12:39:48.735][73394][debug][http] [source/common/http/conn_manager_impl.cc:1069] [C0][S5294506544605446236] request headers complete (end_stream=true):
':authority', '127.0.0.1:8080'
':path', '/foo?secret'
':method', 'GET'
'user-agent', 'HTTPie/2.6.0'
'accept-encoding', 'gzip, deflate'
'accept', '*/*'
'connection', 'keep-alive'
'cookie', 'sessionid=secret;another=secret'
'authorization', 'Basic am9lOnBhc3N3b3Jk'

[2023-05-24 12:39:48.735][73394][debug][http] [source/common/http/conn_manager_impl.cc:1052] [C0][S5294506544605446236] request end stream
[2023-05-24 12:39:48.735][73394][debug][connection] [./source/common/network/connection_impl.h:98] [C0] current connecting state: false
[2023-05-24 12:39:48.735][73394][debug][router] [source/common/router/router.cc:478] [C0][S5294506544605446236] cluster 'mycluster' match for URL '/foo?secret'
[2023-05-24 12:39:48.735][73394][debug][router] [source/common/router/router.cc:691] [C0][S5294506544605446236] router decoding headers:
':authority', '127.0.0.1:8080'
':path', '/foo?secret'
':method', 'GET'
':scheme', 'http'
'user-agent', 'HTTPie/2.6.0'
'accept-encoding', 'gzip, deflate'
'accept', '*/*'
'cookie', 'sessionid=secret;another=secret'
'authorization', 'Basic am9lOnBhc3N3b3Jk'
'x-forwarded-proto', 'http'
'x-request-id', 'b210cc8f-82ea-4132-9f2b-4d41e69f999c'
'x-envoy-expected-rq-timeout-ms', '15000'


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
