
# Issues
#
# Spaces in environment variables cause configuration problems in filters
# https://github.com/fluent/fluent-bit/issues/1225
#
# [request] Support spaces in Modify filter keys and values
# https://github.com/fluent/fluent-bit/issues/4286



cd build
cmake ..
make -j

bin/fluent-bit -v --config=/home/tsaarni/work/devenvs/fluentbit/fluent-bit.conf
echo '{ "message": "Hello FOO world!" }' > /dev/tcp/localhost/5170


# tests
cmake -DFLB_DEV=On -DFLB_TESTS_RUNTIME=On -DFLB_TESTS_INTERNAL=On ..
make -j
make test


# The following tests FAILED:
#           8 - flb-rt-in_netif (Failed)
#          61 - flb-rt-out_td (Failed)
#         108 - flb-it-signv4 (Failed)
# Errors while running CTest


# to run individual tests

# broken by patch
bin/flb-rt-in_netif
bin/flb-it-signv4

# possibly fails always, patch or not
bin/flb-rt-out_td


# run just single test
make flb-rt-filter_modify
bin/flb-rt-filter_modify operation_with_whitespace

bin/flb-it-utils test_flb_utils_split_quoted









# Same test with tail plugin


cat >fluent-bit.conf <<EOF
[INPUT]
   name tail
   path test.log
   parser json

[OUTPUT]
   name stdout

[FILTER]
   name modify
   match *
   condition key_value_matches message .*FOO.*
   add facility "foo bar"
EOF


cat >parsers.conf <<EOF
[PARSER]
   name   json
   format json
EOF


bin/fluent-bit --config=fluent-bit.conf --parser=parsers.conf

echo '{ "message": "Hello FOO world!" }' >> test.log

