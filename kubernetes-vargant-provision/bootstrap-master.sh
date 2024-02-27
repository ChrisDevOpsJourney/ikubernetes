#!/bin/bash

#echo "[TASK 1] Pull required containers"
#kubeadm config images pull >/dev/null 2>&1
#
#echo "[TASK 2] Initialize Kubernetes Cluster"
#kubeadm init --apiserver-advertise-address=192.168.56.100 --pod-network-cidr=172.16.16.0/16 >> /root/kubeinit.log 2>/dev/null
#
#echo "[TASK 3] Deploy Calico network"
#kubectl --kubeconfig=/etc/kubernetes/admin.conf create -f https://docs.projectcalico.org/v3.18/manifests/calico.yaml >/dev/null 2>&1
#
#echo "[TASK 4] Generate and save cluster join command to /joincluster.sh"
#kubeadm token create --print-join-command > /joincluster.sh 2>/dev/null
#
#echo "[TASK 5] Setup kubeconfig"
#mkdir -p ~/.kube/
#cp /etc/kubernetes/admin.conf ~/.kube/config
#
#echo "[TASK 6] Setup Command Alias"
#echo "alias crictl='crictl -runtime-endpoint unix:///run/containerd/containerd.sock'" >> ~/.bashrc
#echo "alias k='kubectl'" >> ~/.bashrc
#echo "alias kgp='kubectl get pods'" >> ~/.bashrc
#echo "alias kgn='kubectl get nodes'" >> ~/.bashrc
