


go run testserver/main.go



http --stream testserver.127-0-0-101.nip.io


ipython3

import logging
import requests

logging.basicConfig(level=logging.DEBUG)

s = requests.Session()
s.get("http://testserver.127-0-0-101.nip.io")
