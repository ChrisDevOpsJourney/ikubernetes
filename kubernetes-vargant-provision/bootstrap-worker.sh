#!/bin/bash

#echo "[TASK 1] Join node to Kubernetes Cluster"
#apt install -qq -y sshpass >/dev/null 2>&1
#sshpass -p "kubeadmin" scp -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no kmaster.example.com:/joincluster.sh /joincluster.sh 2>/dev/null
#bash /joincluster.sh >/dev/null 2>&1
#
#echo "[TASK 2] Setup Command Alias"
#echo "alias crictl='crictl -runtime-endpoint unix:///run/containerd/containerd.sock'" >> ~/.bashrc
#echo "alias k='kubectl'" >> ~/.bashrc
#echo "alias kgp='kubectl get pods'" >> ~/.bashrc
#echo "alias kgn='kubectl get nodes'" >> ~/.bashrc