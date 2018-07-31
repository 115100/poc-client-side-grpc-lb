# POC for client-side load-balancing with grpc-go

This includes everything required to test client-side LB on your own Kubernetes cluster.

Inspired by the [C# example](https://github.com/jtattermusch/grpc-loadbalancing-kubernetes-examples) by jtattermusch.

## Concepts

- Headless services in Kubernetes return all available backends for a service (https://kubernetes.io/docs/concepts/services-networking/dns-pod-service/).

- Lookups on a DNS record that return multiple A records will be called in a round-robin manner if the `round_robin` Balancer is used alongside the `dns:///` scheme.

## Running:

Ensure you have a sensible context set, e.g. a non-production cluster.

```bash
$ kubectl apply -f manifests/*.yaml
$ kubectl logs -f greeter-client-deployment-[...]
2018/07/31 14:09:19 client: Starting to ping server
2018/07/31 14:09:21 client: got response from server: "my IP is 10.1.0.20"
2018/07/31 14:09:22 client: got response from server: "my IP is 10.1.0.17"
2018/07/31 14:09:23 client: got response from server: "my IP is 10.1.0.18"
2018/07/31 14:09:24 client: got response from server: "my IP is 10.1.0.20"
2018/07/31 14:09:25 client: got response from server: "my IP is 10.1.0.17"
2018/07/31 14:09:26 client: got response from server: "my IP is 10.1.0.18"
[...]
```
