package main

import (
	context "context"
	"fmt"
	"log"
	"net"

	grpc "google.golang.org/grpc"
)

func StartMyMicroservice(ctx context.Context, listenAddr, ACLData string) error {
	lis, err := net.Listen("tcp", listenAddr)
	if err != nil {
		log.Fatalln("cant listet port", err)
	}

	server := grpc.NewServer(
	// grpc.UnaryInterceptor(statistics),
	// grpc.UnaryInterceptor(logging)
	)

	RegisterBizServer(server, NewBiz())
	RegisterAdminServer(server, NewAdmin())

	fmt.Println("starting server at", listenAddr)
	server.Serve(lis)
	return nil
}
