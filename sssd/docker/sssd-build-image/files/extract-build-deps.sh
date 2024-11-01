#!/bin/sh

set -e

# Create a temporary directory to store the source package and extract it
mkdir -p /tmp/extract-build-deps
cd /tmp/extract-build-deps
tar xf $1

# Extract the build dependencies from the debian/control file
#
# Example:
#
# Build-Depends: dpkg-dev (>= 1.22.5),
#  autopoint,
#  bc,
#  check <!nocheck>,
#  cifs-utils,
#  debhelper-compat (= 13),
#  ...
#  xsltproc
# Standards-Version: 4.4.0

awk -f- debian/control <<'EOF'
function dep(pkg) {
    sub(/^[ ]+/, "", pkg)
    sub(/[ ,].*$/, "", pkg)
    print pkg
}
/^Build-Depends:/ { flag=1; sub(/Build-Depends:[ ]*/, ""); dep($0) }
flag && /^[A-Za-z-]+:/ { flag=0 }
flag && /^ / { dep($0) }
EOF

rm -rf /tmp/extract-build-deps
