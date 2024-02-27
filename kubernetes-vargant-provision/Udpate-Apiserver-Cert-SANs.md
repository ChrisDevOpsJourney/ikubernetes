
```shell
# Check valid SANs in the existing cert
openssl x509 -in /etc/kubernetes/pki/apiserver.crt -noout -text
```

```shell
kubectl -n kube-system get configmap kubeadm-config -o jsonpath='{.data.ClusterConfiguration}' > kubeadm.yaml
```

Update kubeadm.yaml
```yaml
apiServer:
  certSANs:
  - "kmaster"
  - "kubernetes"
  - "kubernetes.default"
  - "kubernetes.default.svc"
  - "kubernetes.default.svc.cluster.local"
  - "10.96.0.1"
  - "172.16.0.4"
  - "Another-IP"
  extraArgs:
    authorization-mode: Node,RBAC
  timeoutForControlPlane: 4m0s
apiVersion: kubeadm.k8s.io/v1beta3
certificatesDir: /etc/kubernetes/pki
clusterName: kubernetes
controllerManager: {}
dns: {}
etcd:
  local:
    dataDir: /var/lib/etcd
imageRepository: registry.k8s.io
kind: ClusterConfiguration
kubernetesVersion: v1.27.3
networking:
  dnsDomain: cluster.local
  podSubnet: 192.168.0.0/16
  serviceSubnet: 10.96.0.0/12
scheduler: {}
```

```shell
mv /etc/kubernetes/pki/apiserver.{crt,key} /tmp
```

```shell
kubeadm init phase certs apiserver --config kubeadm.yaml
```
