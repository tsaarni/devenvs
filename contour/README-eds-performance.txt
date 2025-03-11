# issue and analysis about EDS performance problem
# https://github.com/projectcontour/contour/issues/6743


# change to separate snapshot cache for EDS and comment about LinearCache
https://github.com/projectcontour/contour/pull/6250/files#diff-1900033d7ec953a2e318fa5369de7c7d199c4e478873a17e61b7f3776771bb1e


kubectl apply -f manifests/echoserver.yaml



function echoserver() {
    local action=$1
    local num=$2
    cat <<EOF | kubectl $action -f -
apiVersion: projectcontour.io/v1
kind: HTTPProxy
metadata:
    name: echoserver-$num
    annotations:
        date: "$(date)"
spec:
    virtualhost:
        fqdn: echoserver-$num.127-0-0-101.nip.io
    routes:
        - services:
            - name: echoserver-$num
              port: 80
---
apiVersion: v1
kind: Service
metadata:
    name: echoserver-$num
    annotations:
        date: "$(date)"
spec:
    selector:
        app.kubernetes.io/name: echoserver
    ports:
        - protocol: TCP
          port: 80
          targetPort: http-api
---
EOF
kubectl get httpproxies,services -o wide
}



echoserver create 0001
echoserver delete 0001
echoserver apply 0001


kubectl delete pod -l app.kubernetes.io/name=echoserver --wait=false


http http://echoserver-0000.127-0-0-101.nip.io
http http://echoserver-7999.127-0-0-101.nip.io



# contour metrics
http http://localhost:8000/metrics
http http://localhost:8000/metrics | grep EndpointDiscoveryService

# envoy metrics
kubectl -n projectcontour port-forward daemonset/envoy 9001:9001
http http://localhost:9001/stats



kubectl -n projectcontour get secret envoycert -o jsonpath='{..ca\.crt}' | base64 -d > envoy-ca.crt
kubectl -n projectcontour get secret envoycert -o jsonpath='{..tls\.crt}' | base64 -d > envoy-tls.crt
kubectl -n projectcontour get secret envoycert -o jsonpath='{..tls\.key}' | base64 -d > envoy-tls.key


###kubectl -n projectcontour port-forward deployment/contour 8001:8001



go run github.com/projectcontour/contour/cmd/contour cli eds --cafile=envoy-ca.crt --cert-file=envoy-tls.crt --key-file=envoy-tls.key --contour="localhost:8001"

kubectl -n projectcontour rollout restart daemonset envoy

kubectl rollout restart deployment echoserver
kubectl get pod -o wide





go get github.com/tsaarni/grpc-json-sniffer


diff --git a/internal/xds/server.go b/internal/xds/server.go
index 90d0e330..56c66926 100644
--- a/internal/xds/server.go
+++ b/internal/xds/server.go
@@ -16,6 +16,7 @@ package xds
 import (
        grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
        "github.com/prometheus/client_golang/prometheus"
+       sniffer "github.com/tsaarni/grpc-json-sniffer"
        "google.golang.org/grpc"
 )

@@ -29,9 +30,11 @@ func NewServer(registry *prometheus.Registry, opts ...grpc.ServerOption) *grpc.S
                metrics = grpc_prometheus.NewServerMetrics()
                registry.MustRegister(metrics)

+               interceptor, _ := sniffer.NewGrpcJsonInterceptor()
+
                opts = append(opts,
-                       grpc.StreamInterceptor(metrics.StreamServerInterceptor()),
-                       grpc.UnaryInterceptor(metrics.UnaryServerInterceptor()),
+                       grpc.ChainStreamInterceptor(metrics.StreamServerInterceptor(), interceptor.StreamServerInterceptor()),
+                       grpc.ChainUnaryInterceptor(metrics.UnaryServerInterceptor(), interceptor.UnaryServerInterceptor()),
                )
        }




  Id  Message              Stream ID  Version Info                          Nonce    Resource Name            Addresses
----  -----------------  -----------  ------------------------------------  -------  -----------------------  ----------------------------
*** Create new service and httpproxy (echoserver-0001)
  36  DiscoveryRequest             5                                                 default/echoserver-0001
  38  DiscoveryResponse            5  a2e7306d-226b-4322-bb7e-37050da2d639  1        default/echoserver-0001  ['10.244.1.4']
  39  DiscoveryRequest             5  a2e7306d-226b-4322-bb7e-37050da2d639  1        default/echoserver-0001
  49  DiscoveryResponse            5  715dca95-b946-44ff-9b8f-21c653eba2bf  2        default/echoserver-0001  ['10.244.1.4']
  53  DiscoveryRequest             5  715dca95-b946-44ff-9b8f-21c653eba2bf  2        default/echoserver-0001
*** Create new service and httpproxy (echoserver-0002)
  58  DiscoveryRequest             6                                                 default/echoserver-0002
  59  DiscoveryResponse            6  715dca95-b946-44ff-9b8f-21c653eba2bf  1        default/echoserver-0002  ['10.244.1.4']
  60  DiscoveryRequest             6  715dca95-b946-44ff-9b8f-21c653eba2bf  1        default/echoserver-0002
  70  DiscoveryResponse            5  4a792f9b-f8e4-4959-90db-748b2ac47491  3        default/echoserver-0001  ['10.244.1.4']
  71  DiscoveryResponse            6  4a792f9b-f8e4-4959-90db-748b2ac47491  2        default/echoserver-0002  ['10.244.1.4']
*** Create new service and httpproxy (echoserver-0003)
  75  DiscoveryRequest             7                                                 default/echoserver-0003
  76  DiscoveryRequest             5  4a792f9b-f8e4-4959-90db-748b2ac47491  3        default/echoserver-0001
  77  DiscoveryRequest             6  4a792f9b-f8e4-4959-90db-748b2ac47491  2        default/echoserver-0002
  78  DiscoveryResponse            7  4a792f9b-f8e4-4959-90db-748b2ac47491  1        default/echoserver-0003  ['10.244.1.4']
  83  DiscoveryRequest             7  4a792f9b-f8e4-4959-90db-748b2ac47491  1        default/echoserver-0003
*** Restart backend service pod (address changes)
  84  DiscoveryResponse            7  8538dde7-8b88-41f5-b712-2683856f9107  2        default/echoserver-0003  ['10.244.1.4']
  85  DiscoveryResponse            6  8538dde7-8b88-41f5-b712-2683856f9107  3        default/echoserver-0002  ['10.244.1.4', '10.244.1.5']
  86  DiscoveryResponse            5  8538dde7-8b88-41f5-b712-2683856f9107  4        default/echoserver-0001  ['10.244.1.4']
  87  DiscoveryRequest             5  8538dde7-8b88-41f5-b712-2683856f9107  4        default/echoserver-0001
  88  DiscoveryRequest             6  8538dde7-8b88-41f5-b712-2683856f9107  3        default/echoserver-0002
  89  DiscoveryRequest             7  8538dde7-8b88-41f5-b712-2683856f9107  2        default/echoserver-0003
  90  DiscoveryResponse            5  e5d7fe89-4b5e-41d2-a093-e692f6942aec  5        default/echoserver-0001  ['10.244.1.4', '10.244.1.5']
  91  DiscoveryResponse            6  e5d7fe89-4b5e-41d2-a093-e692f6942aec  4        default/echoserver-0002  ['10.244.1.4', '10.244.1.5']
  92  DiscoveryResponse            7  e5d7fe89-4b5e-41d2-a093-e692f6942aec  3        default/echoserver-0003  ['10.244.1.4', '10.244.1.5']
  93  DiscoveryRequest             5  e5d7fe89-4b5e-41d2-a093-e692f6942aec  5        default/echoserver-0001
  94  DiscoveryRequest             6  e5d7fe89-4b5e-41d2-a093-e692f6942aec  4        default/echoserver-0002
  95  DiscoveryRequest             7  e5d7fe89-4b5e-41d2-a093-e692f6942aec  3        default/echoserver-0003
  96  DiscoveryResponse            6  9caf4025-50f4-40d2-9927-06431fa4cd74  5        default/echoserver-0002  ['10.244.1.4', '10.244.1.5']
  97  DiscoveryResponse            7  9caf4025-50f4-40d2-9927-06431fa4cd74  4        default/echoserver-0003  ['10.244.1.5']
  98  DiscoveryResponse            5  9caf4025-50f4-40d2-9927-06431fa4cd74  6        default/echoserver-0001  ['10.244.1.4', '10.244.1.5']
  99  DiscoveryRequest             7  9caf4025-50f4-40d2-9927-06431fa4cd74  4        default/echoserver-0003
 100  DiscoveryRequest             5  9caf4025-50f4-40d2-9927-06431fa4cd74  6        default/echoserver-0001
 101  DiscoveryResponse            7  8cba205f-4d88-4e98-992b-91ef68e72f54  5        default/echoserver-0003  ['10.244.1.5']
 102  DiscoveryRequest             6  9caf4025-50f4-40d2-9927-06431fa4cd74  5        default/echoserver-0002
 103  DiscoveryResponse            5  8cba205f-4d88-4e98-992b-91ef68e72f54  7        default/echoserver-0001  ['10.244.1.5']
 104  DiscoveryResponse            6  8cba205f-4d88-4e98-992b-91ef68e72f54  6        default/echoserver-0002  ['10.244.1.5']
 105  DiscoveryRequest             7  8cba205f-4d88-4e98-992b-91ef68e72f54  5        default/echoserver-0003
 106  DiscoveryRequest             6  8cba205f-4d88-4e98-992b-91ef68e72f54  6        default/echoserver-0002
 107  DiscoveryRequest             5  8cba205f-4d88-4e98-992b-91ef68e72f54  7        default/echoserver-0001
 108  DiscoveryResponse            6  ec08e210-870d-4810-876a-44cd5cc7690c  7        default/echoserver-0002  ['10.244.1.5']
 109  DiscoveryResponse            5  ec08e210-870d-4810-876a-44cd5cc7690c  8        default/echoserver-0001  ['10.244.1.5']
 110  DiscoveryResponse            7  ec08e210-870d-4810-876a-44cd5cc7690c  6        default/echoserver-0003  ['10.244.1.5']
 111  DiscoveryRequest             7  ec08e210-870d-4810-876a-44cd5cc7690c  6        default/echoserver-0003
 112  DiscoveryRequest             5  ec08e210-870d-4810-876a-44cd5cc7690c  8        default/echoserver-0001
 113  DiscoveryRequest             6  ec08e210-870d-4810-876a-44cd5cc7690c  7        default/echoserver-0002
 114  DiscoveryResponse            7  c4658889-9134-4698-a7d8-d147f9eb71c9  7        default/echoserver-0003  ['10.244.1.5']
 115  DiscoveryResponse            5  c4658889-9134-4698-a7d8-d147f9eb71c9  9        default/echoserver-0001  ['10.244.1.5']
 116  DiscoveryResponse            6  c4658889-9134-4698-a7d8-d147f9eb71c9  8        default/echoserver-0002  ['10.244.1.5']
 117  DiscoveryRequest             7  c4658889-9134-4698-a7d8-d147f9eb71c9  7        default/echoserver-0003
 118  DiscoveryRequest             5  c4658889-9134-4698-a7d8-d147f9eb71c9  9        default/echoserver-0001
 119  DiscoveryRequest             6  c4658889-9134-4698-a7d8-d147f9eb71c9  8        default/echoserver-0002


  Id  Message              Stream ID  Version Info                            Nonce    Resource Name            Addresses
----  -----------------  -----------  --------------------------------------  -------  -----------------------  ----------------------------
*** Create new service and httpproxy (echoserver-0001)
  36  DiscoveryRequest             5                                                   default/echoserver-0001
  38  DiscoveryResponse            5  140edf95-b952-455a-a3a6-0f4894086eb1-1  1        default/echoserver-0001  ['10.244.1.5']
  39  DiscoveryRequest             5  140edf95-b952-455a-a3a6-0f4894086eb1-1  1        default/echoserver-0001
*** Create new service and httpproxy (echoserver-0002)
  56  DiscoveryRequest             6                                                   default/echoserver-0002
  57  DiscoveryResponse            6  140edf95-b952-455a-a3a6-0f4894086eb1-2  1        default/echoserver-0002  ['10.244.1.5']
  58  DiscoveryRequest             6  140edf95-b952-455a-a3a6-0f4894086eb1-2  1        default/echoserver-0002
*** Create new service and httpproxy (echoserver-0003)
  75  DiscoveryRequest             7                                                   default/echoserver-0003
  76  DiscoveryResponse            7  140edf95-b952-455a-a3a6-0f4894086eb1-3  1        default/echoserver-0003  ['10.244.1.5']
  77  DiscoveryRequest             7  140edf95-b952-455a-a3a6-0f4894086eb1-3  1        default/echoserver-0003
*** Restart upstream servcive pod (address changes)
  78  DiscoveryResponse            7  140edf95-b952-455a-a3a6-0f4894086eb1-4  2        default/echoserver-0003  ['10.244.1.5', '10.244.1.6']
  79  DiscoveryResponse            5  140edf95-b952-455a-a3a6-0f4894086eb1-5  2        default/echoserver-0001  ['10.244.1.5', '10.244.1.6']
  80  DiscoveryResponse            6  140edf95-b952-455a-a3a6-0f4894086eb1-6  2        default/echoserver-0002  ['10.244.1.5', '10.244.1.6']
  81  DiscoveryRequest             5  140edf95-b952-455a-a3a6-0f4894086eb1-5  2        default/echoserver-0001
  82  DiscoveryRequest             7  140edf95-b952-455a-a3a6-0f4894086eb1-4  2        default/echoserver-0003
  83  DiscoveryRequest             6  140edf95-b952-455a-a3a6-0f4894086eb1-6  2        default/echoserver-0002
  84  DiscoveryResponse            5  140edf95-b952-455a-a3a6-0f4894086eb1-7  3        default/echoserver-0001  ['10.244.1.6']
  85  DiscoveryResponse            6  140edf95-b952-455a-a3a6-0f4894086eb1-8  3        default/echoserver-0002  ['10.244.1.6']
  86  DiscoveryResponse            7  140edf95-b952-455a-a3a6-0f4894086eb1-9  3        default/echoserver-0003  ['10.244.1.6']
  87  DiscoveryRequest             6  140edf95-b952-455a-a3a6-0f4894086eb1-8  3        default/echoserver-0002
  88  DiscoveryRequest             7  140edf95-b952-455a-a3a6-0f4894086eb1-9  3        default/echoserver-0003
  89  DiscoveryRequest             5  140edf95-b952-455a-a3a6-0f4894086eb1-7  3        default/echoserver-0001




  Id  Message              Stream ID  Version Info                            Nonce    Resource Name            Addresses
----  -----------------  -----------  --------------------------------------  -------  -----------------------  --------------
*** Envoy is restarted
 103  DiscoveryRequest            10                                                   default/echoserver-0002
 104  DiscoveryRequest            12                                                   default/echoserver-0001
 105  DiscoveryRequest            11                                                   default/echoserver-0003
 106  DiscoveryResponse           10  140edf95-b952-455a-a3a6-0f4894086eb1-9  1        default/echoserver-0002  ['10.244.1.6']
 107  DiscoveryResponse           12  140edf95-b952-455a-a3a6-0f4894086eb1-9  1        default/echoserver-0001  ['10.244.1.6']
 108  DiscoveryResponse           11  140edf95-b952-455a-a3a6-0f4894086eb1-9  1        default/echoserver-0003  ['10.244.1.6']
 109  DiscoveryRequest            10  140edf95-b952-455a-a3a6-0f4894086eb1-9  1        default/echoserver-0002
 111  DiscoveryRequest            11  140edf95-b952-455a-a3a6-0f4894086eb1-9  1        default/echoserver-0003
 112  DiscoveryRequest            12  140edf95-b952-455a-a3a6-0f4894086eb1-9  1        default/echoserver-0001



  Id  Message              Stream ID  Version Info                            Nonce    Resource Name            Addresses
----  -----------------  -----------  --------------------------------------  -------  -----------------------  --------------
*** Restart upstream service pod (address changes) while Contour is down, then Contour comes up
  10  DiscoveryRequest             4  f6891efb-a619-4052-a406-8ce77e11f285-1           default/echoserver-0001
  11  DiscoveryResponse            4  1709b3a2-f9de-4fd1-8c4f-91dcec68f502-1  1        default/echoserver-0001  ['10.244.1.9']
  12  DiscoveryRequest             4  1709b3a2-f9de-4fd1-8c4f-91dcec68f502-1  1        default/echoserver-0001
  13  DiscoveryRequest             5  f6891efb-a619-4052-a406-8ce77e11f285-1           default/echoserver-0002
  14  DiscoveryResponse            5  1709b3a2-f9de-4fd1-8c4f-91dcec68f502-1  1        default/echoserver-0002  ['10.244.1.9']
  15  DiscoveryRequest             5  1709b3a2-f9de-4fd1-8c4f-91dcec68f502-1  1        default/echoserver-0002
  22  DiscoveryRequest             6  f6891efb-a619-4052-a406-8ce77e11f285-1           default/echoserver-0003
  23  DiscoveryResponse            6  1709b3a2-f9de-4fd1-8c4f-91dcec68f502-1  1        default/echoserver-0003  ['10.244.1.9']
  24  DiscoveryRequest             6  1709b3a2-f9de-4fd1-8c4f-91dcec68f502-1  1        default/echoserver-0003





  Id  Message              Stream ID  Version Info                          Nonce    Resource Name            Addresses
----  -----------------  -----------  ------------------------------------  -------  -----------------------  --------------
*** Contour with SnapshotCache starts, Envoy starts
   3  DiscoveryRequest             2                                                 default/echoserver-0001
   4  DiscoveryRequest             4                                                 default/echoserver-0003
   5  DiscoveryResponse            2  88df7617-455e-4bd6-a116-8f615b1ce9e8  1        default/echoserver-0001  ['10.244.1.9']
   7  DiscoveryRequest             3                                                 default/echoserver-0002
   8  DiscoveryResponse            4  88df7617-455e-4bd6-a116-8f615b1ce9e8  1        default/echoserver-0003  ['10.244.1.9']
   9  DiscoveryResponse            3  88df7617-455e-4bd6-a116-8f615b1ce9e8  1        default/echoserver-0002  ['10.244.1.9']
  10  DiscoveryRequest             2  88df7617-455e-4bd6-a116-8f615b1ce9e8  1        default/echoserver-0001
  11  DiscoveryRequest             4  88df7617-455e-4bd6-a116-8f615b1ce9e8  1        default/echoserver-0003
  12  DiscoveryRequest             3  88df7617-455e-4bd6-a116-8f615b1ce9e8  1        default/echoserver-0002

*** Contour is upgraded to LinearCache
   1  DiscoveryRequest             1  88df7617-455e-4bd6-a116-8f615b1ce9e8             default/echoserver-0001
   2  DiscoveryResponse            1  dc40621f-69e5-4055-83ae-d69648c29add-1  1        default/echoserver-0001  ['10.244.1.9']
   3  DiscoveryRequest             1  dc40621f-69e5-4055-83ae-d69648c29add-1  1        default/echoserver-0001
  13  DiscoveryRequest             5  88df7617-455e-4bd6-a116-8f615b1ce9e8             default/echoserver-0003
  14  DiscoveryResponse            5  dc40621f-69e5-4055-83ae-d69648c29add-1  1        default/echoserver-0003  ['10.244.1.9']
  15  DiscoveryRequest             5  dc40621f-69e5-4055-83ae-d69648c29add-1  1        default/echoserver-0003
  27  DiscoveryRequest             7  88df7617-455e-4bd6-a116-8f615b1ce9e8             default/echoserver-0002
  28  DiscoveryResponse            7  dc40621f-69e5-4055-83ae-d69648c29add-1  1        default/echoserver-0002  ['10.244.1.9']
  29  DiscoveryRequest             7  dc40621f-69e5-4055-83ae-d69648c29add-1  1        default/echoserver-0002




  Id  Message              Stream ID  Version Info                          Nonce    Resource Name            Addresses
----  -----------------  -----------  ------------------------------------  -------  -----------------------  --------------
*** Contour with SnapshotCache starts, Envoy starts
   3  DiscoveryRequest             2                                                 default/echoserver-0001
   4  DiscoveryRequest             4                                                 default/echoserver-0003
   6  DiscoveryResponse            2  ee0e8167-8383-4a04-8304-9a8ba75a7b93  1        default/echoserver-0001  ['10.244.1.9']
   7  DiscoveryRequest             3                                                 default/echoserver-0002
   8  DiscoveryResponse            4  ee0e8167-8383-4a04-8304-9a8ba75a7b93  1        default/echoserver-0003  ['10.244.1.9']
   9  DiscoveryRequest             2  ee0e8167-8383-4a04-8304-9a8ba75a7b93  1        default/echoserver-0001
  10  DiscoveryResponse            3  ee0e8167-8383-4a04-8304-9a8ba75a7b93  1        default/echoserver-0002  ['10.244.1.9']
  11  DiscoveryRequest             4  ee0e8167-8383-4a04-8304-9a8ba75a7b93  1        default/echoserver-0003
  12  DiscoveryRequest             3  ee0e8167-8383-4a04-8304-9a8ba75a7b93  1        default/echoserver-0002

*** Restart upstream service pod (address changes) while Contour is down, then Contour with LinearCache comes up

   4  DiscoveryRequest             2  ee0e8167-8383-4a04-8304-9a8ba75a7b93             default/echoserver-0002
   5  DiscoveryResponse            2  ada7afc4-cdc7-44ae-a229-79e432b3adb1-1  1        default/echoserver-0002  ['10.244.1.13']
   6  DiscoveryRequest             2  ada7afc4-cdc7-44ae-a229-79e432b3adb1-1  1        default/echoserver-0002
  16  DiscoveryRequest             6  ee0e8167-8383-4a04-8304-9a8ba75a7b93             default/echoserver-0001
  17  DiscoveryResponse            6  ada7afc4-cdc7-44ae-a229-79e432b3adb1-1  1        default/echoserver-0001  ['10.244.1.13']
  18  DiscoveryRequest             6  ada7afc4-cdc7-44ae-a229-79e432b3adb1-1  1        default/echoserver-0001
  19  DiscoveryRequest             7  ee0e8167-8383-4a04-8304-9a8ba75a7b93             default/echoserver-0003
  20  DiscoveryResponse            7  ada7afc4-cdc7-44ae-a229-79e432b3adb1-1  1        default/echoserver-0003  ['10.244.1.13']
  21  DiscoveryRequest             7  ada7afc4-cdc7-44ae-a229-79e432b3adb1-1  1        default/echoserver-0003
