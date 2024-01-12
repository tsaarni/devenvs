



# Use case for admin read-only attributes in declarative user profile?
https://github.com/keycloak/keycloak/discussions/25073


# Declarative user profile becoming default
https://github.com/keycloak/keycloak/issues/16741
https://github.com/keycloak/keycloak/issues/23905


# Add options to change behavior on how unmanaged attributes are managed
https://github.com/keycloak/keycloak/pull/24937


# Set to activate
"KC_SPI_USER_PROFILE_DECLARATIVE_USER_PROFILE_ADMIN_READ_ONLY_ATTRIBUTES": "foo",
"KC_SPI_USER_PROFILE_DECLARATIVE_USER_PROFILE_READ_ONLY_ATTRIBUTES": "foo",
"KC_FEATURES": "declarative-user-profile",




# Test using https://github.com/marcospereirampj/python-keycloak

python3 -m venv .venv
. .venv/bin/activate
pip install python-keycloak

ipython3

from keycloak import KeycloakOpenIDConnection
from keycloak import KeycloakAdmin

conn = KeycloakOpenIDConnection(server_url="http://localhost:8080", username="admin", password="admin")
admin = KeycloakAdmin(connection=conn)
admin.create_user({"username": "joe", "enabled": True})
admin.create_user({"username": "jill", "enabled": True, "attributes": {"foo": ["val"], "baz": ["val"]}})


admin.get_users()

joeid = admin.get_user_id('joe')
admin.update_user(user_id=joeid, payload={'attributes': {"foo": ["val"]}})
admin.update_user(user_id=joeid, payload={'attributes': {"baz": ["val"]}})
admin.get_user(joeid)
