package main

import (
	context "context"
	"fmt"
	"strings"
	"time"

	"google.golang.org/grpc/metadata"
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
		event := &Event{
			Timestamp: time.Now().Unix(),
			Consumer:  md["consumer"][0],
			Method:    method,
			Host:      md[":authority"][0][:strings.IndexByte(md[":authority"][0], ':')+1],
		}
		biz.Logs <- event
	}
}

func (biz *Biz) Check(ctx context.Context, nothing *Nothing) (*Nothing, error) {

	biz.sendLogs(ctx)

	fmt.Println("Biz.Check logic")
	return &Nothing{Dummy: true}, nil
}

func (biz *Biz) Add(ctx context.Context, nothing *Nothing) (*Nothing, error) {

	biz.sendLogs(ctx)

	fmt.Println("Biz.Add logic")
	return &Nothing{Dummy: true}, nil
}

func (biz *Biz) Test(ctx context.Context, nothing *Nothing) (*Nothing, error) {

	biz.sendLogs(ctx)

	fmt.Println("Biz.Test logic")
	return &Nothing{Dummy: true}, nil
}
