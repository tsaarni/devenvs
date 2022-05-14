
kubectl apply -f manifests/shell.yaml
kubectl exec -it shell -- ash



cat >server.py <<EOF

import asyncio
import websockets

# create handler for each connection
async def handler(websocket, path):
    counter = 1
    while True:
        data = f'Hello world {counter}!'
        print(f'Sending: {data}')
        await websocket.send(data)
        await asyncio.sleep(1)
        counter = counter + 1

start_server = websockets.serve(handler, '0.0.0.0', 8000)
asyncio.get_event_loop().run_until_complete(start_server)
asyncio.get_event_loop().run_forever()
EOF

python3 server.py



# Client

sudo apt-get install python3-websockets -y


cat >client.py <<EOF
import asyncio
import websockets


async def consumer_handler():
    async with websockets.connect('ws://shell.127-0-0-101.nip.io') as websocket:
        async for data in websocket:
            print(f'Received: {data}')

asyncio.get_event_loop().run_until_complete(consumer_handler())

EOF

python3 client.py



sudo nsenter -t $(pgrep -f "python3 server.py") --net wireshark -f "port 8000" -k



kubectl -n projectcontour port-forward daemonset/envoy 9001

http http://localhost:9001/config_dump?include_eds | jq -C .configs[].dynamic_active_clusters

