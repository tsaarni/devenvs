#!/bin/env python3

import paramiko

client = paramiko.SSHClient()
client.set_missing_host_key_policy(paramiko.AutoAddPolicy())

client.connect('localhost', 12222, username='expires', password="expires", look_for_keys=False, allow_agent=False)

# https://docs.paramiko.org/en/latest/api/transport.html#paramiko.transport.Transport.get_banner
auth_banner = client.get_transport().get_banner()

print(auth_banner.decode('utf-8'))

client.close()
