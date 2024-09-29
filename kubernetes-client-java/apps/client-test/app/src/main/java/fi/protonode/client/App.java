
package fi.protonode.client;

import io.kubernetes.client.openapi.ApiClient;
import io.kubernetes.client.openapi.ApiException;
import io.kubernetes.client.openapi.Configuration;
import io.kubernetes.client.openapi.apis.CoreV1Api;
import io.kubernetes.client.openapi.apis.CoreV1Api.APIlistNamespacedPodRequest;
import io.kubernetes.client.openapi.models.V1Pod;
import io.kubernetes.client.openapi.models.V1PodList;
import io.kubernetes.client.util.Config;

import java.io.IOException;

public class App {

    // New 20.0.1
    //
    private static V1PodList listPods(CoreV1Api api) throws ApiException {
        APIlistNamespacedPodRequest list = api.listNamespacedPod("default");
        return list.execute();
    }

    // Legacy 20.0.1-legacy
    //
    // private static V1PodList listPods(CoreV1Api api) throws ApiException {
    //     return api.listNamespacedPod("default", null, null, null, null, null, null, null, null, null, null, null);
    // }

    public static void main(String[] args) throws IOException, ApiException {
        ApiClient client = Config.defaultClient();
        Configuration.setDefaultApiClient(client);

        CoreV1Api api = new CoreV1Api();
        for (V1Pod pod : listPods(api).getItems()) {
            System.out.println("Pod name: " + pod.getMetadata().getName());
            System.out.println(pod.getSpec());
            System.out.println("----------------");
        }

    }
}
