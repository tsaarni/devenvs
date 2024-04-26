

https://github.com/envoyproxy/envoy/issues/32032


./ci/run_envoy_docker.sh './ci/do_ci.sh release.server_only'
￼
Target //distribution/binary:release up-to-date:
  bazel-bin/distribution/binary/release.tar.zst
INFO: Elapsed time: 3157.172s, Critical Path: 600.64s
INFO: 6234 processes: 36 internal, 1 local, 6197 processwrapper-sandbox.
INFO: Build completed successfully, 6234 total actions
INFO: Analyzed target //test/tools/schema_validator:schema_validator_tool.stripped (6 packages loaded, 108 targets configured).
INFO: Found 1 target...
Target //test/tools/schema_validator:schema_validator_tool.stripped up-to-date:
  bazel-bin/test/tools/schema_validator/schema_validator_tool.stripped
INFO: Elapsed time: 26.067s, Critical Path: 24.78s
INFO: 18 processes: 1 internal, 17 processwrapper-sandbox.
INFO: Build completed successfully, 18 total actions
INFO: Analyzed target //test/tools/router_check:router_check_tool.stripped (24 packages loaded, 185 targets configured).
INFO: Found 1 target...
Target //test/tools/router_check:router_check_tool.stripped up-to-date:
  bazel-bin/test/tools/router_check/router_check_tool.stripped
INFO: Elapsed time: 153.377s, Critical Path: 94.20s
INFO: 64 processes: 1 internal, 63 processwrapper-sandbox.
INFO: Build completed successfully, 64 total actions
Release files created in /build/envoy/x64/bin



ENVOY_DOCKER_IN_DOCKER=1 ./ci/run_envoy_docker.sh './ci/do_ci.sh docker'



*** Problem: failed to solve: failed to compute cache key: failed to calculate checksum of ref ...


docker buildx build --progress=plain  --no-cache --platform linux/amd64 -f ci/Dockerfile-envoy --target envoy-tools --sbom=false --provenance=false --load -t envoyproxy/envoy-tools-dev:2f02d47f64a7cb4106d79e6b87ac10353f5a2cb6


Dockerfile-envoy:56
--------------------
  55 |     ENV TARGETPLATFORM="${TARGETPLATFORM:-linux/amd64}"
  56 | >>> COPY --chown=0:0 --chmod=755 \
  57 | >>>     "${TARGETPLATFORM}/schema_validator_tool" "${TARGETPLATFORM}/router_check_tool" /usr/local/bin/
  58 |
--------------------
ERROR: failed to solve: failed to compute cache key: failed to calculate checksum of ref FE7O:47LL:7CIN:VMVY:QH3S:HTOU:3UXQ:35DF:6B7Z:Z7ZO:A7Z6:JJ7B::814m88oanhq92zeei0z1teall: "/linux/amd64/router_check_tool": not found



Solution:

check .dockerignore and make sure the copied file is not being ignored
