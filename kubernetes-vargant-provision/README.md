# kubernetes-vargant-provision

---

Use to provision Kubernetes cluster


# How to use

---
## Prerequisite
- make sure virtualbox installed
- make sure vagrant installed
  https://developer.hashicorp.com/vagrant/docs/installation
- configure virtualbox network to use 192.168.56.0/24


## Operation

---
- Provision kubernetes on your local
```shell
vagrant up
```

- check VM status
```shell
vagrant status
```

- stop cluster
```shell
vagrant halt
```

- destroy cluster
```shell
vagrant destroy
```

- connect to master node
```shell
ssh root@192.168.56.100
```

- check node status
```shell
kgn
```

- check pods
```shell
kgp -n kube-system
```

### Cluster Node Network

---
- master node: 192.168.56.100
- worker node1: 192.168.56.101
- worker node2: 192.168.56.102


