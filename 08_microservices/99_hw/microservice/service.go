package main

import (
	context "context"
	"fmt"
	"log"
	"net"
	"time"

	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func LoggingInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	start := time.Now()

	md, _ := metadata.FromIncomingContext(ctx)

	reply, err := handler(ctx, req)

	fmt.Printf(`--
	after incoming call=%v
	req=%#v
	reply=%#v
	time=%v
	md=%v
	err=%v
`, info.FullMethod, req, reply, time.Since(start), md, err)

	return reply, err
}

func StartMyMicroservice(ctx context.Context, listenAddr, ACLData string) error {

	lis, err := net.Listen("tcp", listenAddr)
	if err != nil {
		log.Fatalln("cant listen port", err)
	}

	server := grpc.NewServer(
	// grpc.UnaryInterceptor(statistics),
	// grpc.StreamServerInterceptor(logging),
	// grpc.UnaryInterceptor(logging)
	)

	RegisterBizServer(server, NewBiz())
	RegisterAdminServer(server, NewAdmin())

	go func() {
		fmt.Println("starting server at", listenAddr)
		if err := server.Serve(lis); err != nil {
			log.Fatalf("server.Serve error: %s", err.Error())
		}
	}()

	go func() {
		<-ctx.Done()
		server.GracefulStop()
	}()

	return nil
}
