#!/bin/bash -ex

# process templates
eval "cat <<< \"$(<database.ldif)\"" > /data/ldif/database.ldif
eval "cat <<< \"$(<users-and-groups.ldif)\"" > /data/ldif/users-and-groups.ldif

# import processed templates
slapadd -v -n 0 -l /data/ldif/database.ldif -F /data/config
slapadd -v -n 1 -l /data/ldif/users-and-groups.ldif -F /data/config
