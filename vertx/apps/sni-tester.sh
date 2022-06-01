#!/bin/bash
#
# sni-tester.sh localhost:8443 "" hostname1 hostname2

connect=$1

printf "\nConnecting to $connect\n\n"
for ((i=2; i<=$#; i++))
do
  servername=${!i}
  if [[ -z $servername ]]; then
    printf "    no SNI:\n"
    echo Q | openssl s_client -connect $connect 2>/dev/null | openssl x509 -text -noout | grep "Issuer:\|Subject:\|DNS:\|Validity\|Not Before:\|Not After :"
  else
    printf "    sending SNI servername: $servername\n"
    echo Q | openssl s_client -connect $connect -servername $servername 2>/dev/null | openssl x509 -text -noout | grep "Issuer:\|Subject:\|DNS:\|Validity\|Not Before:\|Not After :"
  fi
  printf "\n"
done
