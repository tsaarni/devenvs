

Add following to .vscode/launch.json

                "KC_DB_URL": "jdbc:h2:./quarkus/dist/target/keycloakdb;NON_KEYWORDS=VALUE;AUTO_SERVER=TRUE"



# To launch the h2 web console

# NOTE: Run in workspace root !!!!!!!!


h2_version=$(find  ~/.m2/repository/com/h2database/h2/ -maxdepth 1 | sort -V | tail -n 1)
test -e ./quarkus/dist/target/keycloakdb.mv.db && java -cp $h2_version/*.jar org.h2.tools.Console -url "jdbc:h2:file:./quarkus/dist/target/keycloakdb;AUTO_SERVER=TRUE" -user "" -password "" -properties "h2.consoleTimeout=9999999999"

# Remove old h2 file database
rm ./quarkus/dist/target/keycloakdb*



*** Wrong password

Exception in thread "main" org.h2.jdbc.JdbcSQLInvalidAuthorizationSpecException: Wrong user name or password [28000-232]


# Workaround:  delete database
rm ./quarkus/dist/target/keycloakdb*








**** Start H2 web console automatically




cat <<EOF > quarkus/runtime/src/main/java/org/keycloak/quarkus/runtime/H2ConsoleStarter.java

package org.keycloak.quarkus.runtime;

import io.quarkus.runtime.StartupEvent;

import org.h2.server.web.WebServer;
import org.h2.tools.Server;
import org.h2.util.JdbcUtils;
import org.jboss.logging.Logger;

import jakarta.enterprise.context.ApplicationScoped;
import jakarta.enterprise.event.Observes;

import java.sql.Connection;
import java.sql.SQLException;

@ApplicationScoped
public class H2ConsoleStarter {

    private static final Logger logger = Logger.getLogger(H2ConsoleStarter.class);

    final String url = "jdbc:h2:file:./quarkus/dist/target/keycloakdb;AUTO_SERVER=TRUE";
    final String user = "";
    final String password = "";


    void onStart(@Observes StartupEvent ev) {
        try {
            System.setProperty("h2.consoleTimeout", "99999999999");

            WebServer webServer = new WebServer();
            Server web = new Server(webServer, "-webPort", "8082");
            web.start();
            webServer.setShutdownHandler(web::stop);

            Connection conn = JdbcUtils.getConnection(null, url, user, password);
            String path = webServer.addSession(conn);

            logger.infov("H2 console available at {0}", path);
        } catch (SQLException e) {
            throw new RuntimeException("Failed to start H2 console", e);
        }
    }
}
EOF


The URL will be printed in the logs:

2025-03-13 12:45:44,041 INFO  [org.keycloak.quarkus.runtime.H2ConsoleStarter] (Quarkus Main Thread) H2 console available at http://127.0.1.1:8082/frame.jsp?jsessionid=6fb6a1021cbc027cb248fad165e8a7bf
