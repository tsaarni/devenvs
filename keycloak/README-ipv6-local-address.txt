



mkdir -p .vscode
cp -a ~/work/devenvs/keycloak/configs/launch.json  ~/work/devenvs/keycloak/configs/settings.json .vscode



# create dummy interface
sudo ip link add name dummy0 type dummy
sudo ip link set dev dummy0 up

# assign ULA addresses
sudo ip -6 addr add fc00:2000::2000/64 dev dummy0
sudo ip -6 addr add fd00:3000::3000/64 dev dummy0

# assign global addresses
sudo ip -6 addr add 2001::4000/23 dev dummy0


# Check the addresses and make note of the (automatically assigned) link-local address
ip addr show dummy0

# Make request with link-local address
http --form POST http://[fe80::8852:3ff:fee6:58db%dummy0]:8080/realms/master/protocol/openid-connect/token username=admin password=admin grant_type=password client_id=admin-cli

# Make request with ULA
http --form POST http://[fc00:2000::2000]:8080/realms/master/protocol/openid-connect/token username=admin password=admin grant_type=password client_id=admin-cli
http --form POST http://[fd00:3000::3000]:8080/realms/master/protocol/openid-connect/token username=admin password=admin grant_type=password client_id=admin-cli

# Make request with global address (this will fail)
http --form POST http://[2001::4000]:8080/realms/master/protocol/openid-connect/token username=admin password=admin grant_type=password client_id=admin-cliHTTP/1.1
