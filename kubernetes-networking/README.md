# Linux Namespace
Linux Namespace serves as an abstraction layer over operating system resources. Virtualize a namespace as a container that encapsulates specific system resources, with each type of namespace representing a distinct box.
Namespaces are a feature of the Linux kernel that partitions kernel resources such that one set of processes sees one set of resources and another set of processes sees a different set of resources. The feature works by having the same namespace for a group of resources and processes, but those namespaces refer to distinct resources

There are 7 types of namespaces.
- Cgroup
- IPC(Inter-Process Communication)
- Network
- Mount
- PID
- User
- UTS(Unix Time-Sharing)
