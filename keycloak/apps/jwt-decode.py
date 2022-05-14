#!/bin/env python3
#
# Decodes JWT from stdin
#
# Install dependencies on ubuntu
#   apt-get install python3-jwt
#
# or with virtual env
#   python3 -m venv venv
#   . venv/bin/activate
#   pip install jwt￼
#
import jwt
import time
import sys
import pprint

t = jwt.decode(sys.stdin.read().rstrip(), verify=False)
pprint.pprint(t)

if t.get('iat') is not None:
    print("# iat={}".format(time.ctime(t['iat'])))
if t.get('exp') is not None:
    print("# exp={}".format(time.ctime(t['exp'])))
