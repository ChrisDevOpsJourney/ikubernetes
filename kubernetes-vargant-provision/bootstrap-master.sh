#!/bin/bash

echo "[TASK 1] Pull required containers"
kubeadm config images pull >/dev/null 2>&1

echo "[TASK 2] Initialize Kubernetes Cluster"

# Uncomment for Calico CNI
# kubeadm init --apiserver-advertise-address=192.168.10.100 --pod-network-cidr=172.16.16.0/16 >> /root/kubeinit.log 2>/dev/null

# Use cilium CNI
sudo kubeadm init --apiserver-advertise-address=192.168.10.100 --service-cidr=172.16.20.0/16 --skip-phases=addon/kube-proxy >> /root/kubeinit.log 2>/dev/null

# Uncomment for Calico CNI
# echo "[TASK 3] Deploy Calico network"
# kubectl --kubeconfig=/etc/kubernetes/admin.conf create -f kubectl create -f https://raw.githubusercontent.com/projectcalico/calico/v3.28.2/manifests/tigera-operator.yaml >/dev/null
# kubectl --kubeconfig=/etc/kubernetes/admin.conf create -f https://raw.githubusercontent.com/projectcalico/calico/v3.27.0/manifests/custom-resources.yaml >/dev/null

echo "[TASK 4] Generate and save cluster join command to /joincluster.sh"
kubeadm token create --print-join-command > /joincluster.sh

echo "[TASK 5] Setup kubeconfig"
mkdir -p ~/.kube/
cp /etc/kubernetes/admin.conf ~/.kube/config

echo "[TASK 6] Setup Command Alias"
echo "alias crictl='crictl -runtime-endpoint unix:///run/containerd/containerd.sock'" >> ~/.bashrc
echo "alias k='kubectl'" >> ~/.bashrc

echo "[TASK 7] Setup Helm"
curl -fsSL -o get_helm.sh https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 >/dev/null 2>&1
chmod 700 get_helm.sh
./get_helm.sh

echo "[TASK 8] Setup Cilium"
curl -L --remote-name-all https://github.com/cilium/cilium-cli/releases/latest/download/cilium-linux-amd64.tar.gz{,.sha256sum} >/dev/null 2>&1
sha256sum --check cilium-linux-amd64.tar.gz.sha256sum
sudo tar xzvfC cilium-linux-amd64.tar.gz /usr/local/bin
rm cilium-linux-amd64.tar.gz{,.sha256sum}

echo "[TASK 9] add cilium helm repo"
helm repo add cilium https://helm.cilium.io/

echo "[TASK 10] Install cilium"
helm install cilium cilium/cilium \
    --namespace kube-system \
    --set kubeProxyReplacement=true \
    --set k8sServiceHost=192.168.10.100 \
    --set k8sServicePort=6443 \
    --set ipam.operator.clusterPoolIPv4PodCIDRList="10.18.0.0/16"
