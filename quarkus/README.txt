
export MAVEN_OPTS="-Xmx4g"
./mvnw -Dquickly



the built version will be called:  999-SNAPSHOT





*** Run unit tests in vscode


add following to .vscode/settings.json

{
    "java.test.config": [
        {
            "name": "quarkusConfiguration",
            "vmargs": [ "-Djava.util.logging.manager=org.jboss.logmanager.LogManager", "-Dquarkus.http.test-timeout=600000" ],
        },
    ]
}


Setting property quarkus.http.test-timeout allows stopping at breakpoint without causing: java.net.SocketTimeoutException: Read timed out

