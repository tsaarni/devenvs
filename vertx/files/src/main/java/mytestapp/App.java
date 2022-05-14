package mytestapp;

import io.vertx.core.Handler;
import io.vertx.core.Vertx;
import io.vertx.core.http.HttpServer;
import io.vertx.core.http.HttpServerOptions;
import io.vertx.core.net.PemKeyCertOptions;

public class App {

    public static void main(String[] args) {
        int port = 8443;

        HttpServerOptions options = new HttpServerOptions()
                .setPort(port)
                .setSsl(true)
                .setPemKeyCertOptions(
                        new PemKeyCertOptions()
                                .addCertPath("/home/tsaarni/work/devenvs/vertx/certs/server.pem")
                                .addKeyPath("/home/tsaarni/work/devenvs/vertx/certs/server-key.pem")
                                .addCertPath("/home/tsaarni/work/devenvs/vertx/certs/server2.pem")
                                .addKeyPath("/home/tsaarni/work/devenvs/vertx/certs/server2-key.pem"))
                .setSni(true)
                .addEnabledSecureTransportProtocol("TLSv1.3");

        Vertx vertx = Vertx.vertx();

        vertx.exceptionHandler(new Handler<Throwable>() {
            @Override
            public void handle(Throwable event) {
                System.err.println("Exception: " + event + event.getStackTrace());
            }
        });

        HttpServer server = vertx.createHttpServer(options);

        server.requestHandler(
                request -> {
                    request.response().end("Hello world");
                })
                .listen(res -> {
                    if (res.succeeded()) {
                        System.out.println("listening :" + port);
                    } else {
                        System.out.println("Failed: " + res.cause());
                    }
                });
    }
}

