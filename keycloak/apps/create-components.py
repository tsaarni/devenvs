#!/bin/env python3
#
# Takes a list of JSON files containing REST requests to create components in Keycloak.
#
# Usage: create-components.py --server <server> --user <user> --password <password> <rest-request-file>...
#
# Example:
#   create-components.py rest-requests/create-ldap-simple-auth-provider.json
#

import getopt
import logging
import os
import sys

import requests
import json

# Default values
user = "admin"
password = "admin"
server = os.getenv("KEYCLOAK_SERVER_URL", "http://localhost:8080")

# enable info logging
logging.basicConfig(level=logging.INFO)

# Parse command line arguments
long_options = ["help", "server=", "user=", "password="]

try:
    opts, args = getopt.getopt(sys.argv[1:], "", long_options)
except getopt.GetoptError as err:
    print(err)
    sys.exit(1)

for opt, arg in opts:
    if opt == "--help":
        print(
            f"Usage: {sys.argv[0]} --server <server> --user <user> --password <password> <rest-request-file>..."
        )
        sys.exit(0)
    if opt == "--server":
        server = arg
    if opt == "--user":
        user = arg
    if opt == "--password":
        password = arg

# Remaining arguments are REST request to create component(s)
request_files = args


# Authenticate as admin
logging.info(f"Authenticating as {user} on {server}")
res = requests.post(
    f"{server}/realms/master/protocol/openid-connect/token",
    data={
        "username": user,
        "password": password,
        "grant_type": "password",
        "client_id": "admin-cli",
    },
    verify=False,
)
res.raise_for_status()

token = res.json()["access_token"]

# Create components
for f in request_files:
    logging.info(f"Creating component from {f}")
    with open(f, "r") as fd:
        res = requests.post(
            f"{server}/admin/realms/master/components",
            headers={
                "Authorization": f"Bearer {token}",
            },
            json=json.load(fd),
            verify=False,
        )
        res.raise_for_status()
        if res.ok:
            res = requests.get(
                f"{res.headers['Location']}",
                headers={"Authorization": f"Bearer {token}"},
                verify=False,
            )
            res.raise_for_status()
            json = res.json()
            logging.info(
                f"Successfully created: providerType={json['providerType']} id={json['id']}, name={json['name']}"
            )
        else:
            logging.error(f"Error {res.status_code}")
            logging.error(res.text)
