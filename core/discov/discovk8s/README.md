# 设计思路

支持通过Kubernetes Service进行服务发现，方便集成现有服务。<br />
每个k8s Service都关联有唯一的Endpoints对象，保存了所有ready和not ready的Pod，
deployment扩缩容时，会实时更新Endpoints下的地址列表，使用k8s的informer sdk 
来watch我们感兴趣的Endpoints，再更新到gRPC。



## 配置方法

复用Endpoints配置项，要求数组长度为1且格式如下
```yaml
Transform:
   Endpoints:
   - k8s:///transform-svc.ns:8081
```

## RBAC

Pod需要具有读取Endpoints权限，需执行以下配置

```yaml
# 一个K8s集群仅需配置一次
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: endpoints-reader-role
rules:
- apiGroups: [""]
  resources: ["endpoints"]
  verbs: ["get", "watch", "list"]

--- 

#  分namespace执行，以default举例
apiVersion: v1
kind: ServiceAccount
metadata:
  name: endpoints-reader-sa
  namespace: default

---

#  分namespace执行，以default举例
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: endpoints-reader-binding
subjects:
- kind: ServiceAccount
  name: endpoints-reader-sa
  namespace: default
roleRef:
  kind: ClusterRole
  name: endpoints-reader-role
  apiGroup: rbac.authorization.k8s.io

---

# deployments配置，注意其中serviceAccountName字段，此时该Deployment下的Pod便具有访问Endpoints的权限了
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
  labels:
    app: nginx
spec:
  replicas: 1
  selector:
    matchLabels:
      app: nginx
  template:
    metadata:
      labels:
        app: nginx
    spec:
      serviceAccountName: endpoints-reader-sa
      containers:
      - name: nginx
        image: nginx:1.14.2
        ports:
        - containerPort: 80
```