# Test application for service account token renewal


This is a test application for https://github.com/fabric8io/kubernetes-client to test that service account token is reloaded after the token expires and authentication fails.

## Test procedure

1. Clone the pull request from https://github.com/fabric8io/kubernetes-client/pull/3445, build an install it to local maven cache.

2. Build the test application by running `gradle build`

3. Create docker container by running `docker build . -t fabric8io-kubernetes-client-app:latest`

4. Setup a Kubernetes cluster with expiring service account token.
See instructions at https://github.com/tsaarni/cloud-playground/tree/master/kubernetes/sa-token-expiration.

5. Upload the docker image to the cluster by running `kind load docker-image fabric8io-kubernetes-client-app:latest --name exptest`.

6. Deploy the test application by runnig `kubectl apply -f manifests/manifests/clientapp.yaml`

7. The test application will list pods in `default` namespace once every second.
Observe the logs for >10 minutes to see that the token is reloaded `kubectl logs fabric8io-kubernetes-client-app`


Following debug logs show that token has been successfully reloaded after authentication fails.

```console
waiting for 1 second...
Pod: fabric8io-kubernetes-client-app
waiting for 1 second...
[main] DEBUG io.fabric8.kubernetes.client.Config - Trying to configure client from Kubernetes config...
[main] DEBUG io.fabric8.kubernetes.client.Config - Did not find Kubernetes config at: [/root/.kube/config]. Ignoring.
[main] DEBUG io.fabric8.kubernetes.client.Config - Trying to configure client from service account...
[main] DEBUG io.fabric8.kubernetes.client.Config - Found service account host and port: 10.96.0.1:443
[main] DEBUG io.fabric8.kubernetes.client.Config - Found service account ca cert at: [/var/run/secrets/kubernetes.io/serviceaccount/ca.crt}].
[main] DEBUG io.fabric8.kubernetes.client.Config - Found service account token at: [/projected/token].
[main] DEBUG io.fabric8.kubernetes.client.Config - Trying to configure client namespace from Kubernetes service account namespace path...
[main] DEBUG io.fabric8.kubernetes.client.Config - Found service account namespace at: [/var/run/secrets/kubernetes.io/serviceaccount/namespace].
[main] DEBUG io.fabric8.kubernetes.client.Config - Trying to configure client from Kubernetes config...
[main] DEBUG io.fabric8.kubernetes.client.Config - Did not find Kubernetes config at: [/root/.kube/config]. Ignoring.
[main] DEBUG io.fabric8.kubernetes.client.Config - Trying to configure client from service account...
[main] DEBUG io.fabric8.kubernetes.client.Config - Found service account host and port: 10.96.0.1:443
[main] DEBUG io.fabric8.kubernetes.client.Config - Found service account ca cert at: [/var/run/secrets/kubernetes.io/serviceaccount/ca.crt}].
[main] DEBUG io.fabric8.kubernetes.client.Config - Found service account token at: [/projected/token].
[main] DEBUG io.fabric8.kubernetes.client.Config - Trying to configure client namespace from Kubernetes service account namespace path...
[main] DEBUG io.fabric8.kubernetes.client.Config - Found service account namespace at: [/var/run/secrets/kubernetes.io/serviceaccount/namespace].
Pod: fabric8io-kubernetes-client-app
waiting for 1 second...
...
```

## Alternative flow

In this flow, authentication still fails after token has been reloaded

1. Exec into the pod and make a copy of the service account token

2. Start the test application with environment variable pointing to the copy `KUBERNETES_AUTH_SERVICEACCOUNT_TOKEN=/copy-of-token`

3. Wait until the copy of service account token expires and observe the authentication failure.
It will take ~minute after the expiration time has passed.
Following backtrace is printed:

```
Exception in thread "main" io.fabric8.kubernetes.client.KubernetesClientException: Failure executing: GET at: https://10.96.0.1/api/v1/namespaces/default/pods?limit=5. Message: Unauthorized! Configured service account doesn't have access. Service account may have been revoked. Unauthorized.
        at io.fabric8.kubernetes.client.dsl.base.OperationSupport.requestFailure(OperationSupport.java:686)
        at io.fabric8.kubernetes.client.dsl.base.OperationSupport.assertResponseCode(OperationSupport.java:623)
        at io.fabric8.kubernetes.client.dsl.base.OperationSupport.handleResponse(OperationSupport.java:565)
        at io.fabric8.kubernetes.client.dsl.base.OperationSupport.handleResponse(OperationSupport.java:526)
        at io.fabric8.kubernetes.client.dsl.base.OperationSupport.handleResponse(OperationSupport.java:509)
        at io.fabric8.kubernetes.client.dsl.base.BaseOperation.listRequestHelper(BaseOperation.java:137)
        at io.fabric8.kubernetes.client.dsl.base.BaseOperation.list(BaseOperation.java:524)
        at io.fabric8.kubernetes.client.dsl.base.BaseOperation.list(BaseOperation.java:88)
        at fi.protonode.testclient.App.main(App.java:21)
```
