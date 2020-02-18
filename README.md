# grpc-loadbalance-on-kubernetes-and-istio
这个go demo用来测试grpc-gateway在kubernetes-istio这个platform下的负载均衡效果.

## 环境依赖:
这部分需要什么大家就go get吧

## 部署流程：
### proto编译
包括grpc-gateway, 代码中已编译好上传了.

```bash
cd gitlab.bj.sensetime.com/SenseGo/grpc-gateway-demo/proto
protoc --go_out=plugins=grpc:. helloworld.proto
cd gitlab.bj.sensetime.com/SenseGo/grpc-gateway-demo/grpc-gateway
protoc --go_out=plugins=grpc:. gateway.proto
cp gateway.pb.gw.go ../proto/
```

### 服务端编译部署:

```bash
go build -o progress-server main.go
docker build . -t localhost:5000/grpc-server-istio:v0.5
kubectl apply -f grpc-server.yaml
kubectl apply -f grpc-server-svc.yaml
```

### 客户端编译部署:

```bash
go build -o progress-client main.go
docker build . -t localhost:5000/grpc-client-istio:v0.5
kubectl apply -f grpc-client.yaml
```


## 测试流程:
### 一个服务端，一个客户端
很正常
![](https://raw.githubusercontent.com/zhoushuke/BlogPhoto/master/githuboss/grpc-on-kubernetes-loadbalance00.png)

### 服务端个数一个，客户端从一个扩容到三个
这时在服务端显示有3个客户端的请求
![](https://raw.githubusercontent.com/zhoushuke/BlogPhoto/master/githuboss/grpc-on-kubernetes-loadbalance01.png)

### 客户端个数三个, 服务端从一个扩容至五个
这下图可以从客户端发现，请求被均匀地分配给了服务端的5个地址, 说明，通过istio是可以实现grpc服务的负载均衡效果的
![](https://raw.githubusercontent.com/zhoushuke/BlogPhoto/master/githuboss/grpc-on-kubernetes-loadbalance03.png)

### 服务端从五个缩容至一个
客户端也能实时地发现服务的变动, 如上图下半部分所示
