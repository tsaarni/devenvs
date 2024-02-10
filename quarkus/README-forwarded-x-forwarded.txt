
https://github.com/quarkusio/quarkus/issues/35751

export MAVEN_OPTS="-Xmx4g"
./mvnw -Dquickly


./mvnw test -f extensions/vertx-http/deployment -Dtest=io.quarkus.vertx.http.AllowBothForwardedHeadersTest


mvn verify -Dtest=io.quarkus.vertx.http.*Forwarded*
mvn verify -Dtest=io.quarkus.vertx.http.AllowBothForwardedHeadersTest
AllowOnlyForwardedHeaderTest
AllowOnlyXForwardedHeaderTest
AllowOnlyXForwardedHeaderUsingDefaultConfigTest
ForwardedHeaderTest
