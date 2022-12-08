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

func NewBiz(acl map[string][]string, logs chan *Event) *Biz {
	return &Biz{ServiceName: "Some buisness logic", ACL: acl, Logs: logs}
}

func (biz *Biz) sendLogs(ctx context.Context) {

	md, _ := metadata.FromIncomingContext(ctx)

	// host := strings.Split(md[":authority"][0], ":")[0]

	if method, ok := ctx.Value("method").(string); ok {
		event := Event{
			Timestamp: time.Now().Unix(),
			Consumer:  md["consumer"][0],
			Method:    method,
			Host:      md[":authority"][0][:strings.IndexByte(md[":authority"][0], ':')+1],
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

	if method, ok := ctx.Value("method").(string); ok {

		allowIn := false
		for _, allowedMethod := range biz.ACL[consumer[0]] {

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
	}

	return nil
}

func (biz *Biz) Check(ctx context.Context, nothing *Nothing) (*Nothing, error) {

	err := biz.checkACL(ctx)
	if err != nil {
		return nil, err
	}

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
