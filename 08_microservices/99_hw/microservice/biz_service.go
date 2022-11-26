package main

import (
	context "context"
	"fmt"
)

type Biz struct {
	UnimplementedBizServer
	ServiceName string
}

func NewBiz() *Biz {
	return &Biz{ServiceName: "Some buisness logic"}
}

func (biz *Biz) Check(ctx context.Context, nothing *Nothing) (*Nothing, error) {
	fmt.Println("Biz.Check logic")
	return nil, nil
}

func (biz *Biz) Add(ctx context.Context, nothing *Nothing) (*Nothing, error) {
	fmt.Println("Biz.Add logic")
	return nil, nil
}

func (biz *Biz) Test(ctx context.Context, nothing *Nothing) (*Nothing, error) {
	fmt.Println("Biz.Test logic")
	return nil, nil
}
