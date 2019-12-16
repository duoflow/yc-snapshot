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
	conf, vms, _ := config.ReadConfiguration(ctx)
	//fmt.Println(conf)
	// get new IAM token
	token.GetIAMToken(&conf)
	//---
	log.Printf("%#v", vms)
	vm01 := instance.New(&conf)
	vm01.Get(ctx, "ef3cbgorepe1bquo0efr")
	snap01 := snapshot.New(&conf)
	log.Println(snap01.Name)
	snap01.List(ctx)
}
