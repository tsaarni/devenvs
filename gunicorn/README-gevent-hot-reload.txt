
# Create virtual environment
cd ~/work/devenvs/gunicorn
python3 -m venv venv


# Activate virtual environment
. venv/bin/activate


# Create certs
rm -rf certs
mkdir -p certs
certyaml -d certs


# Install gunicorn and dependencies to virtual environment
cd ~/work/gunicorn
pip install -r requirements_test.txt
pip install -r requirements_dev.txt
pip install --editable .


# Start gunicorn with gevent worker
cd ~/work/devenvs/gunicorn
gunicorn --worker-class=gevent --workers=1 --certfile=certs/server.pem --keyfile=certs/server-key.pem --config configs/gunicorn-hotreload-conf.py myapp:app


http --verify=certs/ca.pem https://localhost:8000

# Show server certificate and check the not before / not after dates
openssl s_client -connect localhost:8000 -showcerts </dev/null 2>/dev/null | openssl x509 -text


# Rotate certs
rm certs/server*.pem; certyaml -d certs



# Then check the server certificate again





### .vscode/launch.json

{
    // Use IntelliSense to learn about possible attributes.
    // Hover to view descriptions of existing attributes.
    // For more information, visit: https://go.microsoft.com/fwlink/?linkid=830387
    "version": "0.2.0",
    "configurations": [
        {
            "name": "Gunicorn",
            "type": "debugpy",
            "request": "launch",
            "justMyCode": false,
            "program": "${workspaceFolder}/venv/bin/gunicorn",
            "args": [
                "myapp:app",
                "--workers=1",
                "--keyfile=certs/server-key.pem",
                "--certfile=certs/server.pem",
                "--config=configs/gunicorn-hotreload-conf.py",
            ],
            "console": "integratedTerminal",
        }
    ]
}
