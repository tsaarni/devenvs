
gradle run --args referral
gradle run --args starttls
gradle run --args anonymous


cd ~/work/devenvs/keycloak
docker compose up openldap

wireshark -k -i lo -f "port 389 or port 636" -Y "ldap"
