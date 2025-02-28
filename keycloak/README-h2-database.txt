

Add following to .vscode/launch.json

                "KC_DB_URL": "jdbc:h2:./quarkus/dist/target/keycloakdb;NON_KEYWORDS=VALUE;AUTO_SERVER=TRUE"



# To launch the h2 web console

# NOTE: Run in workspace root !!!!!!!!


h2_version=$(find  ~/.m2/repository/com/h2database/h2/ -maxdepth 1 | sort -V | tail -n 1)
java -cp $h2_version/*.jar org.h2.tools.Console -url "jdbc:h2:file:./quarkus/dist/target/keycloakdb;AUTO_SERVER=TRUE" -user "" -password ""

# Remove old h2 file database
rm ./quarkus/dist/target/keycloakdb*



*** Wrong password

Exception in thread "main" org.h2.jdbc.JdbcSQLInvalidAuthorizationSpecException: Wrong user name or password [28000-232]


# Workaround:  delete database
rm ./quarkus/dist/target/keycloakdb*
