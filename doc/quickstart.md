# 快速开始

0.配置

```yaml
app:
  name: 服务名称
  version: v0.0.1
server:
  http:
    addr: :8080
    readTimeout: 1s
    writeTimeout: 1s
  grpc:
    addr: :9090
client:
  srv-dict:
    debug: true
    address: :9090
    block: false
    dialTimeout: 3s
    readTimeout: 1s
```

1.创建一个Application

```go
type MyApplication struct {
    frame.Application
}

func main() {
    app := &MyApplication{}
    app.Startup()
    app.Run()
}
```

2.为Application添加一个http server

```go
func(app *MyApplication) startHTTPServer() error {
    server := xmux.Build()
    server.Handle("/list", httpTransport.NewServer(ListEndpoint, decodeListRequest, encode.JsonResponse, opts...)).Methods(http.MethodGet)
    return app.Serve(server)
}

func main() {
    app := &MyApplication{}
    app.Startup(
        app.startHTTPServer,
    )
    app.Run()
}
```

3.为Application添加一个grpc server

```go
func(app *MyApplication) startGRPCServer() error {
    server := xgrpc.Build()
    helloworld.RegisterGreeterServer(server.Server, new(greeter.Greeter))
    return app.Serve(server)
}

func main() {
    app := &MyApplication{}
    app.Startup(
        app.startGRPCServer,
    )
    app.Run()
}
```


4.创建一个grpc客户端

直连服务器:
```go
import (
	"context"
	pinyinClient "esfgit.leju.com/golang/dict/pkg/pinyin/client"
	clientgrpc "github.com/Remember9/frame/client/grpc"
)

func dialServer() {
    conn := clientgrpc.Build("srv-dict")

    client := pinyinClient.NewPinyinGRPCClient(conn)
    resp, err := client.Get(context.Background(), "汉字测试", "", "")

}
```

5.监控服务
```
server:
  govern:
    isuse: true
    addr: :8090
```

所有可访问的路由 `http://127.0.0.1:8090/routes`
实时监控 `http://127.0.0.1:8090/debug/statsviz/`





