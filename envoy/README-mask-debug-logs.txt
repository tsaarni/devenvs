bazel build -c dbg //source/exe:envoy-static


# test envoy
python3 -m http.server --bind 127.0.0.1 8081
bazel-bin/source/exe/envoy-static -c ~/work/devenvs/envoy/configs/envoy-static-virtualhost.yaml --log-level debug
http http://127.0.0.1:8080/
http "http://127.0.0.1:8080/foo?secret"






# To run under debugger, create vscode launch.json entry for running

    {
      "name": "Run envoy under gdb",
      "request": "launch",
      "type": "gdb",
      "target": "${workspaceFolder}/bazel-bin/source/exe/envoy-static",
      "arguments": "-c /home/tsaarni/work/devenvs/envoy/configs/envoy-static-virtualhost.yaml --log-level debug",
      "cwd": "${workspaceFolder}",
      "valuesFormatting": "disabled"
    }
