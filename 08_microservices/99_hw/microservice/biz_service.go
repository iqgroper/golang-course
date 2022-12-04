package main

import (
	context "context"
	"fmt"
	"strings"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	status "google.golang.org/grpc/status"
)

type Biz struct {
	UnimplementedBizServer
	Logs        chan *Event
	Stats       chan *Stat
	ACL         map[string][]string
	Host        string
	ServiceName string
}

func NewBiz(host string, acl map[string][]string, logs chan *Event) *Biz {
	return &Biz{ServiceName: "Some buisness logic", Host: host, ACL: acl, Logs: logs}
}

func (biz *Biz) sendLogs(ctx context.Context) {

	md, _ := metadata.FromIncomingContext(ctx)
	fmt.Println("CONSUMER", md["consumer"][0])

	if method, ok := ctx.Value("method").(string); ok {
		event := Event{
			Timestamp: time.Now().Unix(),
			Consumer:  md["consumer"][0],
			// Consumer: "consumer",
			Method: method,
			Host:   biz.Host,
		}
		biz.Logs <- &event
	}
}

func (biz *Biz) checkACL(ctx context.Context) error {
	md, _ := metadata.FromIncomingContext(ctx)
	consumer := md["consumer"]
	if len(consumer) == 0 {
		return status.Error(codes.Unauthenticated, "no consumer in metadata")
	}

	if _, ok := biz.ACL[consumer[0]]; !ok {
		return status.Error(codes.Unauthenticated, "not allowed")
	}

	allowIn := false
	for _, allowedMethod := range biz.ACL[consumer[0]] {
		allowIn = true
		if strings.Contains("*", allowedMethod) {
			// если содержится звездочка - проверяем вхождение подстроки в урле, иначе можно просто втупую сравнить
			allowIn = true
			break
		}

	}
	if !allowIn {
		return status.Error(codes.Unauthenticated, "not allowed")
	}

	return nil
}

func (biz *Biz) Check(ctx context.Context, nothing *Nothing) (*Nothing, error) {

	err := biz.checkACL(ctx)
	if err != nil {
		return nil, err
	}
	// fmt.Println("acl in test", biz.ACL)

	biz.sendLogs(ctx)

	fmt.Println("Biz.Check logic")
	return &Nothing{Dummy: true}, nil
}

func (biz *Biz) Add(ctx context.Context, nothing *Nothing) (*Nothing, error) {

	err := biz.checkACL(ctx)
	if err != nil {
		return nil, err
	}

	biz.sendLogs(ctx)

	fmt.Println("Biz.Add logic")
	return &Nothing{Dummy: true}, nil
}

func (biz *Biz) Test(ctx context.Context, nothing *Nothing) (*Nothing, error) {

	err := biz.checkACL(ctx)
	if err != nil {
		return nil, err
	}

	biz.sendLogs(ctx)

	fmt.Println("Biz.Test logic")
	return &Nothing{Dummy: true}, nil
}
