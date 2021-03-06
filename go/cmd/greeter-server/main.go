package main

import (
	"context"
	"errors"
	"log"
	"net"
	"os"
	"time"

	"github.com/115100/poc-client-side-grpc-lb/go/greeterpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

func main() {
	if err := run(); err != nil {
		log.Fatalln(err)
	}
}

func run() error {
	// assume pod IP is static
	ip := os.Getenv("MY_POD_IP")
	if ip == "" {
		return errors.New("run: MY_POD_IP is not set")
	}

	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		return err
	}
	defer lis.Close()

	grpcServer := grpc.NewServer(grpc.KeepaliveParams(
		keepalive.ServerParameters{
			// this is needed to force the client to refresh DNS
			// lookups every ten minutes. It will be 30m otherwise.
			// See: https://github.com/grpc/grpc-go/issues/1663#issuecomment-358698804
			MaxConnectionAge: 10 * time.Minute,
		},
	))
	greeterpb.RegisterGreeterServer(grpcServer, &greeterServer{ip: ip})

	log.Println("server: Starting to run port :8080")
	return grpcServer.Serve(lis)
}

type greeterServer struct{ ip string }

func (s *greeterServer) Greet(_ context.Context, req *greeterpb.GreetRequest) (*greeterpb.GreetReply, error) {
	return &greeterpb.GreetReply{
		Message: "my IP is " + s.ip,
	}, nil
}
