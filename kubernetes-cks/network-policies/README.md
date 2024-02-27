# Network Policies

- Firewall Rule in Kubernetes
- Implemented by the network plugin CNI
- Namespace level
- Restrict the Ingress/Egress for a group of pods based on certain rules and conditions



# Without Network Policies
- By default, every pod can access every pod
- Pods are not isolated 

# Example Code 
- Pod Selector
```yaml
apiVersion: v1
kind: NetworkPolicy
metadata:
  name: pod-selector-example
  namespace: default
spec:
  podSelector:
    matchLabels:
      id: frontend
    policyType:
      - Egress
```
What does this policy do?
- Implicitly deny all outgoing traffic from pods with label id=frontend in namespace default

```yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: example
  namespace: default
spec:
  podSelector:
    matchLabels:
      id: frontend
  policyTypes:
    - Egress
  egress:
    - to: 
        - namespaceSelector:
          - matchLabels:
              id: ns1
      ports:
          - protocol: TCP
            port: 80
    - to:
        - podSelector:
            matchLabels:
              id: backend
```
What does this policy do?
- Explicitly allow the pods with label id=frontend to communicate with pods in namespace which it labels id=ns1 through port 80
- Explicitly allow the pods with label id=frontend to communicate with pods which it labels id=backend in the **same namespace**.

# Multiple Network Policies
- It is possible to have multiple NPs selecting the same pods
- If the pod has more than one NP
  - union of all NPs are applied
  - order does not affect a policy result.

# Practice - Default Deny
```shell
kubectl run frontend --image=nginx
kubectl run backend --image=nginx

kubectl expose pod frontend --port 80
kubectl expose pod backedn --port 80

kubectl exec frontend curl backend  # You will be able to access backend from frontend
kubectl exec backend curl frontend  # You will be able to access frontend from backend
```
Let's apply network policy
```yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: default-deny
  namespace: default
spec:
  podSelector: {}
  policyTypes:
    - Ingress
    - Egress
```
# Practice — Allow frontend pods to talk to backend pods based on pod selector
```yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: frontend
  namespace: default
spec:
  podSelector:
    matchLabels:
      run: frontend
  policyTypes:
    - Egress
  egress:
    - to:
        - podSelector:
            matchLabels:
              run: backend
---
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: backend
  namespace: default
spec:
  podSelector:
    matchLabels:
      run: backend
  policyTypes:
    - Ingress
  ingress:
    - from:
        - podSelector:
            matchLabels:
              run: frontend
```
when we try to access backend pods from frontend. We are still unable to access. 
This is because the default denied policy also denies DNS traffic on port 53.

To Allow DNS Resolutions
```yaml
# deny all incoming and outgoing traffic from all pods in namespace default
# but allows DNS traffic. 
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: deny
  namespace: default
spec:
  podSelector: {}
  policyTypes:
    - Ingress
    - Egress
  egress:
    - to:
      ports:
        - port: 53
          protocol: TCP
        - port: 53
          protocol: UDP
```