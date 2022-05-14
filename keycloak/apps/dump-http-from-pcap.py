#!/bin/env python3
#
# Dumps HTTP requests and responses in textual format from .pcap file
#
# Usage
#   dump-http-from-pcap.py traffic.cap
#
# To install dependencies into virtual env
#   python3 -m venv venv
#   . venv/bin/activate
#   pip install scapy
#

import sys
import datetime
import codecs
from scapy.all import *
from scapy.layers.http import *

# Generate unique steam id for (full duplex) TCP stream
def stream_id(p):
    return str(sorted([p[IP].src, p[TCP].sport, p[IP].dst, p[TCP].dport], key=str))

def main(capture_filename):
    #cap = sniff(offline=capture_filename, session=TCPSession, filter='tcp')
    # with filter it throws:
    #    AttributeError: 'PcapNgReader' object has no attribute 'linktype'
    # in scapy/utils.py:2138 (scape 2.4.5)
    cap = sniff(offline=capture_filename, session=TCPSession)

    # Separate the unique TCP streams out from the packets in the capture file.
    tcp_streams = cap.sessions(stream_id)

    # Sort the streams by the timestamp of the first packet.
    sorted_sessions = sorted(tcp_streams.values(), key=lambda plist: plist[0].time)

    for pkgs in sorted_sessions:
        for p in pkgs:
            # Process only packets that have HTTP layer.
            if HTTP in p:
                # Print header for each request with some metadata.
                if HTTPRequest in p:
                    zt = datetime.utcfromtimestamp(float(p.time)).strftime('%Y-%m-%dT%H:%M:%S.%fZ')
                    print(f'''
################################################################################
#
# {zt} {p[IP].src}:{p[TCP].sport} -> {p[IP].dst}:{p[TCP].dport}
#
################################################################################
''')
                if HTTPResponse in p:
                    print('\n# Response:')
                    print(p.summary())
                print(codecs.decode(bytes(p[HTTP].payload), 'utf-8', errors='ignore'))

if __name__ == '__main__':
    main(sys.argv[1])
