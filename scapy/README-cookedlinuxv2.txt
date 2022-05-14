

python3 -m venv venv
. venv/bin/activate

pip install -r requirements.txt


cd scapy/tests

./run_tests -h                             # help
./run_tests                                # run all unit tests
./run_tests -t scapy/layers/l2.uts -F
