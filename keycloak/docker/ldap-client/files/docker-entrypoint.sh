#!/bin/bash -ex
chmod 600 /etc/sssd/sssd.conf
mkdir -p /run/sshd
systemctl init ssh sssd
