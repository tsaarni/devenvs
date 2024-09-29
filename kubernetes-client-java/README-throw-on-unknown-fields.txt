

# client-java throws an exception on unknown fields

https://github.com/kubernetes-client/java/issues/3428




# To reproduce, use 20.0.1 version of client-java and run apps/client-test

# contains field that is unknown to java-client
kubectl apply -f manifests/apparmor.yaml

cd apps/client-test/
gradle run

# modify app/src/main/java/fi/protonode/client/App.java for 20.0.1-legacy to compare behavior on unknown fields


