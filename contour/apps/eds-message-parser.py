#!/bin/env python3
#
# This script reads a JSON file containing EDS messages and generates a table of DiscoveryRequest
# and DiscoveryResponse messages.
#
# The JSON file is generated by https://github.com/tsaarni/grpc-json-sniffer
#
# Run the script to generate a table of EDS messages
#
#    apps/grpc-eds-sequencer.py grpc_capture.json
#    apps/grpc-eds-sequencer.py grpc_capture.json <min_message_id>  # analyze messages with message_id >= min_message_id
#

import json
import sys

class EdsMessageParser:

    def __init__(self):
        self.rows = []

    def set_min_message_id(self, min_message_id):
        self.min_message_id = int(min_message_id)

    def process_discovery_request(self, data, message_id, stream_id, content):
        if data.get("error") is not None:
            if "context canceled" in data.get("error"):
                self.rows.append([message_id, "DiscoveryRequest", stream_id, "<context canceled>", "", "<context canceled>f", ""])
                return
            raise ValueError(f"Unknown error in DiscoveryRequest: {content.get('error')}")
        version_info = content.get("versionInfo")
        resource_name = ", ".join(content.get("resourceNames"))
        response_nonce = content.get("responseNonce")
        self.rows.append([message_id, "DiscoveryRequest", stream_id, version_info, response_nonce, resource_name, ""])

    def process_discovery_response(self, data, message_id, stream_id, content):
        version_info = content.get("versionInfo")
        nonce = content.get("nonce")
        resources = content.get("resources")

        if len(resources) == 0:
            self.rows.append([message_id, "DiscoveryResponse", stream_id, version_info, nonce, "", ""])
        elif len(resources) == 1:
            cluster_name = resources[0].get("clusterName")
            endpoints = resources[0].get("endpoints")
            if len(endpoints) == 0:
                self.rows.append([message_id, "DiscoveryResponse", stream_id, version_info, nonce, cluster_name, ""])
                return
            lb_endpoints = endpoints[0].get("lbEndpoints")
            addresses = [
                lb_endpoint.get("endpoint").get("address").get("socketAddress").get("address")
                for lb_endpoint in lb_endpoints
            ]
            self.rows.append([message_id, "DiscoveryResponse", stream_id, version_info, nonce, cluster_name, ", ".join(addresses)])
        else:
            raise ValueError("Cannot process multiple resources in DiscoveryResponse")

    def process(self, line):
        data = json.loads(line)
        message_id = data.get("message_id")

        if message_id < self.min_message_id:
            return

        stream_id = data.get("stream_id")
        content = data.get("content", {})

        if data.get("method") == "/envoy.service.endpoint.v3.EndpointDiscoveryService/StreamEndpoints":
            if data.get("message") == "envoy.service.discovery.v3.DiscoveryRequest":
                self.process_discovery_request(data, message_id, stream_id, content)
            elif data.get("message") == "envoy.service.discovery.v3.DiscoveryResponse":
                self.process_discovery_response(data, message_id, stream_id, content)

    def print_table(self):
        headers = ["Id", "Message", "Stream ID", "Version Info", "Nonce", "Resource Name", "Addresses"]
        column_widths = [max(len(str(row[i])) for row in [headers] + self.rows) for i in range(len(headers))]

        def format_row(row):
            return "   ".join(str(row[i]).ljust(column_widths[i]) for i in range(len(row)))

        print(format_row(headers))
        print("- -".join("-" * width for width in column_widths))
        for row in self.rows:
            print(format_row(row))

def main():

    min_message_id = 0
    if len(sys.argv) > 2:
        min_message_id = sys.argv[2]

    table = EdsMessageParser()
    table.set_min_message_id(min_message_id)

    with open(sys.argv[1], "r") as f:
        try:
            for line in f:
                table.process(line)
        except Exception:
            print(f"Error while processing line:\n{line}")
            import traceback
            traceback.print_exc()
            print("Exiting")
            exit(1)


    table.print_table()


if __name__ == "__main__":
    main()
