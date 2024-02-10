
https://github.com/keycloak/keycloak/issues/11716




# NOTE: comment out the command from postgres service in docker-compose.yml to disable TLS
docker-compose rm -f
docker-compose up postgres
docker-compose up pgadmin    # http://localhost:8081/browser/


docker-compose rm -v postgres  # remove volumes

http://keycloak.127.0.0.1.nip.io:8080






# Not resolved
Support dynamic custom read/write Strings from @ColumnTransformer annotation for encryption safety/best practices
https://hibernate.atlassian.net/browse/HHH-13827









# example
--- a/model/jpa/src/main/java/org/keycloak/models/jpa/entities/ComponentConfigEntity.java
+++ b/model/jpa/src/main/java/org/keycloak/models/jpa/entities/ComponentConfigEntity.java
@@ -17,6 +17,7 @@

 package org.keycloak.models.jpa.entities;

+import org.hibernate.annotations.ColumnTransformer;
 import org.hibernate.annotations.Nationalized;

 import jakarta.persistence.Access;
@@ -50,6 +51,7 @@ public class ComponentConfigEntity {
     protected String name;
     @Nationalized
     @Column(name = "VALUE")
+    @ColumnTransformer(read = "pgp_sym_decrypt(dearmor(VALUE), 'password')", write = "armor(pgp_sym_encrypt(?, 'password'))")
     protected String value;

     public String getId() {



# enable pgcrypto
CREATE EXTENSION IF NOT EXISTS pgcrypto;



# how to use gpg to decrypt columns manually
https://www.crunchydata.com/blog/postgres-pgcrypto


# how to change annotations programmatically
https://stackoverflow.com/questions/19541252/i-would-like-to-know-if-there-is-a-way-with-hibernate-to-perform-a-programmatic


# store key in postgres.conf
https://stackoverflow.com/questions/42437840/how-to-encrypt-a-column-in-postgres-using-hibernate-columntransformer

pgp_sym_decrypt(" +
            "    test, " +
            "    current_setting('encrypt.key')" +
            ")",
    write = "pgp_sym_encrypt( " +
            "    ?, " +
            "    current_setting('encrypt.key')" +
            ") "

https://www.postgresql.org/docs/current/functions-admin.html


# store key in session variables
https://blog.mclaughlinsoftware.com/2019/09/28/session-variables/

SET SESSION "my.password" = "password";
SELECT pgp_sym_decrypt(dearmor(VALUE), current_setting('my.password')) FROM component_config;
