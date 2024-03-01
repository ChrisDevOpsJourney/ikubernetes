Environment Info
- Ubuntu22.04 - 5.15.0-97-generic
- kubernetes 1.29
- ipvsadm v1.31


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
kubeadm config print init-defaults > kubeadm.yaml

# Needs to do some modification before apply, can refer to kubeadm.yaml that in the repo
sudo kubeadm init  --config kubeadm.yaml

```

```shell
# setup kubeconfig
mkdir -p $HOME/.kube
sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
sudo chown $(id -u):$(id -g) $HOME/.kube/config
```

```shell
#Download Calico CNI
curl https://raw.githubusercontent.com/projectcalico/calico/v3.27.2/manifests/calico.yaml -O
#Apply Calico CNI
kubectl apply -f ./calico.yaml
```

```shell
#Get cluster info
kubectl cluster-info
```

```shell
#List the virtual server table 
sudo ipvsadm -L
```

```shell
#Schedule a Kubernetes deployment using a container from Google samples
kubectl create deployment hello-world --image=gcr.io/google-samples/hello-app:1.0

#View all Kubernetes deployments
kubectl get deployments

kubectl expose deployment hello-world --port=8090 --target-port=8080 

kubectl get services

kubectl get pods -o wide

#List the virtual server table
sudo ipvsadm -L

kubectl scale --replicas=4 deployment/hello-world

#List the virtual server table
sudo ipvsadm -L
```
