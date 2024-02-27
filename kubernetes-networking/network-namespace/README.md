# Understanding Container Networking using Linux Network Namespaces to isolate the server

## What is network namespace
network namespace is a Linux kernel feature that allows us to isolate network environments through virtualization. for example, using a network namespace you can create separate network interface and routing tables that are isolated from the rest of the system and operate independently.
Network Namespace is a core component of the docker network.

network namespace provides isolation of the system resources associated with networking
- network devices
- ipv4/ipv6 protocol stacks
- IP routing tables
- Firewall Rules
- the /proc/net directory
- the /sys/class/net directory
- various files under /proc/sys/net
- port numbers(sockets)

## Virtual Interfaces and Bridges:
virtual interfaces provide us with virtualized representations of physical network interfaces and the Bridge gives the virtual equivalent of a Switch.

In this demo, we are going to create 
- 2 network namespaces(like 2 isolated servers)
- 2 veth pairs(like 2 physical ethernet cables)
- 1 bridge (for routing traffic between namespaces)

Demo Objective
- configure the bridge to allow 2 namespaces to communicate with each other
- connect the bridge to the host and internet
- configure for incoming traffic(outside) to the namespace

## Demo
- List all network interface
```sh
sudo ip link show
```

Step 1.1: Create 2 network namespaces
```sh
# Add two network namespaces using "ip netns" command
sudo ip netns add NS1
sudo ip netns add NS2
```

```sh
# List the created network namespaces
sudo ip netns show

# By convention, network namespace handles created by
# iproute2 live under `/var/run/netns`
sudo ls /var/run/netns/
```

Step 1.2: create veth pairs
```sh
sudo ip link add veth10 type veth peer name veth11
sudo ip link add veth20 type veth peer name veth21

ip link show type veth
```

Step 1.3: add the veth paris to the namespace
```sh
sudo ip link set veth11 netns NS1
sudo ip link set veth21 netns NS2
```

Step 1.4: Configure the interfaces in the network namespaces with IP address
```sh
sudo ip netns exec NS1 ip addr add 172.16.0.2/24 dev veth11 
sudo ip netns exec NS2 ip addr add 172.16.0.3/24 dev veth21
```

Step 1.5: Enable the interfaces in the network namespaces
```sh
sudo ip netns exec NS1 ip link set dev veth11 up
sudo ip netns exec NS2 ip link set dev veth21 up
```

Step 1.5: Create a bridge
```sh
sudo ip link add br0 type bridge

sudo ip link show type bridge
sudo ip link show br0
```

Step 1.6: Add the network interface to the bridge
```sh
sudo ip link set dev veth10 master br0
sudo ip link set dev veth20 master br0
```

Step 1.7: Configure IP for the bridge
```sh
sudo ip addr add 172.16.0.1/24 dev br0
```

Step 1.8: Enable the Bridge and the interface connected to the bridge
```sh
sudo ip link set dev br0 up
sudo ip link set dev veth10 up
sudo ip link set dev veth20 up
```

Step 1.9: Enable the loopback interfaces in the network namespaces
```sh
sudo ip netns exec NS1 ip link set lo up
sudo ip netns exec NS2 ip link set lo up

sudo ip netns exec NS1 ip a
sudo ip netns exec NS2 ip a
```

## Verify connectivity between two netns
```s
# We can log in to netns environment using the below; It will be isolated from any other network
sudo nsenter --net=/var/run/netns/NS1

#Ping adaptor attached to NS1
ping -W 1 -c 2 172.16.0.2

#Ping the bridge
ping -W 1 -c 2 172.16.0.1

#Ping the adaptor in NS2
ping -W 1 -c 2 172.16.0.3
```

## Verify connectivity from netns to internet
```sh
sudo ip netns exec NS1 ping -W 1 -c 2 10.20.3.7
```
> ping: connect: Network is unreachable

```sh
# Check the route inside NS1
sudo ip netns exec NS1 route -n
```
![image](https://github.com/ChrisDevOpsOrg/linux-namespace/assets/54896350/096f13ed-fde3-49ef-ba25-0597a5b63177)
> As we can see, no route is defined to carry other traffic than 172.16.0.0/24
> We can fix this by adding a default route


```sh
# Setting the default route in the network namespaces
sudo ip netns exec NS1 ip route add default via 172.16.0.1 dev veth11
sudo ip netns exec NS2 ip route add default via 172.16.0.1 dev veth21

sudo ip netns exec NS1 route -n
```
![image](https://github.com/ChrisDevOpsOrg/linux-namespace/assets/54896350/c04c0287-1cbc-49f9-8005-ee0400d16e58)


```sh
# Now first ping the host machine eth0
sudo ip netns exec NS1 ping -W 1 -c 2 10.20.3.7

# Now ping 8.8.8.8 again
sudo ip netns exec NS1 ping 8.8.8.8
# still unreachable
```
```sh
# open tcpdump in eth0 to see the packet
sudo tcpdump -i eth0 icmp
```
![image](https://github.com/ChrisDevOpsOrg/linux-namespace/assets/54896350/831a39d4-c5ea-45c8-8361-28d508758b70)
> no packet captured, let's capture traffic for br0

```sh
sudo tcpdump -i br0 icmp
```
![image](https://github.com/ChrisDevOpsOrg/linux-namespace/assets/54896350/b4251edf-bca1-4f89-9639-dfbb4b9f9712)
> We can see the traffic at br0 but we don't get a response from eth0.
> It's because of an IP forwarding issue
> enabling ip forwarding by changing value 0 to 1

```sh
sudo sysctl -w net.ipv4.ip_forward=1
sudo cat /proc/sys/net/ipv4/ip_forward
sudo tcpdump -i eth0 icmp
```
![image](https://github.com/ChrisDevOpsOrg/linux-namespace/assets/54896350/8b3984ac-0bf8-4443-86c4-bdda0e4a830c)
>|  As we can see how we are getting responses eth0 but ping 8.8.8.8 still not working. Although the network is now reachable, there’s no way that we can have responses back - cause packets from external networks 
can’t be sent directly to our `172.16.0.0/24` network.

```sh
sudo iptables \
        -t nat \
        -A POSTROUTING \
        -s 172.16.0.0/24 ! -o br0 \
        -j MASQUERADE
# -t specifies the table to which the commands
# should be directed to. By default, it's `filter`.
# -A specifies that we're appending a rule to the
# chain then we tell the name after it;
# -s specifies a source address (with a mask in this case).
# -j specifies the target to jump to (what action to take).

# Now we're getting a response from Google DNS
sudo ip netns exec NS1 ping -c 2 8.8.8.8
```
![image](https://github.com/ChrisDevOpsOrg/linux-namespace/assets/54896350/3791022f-9db4-4c4f-b774-5759586c1fbd)

## open a service in one of the namespaces and try to get a response from outside
```sh
sudo nsenter --net=/var/run/netns/NS1
python3 -m http.server --bind 172.16.0.2 3000
```
Access server via host public IP
```sh
telnet YOUR_HOST_PUBLIC_IP:3000
```
![image](https://github.com/ChrisDevOpsOrg/linux-namespace/assets/54896350/fc362ffd-420a-47bc-8e84-aaea4a455e3e)
>| As we can see we can't reach the destination. Because we didn't tell the Host machine where to put the incoming traffic. We have to NAT again, this time we will define the destination

```sh
sudo iptables \
        -t nat \
        -A PREROUTING \
        -d 10.20.3.7 \
        -p tcp -m tcp --dport 3000 \
        -j DNAT --to-destination 172.16.0.2:3000
```
![image](https://github.com/ChrisDevOpsOrg/linux-namespace/assets/54896350/452258d0-d37c-4a62-bfde-1ddef09b9012)





