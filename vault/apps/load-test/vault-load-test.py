#!/usr/bin/env python3

import argparse
import asyncio
import json
import logging
import os
import ssl
import time
import urllib.parse

logging.basicConfig(level=logging.INFO, format="%(asctime)s %(levelname)s: %(message)s")
logger = logging.getLogger("vault-load-test")

class KubernetesLoginTester:

    def __init__(self, vault_addr: str, concurrent_requests: int, max_failed_logins: int, role: str, token: str, ca_data: str|None) -> None:
        self.url = urllib.parse.urlparse(vault_addr)

        # Number of concurrent login requests to make.
        self.num_workers = concurrent_requests

        # When the number of failed logins reaches this value, the test will stop.
        self.max_failed_logins = max_failed_logins

        self.ca_data = ca_data

        # Lock to protect the counters.
        self.lock = asyncio.Lock()
        self.failed_logins = 0
        self.successful_logins = 0

        # Prepare the request to be sent.
        request_body = json.dumps({
            "role": role,
            "jwt": token
        })

        request = (
            f"POST /v1/auth/kubernetes/login HTTP/1.1\r\n"
            f"Host: {self.url.hostname}\r\n"
            "Content-Type: application/json\r\n"
            f"Content-Length: {len(request_body)}\r\n"
            "\r\n"
            f"{request_body}\r\n"
            "\r\n"
        )

        logger.debug(f"Request: {request}")
        self.encoded_request = request.encode()



    # Create tasks (num_workers) that will send requests concurrently.
    async def start(self) -> None:
        logger.info(f"Starting logins with {self.num_workers} workers towards {self.url.geturl()}")

        # Create task that will report the number of successful and failed logins periodically.
        asyncio.create_task(self.status_reporter())

        self.start_time = time.time()

        # Create a number of workers that will send login requests to Vault concurrently.
        await asyncio.gather(*[self.worker(i) for i in range(self.num_workers)])


    def print_stats(self) -> None:
        end_time = time.time()
        elapsed_time = end_time - self.start_time
        logger.info(f"Test completed in {elapsed_time:.2f} seconds")
        logger.info(f"Successful logins: {self.successful_logins}, Failed logins: {self.failed_logins}")
        logger.info(f"Average successful logins per second: {self.successful_logins // elapsed_time}")
        logger.info(f"Average failed logins per second: {self.failed_logins // elapsed_time}")
        if self.successful_logins > 0:
            logger.info(f"Average time per successful login: {elapsed_time / self.successful_logins * 1000:.2f} ms")
        if self.failed_logins > 0:
            logger.info(f"Average time per failed login: {elapsed_time / self.failed_logins * 1000:.2f} ms")

    # Reporter function that prints the number of successful and failed logins periodically.
    async def status_reporter(self) -> None:
        logger.debug("Starting reporter")
        start_time = asyncio.get_event_loop().time()
        while True:
            await asyncio.sleep(1)
            async with self.lock:
                logger.info(f"Successful logins: {self.successful_logins}, Failed logins: {self.failed_logins}, Runtime: {asyncio.get_event_loop().time() - start_time:.0f} seconds")
                if self.failed_logins >= self.max_failed_logins:
                    break


    # Worker function that sends login requests to Vault until the number of failed logins reaches the maximum.
    async def worker(self, worker_id: int) -> None:
        ssl_context = None
        if self.url.scheme == "https" and self.ca_data is not None:
            ssl_context = ssl.create_default_context(cadata=self.ca_data)
        elif self.url.scheme == "https":
            ssl_context = ssl.create_default_context()

        while True:
            result = await self.login(ssl_context)
            async with self.lock:
                if result:
                    self.successful_logins += 1
                else:
                    self.failed_logins += 1
                    if self.failed_logins >= self.max_failed_logins:
                        logger.debug(f"Worker {worker_id} reached maximum of {self.max_failed_logins} failed logins.")
                        break

    async def login(self, ssl_context: ssl.SSLContext|None) -> bool:
        logger.debug(f"Sending request to {self.url.geturl()}")

        sock = asyncio.open_connection(self.url.hostname, self.url.port, ssl=ssl_context)
        reader, writer = await sock

        writer.write(self.encoded_request)
        await writer.drain()

        # Read full response from the server.
        buffer = await reader.read(4096)
        response = buffer.decode()
        logger.debug(response)

        writer.close()
        await writer.wait_closed()

        if response.startswith("HTTP/1.1 200 OK"):
            return True
        else:
            return False

if __name__ == "__main__":

    parser = argparse.ArgumentParser(description="Vault load test")
    parser.add_argument("--concurrency", type=int, help="Number of concurrent requests", default=500)
    parser.add_argument("--max-failures", type=int, help="Maximum number of failed logins", default=1000)
    parser.add_argument("--vault-addr", type=str, help="Vault address", default=os.getenv("VAULT_ADDR", "http://localhost:8200"))
    parser.add_argument("--service-account", type=argparse.FileType('r'), help="Kubernetes service account path", default="/var/run/secrets/kubernetes.io/serviceaccount/token")
    parser.add_argument("--ca-cert", type=argparse.FileType('r'), help="CA certificate path")
    parser.add_argument("--role", type=str, help="Kubernetes role", required=True)
    parser.add_argument("--verbose", action="store_true", help="Enable verbose logging")
    args = parser.parse_args()

    ca_data = None
    if args.ca_cert:
        ca_data = args.ca_cert.read()

    if args.service_account:
        service_account = args.service_account.read()

    if args.verbose:
        logger.setLevel(logging.DEBUG)

    tester = KubernetesLoginTester(args.vault_addr, args.concurrency, args.max_failures, args.role, service_account, ca_data)
    try:
        asyncio.run(tester.start())
    except KeyboardInterrupt:
        logger.info("Received KeyboardInterrupt, stopping test")

    tester.print_stats()
