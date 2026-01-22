#!/usr/bin/env python3
#
# Copyright 2024 tero.saarni@est.tech
#
# Licensed to the Apache Software Foundation (ASF) under one or more contributor license agreements; and to You under the Apache License, Version 2.0.
#
# jsonlogparse.py - Parse and format JSON log files into human-readable format.
#
# Usage:
#  python jsonlogparse.py [--nocolor] [--outfile OUTFILE] [INFILE]
#
# By default it reads from STDIN and writes to STDOUT.
#
# Parses following fields from JSON log file:
# - timestamp
# - severity
# - message
#
# Optionally colorizes output based on severity and highlights keywords in message.
#

import json
import sys
import argparse
import re
import codecs

# Custom JSON decoder to handle invalid escape sequences
def raw_decode_with_invalid_escapes(s, *args, **kwargs):
    try:
        return json.loads(s, *args, **kwargs)
    except json.JSONDecodeError as e:
        if 'Invalid \\escape' in str(e):
            # Replace the problematic string with a valid JSON string
            # by escaping all backslashes
            s = s.replace('\\', '\\\\')
            # But fix double escaping of valid escape sequences
            for esc in ['\\"', '\\/', '\\b', '\\f', '\\n', '\\r', '\\t']:
                s = s.replace('\\\\' + esc[1:], esc)
            return json.loads(s, *args, **kwargs)
        raise

SEVERITY_COLORS = {
    "DEBUG": "\033[0;37m",
    "INFO": "\033[0;32m",
    "WARNING": "\033[0;33m",
    "ERROR": "\033[0;31m",
    "CRITICAL": "\033[0;31m"
}

# highlighted keywords in message
KEYWORDS = ["error", "warning", "timeout", "timed out", "failed", "failure", "invalid", "exception", "caused by", "stacktrace"]

def highlight_keyword(match):
    return f"\033[0;31m{match.group(0)}\033[0m"

def format_log_line(line, colorize):
    # handle escapes
    try:
        jsondoc = raw_decode_with_invalid_escapes(line)
        timestamp = jsondoc.get("timestamp", "")
        severity = jsondoc.get("severity", "").upper()
        message = jsondoc.get("message", "")
    except Exception as e:
        return f"Error parsing log: {str(e)}\n"

    if colorize:
        timestamp = f"\033[0;36m{timestamp}\033[0m"
        severity = f"{SEVERITY_COLORS.get(severity, '')}{severity}\033[0m"
        for keyword in KEYWORDS:
            message = re.sub(rf"\b{keyword}\b", highlight_keyword, message, flags=re.IGNORECASE)

    return f"{timestamp} {severity}  {message}\n"

def process_log_file(infile, outfile, colorize):
    for line in infile:
        outfile.write(format_log_line(line, colorize))

def main():
    parser = argparse.ArgumentParser(description="Parse and format JSON log files")
    parser.add_argument("--nocolor", help="Do not colorize output", action="store_true", default=False)
    parser.add_argument("--outfile", help="Path to output file (default is stdout)", nargs='?', type=argparse.FileType('w'), default=sys.stdout)
    parser.add_argument("infile", help="Path to JSON log file (default is stdin)", nargs='?', type=argparse.FileType('r'), default=sys.stdin)
    args = parser.parse_args()

    colorize = not args.nocolor
    process_log_file(args.infile, args.outfile, colorize)

if __name__ == "__main__":
    try:
        main()
    except (BrokenPipeError, KeyboardInterrupt):
        pass
