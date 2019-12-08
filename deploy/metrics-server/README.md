# metrics-server for kutti

Kutti uses kubeadm as its setup machanism. The [metrics-server](https://kubernetes.io/docs/tasks/debug-application-cluster/resource-metrics-pipeline/) story for kubeadm is far from perfect as of December 2019 - so kutti has to cope.

As of now, kutti tracks the [metrics server repo](https://github.com/kubernetes-incubator/metrics-server), and copies the **1.8+** manifest files from there. The following additional command-line arguments are added to the deployment:
* --kubelet-insecure-tls
* --kubelet-preferred-address-types=InternalIP
