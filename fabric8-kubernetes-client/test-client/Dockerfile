# docker build . -t fabric8io-kubernetes-client-app:latest
# kind load docker-image fabric8io-kubernetes-client-app:latest --name exptest

FROM openjdk:11-slim

COPY app/build/distributions/app.tar /
RUN tar xvf app.tar

CMD ["/app/bin/app"]
