## 安装工具库
```shell
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install github.com/go-micro/generator/cmd/protoc-gen-micro@latest
```

## 创建项目并安装依赖
```shell
mkdir -p example
cd example
go mod init github.com/killlowkey/go-micro-example

go get google.golang.org/protobuf
go get go-micro.dev/v4
```

## 编写 proto 文件
在 proto 目录下创建 greeter.proto 文件，并填充如下内容
```proto
syntax = "proto3";

package pb;
option go_package = "pkg/proto;pb";

service Greeter {
	rpc Hello(Request) returns (Response) {}
}

message Request {
	string name = 1;
}

message Response {
	string msg = 1;
}
```

## 生成代码
```shell
protoc --proto_path=. --micro_out=. --go_out=. proto/greeter.proto
```

## 服务端代码
cmd/server/main.go
```go
package main

import (
	"context"
	"github.com/killlowkey/go-micro-example/pkg/proto"
	"go-micro.dev/v4"
	"google.golang.org/grpc"
	"log"
	"time"
)

type Greeter struct{}

func (g *Greeter) Hello(ctx context.Context, req *pb.Request, rsp *pb.Response) error {
	rsp.Msg = "Hello " + req.Name
	return nil
}

func main() {
	go func() {
		for {
			grpc.DialContext(context.TODO(), "127.0.0.1:9091")
			time.Sleep(time.Second)
		}
	}()

	service := micro.NewService(
		micro.Name("go.micro.srv.greeter"),
	)

	// optionally setup command line usage
	service.Init()

	// Register Handlers
	if err := pb.RegisterGreeterHandler(service.Server(), &Greeter{}); err != nil {
		log.Fatal(err)
	}

	// Run server
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
```

## 客户端代码
```go
package main

import (
	"context"
	"fmt"

	"github.com/killlowkey/go-micro-example/pkg/proto"
	"go-micro.dev/v4"
)

func main() {
	// create a new service
	service := micro.NewService()

	// parse command line flags
	service.Init()

	// Use the generated client stub
	cl := pb.NewGreeterService("go.micro.srv.greeter", service.Client())

	// Make request
	rsp, err := cl.Hello(context.Background(), &pb.Request{
		Name: "John",
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(rsp.Msg)
}
```
## 运行服务端
> 运行之前执行 go mod tidy
```shell
go run cmd/server/main.go
```

## 运行客户端
> 运行之前执行 go mod tidy
```shell
go run cmd/client/main.go
```