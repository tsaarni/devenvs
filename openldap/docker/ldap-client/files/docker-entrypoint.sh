#!/bin/bash -ex

mkdir -p /run/sshd
systemctl init ssh sssd
