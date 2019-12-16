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
	conf, _ := config.ReadConfigurationV2(ctx)
	//fmt.Println(conf)
	// get new IAM token
	token.GetIAMToken(&conf)
	//---
	vm01 := instance.New(&conf)
	vm01.Get(ctx, "ef3cbgorepe1bquo0efr")
	snap01 := snapshot.New(&conf)
	log.Println(snap01.Name)
	snap01.List(ctx)
}
