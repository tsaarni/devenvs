#!/bin/bash -e

mkdir -p /run/sshd
echo "NOTE: Wait for a while for sssd to start"
systemctl init ssh sssd
