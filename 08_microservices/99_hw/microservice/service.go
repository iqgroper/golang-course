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
	// grpc.UnaryInterceptor(authInterceptor),
	// grpc.InTapHandle(rateLimiter),
	)

	// RegisterAdminServer(server, )
	RegisterBizServer(server, NewBiz())

	fmt.Println("starting server at", listenAddr)
	server.Serve(lis)
	return nil
}

type Biz struct {
	UnimplementedBizServer
	ServiceName string
}

func NewBiz() *Biz {
	return &Biz{ServiceName: "Some buisness logic"}
}

func (biz *Biz) Check(context.Context, *Nothing) (*Nothing, error) {
	fmt.Println("Biz.Check logic")
	return nil, nil
}

func (biz *Biz) Add(context.Context, *Nothing) (*Nothing, error) {
	fmt.Println("Biz.Add logic")
	return nil, nil
}

func (biz *Biz) Test(context.Context, *Nothing) (*Nothing, error) {
	fmt.Println("Biz.Test logic")
	return nil, nil
}
