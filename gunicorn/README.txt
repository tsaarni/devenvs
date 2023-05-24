

cd ~/work/devenvs/gunicorn
python3 -m venv venv
. venv/bin/activate

mkdir -p certs
certyaml -d certs


cd ~/work/gunicorn

pip install -r requirements_test.txt
pip install -r requirements_dev.txt

python setup.py install
python setup.py test


cd ~/work/devenvs/gunicorn
. venv/bin/activate
gunicorn --worker-class=sync --workers=4 --certfile=certs/server.pem --keyfile=certs/server-key.pem myapp:app

http --verify=certs/ca.pem https://localhost:8000

gunicorn --worker-class=sync --workers=4 --certfile=certs/server.pem --keyfile=certs/server-key.pem --ca-certs=certs/ca.pem --cert-reqs 2 myapp:app
http --verify=certs/ca.pem --cert=certs/client.pem --cert-key=certs/client-key.pem  https://localhost:8000


gunicorn --worker-class=sync --workers=4 --certfile=certs/server.pem --keyfile=certs/server-key.pem --config configs/gunicorn-mintls-conf.py myapp:app


sslyze localhost:8000

openssl s_client -quiet -connect localhost:8000 -servername foo
openssl s_client -quiet -connect localhost:8000 -servername foo.127.0.0.1.nip.io


# start reading stream but block
import requests
r = requests.get("https://localhost:8000/stream", stream=True, verify="ca.pem")


# long-lived keep-alive connection
import requests, logging
logging.basicConfig(level=logging.DEBUG)
s = requests.Session()
s.get("https://localhost:8000", verify="ca.pem")




# regenerate .rst docs
#   note: it will overwrite some defaults, edit them back manually
make -C docs html
xdg-open docs/build/html/index.html

summary for timeout behavior

sync

  - connection is closed after idle for 30 seconds when
    - client connects but does not send HTTP request
    - client connects with HTTPS but does not proceed with TLS handshake

  - connection is closed immediately when
    - client sends single HTTP request with keep-alive

gthread

  - connection remains established forever when
    - client connects but does not send HTTP request
    - client connects with HTTPS but does not proceed with TLS handshake

  - connection is closed after idle for 2 seconds when
    - client sends single HTTP request with keep-alive

tornado

  - connection remains established even if:
    - client connects but does not send HTTP request
    - client sends single HTTP request with keep-alive and never sends further requests
    - client connects with HTTPS but does not proceed with TLS handshake


eventlet and gevent

  - connection is closed after idle for 2 seconds when
    - client connects but does not send HTTP request
    - client sends single HTTP request with keep-alive and never sends further requests
    - client connects with HTTPS but does not proceed with TLS handshake
