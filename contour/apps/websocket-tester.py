#!/usr/bin/env python

import asyncio
import websockets


async def consumer_handler():
    async with websockets.connect('ws://shell.127-0-0-101.nip.io') as websocket:
        async for data in websocket:
            print(f'Received: {data}')

asyncio.get_event_loop().run_until_complete(consumer_handler())
