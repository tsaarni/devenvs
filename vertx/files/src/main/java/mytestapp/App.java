package mytestapp;

import io.vertx.core.Vertx;
import io.vertx.core.http.HttpServer;
import io.vertx.core.http.HttpServerOptions;
import io.vertx.core.net.PemKeyCertOptions;
import io.vertx.core.net.PfxOptions;

public class App {

  static final int port = 8443;

  public static void main(String[] args) {

    System.setProperty("javax.net.debug", "keymanager");
    System.setProperty("vertxweb.environment", "development");

    Vertx vertx = Vertx.vertx();

    HttpServerOptions options = new HttpServerOptions()
        .setPort(port)
        .setSsl(true)
        .setSni(true)
        .addEnabledSecureTransportProtocol("TLSv1.3");
    configurePkcs12(options);
    //configurePem(options);

    HttpServer server = vertx.createHttpServer(options);

    server.requestHandler(request -> request.response().end("Hello world"))
        .exceptionHandler(e -> System.err.println("Exception: " + e.getStackTrace()))
        .listen(r -> {
          if (r.failed()) {
            r.cause().printStackTrace();
          }
        });
  }

  static void configurePkcs12(HttpServerOptions options) {
    options.setPfxKeyCertOptions(
        new PfxOptions()
            //.setPath("/home/tsaarni/work/devenvs/vertx/certs/server.p12")
            //.setPassword("secret"));
            .setPath("/home/tsaarni/work/vert.x/src/test/resources/tls/sni-keystore.jks")
            .setPassword("wibble"));

  }

  static void configurePem(HttpServerOptions options) {
    options.setPemKeyCertOptions(
        new PemKeyCertOptions()
            .addCertPath("/home/tsaarni/work/devenvs/vertx/certs/server.pem")
            .addKeyPath("/home/tsaarni/work/devenvs/vertx/certs/server-key.pem")
            .addCertPath("/home/tsaarni/work/devenvs/vertx/certs/server2.pem")
            .addKeyPath("/home/tsaarni/work/devenvs/vertx/certs/server2-key.pem"));

  }

}
