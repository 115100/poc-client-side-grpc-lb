package main

import (
	"context"
	"log"
	"time"

	"github.com/115100/poc-client-side-grpc-lb/go/greeterpb"
	"google.golang.org/grpc"
)

func main() {
	if err := run(); err != nil {
		log.Fatalln(err)
	}
}

func run() error {
	ctx := context.Background()

	// Requires grpc >= 1.14.0 for proper DNS resolution behaviour.
	cc, err := grpc.DialContext(
		ctx,
		// headless service exposed here; always in fmt dns:///<svc>.<namespace>.svc.cluster.local:<port>
		"dns:///greeter-server.default.svc.cluster.local:8080",
		// round_robin is a built-in balancer
		grpc.WithBalancerName("round_robin"),
		grpc.WithInsecure(),
	)
	if err != nil {
		return err
	}
	defer cc.Close()

	log.Println("client: Starting to ping server")
	client := greeterpb.NewGreeterClient(cc)
	for {
		resp, err := client.Greet(ctx, &greeterpb.GreetRequest{Name: "ping"})
		if err != nil {
			return err
		}

		log.Printf("client: got response from server: \"%s\"\n", resp.Message)
		time.Sleep(time.Second)
	}

	return nil
}
