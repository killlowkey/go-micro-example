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
