Environment Info
- Ubuntu22.04 - 5.15.0-97-generic
- kubernetes 1.29
- ipvsadm v1.31

```shell
export master='10.20.3.7'
export node1='10.20.3.8'
```

```shell
# update the server
sudo apt-get update -y
sudo apt-get upgrade -y
sudo apt-get install -y apt-transport-https ca-certificates curl gpg
```

```shell
# If the folder `/etc/apt/keyrings` does not exist, it should be created before the curl command, read the note below.
# sudo mkdir -p -m 755 /etc/apt/keyrings
curl -fsSL https://pkgs.k8s.io/core:/stable:/v1.29/deb/Release.key | sudo gpg --dearmor -o /etc/apt/keyrings/kubernetes-apt-keyring.gpg

```

```shell
# This overwrites any existing configuration in /etc/apt/sources.list.d/kubernetes.list
echo 'deb [signed-by=/etc/apt/keyrings/kubernetes-apt-keyring.gpg] https://pkgs.k8s.io/core:/stable:/v1.29/deb/ /' | sudo tee /etc/apt/sources.list.d/kubernetes.list

```
```shell
#Install all of the necessary Kubernetes components with the command:
sudo apt-get update
sudo apt-get install kubeadm kubelet kubectl ipvsadm containerd -y
```
```shell
#Configure containerd and start the service
sudo mkdir -p /etc/containerd
sudo su -
    containerd config default > /etc/containerd/config.toml
    sed -E -i 's/(SystemdCgroup = ).*/\1true/g' config.toml
    systemctl restart containerd
exit
```

```shell
# Forwarding IPv4 and letting iptables see bridged traffic
cat <<EOF | sudo tee /etc/modules-load.d/k8s.conf
overlay
br_netfilter
EOF

sudo modprobe overlay
sudo modprobe br_netfilter
```

```shell
# sysctl params required by setup, params persist across reboots
cat <<EOF | sudo tee /etc/sysctl.d/k8s.conf
net.bridge.bridge-nf-call-iptables  = 1
net.bridge.bridge-nf-call-ip6tables = 1
net.ipv4.ip_forward                 = 1
EOF
```

```shell
# Apply sysctl params without reboot
sudo sysctl --system
```

```shell
# Verify that the br_netfilter, overlay modules are loaded
lsmod | grep br_netfilter
lsmod | grep overlay
```

```shell
#Pull the necessary containers with the command:
sudo kubeadm config images pull
```

**########## This section must be run only on master node ######### **
```shell
#Make sure "kube-proxy" is not installed, we want cilium to use the new "eBPF" based proxy
sudo kubeadm init --skip-phases=addon/kube-proxy

```

```shell
# setup kubeconfig
mkdir -p $HOME/.kube
sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
sudo chown $(id -u):$(id -g) $HOME/.kube/config
```

```shell
#Install cilium CLI
curl -L --remote-name-all https://github.com/cilium/cilium-cli/releases/latest/download/cilium-linux-amd64.tar.gz{,.sha256sum}
sha256sum --check cilium-linux-amd64.tar.gz.sha256sum
sudo tar xzvfC cilium-linux-amd64.tar.gz /usr/local/bin
rm cilium-linux-amd64.tar.gz{,.sha256sum}
```

```shell
#Setup Helm repository
helm repo add cilium https://helm.cilium.io/

#Deploy Cilium release via Helm:
helm install cilium cilium/cilium 
    --namespace kube-system \
    --set kubeProxyReplacement=strict \
    --set k8sServiceHost=@$master \
    --set k8sServicePort=6443
```

```shell
kubectl -n kube-system get pods -l k8s-app=cilium -o wide
MASTER_CILIUM_POD=$(kubectl -n kube-system get pods -l k8s-app=cilium -o wide |  grep master | awk '{ print $1}' )
echo $MASTER_CILIUM_POD

#validate that the Cilium agent is running in the desired mode (non kube-proxy)
kubectl exec -it -n kube-system $MASTER_CILIUM_POD -- cilium status | grep KubeProxyReplacement

#Validate that Cilium installation
cilium status --wait

#Review what network interfaces Cilium has created
ip link show
```

```shell
#Schedule a Kubernetes deployment using a container from Google samples
kubectl create deployment hello-world --image=gcr.io/google-samples/hello-app:1.0

#Scale up the replica set to 4
kubectl scale --replicas=4 deployment/hello-world

#Get pod info
kubectl get pods -o wide

kubectl expose deployment hello-world --port=8080 --target-port=8080 --type=NodePort
kubectl get service hello-world

kubectl exec -it -n kube-system $MASTER_CILIUM_POD -- cilium service list

#Verify that iptables are not used
sudo iptables-save | grep KUBE-SVC

export CLUSTERIP=$(kubectl get service hello-world  -o jsonpath='{ .spec.clusterIP }')
echo $CLUSTERIP

PORT=$( kubectl get service hello-world  -o jsonpath='{.spec.ports[0].port}')
echo $PORT

curl http://$CLUSTERIP:$PORT

NODEPORT=$( kubectl get service hello-world  -o jsonpath='{.spec.ports[0].nodePort}')
echo $NODEPORT

curl http://$master:$NODEPORT

```

**########## Setup Hubble**########## 

```shell
cilium hubble enable

#Enabling Hubble requires the TCP port 4245 to be open on all nodes running Cilium. This is required for Relay to operate correctly.

cilium status

#In order to access the observability data collected by Hubble, install the Hubble CL
export HUBBLE_VERSION=$(curl -s https://raw.githubusercontent.com/cilium/hubble/master/stable.txt)
curl -L --remote-name-all https://github.com/cilium/hubble/releases/download/$HUBBLE_VERSION/hubble-linux-amd64.tar.gz{,.sha256sum}
sha256sum --check hubble-linux-amd64.tar.gz.sha256sum
sudo tar xzvfC hubble-linux-amd64.tar.gz /usr/local/bin
rm hubble-linux-amd64.tar.gz{,.sha256sum}

#In order to access the Hubble API, create a port forward to the Hubble service from your local machine
cilium hubble port-forward&

hubble status
#If you get "Unavailable Nodes: ", follow below troubleshooting:
######Hubbel trouble shooting####

    #Get resolution from: https://github.com/cilium/hubble/issues/599
    kubectl delete secrets -n kube-system cilium-ca
    kubectl get secrets -n kube-system hubble-ca-secret -o yaml | sed -e 's/name: hubble-ca-secret/name: cilium-ca/;/\(resourceVersion\|uid\)/d' | kubectl apply -f -
    cilium hubble disable
    cilium hubble enable
    #Please note that the next time the hubble-generate-certs CronJob runs, 
    #it will override the TLS certificates for both Hubble and Relay signing them with hubble-ca-secret (i.e. not ciliium-ca). 
    #Relay should continue to work, but this could bring more incompatibility with the CLI (e.g. if you were to disable then re-enable Hubble again through the CLI).
    cilium hubble port-forward&
    hubble status
    hubble observe

#Setup Hubble UI
cilium hubble enable --ui

cilium hubble ui
```
