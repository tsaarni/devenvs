

docker run --rm -it \
    --env USER_UID=$(id -u) \
    --env USER_GID=$(id -g) \
    --volume "$(pwd)":/envoy \
    --workdir /envoy \
    ubuntu:25.04 /bin/bash

# Install the compilation dependencies.
export DEBIAN_FRONTEND=noninteractive
apt-get update
apt-get install -y --no-install-recommends build-essential python3 curl ca-certificates git


# Check that gold is NOT installed
gold --version

-bash: gold: command not found

curl --location --output /usr/local/bin/bazel https://github.com/bazelbuild/bazelisk/releases/latest/download/bazelisk-linux-amd64
chmod +x /usr/local/bin/bazel



userdel ubuntu
groupadd --gid $USER_GID envoybuild
useradd --uid $USER_UID --gid $USER_GID --create-home --shell /bin/bash envoybuild


su - envoybuild


cd /envoy
cat > user.bazelrc <<EOF
build --config=gcc
build --local_resources=cpu=HOST_CPUS-2
build --local_resources=memory=HOST_RAM*.5
EOF



bazel build //source/exe:envoy-static







diff --git a/bazel/rbe/toolchains/configs/linux/gcc/cc/cc_toolchain_config.bzl b/bazel/rbe/toolchains/configs/linux/gcc/cc/cc_toolchain_config.bzl
index e65754720c..06c4624dc9 100755
--- a/bazel/rbe/toolchains/configs/linux/gcc/cc/cc_toolchain_config.bzl
+++ b/bazel/rbe/toolchains/configs/linux/gcc/cc/cc_toolchain_config.bzl
@@ -820,13 +820,6 @@ def _impl(ctx):
                     flag_group(
                         iterate_over = "libraries_to_link",
                         flag_groups = [
-                            flag_group(
-                                flags = ["-Wl,--start-lib"],
-                                expand_if_equal = variable_with_value(
-                                    name = "libraries_to_link.type",
-                                    value = "object_file_group",
-                                ),
-                            ),
                             flag_group(
                                 flags = ["-Wl,-whole-archive"],
                                 expand_if_true =
@@ -879,13 +872,6 @@ def _impl(ctx):
                                 flags = ["-Wl,-no-whole-archive"],
                                 expand_if_true = "libraries_to_link.is_whole_archive",
                             ),
-                            flag_group(
-                                flags = ["-Wl,--end-lib"],
-                                expand_if_equal = variable_with_value(
-                                    name = "libraries_to_link.type",
-                                    value = "object_file_group",
-                                ),
-                            ),
                         ],
                         expand_if_available = "libraries_to_link",
                     ),
