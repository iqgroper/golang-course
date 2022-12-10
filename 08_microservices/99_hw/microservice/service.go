package main

import (
	context "context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strings"

	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	status "google.golang.org/grpc/status"
)

type ACL struct {
	Logger    []string
	Stat      []string
	Biz_user  []string
	Biz_admin []string
}

func bizInterceptor(acl map[string][]string) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {

		err := checkACL(acl, info.FullMethod, ctx)
		if err != nil {
			return nil, err
		}

		ctx = context.WithValue(ctx, "method", info.FullMethod)

		reply, err := handler(ctx, req)

		return reply, err
	}
}

func adminInterceptor(acl map[string][]string) grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		stream grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		err := checkACL(acl, info.FullMethod, stream.Context())
		if err != nil {
			return err
		}
		handler(srv, stream)
		return nil
	}
}

func checkACL(acl map[string][]string, method string, ctx context.Context) error {
	md, _ := metadata.FromIncomingContext(ctx)
	consumer := md["consumer"]
	if len(consumer) == 0 {
		return status.Error(codes.Unauthenticated, "no consumer in metadata")
	}

	if _, ok := acl[consumer[0]]; !ok {
		return status.Error(codes.Unauthenticated, "not allowed")
	}

	allowIn := false
	for _, allowedMethod := range acl[consumer[0]] {

		if strings.Contains(allowedMethod, "*") {
			if strings.HasPrefix(method, strings.TrimRight(allowedMethod, "*")) {
				allowIn = true
				break
			}

		} else {
			if method == allowedMethod {
				allowIn = true
				break
			}
		}

	}
	if !allowIn {
		return status.Error(codes.Unauthenticated, "not allowed")
	}

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
		grpc.UnaryInterceptor(bizInterceptor(acl)),
		grpc.StreamInterceptor(adminInterceptor(acl)),
	)

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
