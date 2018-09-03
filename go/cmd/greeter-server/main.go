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
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/keepalive"
)

func main() {
	if err := run(); err != nil {
		log.Fatalln(err)
	}
}

func run() error {
	ctx := context.Background()

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
	srv := &greeterServer{ip: ip}
	greeterpb.RegisterGreeterServer(grpcServer, srv)

	hs := health.NewServer()
	go updateHealthServer(ctx, hs, srv)
	healthpb.RegisterHealthServer(grpcServer, hs)

	log.Println("server: Starting to run port :8080")
	return grpcServer.Serve(lis)
}

// -----------------------------------------------------------------------------

type greeterServer struct{ ip string }

func (s *greeterServer) Greet(_ context.Context, req *greeterpb.GreetRequest) (*greeterpb.GreetReply, error) {
	return &greeterpb.GreetReply{
		Message: "my IP is " + s.ip,
	}, nil
}

// -----------------------------------------------------------------------------

func updateHealthServer(ctx context.Context, hs *health.Server, srv *greeterServer) {
	c := time.NewTicker(5 * time.Second) // Adjust/make configurable as needed
	defer c.Stop()

	for {
		select {
		case <-c.C:
			_ = srv // This is where you check the server health with custom things like DB pings
			if true {
				hs.SetServingStatus("greeter-server", healthpb.HealthCheckResponse_SERVING)
			} else {
				hs.SetServingStatus("greeter-server", healthpb.HealthCheckResponse_NOT_SERVING)
			}
		case <-ctx.Done():
			break
		}
	}
}
