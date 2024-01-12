


# This script modifies the envoy build filesto allow building the postgres proxy filter from contrib.

# Add the postgres proxy filter to the envoy build
# - see https://github.com/envoyproxy/envoy/blob/main/bazel/README.md#customize-extension-build-config
#
# Notes: 
#
# - Instead of creating new Bazel workspace like instructed in README we modify the existing one.
#
#   Rationale: Own workspace should contain copy of source/extensions/extensions_build_config.bzl which
#   we would need to maintain. We can just modify the existing one by using script.
#
# - Add the postgres filter to the list of extensions to build in //source/exe:envoy-static target
#   instead of using //contrib/exe:envoy-static.
#
#   Reationale: It is easier to add single contrib filter to normal build than to remove all other 
#   contrib filters from contrib build.
#   This is also documented as supported in the Bazel README:
#   "no need to specifically perform a contrib build to include a contrib extension"
#
sed -i '/EXTENSIONS = {/a \    "envoy.filters.network.postgres_proxy": "//contrib/postgres_proxy/filters/network/source:config",' source/extensions/extensions_build_config.bzl

# Change the visibility to public to allow the filter to compile
# - see https://github.com/envoyproxy/envoy/blob/main/bazel/README.md#extra-extensions
sed -i 's|^CONTRIB_EXTENSION_PACKAGE_VISIBILITY .*|CONTRIB_EXTENSION_PACKAGE_VISIBILITY = ["//visibility:public"]|' source/extensions/extensions_build_config.bzl

