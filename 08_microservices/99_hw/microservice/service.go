package main

import (
	context "context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"time"

	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func BizLoggingInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	start := time.Now()

	md, _ := metadata.FromIncomingContext(ctx)
	ctx = context.WithValue(ctx, "method", info.FullMethod)

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

type ACL struct {
	Logger    []string
	Stat      []string
	Biz_user  []string
	Biz_admin []string
}

// type myStream struct {
// 	grpc.ServerStream
// 	method string
// }

// func (s *myStream) Context() context.Context {
// 	return context.WithValue(s.ServerStream.Context(), "method", s.method)
// }

func AdminLoggingInterceptor(
	srv interface{},
	stream grpc.ServerStream,
	info *grpc.StreamServerInfo,
	handler grpc.StreamHandler) error {

	// newStream := &myStream{
	// 	method: info.FullMethod,
	// }

	handler(srv, stream)

	return nil
}

func StartMyMicroservice(ctx context.Context, listenAddr, ACLData string) error {

	var acl map[string][]string
	errorMarshal := json.Unmarshal([]byte(ACLData), &acl)
	if errorMarshal != nil {
		return errorMarshal
	}

	lis, err := net.Listen("tcp", listenAddr)
	if err != nil {
		log.Fatalln("cant listen port", err)
	}

	server := grpc.NewServer(
		grpc.UnaryInterceptor(BizLoggingInterceptor),
		grpc.StreamInterceptor(AdminLoggingInterceptor),
	)

	// logs := &[]Event

	logs := make(chan *Event, 2)
	RegisterBizServer(server, NewBiz(acl, logs))
	RegisterAdminServer(server, NewAdmin(ctx, acl, logs))

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
