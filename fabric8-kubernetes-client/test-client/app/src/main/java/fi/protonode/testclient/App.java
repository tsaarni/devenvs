package fi.protonode.testclient;

import io.fabric8.kubernetes.api.model.ListOptionsBuilder;
import io.fabric8.kubernetes.api.model.PodList;
import io.fabric8.kubernetes.client.ConfigBuilder;
import io.fabric8.kubernetes.client.DefaultKubernetesClient;
import io.fabric8.kubernetes.client.KubernetesClient;
import java.util.concurrent.TimeUnit;

/**
 * Hello world!
 *
 */
public class App
{
    public static void main( String[] args ) throws InterruptedException
    {
        final ConfigBuilder configBuilder = new ConfigBuilder();
        try (KubernetesClient client = new DefaultKubernetesClient(configBuilder.build())) {
            while (true) {
                PodList podList = client.pods().inNamespace("default").list(new ListOptionsBuilder().withLimit(5L).build());
                for (var pod : podList.getItems()) {
                    System.out.printf("Pod: %s\n", pod.getMetadata().getName());
                }

                System.out.println("waiting for 1 second...");
                TimeUnit.SECONDS.sleep(1);
            }
        }
    }
}
