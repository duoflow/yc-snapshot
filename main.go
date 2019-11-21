package main

import (
	"context"
	"log"

	"github.com/duoflow/yc-snapshot/config"
	"github.com/duoflow/yc-snapshot/instance"
	"github.com/duoflow/yc-snapshot/snapshot"
	"github.com/duoflow/yc-snapshot/token"
)

func main() {
	ctx := context.Background()
	conf, _ := config.ReadConfiguration(ctx)
	// get new IAM token
	conf.Token = token.GetIAMToken(&conf)
	//fmt.Println("New IAMToken: " + conf.Token)
	vm01 := instance.New(&conf)
	vm01.Get(ctx, "epdnvonl4vd6ik8ngood")
	vm01.Stop(ctx, "epdnvonl4vd6ik8ngood")
	//
	snap01 := snapshot.New(&conf)
	log.Println(snap01.Name)
	snap01.List(ctx)
}
