package main

import (
	"context"
	"log"
	"time"

	"github.com/115100/poc-client-side-grpc-lb/greeterpb"
	"google.golang.org/grpc"
)

func main() {
	if err := run(); err != nil {
		log.Fatalln(err)
	}
}

func run() error {
	ctx := context.Background()

	// TODO: As of 2018-07-31, due to
	// https://github.com/grpc/grpc-go/pull/2201, this is bugged if the
	// dns resolver returns an empty set and will not fetch DNS records
	// for another 30 minutes (defaultFreq). Until the next tagged release,
	// you should use master branch.
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
